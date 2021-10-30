package dao

import (
	"database/sql"
	"log"
	//~ "fmt"
	//~ "os"
	
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type CornerDao struct{
	Db *sql.DB
	Logger *log.Logger 
}

type CornersData struct {
	Id 			int			`json:"id"`
	Name 		string		`json:"name"`
	Avatar 		string		`json:"avatar"`
	Visible		int			`json:"visible"`
	Count		int			`json:"count"`
}

type CornersItem struct {
	Active		bool			`json:"active"`
	Title		string			`json:"title"`
	Avatar		string			`json:"avatar"`
}

type CornersDataItems struct {
	Id 			int				`json:"id"`
	Name 		string			`json:"name"`
	Avatar 		string			`json:"avatar"`
	Visible		bool			`json:"visible"`
	Count		int			`json:"count"`
	Items		[]CornersItem	`json:"items"`
}

type CornerVisibilityItems struct {
	Arr []CornersDataItems	`json:"arr"`
}

func (cornerDao CornerDao) FetchAllCorners() []CornersData {
	corners := []CornersData{}
	
	//~ path, err := os.Getwd()
	//~ if err != nil {
		//~ fmt.Println(err)
	//~ }
	//~ fmt.Println(path)
	
	row, err := cornerDao.Db.Query("SELECT * FROM corner")
	
	if err != nil {
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	defer row.Close()
	
	for row.Next() {
		var cdata CornersData
		row.Scan(&cdata.Id, &cdata.Name, &cdata.Avatar, &cdata.Visible, &cdata.Count)
		
		corners = append(corners, cdata)
	}
	
	return corners
}

func (cornerDao CornerDao) UpdateVisibility(cornerVisibilityItems CornerVisibilityItems) string {
	var visible int
	
	// Persist data	
	//cornerDao.Db.Serialize(func() {
	tx, err := cornerDao.Db.Begin() // BEGIN TRANSACTION
	if err != nil {
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	for _, corner := range cornerVisibilityItems.Arr {
		if corner.Visible == true {
			visible = 1
		} else {
			visible = 0
		}
		
		_, execErr := tx.Exec(`UPDATE corner SET visible = ? WHERE id = ?`, visible, corner.Id)
		if execErr != nil {
			_ = tx.Rollback()
			cornerDao.Logger.Println(err)
			panic(execErr)
		}
	}
	
	if err := tx.Commit(); err != nil {
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	return "ok"
}

func (cornerDao CornerDao) InsertCorner(cornersDataItem CornersDataItems) string {
	tx, err := cornerDao.Db.Begin() // BEGIN TRANSACTION
	if err != nil {
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	stmt, err := tx.Prepare("INSERT INTO corner VALUES(?,?,?,?,?)")
	if err != nil {
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	_, err = stmt.Exec(cornersDataItem.Id, cornersDataItem.Name, cornersDataItem.Avatar, cornersDataItem.Visible, cornersDataItem.Count)
	if err != nil {
		_ = tx.Rollback()
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	tx.Commit()
	
	return "ok";
}

func (cornerDao CornerDao) DeleteCorner(cornersDataItem CornersDataItems) string {
	tx, err := cornerDao.Db.Begin() // BEGIN TRANSACTION
	if err != nil {
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	stmt, err := tx.Prepare("DELETE FROM corner WHERE id = ?")
	if err != nil {
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	_, err = stmt.Exec(cornersDataItem.Id)
	if err != nil {
		_ = tx.Rollback()
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	tx.Commit()
	
	return "ok";
}

func (cornerDao CornerDao) UpdateCorner(cornersDataItem CornersDataItems) string {
	tx, err := cornerDao.Db.Begin() // BEGIN TRANSACTION
	if err != nil {
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	stmt, err := tx.Prepare("UPDATE corner SET avatar = ?, count = ? WHERE id = ?")
	if err != nil {
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	_, err = stmt.Exec(cornersDataItem.Avatar, cornersDataItem.Count, cornersDataItem.Id)
	if err != nil {
		_ = tx.Rollback()
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	tx.Commit()
	
	return "ok";
}

func (cornerDao CornerDao) DeleteAllCorners() string {
	tx, err := cornerDao.Db.Begin() // BEGIN TRANSACTION
	if err != nil {
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	stmt, err := tx.Prepare("DELETE FROM corner WHERE id > 1")
	if err != nil {
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	_, err = stmt.Exec()
	if err != nil {
		_ = tx.Rollback()
		cornerDao.Logger.Println(err)
		panic(err)
	}
	
	tx.Commit()
	
	return "ok";
}