// Copyright 2014 The Transporter Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// seed is a reimagining of the seed mongo to mongo tool

package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"

	"github.com/compose/transporter/pkg/events"
	"github.com/compose/transporter/pkg/transporter"
)

var (
	sourceUri string			= os.Getenv("SOURCE_MONGO_URL")
	destUri string				= os.Getenv("DESTINATION_MONGO_URL")
	sourceDB string				= os.Getenv("SOURCE_DB")
	destinationDB string	= os.Getenv("DEST_DB")
	tail									= os.Getenv("TAIL")
	debug									= os.Getenv("DEBUG")
)

// Connect to source URI

sess, err := mgo.Dial(sourceURI)
	if err != nil {
	  fmt.Println("can't connect " + err.Error())
	}

// Get collection names from source DB

names, err := sess.DB(sourceDB).CollectionNames()
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}

for _, name := range names {
	func main() {
		source :=
			transporter.NewNode("source", "mongo", map[string]interface{}{"uri": sourceUri, "namespace": sourceDB + "." + name, "tail": tail}).
				Add(transporter.NewNode("out", "mongo", map[string]interface{}{"uri": destUri, "namespace": destinationDB + "." + name}))

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
