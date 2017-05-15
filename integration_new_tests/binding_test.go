// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_new_tests

import (
	"fmt"
	"math/rand"
	"net/http"

	. "github.com/onsi/ginkgo"
	"github.com/pivotal-cf/on-demand-service-broker/boshdirector"
	"github.com/pivotal-cf/on-demand-service-broker/broker"
	"github.com/pivotal-cf/on-demand-service-broker/serviceadapter"
	sdk "github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

var _ = Describe("binding service instances", func() {
	It("binds a service to an application instance", func() {
		When(creatingNewBinding).
			With(NoCredhub, serviceAdapterReturnsBinding, boshHasVMsForServiceInstance).
			theBroker(
				RespondsWith(http.StatusCreated, BindingResponse),
				Logs(fmt.Sprintf("create binding with ID %s", bindingGUIDfromCF)))
	})

	It("sends login details to credhub when credhub configured", func() {
		aCredhub := NewCredhub()
		boshHasDeploymentWithCredhub := func(env *BrokerEnvironment) {
			boshHasVMsForServiceInstance(env)
			aCredhub.WillReceiveCredentials(env.serviceInstanceID)
		}

		When(creatingNewBinding).
			With(aCredhub, serviceAdapterReturnsBinding, boshHasDeploymentWithCredhub).
			theBroker(
				RespondsWith(http.StatusCreated, BindingResponse),
				Logs(fmt.Sprintf("create binding with ID %s", bindingGUIDfromCF)),
			)
	})

	It("fails when rejected by adapter", func() {
		stderrMessage := fmt.Sprintf("binding stderr message-%d", rand.Int())
		serviceAdapterFails := func(sa *ServiceAdapter, id ServiceInstanceID) {
			sa.FailsToBindBecause(sdk.BindingAlreadyExistsErrorExitCode, stderrMessage)
		}

		When(creatingNewBinding).
			With(NoCredhub, serviceAdapterFails, boshHasVMsForServiceInstance).
			theBroker(
				RespondsWith(http.StatusConflict, errorBody(serviceadapter.BindingAlreadyExistsMessage)),
				Logs(stderrMessage),
			)
	})

	It("fails when bosh is unreachable", func() {
		When(creatingNewBinding).
			With(NoCredhub, noServiceAdapter, boshConnectionFails).
			theBroker(
				RespondsWith(http.StatusInternalServerError, errorBody("Currently unable to bind service instance, please try again later")),
				Logs(boshdirector.UnreachableMessage),
			)
	})

	It("fails when bosh deployment doesn't exist", func() {
		When(creatingNewBinding).
			With(NoCredhub, noServiceAdapter, boshHasNoVMs).
			theBroker(
				RespondsWith(http.StatusNotFound, errorBody("instance does not exist")),
				LogsWithServiceId("instance %s, not found"),
			)
	})
})

var creatingNewBinding = func(env *BrokerEnvironment) *http.Request {
	return env.Broker.CreateBindingRequest(env.serviceInstanceID)
}
var boshConnectionFails = func(env *BrokerEnvironment) { env.Bosh.Close() }
var boshHasVMsForServiceInstance = func(env *BrokerEnvironment) {
	deploymentName := broker.DeploymentNameFrom(string(env.serviceInstanceID))
	env.Bosh.HasVMsFor(deploymentName)
	env.Bosh.HasManifestFor(deploymentName)
}
var boshHasNoVMs = func(env *BrokerEnvironment) {
	env.Bosh.HasNoVMsFor(broker.DeploymentNameFrom(string(env.serviceInstanceID)))
}
var serviceAdapterReturnsBinding = func(sa *ServiceAdapter, id ServiceInstanceID) { sa.ReturnsBinding() }
var noServiceAdapter = func(sa *ServiceAdapter, id ServiceInstanceID) {}

func errorBody(message string) string {
	return fmt.Sprintf(`{"description": "%s"}`, message)
}
