// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_tests

import (
	"fmt"
	"math/rand"
	"net/http"

	. "github.com/onsi/ginkgo"
	"github.com/pivotal-cf/on-demand-service-broker/adapterclient"
	"github.com/pivotal-cf/on-demand-service-broker/boshclient"
	sdk "github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

var _ = Describe("binding service instances", func() {
	It("binds a service to an application instance", func() {
		When(WithoutCredhub, serviceAdapterReturnsBinding, boshHasDeploymentForServiceInstance).
			brokerRespondsWith(http.StatusCreated, BindingResponse, fmt.Sprintf("create binding with ID %s", bindingId))
	})

	It("sends login details to credhub when credhub configured", func() {
		mockCredhub := NewCredhub()
		boshHasDeploymentWithCredhub := func(env *BrokerEnvironment, id ServiceInstanceID) {
			boshHasDeploymentForServiceInstance(env, id)
			mockCredhub.WillReceiveCredentials(id)
		}

		When(mockCredhub, serviceAdapterReturnsBinding, boshHasDeploymentWithCredhub).
			brokerRespondsWith(http.StatusCreated, BindingResponse, fmt.Sprintf("create binding with ID %s", bindingId))
	})

	It("fails when rejected by adapter", func() {
		stderrMessage := fmt.Sprintf("binding stderr message-%d", rand.Int())
		serviceAdapterFails := func(sa *ServiceAdapter) { sa.FailsToBindBecause(sdk.BindingAlreadyExistsErrorExitCode, stderrMessage) }

		When(WithoutCredhub, serviceAdapterFails, boshHasDeploymentForServiceInstance).
			brokerRespondsWith(http.StatusConflict, errorBody(adapterclient.BindingAlreadyExistsMessage), stderrMessage)
	})

	It("fails when bosh is unreachable", func() {
		When(WithoutCredhub, noServiceAdapter, boshConnectionFails).
			brokerRespondsWith(http.StatusInternalServerError, errorBody("Currently unable to bind service instance, please try again later"), boshclient.UnreachableMessage)
	})

	It("fails when bosh deployment doesn't exist", func() {
		When(WithoutCredhub, noServiceAdapter, boshHasNoDeployment).
			brokerRespondsWith(http.StatusNotFound, errorBody("instance does not exist"), "not found") // TODO Where to get service instance ID?
	})
})

var boshConnectionFails = func(env *BrokerEnvironment, id ServiceInstanceID) { env.Bosh.Close() }
var boshHasDeploymentForServiceInstance = func(env *BrokerEnvironment, id ServiceInstanceID) { env.Bosh.HasDeploymentFor(id) }
var boshHasNoDeployment = func(env *BrokerEnvironment, id ServiceInstanceID) { env.Bosh.HasNoDeploymentFor(id) }
var serviceAdapterReturnsBinding = func(sa *ServiceAdapter) { sa.ReturnsBinding() }
var noServiceAdapter = func(sa *ServiceAdapter) {}

func errorBody(message string) string {
	return fmt.Sprintf(`{"description": "%s"}`, message)
}
