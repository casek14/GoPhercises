package storyTelling

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func init() {
	tpl = template.Must(template.New("").ParseFiles("templates/story.html"))
}

var tpl *template.Template

func NewStoriesRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HomePageHandler)
	r.HandleFunc("/story/{story}/{chapter}", StoriesHandler)
	return r
}

// Handle homepage for the application
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
	}
	var stories Stories
	stories.StoriesList = GetAvailableStories(dirPath)
	t.Execute(w, stories)

}

func StoriesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	story, err := GetRenderStory(vars["story"], vars["chapter"])
	if err != nil {
		fmt.Fprintf(w, "Unable to find story %s. ERR: %s", vars["story"], err)
	}
	t, err := template.ParseFiles("templates/story.html")
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
	}

	t.Execute(w, story)
}
