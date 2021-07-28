package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/landscape", landscapeHandler)
	http.Handle("/", http.FileServer(http.Dir("./web/static")))

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func landscapeHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := getCicdProjects()
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode([]Project{
			{
				Name: "project could not be find",
			},
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
	return
}
