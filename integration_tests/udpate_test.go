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
	"github.com/pivotal-cf/on-demand-service-broker/broker"
)

var _ = Describe("updating a service instance", func() {
	It("returns tracking data for an update operation", func() {
		updateTaskID := rand.Int()
		boshDeploysUpdatedManifest := func(env *BrokerEnvironment, id ServiceInstanceID) {
			deploymentName := broker.DeploymentNameFrom(string(id))

			env.Bosh.HasNoTasksFor(deploymentName)
			env.Bosh.HasManifestFor(deploymentName)
			env.Bosh.DeploysWithoutContextId(deploymentName, updateTaskID)
		}

		When(updatingServiceInstance).
			with(NoCredhub, serviceAdapterGeneratesManifest, boshDeploysUpdatedManifest).
			brokerRespondsWith(
				http.StatusAccepted,
				fmt.Sprintf(`{"operation":{"OperationType":"update", "BoshTaskID": %d}`, updateTaskID),
				"foo",
			)
	})
})

var updatingServiceInstance = func(env *BrokerEnvironment, id ServiceInstanceID) *http.Request {
	return env.Broker.UpdateServiceInstanceRequest(id)
}
var serviceAdapterGeneratesManifest = func(sa *ServiceAdapter, id ServiceInstanceID) {
	sa.adapter.GenerateManifest().ToReturnManifest(rawManifestWithDeploymentName(id))
}

func rawManifestWithDeploymentName(id ServiceInstanceID) string {
	return "name: " + broker.DeploymentNameFrom(string(id))
}
