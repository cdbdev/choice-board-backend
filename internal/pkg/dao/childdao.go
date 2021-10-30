package dao

import (
	"database/sql"
	"time"
	"log"
	//~ "fmt"
	//~ "os"
	
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type ChildDao struct{
	Db *sql.DB
	Logger *log.Logger 
}

type ChildData struct {
	Id 			int			`json:"id"`
	Lastname 	string		`json:"lastname"`
	Firstname 	string		`json:"firstname"`
	Avatar		string		`json:"avatar"`
	Position	int			`json:"position"`
}

type ChildDataItem struct {
	Id 			int			`json:"id"`
	Lastname 	string		`json:"lastname"`
	Firstname 	string		`json:"firstname"`
	Fullname	string		`json:"fullName"`
	Avatar		string		`json:"avatar"`
	Position	int			`json:"position"`
}

type ChildDataItems struct {
	Arr []ChildDataItem	`json:"arr"`
}

type ChildActivity struct {
	Child		int
	Timestamp   string
	Corner		int
	Name		string
}

func (childDao ChildDao) FetchAllChildren() []ChildData {
	children := []ChildData{}
	
	row, err := childDao.Db.Query("SELECT * FROM child")
	
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	defer row.Close()
	
	for row.Next() {
		var childdata ChildData
		row.Scan(&childdata.Id, &childdata.Lastname, &childdata.Firstname, &childdata.Avatar, &childdata.Position)
		
		children = append(children, childdata)
	}
	
	return children
}

func (childDao ChildDao) InsertChild(childItem ChildDataItem) string {
	tx, err := childDao.Db.Begin() // BEGIN TRANSACTION
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	stmt, err := tx.Prepare("INSERT INTO child VALUES(?,?,?,?,?)")
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	_, err = stmt.Exec(childItem.Id, childItem.Lastname, childItem.Firstname, childItem.Avatar, childItem.Position)
	if err != nil {
		_ = tx.Rollback()
		childDao.Logger.Println(err)
		panic(err)
	}
	
	tx.Commit()
	
	return "ok";
}

func (childDao ChildDao) DeleteChild(childItem ChildDataItem) string {
	tx, err := childDao.Db.Begin() // BEGIN TRANSACTION
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	stmt, err := tx.Prepare("DELETE FROM child WHERE id = ?")
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	_, err = stmt.Exec(childItem.Id)
	if err != nil {
		_ = tx.Rollback()
		childDao.Logger.Println(err)
		panic(err)
	}
	
	// Also delete activities
	stmt, err = tx.Prepare("DELETE FROM childactivity WHERE child = ?")
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	_, err = stmt.Exec(childItem.Id)
	if err != nil {
		_ = tx.Rollback()
		childDao.Logger.Println(err)
		panic(err)
	}
	
	
	tx.Commit()
	
	return "ok";
}

func (childDao ChildDao) UpdateChild(childItem ChildDataItem) string {
	tx, err := childDao.Db.Begin() // BEGIN TRANSACTION
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	stmt, err := tx.Prepare("UPDATE child SET lastname = ?, firstname = ?, avatar = ?, position = ? WHERE id = ?")
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	_, err = stmt.Exec(childItem.Lastname, childItem.Firstname, childItem.Avatar, childItem.Position, childItem.Id)
	if err != nil {
		_ = tx.Rollback()
		childDao.Logger.Println(err)
		panic(err)
	}
	
	tx.Commit()
	
	return "ok";
}

func (childDao ChildDao) UpdateChildposition(childItemsArr ChildDataItems) string {
	tx, err := childDao.Db.Begin() // BEGIN TRANSACTION
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	childItems := childItemsArr.Arr
	
	for i := 0; i < len(childItems); i++ {
		childItem := childItems[i]
		
		// Update position
		stmt, err := tx.Prepare("UPDATE child SET position = ? WHERE id = ?")
		if err != nil {
			childDao.Logger.Println(err)
			panic(err)
		}
		
		_, err = stmt.Exec(childItem.Position, childItem.Id)
		if err != nil {
			_ = tx.Rollback()
			childDao.Logger.Println(err)
			panic(err)
		}
		
		// Insert activity
		timestamp := getDefaultDateFormat()
		
		stmt, err = tx.Prepare("INSERT INTO childactivity VALUES(?,?,?)")
		if err != nil {
			childDao.Logger.Println(err)
			panic(err)
		}
		
		_, err = stmt.Exec(childItem.Id, timestamp, childItem.Position)
		if err != nil {
			_ = tx.Rollback()
			childDao.Logger.Println(err)
			panic(err)
		}
	}
	
	tx.Commit()
	
	return "ok";
}

func (childDao ChildDao) GetAllChildActivities(childItem ChildDataItem) []ChildActivity {
	var childActivities []ChildActivity
	
	row, err := childDao.Db.Query("SELECT ch.*, co.name FROM childactivity as ch inner join corner as co on co.id = ch.corner where ch.child = ?", childItem.Id)
	
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	defer row.Close()
	
	for row.Next() {
		var childactivity ChildActivity
		row.Scan(&childactivity.Child, &childactivity.Timestamp, &childactivity.Corner, &childactivity.Name)
		
		childActivities = append(childActivities, childactivity)
	}
	
	return childActivities
}

func (childDao ChildDao) DeleteAllChildren() string {
	tx, err := childDao.Db.Begin() // BEGIN TRANSACTION
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	stmt, err := tx.Prepare("DELETE FROM child")
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	_, err = stmt.Exec()
	if err != nil {
		_ = tx.Rollback()
		childDao.Logger.Println(err)
		panic(err)
	}
	
	tx.Commit()
	
	return "ok";
}

func (childDao ChildDao) DeleteAllActivities() string {
	tx, err := childDao.Db.Begin() // BEGIN TRANSACTION
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	stmt, err := tx.Prepare("DELETE FROM childactivity")
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	_, err = stmt.Exec()
	if err != nil {
		_ = tx.Rollback()
		childDao.Logger.Println(err)
		panic(err)
	}
	
	tx.Commit()
	
	return "ok";
}

func (childDao ChildDao) GetAllActivities() []ChildActivity {
	var childActivities []ChildActivity
	
	row, err := childDao.Db.Query("SELECT ch.*, co.name FROM childactivity as ch inner join corner as co on co.id = ch.corner")
	
	if err != nil {
		childDao.Logger.Println(err)
		panic(err)
	}
	
	defer row.Close()
	
	for row.Next() {
		var childactivity ChildActivity
		row.Scan(&childactivity.Child, &childactivity.Timestamp, &childactivity.Corner, &childactivity.Name)
		
		childActivities = append(childActivities, childactivity)
	}
	
	return childActivities
}

func getDefaultDateFormat() string{
	currentTime := time.Now()	
	timestamp := currentTime.Format("2006-02-01 15:04:05")
	return timestamp
}