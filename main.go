package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	tagName         = flag.String("tagName", "", "Please provide a tag name for the release")
	targetCommitish = flag.String("targetCommitish", "", "Please provide a targetCommitish value namely master")
	name            = flag.String("name", "", "A name for the Github Release")
	body            = flag.String("body", "", "Provide a description of the release")
	username        = flag.String("username", "", "Enter your github username/organization")
	repo            = flag.String("repo", "", "Enter your github repository")
	accessToken     = flag.String("accessToken", "", "Enter your Github Access Token")
)

type Release struct {
	URL             string    `json:"url"`
	HTMLURL         string    `json:"html_url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	ID              int       `json:"id"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Body            string    `json:"body"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Assets []interface{} `json:"assets"`
}

func checkArgs() {
	if *tagName == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}
	if *targetCommitish == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}
	if *name == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}
	if *body == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}
	if *username == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}
	if *repo == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}
	if *accessToken == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}
}

func main() {
	flag.Parse()

	checkArgs()

	body := struct {
		TagName         string `json:"tag_name"`
		TargetCommitish string `json:"target_commitish"`
		Name            string `json:"name"`
		Body            string `json:"body"`
		Draft           bool   `json:"draft"`
		Prerelease      bool   `json:"prerelease"`
	}{
		*tagName,
		*targetCommitish,
		*name,
		*body,
		false,
		false,
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalln("An issue occurred marshalling this data ", err)
	}

	client := http.Client{}

	req, err := http.NewRequest(
		"POST",
		"https://api.github.com/repos/"+*username+"/"+*repo+"/releases",
		bytes.NewBuffer(postBody),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+*accessToken)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != http.StatusCreated {
		log.Fatalln("Should receive a status code of 201 created")
	}

	release := new(Release)

	json.NewDecoder(resp.Body).Decode(&release)

	json, err := json.MarshalIndent(release, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(json))
}
