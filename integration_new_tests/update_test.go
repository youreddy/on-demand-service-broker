// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_new_tests

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf/on-demand-service-broker/broker"
)

var _ = Describe("updating a service instance", func() {
	It("returns tracking data for an update operation", func() {
		updateTaskID := rand.Int()
		boshDeploysUpdatedManifest := func(env *BrokerEnvironment) {
			deploymentName := env.DeploymentName()

			env.Bosh.HasNoTasksFor(deploymentName)
			env.Bosh.HasManifestFor(deploymentName)
			env.Bosh.DeploysWithoutContextId(deploymentName, updateTaskID)
		}

		When(updatingServiceInstance).
			With(NoCredhub, serviceAdapterGeneratesManifest, boshDeploysUpdatedManifest).
			theBroker(
				RespondsWith(http.StatusAccepted, matchingUpdateOperationWith(updateTaskID)),
				LogsWithServiceId("updating instance %s"),
				LogsWithDeploymentName(fmt.Sprintf("Bosh task ID for update deployment %%s is %d", updateTaskID)),
			)
	})
})

var updatingServiceInstance = func(env *BrokerEnvironment) *http.Request {
	return env.Broker.UpdateServiceInstanceRequest(env.serviceInstanceID)
}
var serviceAdapterGeneratesManifest = func(sa *ServiceAdapter, id ServiceInstanceID) {
	sa.adapter.GenerateManifest().ToReturnManifest(rawManifestWithDeploymentName(id))
}

func rawManifestWithDeploymentName(id ServiceInstanceID) string {
	return "name: " + broker.DeploymentNameFrom(string(id))
}

func matchingUpdateOperationWith(updateTaskId int) types.GomegaMatcher {
	return &updateOperation{
		expected: MatchJSON(fmt.Sprintf(`{"BoshTaskID": %d, "OperationType": "update"}`, updateTaskId)),
	}
}

type updateOperation struct {
	expected types.GomegaMatcher
}

func (uo *updateOperation) Match(actual interface{}) (bool, error) {
	return uo.expected.Match(asOperationData(actualBytes(actual)))
}

func (uo *updateOperation) FailureMessage(actual interface{}) (message string) {
	return uo.expected.FailureMessage(asOperationData(actualBytes(actual)))
}

func (uo *updateOperation) NegatedFailureMessage(actual interface{}) (message string) {
	return uo.expected.NegatedFailureMessage(asOperationData(actualBytes(actual)))
}

func asOperationData(source []byte) []byte {
	var updateResponse brokerapi.UpdateResponse
	err := json.Unmarshal(source, &updateResponse)
	Expect(err).NotTo(HaveOccurred())

	var operationData broker.OperationData
	err = json.Unmarshal([]byte(updateResponse.OperationData), &operationData)
	Expect(err).NotTo(HaveOccurred())
	return serialized(&operationData)
}

func actualBytes(actual interface{}) []byte {
	bytes, isBytes := actual.([]byte)
	Expect(isBytes).To(BeTrue(), fmt.Sprintf("converting actual to string: %s", actual))
	return bytes
}

func serialized(source *broker.OperationData) []byte {
	bytes, err := json.Marshal(source)
	Expect(err).NotTo(HaveOccurred())
	return bytes
}
