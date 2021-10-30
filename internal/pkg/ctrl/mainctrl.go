package ctrl

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"database/sql"
	"bytes"
	"encoding/json"
	"encoding/base64"
	"os"
	"log"
	"corner-backend/internal/pkg/dao"
)

type ctxKey struct{}

type Cleanup struct {
	Children bool
	Corners bool
	Activities bool
}

type Image struct {
	Base64str	string
	File		string
}

type MainController struct{ 
	Db *sql.DB
	Logger *log.Logger
}

var (
	HomeRe = regexp.MustCompile(`^\/index.html`)
	AdminRe = regexp.MustCompile(`^\/admin.*`)
	CleanupRe = regexp.MustCompile(`^\/cleanup.*`)
	ImageRe = regexp.MustCompile(`^\/image.*`)
	FileServerRe = regexp.MustCompile(`^\/.*`)
)

func (c MainController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE");
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	
	w.Header().Set("Content-Type", "text/html")
	
	switch{
		case r.Method == http.MethodGet && HomeRe.MatchString(r.URL.Path):
			c.Home(w,r)
			return
		case r.Method == http.MethodGet && AdminRe.MatchString(r.URL.Path):
			c.Admin(w,r)
			return
		case r.Method == http.MethodPost && CleanupRe.MatchString(r.URL.Path):
			c.Cleanup(w,r)
			return
		case r.Method == http.MethodPost&& ImageRe.MatchString(r.URL.Path):
			c.ImageCreator(w,r)
			return			
		case r.Method == http.MethodGet && FileServerRe.MatchString(r.URL.Path):
			c.FileServer(w,r)
		case r.Method == "OPTIONS":
			w.WriteHeader(http.StatusOK)
			return
		default: 
			http.NotFound(w,r)
			return
	}
}

func (c MainController) Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web\\index.html")
}

func (c MainController) Admin(w http.ResponseWriter, r *http.Request) {
	slug := "/admin"
	regexpres := regexp.MustCompile("/admin/.*")	
	matches := regexpres.FindStringSubmatch(r.URL.Path)
	
	if len(matches) > 0 {
		ctx := context.WithValue(r.Context(), ctxKey{}, matches[0:])		
		slug = ctx.Value(ctxKey{}).([]string)[0]
	}
	
	if slug == "/admin" {				
		fmt.Fprintf(w, buildAdminPage(r.Host))
	} else if slug == "/admin/corner" {
		fmt.Fprintf(w, buildCornerQueryPage(c))
	} else if slug == "/admin/child" {		
		fmt.Fprintf(w, buildChildQueryPage(c))
	} else if slug == "/admin/activities" {
		fmt.Fprintf(w, buildActivitiesQueryPage(c))
	} else {
		http.NotFound(w,r)
	}
}

func (c MainController) Cleanup(w http.ResponseWriter, r *http.Request) {
	// input: {   "children": false,  "corners": false,  "activities": true  }
	var cleanup Cleanup
	
	cornerdao := dao.CornerDao{Db: c.Db}
	childdao := dao.ChildDao{Db: c.Db}
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cleanup)
	if err != nil {
		c.Logger.Println(err)
		panic(err)
	}
	
	// Check which option should be cleaned
	resp := ""
	if cleanup.Children {
		resp += "Status cleaning children [" + childdao.DeleteAllChildren() + "] "
	} 
	
	if cleanup.Corners {
		resp += "Status cleaning corners [" + cornerdao.DeleteAllCorners() + "] "
	} 
	
	if cleanup.Activities {
		resp += "Status cleaning activities [" + childdao.DeleteAllActivities() + "] "
	} 
	
	if resp == "" {
		resp += "Nothing to clean."
	}
	
	fmt.Fprintf(w, resp)
}

func (c MainController) ImageCreator(w http.ResponseWriter, r *http.Request) {
	// Example: {   "base64str":  "/9j/4AAQSkZJRgABAQEBLAEsAAD/4gJASUN...",   "file": "static/corners/heb-je-een-afspraak-met-fluvius.jpg"}
	var image Image
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&image)
	if err != nil {
		c.Logger.Println(err)
		panic(err)
	}
	
	// Decode base64
	base64dec, err := base64.StdEncoding.DecodeString(image.Base64str)
    if err != nil {
		c.Logger.Println(err)
        panic(err)
    }
	
	f, err := os.Create("web\\srv\\" + image.File)
    if err != nil {
		c.Logger.Println(err)
        panic(err)
    }
    defer f.Close()

    if _, err := f.Write(base64dec); err != nil {
		c.Logger.Println(err)
        panic(err)
    }
	
    if err := f.Sync(); err != nil {
		c.Logger.Println(err)
        panic(err)
    }
	
	fmt.Fprintf(w, "Picture saved")
}

func (c MainController) FileServer(w http.ResponseWriter, r *http.Request) {	
	regexpres := regexp.MustCompile("/.*")	
	matches := regexpres.FindStringSubmatch(r.URL.Path)
	ctx := context.WithValue(r.Context(), ctxKey{}, matches[0:])

	urlreq := ctx.Value(ctxKey{}).([]string)
	p := "web\\" + urlreq[0]
	http.ServeFile(w, r, p)
}

func buildAdminPage(host string) string {
	var html bytes.Buffer
		
	linkcorner := "http://" + host + "/admin/corner"
	linkchild := "http://" + host + "/admin/child"
	linkactivities := "http://" + host + "/admin/activities"
	
	html.WriteString(`<html><head></head><body><h1>Database info</h1><h3>Available tables</h3>`)
	urls := fmt.Sprintf("<a href='%s' target='_blank'>Corner</a></br><a href='%s' target='_blank'>Child</a></br><a href='%s' target='_blank'>Activities</a></br>", linkcorner, linkchild, linkactivities)
	html.WriteString(urls)
	
	// Get current database info
	mydir, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
    }
	
	dbinfo := fmt.Sprintf("</br></br><i>Running database from: <b>%s\\db\\corner_data.db</b></i>", mydir)
	html.WriteString(dbinfo)
	
	return html.String()
}

func buildCornerQueryPage(c MainController) string {
	var html bytes.Buffer
		
	cornerdao := dao.CornerDao{Db: c.Db}
	corners := cornerdao.FetchAllCorners()
	
	html.WriteString(`<html><head></head><body><table border = "1"><tr><th>ID</th><th>Name</th><th>Avatar</th><th>Visible</th></tr>`)
	
	for i := 0; i < len(corners); i++ {
		corner := corners[i]
		
		record := fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%d</td></tr>", 
								corner.Id, corner.Name, corner.Avatar, corner.Visible)
		html.WriteString(record)
	}
	
	html.WriteString(`</table></body></html>`)
	return html.String()
}

func buildChildQueryPage(c MainController) string {		
	var html bytes.Buffer
		
	childdao := dao.ChildDao{Db: c.Db}
	children := childdao.FetchAllChildren()
	
	html.WriteString(`<html><head></head><body><table border = "1"><tr><th>ID</th><th>Lastname</th><th>Firstname</th><th>Avatar</th><th>Position</th></tr>`)
	
	for i := 0; i < len(children); i++ {
		child := children[i]
		
		record := fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td></tr>", 
								child.Id, child.Lastname, child.Firstname, child.Avatar, child.Position)
		html.WriteString(record)
	}
	
	html.WriteString(`</table></body></html>`)
	return html.String()
}

func buildActivitiesQueryPage(c MainController) string {
	var html bytes.Buffer
		
	childdao := dao.ChildDao{Db: c.Db}
	activities := childdao.GetAllActivities()
	
	html.WriteString(`<html><head></head><body><table border = "1"><tr><th>Child</th><th>Timestamp</th><th>Corner</th><th>Corner name</th></tr>`)
	
	for i := 0; i < len(activities); i++ {
		activity := activities[i]
		
		record := fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%d</td><td>%s</td></tr>", 
								activity.Child, activity.Timestamp, activity.Corner, activity.Name)
		html.WriteString(record)
	}
	
	html.WriteString(`</table></body></html>`)
	return html.String()
}