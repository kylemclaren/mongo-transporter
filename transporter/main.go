// Copyright 2014 The Transporter Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// seed is a reimagining of the seed mongo to mongo tool

// users and indexes

package main

import (
	"fmt"
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
)

func main() {

	var (
		tail  bool
		debug bool
	)

	tail = (strings.ToLower(envTail) == "true")
	debug = (strings.ToLower(envDebug) == "true")

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

		pipeline.Run()
	}
}
