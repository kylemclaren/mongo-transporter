// Copyright 2014 The Transporter Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/compose/transporter/pkg/events"
	"github.com/compose/transporter/pkg/transporter"
)

var (
	sourceUri     string = os.Getenv("SOURCE_MONGO_URL")
	destUri       string = os.Getenv("DESTINATION_MONGO_URL")
	sourceDB      string = os.Getenv("SOURCE_DB")
	destinationDB string = os.Getenv("DEST_DB")
	envTail              = os.Getenv("TAIL")
	envDebug             = os.Getenv("DEBUG")
	slackNotify          = os.Getenv("SLACK_NOTIFY")
	slackToken    string = os.Getenv("SLACK_TOKEN")
	slackChannel  string = os.Getenv("SLACK_CHANNEL")
	slackApp      string = os.Getenv("SLACK_APP_NAME")
)

func main() {

	var (
		tail  bool
		debug bool
		slack bool
	)

	tail = (strings.ToLower(envTail) == "true")
	debug = (strings.ToLower(envDebug) == "true")
	slack = (strings.ToLower(slackNotify) == "true")

	// Connect to source URI

	sess, err := mgo.Dial(sourceUri)
	if err != nil {
		fmt.Println("Can't connect: " + err.Error())
		os.Exit(1)
	}

	// Get collection names from source DB

	names, err := sess.DB(sourceDB).CollectionNames()
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}

	// Iterate over collection names and run a pipeline for each

	for _, name := range names {

		if strings.HasPrefix(name, "system.") {
			continue
		}

		srcNamespace := fmt.Sprintf("%s.%s", sourceDB, name)
		destNamespace := fmt.Sprintf("%s.%s", destinationDB, name)

		fmt.Println("Copying from " + srcNamespace + " to " + destNamespace)

		slackMessage := "Copying from " + srcNamespace + " to " + destNamespace

		if slack == true {

			apiUrl := "https://slack.com/api"
			resource := "/chat.postMessage/"
			data := url.Values{}
			data.Set("token", slackToken)
			data.Add("channel", slackChannel)
			data.Add("username", slackApp)
			data.Add("icon_url", "https://raw.githubusercontent.com/kylemclaren/mongo-transporter/notifications/slack_icon.png")
			data.Add("text", slackMessage)

			u, _ := url.ParseRequestURI(apiUrl)
			u.Path = resource
			u.RawQuery = data.Encode()
			urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/?name=foo&surname=bar"

			client := &http.Client{}
			r, _ := http.NewRequest("POST", urlStr, nil)
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

			resp, _ := client.Do(r)

		}

		source :=
			transporter.NewNode("source", "mongo", map[string]interface{}{"uri": sourceUri, "namespace": srcNamespace, "tail": tail}).
				Add(transporter.NewNode("out", "mongo", map[string]interface{}{"uri": destUri, "namespace": destNamespace}))

		if debug == true {
			source.Add(transporter.NewNode("out", "file", map[string]interface{}{"uri": "stdout://"}))
		}

		pipeline, err := transporter.NewPipeline(source, events.NewLogEmitter(), 1*time.Second)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		go pipeline.Run()
	}

	c := make(chan bool)
	<-c

}
