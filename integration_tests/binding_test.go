// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_tests

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/on-demand-service-broker/adapterclient"
	"github.com/pivotal-cf/on-demand-service-broker/boshclient"
	sdk "github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

var (
	brokerPath         = NewBinary("github.com/pivotal-cf/on-demand-service-broker/cmd/on-demand-service-broker")
	serviceAdapterPath = NewBinary("github.com/pivotal-cf/on-demand-service-broker/integration_tests/mock/adapter")
)

type Test struct {
	expectedStatusCode int
	expectedBody       string
	expectedLogMessage string
}

func NewTest(expectedStatusCode int, expectedBody, expectedLogMessage string) *Test {
	return &Test{
		expectedStatusCode: expectedStatusCode,
		expectedBody:       expectedBody,
		expectedLogMessage: expectedLogMessage,
	}
}

var _ = Describe("binding service instances", func() {
	It("binds a service to an application instance", func() {
		serviceInstanceID := AServiceInstanceID()

		env := NewBrokerEnvironment(NewBosh(), NewCloudFoundry(), NewServiceAdapter(serviceAdapterPath.Path()), NoopCredhub, brokerPath.Path())
		defer env.Close()

		env.ServiceAdapter.ReturnsBinding()

		env.Start()

		env.Bosh.HasDeploymentFor(serviceInstanceID)

		response := responseTo(env.Broker.CreateBindingRequest(serviceInstanceID))
		Expect(response.StatusCode).To(Equal(http.StatusCreated))
		Expect(bodyOf(response)).To(MatchJSON(BindingResponse))
		env.Broker.HasLogged(fmt.Sprintf("create binding with ID %s", bindingId))

		env.Verify()
	})

	It("sends login details to credhub when credhub configured", func() {
		serviceInstanceID := AServiceInstanceID()
		mockCredhub := NewCredhub()
		env := NewBrokerEnvironment(NewBosh(), NewCloudFoundry(), NewServiceAdapter(serviceAdapterPath.Path()), mockCredhub, brokerPath.Path())
		defer env.Close()

		env.ServiceAdapter.ReturnsBinding()

		env.Start()

		env.Bosh.HasDeploymentFor(serviceInstanceID)
		mockCredhub.WillReceiveCredentials(serviceInstanceID)

		response := responseTo(env.Broker.CreateBindingRequest(serviceInstanceID))
		Expect(response.StatusCode).To(Equal(http.StatusCreated))
		Expect(bodyOf(response)).To(MatchJSON(BindingResponse))
		env.Broker.HasLogged(fmt.Sprintf("create binding with ID %s", bindingId))

		env.Verify()
	})

	It("fails when rejected by adapter", func() {
		stderrMessage := fmt.Sprintf("binding stderr message-%d", rand.Int())

		serviceInstanceID := AServiceInstanceID()
		env := NewBrokerEnvironment(NewBosh(), NewCloudFoundry(), NewServiceAdapter(serviceAdapterPath.Path()), NoopCredhub, brokerPath.Path())
		defer env.Close()

		env.ServiceAdapter.FailsToBindBecause(sdk.BindingAlreadyExistsErrorExitCode, stderrMessage)

		env.Start()

		env.Bosh.HasDeploymentFor(serviceInstanceID)

		response := responseTo(env.Broker.CreateBindingRequest(serviceInstanceID))
		Expect(response.StatusCode).To(Equal(http.StatusConflict))
		Expect(bodyOf(response)).To(MatchJSON(errorResponse(adapterclient.BindingAlreadyExistsMessage)))
		env.Broker.HasLogged(stderrMessage)

		env.Verify()
	})

	It("fails when bosh is unreachable", func() {
		serviceInstanceID := AServiceInstanceID()
		env := NewBrokerEnvironment(NewBosh(), NewCloudFoundry(), NewServiceAdapter(serviceAdapterPath.Path()), NoopCredhub, brokerPath.Path())
		defer env.Close()

		env.Start()

		env.Bosh.Close()

		response := responseTo(env.Broker.CreateBindingRequest(serviceInstanceID))
		Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
		Expect(bodyOf(response)).To(MatchJSON(errorResponse("Currently unable to bind service instance, please try again later")))
		env.Broker.HasLogged(boshclient.UnreachableMessage)

		env.Verify()
	})

	It("fails when bosh deployment doesn't exist", func() {
		serviceInstanceID := AServiceInstanceID()
		env := NewBrokerEnvironment(NewBosh(), NewCloudFoundry(), NewServiceAdapter(serviceAdapterPath.Path()), NoopCredhub, brokerPath.Path())
		defer env.Close()

		env.Start()

		env.Bosh.HasNoDeploymentFor(serviceInstanceID)

		response := responseTo(env.Broker.CreateBindingRequest(serviceInstanceID))
		Expect(response.StatusCode).To(Equal(http.StatusNotFound))
		Expect(bodyOf(response)).To(MatchJSON(errorResponse("instance does not exist")))
		env.Broker.HasLogged(fmt.Sprintf("binding: instance %s, not found", serviceInstanceID))

		env.Verify()
	})

})

func responseTo(request *http.Request) *http.Response {
	response, err := http.DefaultClient.Do(request)
	Expect(err).ToNot(HaveOccurred())
	return response
}

func errorResponse(message string) string {
	return fmt.Sprintf(`{"description": "%s"}`, message)
}

func bodyOf(response *http.Response) []byte {
	body, err := ioutil.ReadAll(response.Body)
	Expect(err).NotTo(HaveOccurred())
	return body
}
