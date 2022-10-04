package main

import (
	"fmt"
	"net/http"
)

func handleSlash(p_w http.ResponseWriter, p_req *http.Request) {
	if p_req.URL.Path != "/" {
		http.Error(p_w, "404 Not found", http.StatusNotFound)

		return
	}

	p_w.Write([]byte("<h1>Hello from golang!</h1>"))
}

func handleSlashTest(p_w http.ResponseWriter, p_req *http.Request) {
	if p_req.URL.Path != "/test" {
		http.Error(p_w, "404 Not found", http.StatusNotFound)

		return
	}

	p_w.Write([]byte("<h1>Hello from golang!</h1><h1>Welcome to the test page.</h1>"))
}

func main() {
	http.HandleFunc("/",     handleSlash)
	http.HandleFunc("/test", handleSlashTest)

	fmt.Println("Starting a server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
