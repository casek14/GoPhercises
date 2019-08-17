package main

import (
	"github.com/casek14/GoPhercises/StoryTelling"
	"net/http"
)

func main() {

	r := storyTelling.NewStoriesRouter()
	http.Handle("/",r)
	http.ListenAndServe(":8080", nil)
}
