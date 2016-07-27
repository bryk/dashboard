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

package handler

import (
	"fmt"
	"net/http"
	"log"
	"strconv"

	restful "github.com/emicklei/go-restful"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

const (
	// RequestLogString is a template for request log message.
	RequestLogString = "[%s] Incoming %s %s %s request from %s"

	// ResponseLogString is a template for response log message.
	ResponseLogString = "[%s] Outcoming response to %s with %d status code"
)

// ApiHandler is a representation of API handler. Structure contains client, Heapster client and
// client configuration.
type ApiHandler struct {
	client         *client.Client
	clientConfig   clientcmd.ClientConfig
}

type Scale struct {

}

type ScaleSpec struct {}

func CreateHttpApiHandler(client *client.Client,
	clientConfig clientcmd.ClientConfig) http.Handler {

log.Printf("Hi")
	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	apiHandler := ApiHandler{client, clientConfig}

	apiV1Ws := new(restful.WebService)
	apiV1Ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	wsContainer.Add(apiV1Ws)
	apiV1Ws.Route(
		apiV1Ws.GET("/scale/{namespace}/{name}/{replicas}").
			To(apiHandler.handleScale).
			Writes(Scale{}))

		return wsContainer
}

func (apiHandler *ApiHandler) handleScale(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	namespace := request.PathParameter("namespace")
	replicas := request.PathParameter("replicas")
	log.Printf("Scale request: %s %s %s", name, namespace, replicas)
	replicasInt, err := strconv.Atoi(replicas)
	if err != nil {
		fmt.Printf("Replicas is not an int")
		return
	}

	result, err := apiHandler.client.Deployments(namespace).Get(name)
	if err != nil {
		fmt.Printf("%#v, %#v\n", response, err)
		return
	}

	result.Spec.Replicas = int32(replicasInt)
	apiHandler.client.Deployments(namespace).Update(result)

	out := fmt.Sprintf("Scaled deployment %s in namespace %s to %s", name, namespace, replicas)
	response.WriteHeaderAndEntity(http.StatusCreated, out)
}
