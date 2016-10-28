// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/pflag"
)

var (
	argPort          = pflag.Int("port", 9090, "The port to listen to for incoming HTTP requests")
	argApiserverHost = pflag.String("database-host", "", "The address of the Kubernetes Apiserver "+
		"to connect to in the format of protocol://address:port, e.g., "+
		"http://localhost:8080. If not specified, the assumption is that the binary runs inside a"+
		"Kubernetes cluster and local discovery is attempted.")
	argHeapsterHost = pflag.String("login-service-host", "", "The address of the Heapster Apiserver "+
		"to connect to in the format of protocol://address:port, e.g., "+
		"http://localhost:8082. If not specified, the assumption is that the binary runs inside a"+
		"Kubernetes cluster and service proxy will be used.")
	abFlag = os.Getenv("AB_EXPERIMENT_TURNED_ON");
	reqs   = 0
)

func handler(w http.ResponseWriter, r *http.Request) {
	img := "t1"
	reqs++
	log.Printf("Handling request to review app homepage, requests so far: %d", reqs)
	if abFlag == "true" && (reqs%3 == 0) {
		img = "t2"
		log.Printf("Using alternate versions: %s\n", img)
	}

	html := fmt.Sprintf("<html><body style=\"background-color: #97b1c7; display: flex;\"><img "+
		"src=\"public/en/assets/images/%s.png\" style=\"margin: 0 auto;\"></body></html>", img)
	log.Printf("HTML: %s\n", html)

	fmt.Fprintf(w, html)
}

func main() {
	// Set logging output to standard console out
	log.SetOutput(os.Stdout)

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	flag.CommandLine.Parse(make([]string, 0)) // Init for glog calls in kubernetes packages

	log.Printf("Starting t-shirt review app on port: %d", *argPort)
	log.Printf("Connected to database at %s", *argApiserverHost)
	log.Printf("Connected to login service at %s", *argHeapsterHost)
	log.Printf("A/B experiment config: %s", abFlag)

	http.HandleFunc("/", handler)
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	log.Print(http.ListenAndServe(fmt.Sprintf(":%d", *argPort), nil))
}
