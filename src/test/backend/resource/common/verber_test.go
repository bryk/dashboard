package common

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/client/restclient"
)

type clientFunc func(req *http.Request) (*http.Response, error)

func (f clientFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type FakeRESTClient struct {
	response *http.Response
	err      error
}

func (c *FakeRESTClient) Delete() *restclient.Request {
	return restclient.NewRequest(clientFunc(func(req *http.Request) (*http.Response, error) {
		fmt.Printf("%#v\n", req.URL)
		return nil, errors.New("err")
	}), "DELETE", nil, "/api/v1", restclient.ContentConfig{}, nil, nil)
}

func TestDeleteShouldPropagateErrors(t *testing.T) {
	verber := ResourceVerber{client: &FakeRESTClient{}}

	err := verber.Delete("replicaset", "bar", "baz")

	if !reflect.DeepEqual(err, errors.New("err")) {
		t.Fatalf("Expected error on verber delete but got %#v", err)
	}
	t.Fatalf("foo")
}
