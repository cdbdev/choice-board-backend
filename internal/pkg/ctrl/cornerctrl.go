package ctrl

//~ ---> TODO add following:
		//~ http://localhost:3500/corner( delete, put)

import (
	"corner-backend/internal/pkg/dao"
	"encoding/json"
	"net/http"
	"bytes"
	"database/sql"
	"log"
	"regexp"
)

var (
	getAllCornersRe = regexp.MustCompile(`^\/cornerData*$`)
	CornerRe = regexp.MustCompile(`^\/corner*$`)
	CornerVisibilityRe = regexp.MustCompile(`^\/cornerVisibility*$`)
	ChartRe = regexp.MustCompile(`^\/chart*$`)
	ChartOverviewRe = regexp.MustCompile(`^\/chartOverview*$`)
)

type CornerController struct{
	Db *sql.DB
	Logger *log.Logger 
}

type CollectedData struct {
	Corners		[]dao.CornersData	`json:"hoeken"`
	Children	[]dao.ChildData		`json:"kleuters"`
}

type ResultsData struct {
	Results CollectedData		`json:"results"`
}

type Chart struct {
	Corner 		int				`json:"corner"`
	CornerName	string			`json:"cornerName"`
	Count		int				`json:"count"`
	Percent		float64 		`json:"percent"`
}

func (c CornerController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE");
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	
	w.Header().Set("Content-Type", "application/json")
	
	switch{
		case r.Method == http.MethodGet && getAllCornersRe.MatchString(r.URL.Path):
			c.GetAllData(w,r)
			return
		case r.Method == http.MethodPost && CornerRe.MatchString(r.URL.Path):
			c.AddCorner(w,r)
			return
		case r.Method == http.MethodPut && CornerRe.MatchString(r.URL.Path):
			c.ChangeCorner(w,r)
			return
		case r.Method == http.MethodDelete && CornerRe.MatchString(r.URL.Path):
			c.RemoveCorner(w,r)
			return
		case r.Method == http.MethodPost && CornerVisibilityRe.MatchString(r.URL.Path):
			c.UpdateCornerVisibility(w,r)
			return
		case r.Method == http.MethodPost && ChartRe.MatchString(r.URL.Path):
			c.ChartDetails(w,r)
			return
		case r.Method == http.MethodGet && ChartOverviewRe.MatchString(r.URL.Path):
			c.ChartOverview(w,r)
			return
		case r.Method == "OPTIONS":
			w.WriteHeader(http.StatusOK)
			return
		default: 
			notFound(w,r)
			return
	}
}

func (c CornerController) GetAllData(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
    //w.WriteHeader(http.StatusCreated)
	
	// Get DAO's
	cornerdao := dao.CornerDao{Db: c.Db, Logger: c.Logger}
	childdao := dao.ChildDao{Db: c.Db, Logger: c.Logger}
	
	// Get all corners
	cornersData := cornerdao.FetchAllCorners()
	
	/// get all children
	childData := childdao.FetchAllChildren()
		
	resultsData := &ResultsData{ Results: CollectedData{ Corners: cornersData, Children: childData } }
	
	data, err := json.Marshal(resultsData)
	if err != nil {
		c.Logger.Println(err)
		panic(err)
	}
	
	w.Write(data)
}

func (c CornerController) AddCorner(w http.ResponseWriter, r *http.Request) {
	var cornersDataItems dao.CornersDataItems
	
	cornerdao := dao.CornerDao{Db: c.Db, Logger: c.Logger}
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cornersDataItems)
	if err != nil {
		c.Logger.Println(err)
		panic(err)
	}
	
	data := GenerateResult(cornerdao.InsertCorner(cornersDataItems) )	
	w.Write(data)
}

func (c CornerController) ChangeCorner(w http.ResponseWriter, r *http.Request) {
	var cornersDataItems dao.CornersDataItems
	
	cornerdao := dao.CornerDao{Db: c.Db, Logger: c.Logger}
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cornersDataItems)
	if err != nil {
		c.Logger.Println(err)
		panic(err)
	}
	
	data := GenerateResult(cornerdao.UpdateCorner(cornersDataItems) )	
	w.Write(data)
}

func (c CornerController) RemoveCorner(w http.ResponseWriter, r *http.Request) {
	var cornersDataItems dao.CornersDataItems
	
	cornerdao := dao.CornerDao{Db: c.Db, Logger: c.Logger}
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cornersDataItems)
	if err != nil {
		c.Logger.Println(err)
		panic(err)
	}
	
	data := GenerateResult(cornerdao.DeleteCorner(cornersDataItems) )	
	w.Write(data)
}

func (c CornerController) UpdateCornerVisibility(w http.ResponseWriter, r *http.Request) {
	var cornerVisibilityItems dao.CornerVisibilityItems
	
	cornerdao := dao.CornerDao{Db: c.Db, Logger: c.Logger}
	
	buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
    newStr := buf.String()
	
	err := json.Unmarshal([]byte(newStr), &cornerVisibilityItems)
	if err != nil {
		c.Logger.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	data := GenerateResult(cornerdao.UpdateVisibility(cornerVisibilityItems) )	
	w.Write(data)
}

func (c CornerController) ChartDetails(w http.ResponseWriter, r *http.Request) {	
	// example input: {"id":1,"firstname":"Fran","lastname":"De Boeck","fullName":"Fran De Boeck","avatar":"/srv/static/children/fran.jpg","position":1}
	// example output: {corner: act.corner, cornerName: act.name, count: 1, percent: 0}
	
	var childDataItem dao.ChildDataItem
	var charts []Chart
	var charts_tmp []Chart
	var b []byte
	
	childdao := dao.ChildDao{Db: c.Db, Logger: c.Logger}
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&childDataItem)
	if err != nil {
		c.Logger.Println(err)
		panic(err)
	}
	
	// result of all existing corners where id <> 1
	childActivities := childdao.GetAllChildActivities(childDataItem)
	sumOfActivities := 0.0
	
	for i := 0; i < len(childActivities); i++ {
		childActivity := childActivities[i]
		
		if childActivity.Corner == 1 {
			continue
		}
		
		sumOfActivities += 1
		
		// Collect activities and do the count
		currentChartIndex := -1
		
		for i := range charts_tmp {
			if charts_tmp[i].Corner == childActivity.Corner {
				currentChartIndex = i
				break
			}
		}
		
		// Search existing chartdata 
		if currentChartIndex == -1 {			
			newChartData := Chart{childActivity.Corner,childActivity.Name, 1, 0}
			charts_tmp = append(charts_tmp, newChartData)
		} else {
			// add 1 to existing count of found records
			charts_tmp[currentChartIndex].Count = charts_tmp[currentChartIndex].Count + 1
		}
	}
	
	// Calculate percentages
	for i := 0; i < len(charts_tmp); i++ {
		currentChart := charts_tmp[i]
		
		percent := float64(100 / ( sumOfActivities / float64(currentChart.Count) ))
		currentChart.Percent = percent;
		charts = append(charts, currentChart)
	}
	
	if len(charts) > 0 {
		b, err = json.Marshal(charts)
		if err != nil {
			c.Logger.Println(err)
			panic(err)
		} 
	} else {
		b, err = json.Marshal([]Chart{})
		if err != nil {
			c.Logger.Println(err)
			panic(err)
		} 
	}
	
	w.Write(b)	
}

func (c CornerController) ChartOverview(w http.ResponseWriter, r *http.Request) {
	var charts []Chart
	var charts_tmp []Chart
	
	childdao := dao.ChildDao{Db: c.Db, Logger: c.Logger}	
	
	// result of all existing corners where id <> 1
	childActivities := childdao.GetAllActivities()
	sumOfActivities := 0.0
	
	for i := 0; i < len(childActivities); i++ {
		childActivity := childActivities[i]
		
		if childActivity.Corner == 1 {
			continue
		}
		
		sumOfActivities += 1
		
		// Collect activities and do the count
		currentChartIndex := -1
		
		for i := range charts_tmp {
			if charts_tmp[i].Corner == childActivity.Corner {
				currentChartIndex = i
				break
			}
		}
		
		// Search existing chartdata 
		if currentChartIndex == -1 {			
			newChartData := Chart{childActivity.Corner,childActivity.Name, 1, 0}
			charts_tmp = append(charts_tmp, newChartData)
		} else {
			// add 1 to existing count of found records
			charts_tmp[currentChartIndex].Count = charts_tmp[currentChartIndex].Count + 1
		}
	}
	
	// Calculate percentages
	for i := 0; i < len(charts_tmp); i++ {
		currentChart := charts_tmp[i]
		
		percent := float64(100 / ( sumOfActivities / float64(currentChart.Count) ))
		currentChart.Percent = percent;
		charts = append(charts, currentChart)
	}
	
	if len(charts) > 0 {
		b, err := json.Marshal(charts)
		if err != nil {
			c.Logger.Println(err)
			panic(err)
		} 
		w.Write(b)
	} else {
		b, err := json.Marshal([]Chart{})
		if err != nil {
			c.Logger.Println(err)
			panic(err)
		} 
		w.Write(b)
	}		
}