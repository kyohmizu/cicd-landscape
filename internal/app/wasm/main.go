package main

import (
	"encoding/json"
	"fmt"
	"github.com/kyohmizu/cicd-landscape/internal/pkg/common"
	"net/http"
	"strconv"
	"syscall/js"
)

var document = js.Global().Get("document")
var app = js.Global().Get("document").Call("getElementById", "app")

type Project common.Project
type Projects []Project

func main() {
	c := make(chan bool)
	go load()
	<-c
}

func load() {
	projects, err := getProjects()
	if err != nil {
		landscapeDiv := document.Call("createElement", "div")
		landscapeDiv.Get("classList").Call("add", "landscape")
		app.Call("appendChild", landscapeDiv)

		errorDiv := document.Call("createElement", "div")
		errorDiv.Set("textContent", "error")
		landscapeDiv.Call("appendChild", errorDiv)
		return
	}

	//pmap := splitByRelation(projects)
	loadByRelation("CNCF Graduated Projects", filterByRelation(projects, "graduated"))
	loadByRelation("CNCF Incubating Projects", filterByRelation(projects, "incubating"))
	loadByRelation("CNCF Sandbox Projects", filterByRelation(projects, "sandbox"))
	loadByRelation("Others", filterByRelation(projects, ""))
}

func loadByRelation(rel string, projects Projects) {
	if len(projects) == 0 {
		return
	}

	relationDiv := document.Call("createElement", "div")
	relationDiv.Get("classList").Call("add", "relation")
	relationDiv.Set("textContent", rel + " (" + strconv.Itoa(len(projects)) + ")")
	app.Call("appendChild", relationDiv)

	landscapeDiv := document.Call("createElement", "div")
	landscapeDiv.Get("classList").Call("add", "landscape")
	app.Call("appendChild", landscapeDiv)

	for _, proj := range projects {
		landscapeDiv.Call("appendChild", createItem(proj))
	}
}

func getProjects() (Projects, error) {
	resp, err := http.Get("/landscape")
	if err != nil {
		return Projects{}, err
	}
	defer resp.Body.Close()

	projects := Projects{}
	err = json.NewDecoder(resp.Body).Decode(&projects)
	return projects, err
}

func filterByRelation(projects Projects, relation string) Projects {
	var filtered Projects
	for _, proj := range projects {
		if proj.Project == relation {
			filtered = append(filtered, proj)
		}
	}
	return filtered
}

//func splitByRelation(projects Projects) map[string]Projects {
//	var pmap map[string]Projects
//	for _, proj := range projects {
//		switch proj.Project {
//		case "graduated":
//			addProjectToMap(pmap, "graduated", proj)
//		case "incubating":
//			addProjectToMap(pmap, "incubating", proj)
//		case "sandbox":
//			addProjectToMap(pmap, "sandbox", proj)
//		default:
//			addProjectToMap(pmap, "others", proj)
//		}
//	}
//	return pmap
//}

//func addProjectToMap(pmap map[string]Projects, relation string, proj Project) {
//	v, ok := pmap[relation]
//	if ok {
//		v = append(v, proj)
//	} else {
//		pmap[relation] = Projects{proj}
//	}
//}

func createItem(proj Project) js.Value {
	item := document.Call("createElement", "div")
	item.Get("classList").Call("add", "item")

	text := proj.RepoUrl
	if text == "" {
		text = proj.HomepageUrl
	}
	itemLink := document.Call("createElement", "a")
	itemLink.Set("href", text)

	itemTitle := document.Call("createElement", "div")
	itemTitle.Set("textContent", fmt.Sprintf("%s", proj.Name))
	itemTitle.Get("classList").Call("add", "item-title")

	text = proj.Description
	if text == "" {
		text = "No Description"
	}
	itemDescription := document.Call("createElement", "div")
	itemDescription.Set("textContent", text)
	itemDescription.Get("classList").Call("add", "item-description")

	item.Call("appendChild", itemLink)
	itemLink.Call("appendChild", itemTitle)
	itemLink.Call("appendChild", itemDescription)
	return item
}
