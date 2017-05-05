// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_tests

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/on-demand-service-broker/boshclient"
	"github.com/pivotal-cf/on-demand-service-broker/integration_tests/mock"
	"github.com/pivotal-cf/on-demand-service-broker/mockbosh"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
)

var (
	brokerPath         = NewBinary("github.com/pivotal-cf/on-demand-service-broker/cmd/on-demand-service-broker")
	serviceAdapterPath = NewBinary("github.com/pivotal-cf/on-demand-service-broker/integration_tests/mock/adapter")
	adapter            = mock.Adapter{}
)

var _ = Describe("binding service instances", func() {
	It("binds a service to an application instance", func() {
		manifestForFirstDeployment := bosh.BoshManifest{
			Name:           deploymentName(instanceID),
			Releases:       []bosh.Release{},
			Stemcells:      []bosh.Stemcell{},
			InstanceGroups: []bosh.InstanceGroup{},
		}

		adapter.New()
		adapter.CreateBinding().ReturnsBinding(`{
					"credentials": {"secret": "dont-tell-anyone"},
					"syslog_drain_url": "syslog-url",
					"route_service_url": "excellent route"
					}`)

		withBroker(func(b *Broker) {
			b.Bosh.Director.VerifyAndMock(
				mockbosh.VMsForDeployment(deploymentName(instanceID)).RedirectsToTask(2015),
				mockbosh.Task(2015).RespondsWithTaskContainingState(boshclient.BoshTaskDone),
				mockbosh.TaskOutput(2015).RespondsWithVMsOutput([]boshclient.BoshVMsOutput{{IPs: []string{"ip.from.bosh"}, InstanceGroup: "some-instance-group"}}),
				mockbosh.GetDeployment(deploymentName(instanceID)).RespondsWithManifest(manifestForFirstDeployment),
			)

			response, err := http.DefaultClient.Do(b.CreationRequest())
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusCreated))

			// request a new binding from service to application application

			// responds with Created and the binding details
			// logs the bind request with an ID
		})
	})

})

func withBroker(test func(*Broker)) {
	broker := NewBroker(NewBosh(), NewCloudFoundry(), NewServiceAdapter(serviceAdapterPath.Path()), brokerPath.Path())
	defer broker.Close()
	broker.Start()
	test(broker)
	broker.Verify()
}

func deploymentName(instanceID string) string {
	return "service-instance_" + instanceID
}
