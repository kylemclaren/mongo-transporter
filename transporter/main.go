// Copyright 2014 The Transporter Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/franela/goreq"

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
	webhookUrl    string = os.Getenv("SLACK_WEBHOOK_URL")
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

		var slackMessage string = "Copying from " + srcNamespace + " to " + destNamespace

		if slack == true {

			type notification struct {
				channel  string
				username string
				icon_url string
				text     string
			}

			notification = Notification{channel: slackChannel, username: slackApp, icon_url: "https://raw.githubusercontent.com/kylemclaren/mongo-transporter/notifications/slack_icon.png", text: slackMessage}

			res, err := goreq.Request{
				Method: "POST",
				Uri:    webhookUrl,
				Body:   notification,
			}.Do()

			fmt.Println("Posted to Slack")

		}

		source :=
			transporter.NewNode("source", "mongo", map[string]interface{}{"uri": sourceUri, "namespace": srcNamespace, "tail": tail}).
				Add(transporter.NewNode("out", "mongo", map[string]interface{}{"uri": destUri, "namespace": destNamespace}))

		if debug == true {
			source.Add(transporter.NewNode("out", "file", map[string]interface{}{"uri": "stdout://"}))
		}

		pipeline, err := transporter.NewPipeline(source, events.NewLogEmitter(), 1*time.Second)
		if err != nil {
			fmt.Println("Transporter error: " + err.Error())
			os.Exit(1)
		}

		go pipeline.Run()

	}

	c := make(chan bool)
	<-c

}
