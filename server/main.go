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

	"github.com/compose/transporter/pkg/events"
	"github.com/compose/transporter/pkg/transporter"
)

var (
	sourceUri string = os.Getenv("SOURCE_MONGO_URL")
	sinkUri   string = os.Getenv("SINK_MONGO_URL")
	sourceDB  string = os.Getenv("SOURCE_DB")
	sinkDB    string = os.Getenv("SINK_DB")
	envTail          = os.Getenv("TAIL")
	envDebug         = os.Getenv("DEBUG")
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
		sinkNamespace := fmt.Sprintf("%s.%s", sinkDB, name)
		fmt.Println("Copying from " + srcNamespace + " to " + sinkNamespace)

		source :=
			transporter.NewNode("source", "mongo", map[string]interface{}{"uri": sourceUri, "namespace": srcNamespace, "tail": tail}).
				Add(transporter.NewNode("out", "mongo", map[string]interface{}{"uri": sinkUri, "namespace": sinkNamespace}))

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
