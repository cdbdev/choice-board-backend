package ctrl

import (
	"corner-backend/internal/pkg/dao"
	"encoding/json"
	"net/http"
	"database/sql"
	"log"
	"regexp"
)

var (
	//~ AddChildRe = regexp.MustCompile(`^\/child*$`)
	//~ RemoveChildRe = regexp.MustCompile(`^\/child*$`)
	//~ ChangeChildRe = regexp.MustCompile(`^\/child*$`)
	ChildPosRe = regexp.MustCompile(`^\/child*$`)
	ChangeChildPosRe = regexp.MustCompile(`^\/childposition*$`)
)

type ChildController struct{
	Db *sql.DB
	Logger *log.Logger 
}

func (c ChildController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE");
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	
	w.Header().Set("Content-Type", "application/json")
	
	switch{
		case r.Method == http.MethodPost && ChildPosRe.MatchString(r.URL.Path):
			c.AddChild(w,r)
			return
		case r.Method == http.MethodDelete && ChildPosRe.MatchString(r.URL.Path):
			c.RemoveChild(w,r)
			return
		case r.Method == http.MethodPut && ChildPosRe.MatchString(r.URL.Path):
			c.ChangeChild(w,r)
			return
		case r.Method == http.MethodPost && ChangeChildPosRe.MatchString(r.URL.Path):
			c.ChangeChildposition(w,r)
			return
		case r.Method == "OPTIONS":
			w.WriteHeader(http.StatusOK)
			return
		default: 
			notFound(w,r)
			return
	}
}

func (c ChildController) AddChild(w http.ResponseWriter, r *http.Request) {
	// Example input: {   "id": 4,   "lastname": "123",   "firstname": "Test ",   "fullName": "Test  123",   "avatar": "/srv/static/doc-images/lists/qm_face.jpg",   "position": 1  }
	var childDataItem dao.ChildDataItem
	
	childdao := dao.ChildDao{Db: c.Db, Logger: c.Logger}
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&childDataItem)
	if err != nil {
		c.Logger.Println(err)
		panic(err)
	}
	
	data := GenerateResult(childdao.InsertChild(childDataItem) )	
	w.Write(data)
}

func (c ChildController) RemoveChild(w http.ResponseWriter, r *http.Request) {
	//~ // Example input: {   "id": 4,   "lastname": "123",   "firstname": "Test ",   "fullName": "Test  123",   "avatar": "/srv/static/doc-images/lists/qm_face.jpg",   "position": 1  }
	var childDataItem dao.ChildDataItem
	
	childdao := dao.ChildDao{Db: c.Db, Logger: c.Logger}
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&childDataItem)
	if err != nil {
		c.Logger.Println(err)
		panic(err)
	}
	
	data := GenerateResult(childdao.DeleteChild(childDataItem))	
	w.Write(data)
}

func (c ChildController) ChangeChild(w http.ResponseWriter, r *http.Request) {
	//~ // Example input: {   "id": 4,   "lastname": "123",   "firstname": "Test ",   "fullName": "Test  123",   "avatar": "/srv/static/doc-images/lists/qm_face.jpg",   "position": 1  }
	var childDataItem dao.ChildDataItem
	
	childdao := dao.ChildDao{Db: c.Db, Logger: c.Logger}
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&childDataItem)
	if err != nil {
		c.Logger.Println(err)
		panic(err)
	}
	
	data := GenerateResult(childdao.UpdateChild(childDataItem))	
	w.Write(data)
}

func (c ChildController) ChangeChildposition(w http.ResponseWriter, r *http.Request) {
	// Example input: {   "id": 4,   "lastname": "123",   "firstname": "Test ",   "fullName": "Test  123",   "avatar": "/srv/static/doc-images/lists/qm_face.jpg",   "position": 1  }
	var childDataItems dao.ChildDataItems
	
	childdao := dao.ChildDao{Db: c.Db, Logger: c.Logger}
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&childDataItems)
	if err != nil {
		c.Logger.Println(err)
		panic(err)
	}
	
	data := GenerateResult(childdao.UpdateChildposition(childDataItems))	
	w.Write(data)
}