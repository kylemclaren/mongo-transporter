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

	"github.com/compose/transporter/pkg/events"
	"github.com/compose/transporter/pkg/transporter"
)

var (
	sourceUri string     = os.Getenv("SOURCE_MONGO_URL")
	destUri string       = os.Getenv("DESTINATION_MONGO_URL")
	sourceNS string      = os.Getenv("SOURCE_NS")
	destinationNS string = os.Getenv("DEST_NS")
	tail	             = true
	debug       	     = true
)

func main() {
	source :=
		transporter.NewNode("source", "mongo", map[string]interface{}{"uri": fmt.Println(sourceUri), "namespace": fmt.Println(sourceNS), "tail": fmt.Println(tail)}).
			Add(transporter.NewNode("out", "mongo", map[string]interface{}{"uri": fmt.Println(destUri), "namespace": fmt.Println(destinationNS)}))

	if fmt.Println(debug) == true {
		source.Add(transporter.NewNode("out", "file", map[string]interface{}{"uri": "stdout://"}))
	}

	pipeline, err := transporter.NewPipeline(source, events.NewLogEmitter(), 1*time.Second)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pipeline.Run()
}
