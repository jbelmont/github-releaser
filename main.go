package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
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

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
}
