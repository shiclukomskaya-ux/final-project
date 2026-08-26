package main

import (
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
	//dbPath := filepath.Join(rootPath, "scheduler.db")
	if _, err := os.Stat(webPath); os.IsNotExist(err) {
		webPath = "web"
	}
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	address := fmt.Sprintf(":%s", port)

	http.Handle("/", http.FileServer(http.Dir(webPath)))
	http.ListenAndServe(address, nil)
}
