package api

import (
	"fmt"
	"log"
	"os"
	"net"
	"strconv"
	"os/exec"
	"net/http"
	"corner-backend/internal/pkg/dao"
	"corner-backend/internal/pkg/ctrl"
)

const port = 3500

// Start API service
func Start() error {
	if status, err := Check(port); status == false {
		// Just launch the browser
		LaunchBrowser()
	
		// Print message and exit
		fmt.Println(err)
		return nil // do not return error, the app was probably already started
	}
	
	// Create log file, if the file doesn't exist, create it or append to the file
	logger := InitializeLog()
	
	// Initialize database
	db := dao.InitDB()
	
	defer dao.CloseDB(db)  //TODO: how to cleanly close Database...
	
	// Setup routing	
	mux := http.NewServeMux()
	
	mux.Handle("/", &ctrl.MainController{Db: db, Logger: logger})
	mux.Handle("/index.html", &ctrl.MainController{Db: db, Logger: logger})
	mux.Handle("/corner", &ctrl.CornerController{Db: db, Logger: logger})
	mux.Handle("/cornerData", &ctrl.CornerController{Db: db, Logger: logger})
	mux.Handle("/cornerVisibility", &ctrl.CornerController{Db: db, Logger: logger})
	mux.Handle("/chart", &ctrl.CornerController{Db: db, Logger: logger})
	mux.Handle("/chartOverview", &ctrl.CornerController{Db: db, Logger: logger})
	mux.Handle("/child", &ctrl.ChildController{Db: db, Logger: logger})
	mux.Handle("/childposition", &ctrl.ChildController{Db: db, Logger: logger})
	mux.Handle("/admin", &ctrl.MainController{Db: db, Logger: logger})
	mux.Handle("/cleanup", &ctrl.MainController{Db: db, Logger: logger})
	mux.Handle("/image", &ctrl.MainController{Db: db, Logger: logger})
	mux.Handle("/.*", &ctrl.MainController{Db: db, Logger: logger})
	
	// Start the server
	//~ fmt.Println("Listening on port: ", port)
	//~ http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	
	// Start the server (alternative way to be able to start browser)
	server, err := net.Listen("tcp", "localhost:3500")
	if err != nil {
		log.Fatal(err)
	}
	
	LaunchBrowser()
	
	fmt.Println("Listening on port: 3500")
	fmt.Println("---------------------------------------------------------------------------------------------")
	fmt.Println("BELANGRIJK: Bij het sluiten van dit venster wordt de toepassing 'Digitaal Keuzebord' gestopt.")
	log.Fatal(http.Serve(server, mux)) 
	
	return nil
}

func InitializeLog() *log.Logger {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }	
	logger := log.New(file, "", log.LstdFlags)
	return logger
}

func LaunchBrowser() {
	url := "http://localhost:3500/index.html"
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	if err != nil {
		 log.Println(err)
	}
}

func Check(port int) (status bool, err error) {

	// Concatenate a colon and the port
	host := "localhost:" + strconv.Itoa(port)

	// Try to create a server with the port
	server, err := net.Listen("tcp", host)

	// if it fails then the port is likely taken
	if err != nil {
		return false, err
	}

	// close the server
	server.Close()

	// we successfully used and closed the port
	// so it's now available to be used again
	return true, nil

}