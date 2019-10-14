/*
Copyright (c) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains tests for servers.

package tests

import (
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gorilla/mux"
	cmv1 "github.com/openshift-online/ocm-api-metamodel/tests/api/clustersmgmt/v1"
)

type MyTestRootServer struct{}

func (s *MyTestRootServer) Clusters() cmv1.ClustersServer {
	return &MyTestClustersServer{}
}

type MyTestClustersServer struct{}

func (s *MyTestClustersServer) List(ctx context.Context, request *cmv1.ClustersListServerRequest,
	response *cmv1.ClustersListServerResponse) error {
	items, err := cmv1.NewClusterList().
		Items(cmv1.NewCluster().Name("test-list-clusters")).
		Build()
	if err != nil {
		return err
	}
	response.Items(items)
	response.Page(request.Page())
	response.Size(request.Size())
	response.Total(items.Len())
	return nil
}

func (s *MyTestClustersServer) Add(ctx context.Context, request *cmv1.ClustersAddServerRequest,
	response *cmv1.ClustersAddServerResponse) error {
	return nil
}

func (s *MyTestClustersServer) Cluster(id string) cmv1.ClusterServer {
	return &MyTestClusterServer{}
}

type MyTestClusterServer struct{}

func (s *MyTestClusterServer) Get(ctx context.Context, request *cmv1.ClusterGetServerRequest,
	response *cmv1.ClusterGetServerResponse) error {
	cluster, err := cmv1.NewCluster().Name("test-get-cluster-by-id").Build()
	if err != nil {
		return err
	}
	response.Body(cluster)
	return nil
}

func (s *MyTestClusterServer) Update(ctx context.Context, request *cmv1.ClusterUpdateServerRequest,
	response *cmv1.ClusterUpdateServerResponse) error {
	return nil
}

func (s *MyTestClusterServer) Delete(ctx context.Context, request *cmv1.ClusterDeleteServerRequest,
	response *cmv1.ClusterDeleteServerResponse) error {
	return nil
}

func (s *MyTestClusterServer) Groups() cmv1.GroupsServer {
	return nil
}

func (s *MyTestClusterServer) IdentityProviders() cmv1.IdentityProvidersServer {
	return &MyTestIdentityProvidersServer{}
}

type MyTestIdentityProvidersServer struct{}

func (s *MyTestIdentityProvidersServer) List(ctx context.Context,
	request *cmv1.IdentityProvidersListServerRequest,
	response *cmv1.IdentityProvidersListServerResponse) error {
	items, err := cmv1.NewIdentityProviderList().
		Items(cmv1.NewIdentityProvider().Name("test-list-identity-providers")).
		Build()
	if err != nil {
		return err
	}
	response.Items(items)
	response.Page(1)
	response.Size(1)
	response.Total(1)
	return nil
}

func (s *MyTestIdentityProvidersServer) Add(ctx context.Context,
	request *cmv1.IdentityProvidersAddServerRequest,
	response *cmv1.IdentityProvidersAddServerResponse) error {
	return nil
}

func (s *MyTestIdentityProvidersServer) IdentityProvider(id string) cmv1.IdentityProviderServer {
	return nil
}

var _ = Describe("Server", func() {
	It("Can receive a request and return response", func() {
		myTestRootServer := new(MyTestRootServer)
		rootAdapter := cmv1.NewRootAdapter(myTestRootServer, mux.NewRouter())

		request := httptest.NewRequest(http.MethodGet, "/clusters", nil)
		recorder := httptest.NewRecorder()
		rootAdapter.ServeHTTP(recorder, request)

		Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
	})

	It("Returns a 404 for a path with a trailing slash", func() {
		myTestRootServer := new(MyTestRootServer)
		rootAdapter := cmv1.NewRootAdapter(myTestRootServer, mux.NewRouter())

		request := httptest.NewRequest(http.MethodGet, "/clusters/", nil)
		recorder := httptest.NewRecorder()
		rootAdapter.ServeHTTP(recorder, request)

		Expect(recorder.Result().StatusCode).To(Equal(http.StatusNotFound))
	})

	It("Returns a 404 for an unkown resource", func() {
		myTestRootServer := new(MyTestRootServer)
		rootAdapter := cmv1.NewRootAdapter(myTestRootServer, mux.NewRouter())

		request := httptest.NewRequest(http.MethodGet, "/foo", nil)
		recorder := httptest.NewRecorder()
		rootAdapter.ServeHTTP(recorder, request)

		Expect(recorder.Result().StatusCode).To(Equal(http.StatusNotFound))
	})

	It("Can get a list of clusters", func() {
		myTestRootServer := new(MyTestRootServer)
		rootAdapter := cmv1.NewRootAdapter(myTestRootServer, mux.NewRouter())

		request := httptest.NewRequest(http.MethodGet, "/clusters", nil)
		recorder := httptest.NewRecorder()
		rootAdapter.ServeHTTP(recorder, request)

		expected := `{
			"page": 0,
			"size": 0,
			"total": 1,
			"items": [
				{
					"kind": "Cluster",
					"name": "test-list-clusters"
				}
			]
		}`

		Expect(recorder.Body).To(MatchJSON(expected))
		Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
	})

	It("Can get a list of clusters by page", func() {
		myTestRootServer := new(MyTestRootServer)
		rootAdapter := cmv1.NewRootAdapter(myTestRootServer, mux.NewRouter())

		request := httptest.NewRequest(http.MethodGet, "/clusters?page=2", nil)
		recorder := httptest.NewRecorder()
		rootAdapter.ServeHTTP(recorder, request)

		expected := `{
			"page":2,
			"size":0,
			"total":1,
			"items":[{"kind":"Cluster","name":"test-list-clusters"}]
			}`

		Expect(recorder.Body).To(MatchJSON(expected))
		Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
	})

	It("Can get a list of clusters by size", func() {
		myTestRootServer := new(MyTestRootServer)
		rootAdapter := cmv1.NewRootAdapter(myTestRootServer, mux.NewRouter())

		request := httptest.NewRequest(http.MethodGet, "/clusters?size=2", nil)
		recorder := httptest.NewRecorder()
		rootAdapter.ServeHTTP(recorder, request)

		expected := `{
			"page": 0,
			"size": 2,
			"total": 1,
			"items": [
				{
					"kind": "Cluster",
					"name": "test-list-clusters"
				}
			]
		}`

		Expect(recorder.Body).To(MatchJSON(expected))
		Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
	})

	It("Can get a list of clusters by size and page", func() {
		myTestRootServer := new(MyTestRootServer)
		rootAdapter := cmv1.NewRootAdapter(myTestRootServer, mux.NewRouter())

		request := httptest.NewRequest(http.MethodGet, "/clusters?size=2&page=1", nil)
		recorder := httptest.NewRecorder()
		rootAdapter.ServeHTTP(recorder, request)

		expected := `{
			"page": 1,
			"size": 2,
			"total": 1,
			"items": [
				{
					"kind": "Cluster",
					"name": "test-list-clusters"
				}
			]
		}`

		Expect(recorder.Body).To(MatchJSON(expected))
		Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
	})

	It("Can get a cluster by id", func() {
		myTestRootServer := new(MyTestRootServer)
		rootAdapter := cmv1.NewRootAdapter(myTestRootServer, mux.NewRouter())

		request := httptest.NewRequest(http.MethodGet, "/clusters/123", nil)
		recorder := httptest.NewRecorder()
		rootAdapter.ServeHTTP(recorder, request)

		expected := `{
			"kind": "Cluster",
			"name": "test-get-cluster-by-id"
		}`

		Expect(recorder.Body).To(MatchJSON(expected))
		Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
	})

	It("Can get a cluster sub resource by id", func() {
		myTestRootServer := new(MyTestRootServer)
		rootAdapter := cmv1.NewRootAdapter(myTestRootServer, mux.NewRouter())

		request := httptest.NewRequest(
			http.MethodGet,
			"/clusters/123/identity_providers",
			nil,
		)
		recorder := httptest.NewRecorder()
		rootAdapter.ServeHTTP(recorder, request)

		expected := `{
			"page": 1,
			"size": 1,
			"total": 1,
			"items": [
				{
					"kind": "IdentityProvider",
					"name": "test-list-identity-providers"
				}
			]
		}`

		Expect(recorder.Body).To(MatchJSON(expected))
		Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
	})

	It("Returns a 404 for an unkown sub resource", func() {
		myTestRootServer := new(MyTestRootServer)
		rootAdapter := cmv1.NewRootAdapter(myTestRootServer, mux.NewRouter())

		request := httptest.NewRequest(http.MethodGet, "/clusters/123/foo", nil)
		recorder := httptest.NewRecorder()
		rootAdapter.ServeHTTP(recorder, request)

		Expect(recorder.Result().StatusCode).To(Equal(http.StatusNotFound))
	})
})
