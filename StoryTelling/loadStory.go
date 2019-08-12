package storyTelling

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	dirPath = "/home/casek/go/src/github.com/casek14/GoPhercises/StoryTelling/stories"
)

func GetAvailableStories(path string) []string {
	var availableStories []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("Unable to read dir %s. ERROR: %s\n", path, err)
		return availableStories
	}

	for _, f := range files {
		availableStories = append(availableStories, f.Name())
	}

	return availableStories
}

// Story
type Story struct {
	StoryTitle string    `json:"story_title"`
	Chapters   []Chapter `json:"chapters"`
}

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Link string `json:"link"`
}

func LoadStory(storyName string) (*Story, error) {
	jsonFile, err := os.Open(dirPath + "/" + storyName + ".json")
	if err != nil {
		log.Printf("Unable to open story %s.json, ERROR: %s\n", storyName, err)
		return nil, err
	}
	jsonF, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Printf("Unable to read story: %s\n", err)
		return nil, err
	}
	var story Story
	json.Unmarshal(jsonF, &story)
	return &story, nil
}
