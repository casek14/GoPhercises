package storyTelling

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
		availableStories = append(availableStories, strings.TrimSuffix(f.Name(), ".json"))
	}

	return availableStories
}

type Stories struct {
	StoriesList []string
}

type RenderStory struct {
	StoryTitle string
	Paragraphs []string
	Options    []Option
}

// Story
type Story struct {
	StoryTitle string    `json:"story_title"`
	Chapters   []Chapter `json:"chapters"`
}

func FindChapterByTitle(s string, ch string) (*Chapter, error) {
	story, err := LoadStory(s)
	if err != nil {
		log.Printf("Unable to find story %s, ERR: %s", s, err)
		return nil, err
	}
	for _, c := range story.Chapters {
		if c.Title == ch {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("Chapter with title %s does not exists in story %s", ch, s)
}

func GetRenderStory(story string, chapter string) (*RenderStory, error) {
	ch, err := FindChapterByTitle(story, chapter)
	if err != nil {
		return nil, err
	}
	var renderedStory RenderStory
	renderedStory.StoryTitle = story
	renderedStory.Paragraphs = ch.Paragraphs
	renderedStory.Options = ch.Options

	return &renderedStory, nil
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
