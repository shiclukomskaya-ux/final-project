package main

import (
	"final-project/internal/api"
	"final-project/internal/database"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	rootPath := filepath.Dir(execPath)
	webPath := filepath.Join(rootPath, "web")
	dbPath := filepath.Join(rootPath, "scheduler.db")

	if _, err := os.Stat(webPath); os.IsNotExist(err) {
		webPath = "web"
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbPath = "scheduler.db"
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	address := fmt.Sprintf(":%s", port)

	if initErr := database.Init(dbPath); initErr != nil {
		panic(initErr)
	}

	http.Handle("/", http.FileServer(http.Dir(webPath)))
	api.Init()
	http.ListenAndServe(address, nil)

}
