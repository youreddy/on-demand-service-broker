// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_tests

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
)

var (
	brokerPath         = NewBinary("github.com/pivotal-cf/on-demand-service-broker/cmd/on-demand-service-broker")
	serviceAdapterPath = NewBinary("github.com/pivotal-cf/on-demand-service-broker/integration_tests/mock/adapter")
)

var _ = Describe("binding service instances", func() {
	It("binds a service to an application instance", func() {
		withBroker(func(b *Broker) {
			b.Bosh.ReturnsDeployment()

			response, err := http.DefaultClient.Do(b.CreationRequest())
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusCreated))
			Expect(bodyOf(response)).To(MatchJSON(BindingResponse))

			Expect(b.ServiceAdapter.adapter.CreateBinding().ReceivedID()).To(Equal("Gjklh45ljkhn"))
			Expect(b.ServiceAdapter.adapter.CreateBinding().ReceivedBoshVms()).To(Equal(bosh.BoshVMs{"some-instance-group": []string{"ip.from.bosh"}}))
			Expect(b.ServiceAdapter.adapter.CreateBinding().ReceivedRequestParameters()).To(Equal(map[string]interface{}{
				"plan_id":    bindingPlanID,
				"service_id": bindingServiceID,
				"app_guid":   appGUID,
				"bind_resource": map[string]interface{}{
					"app_guid": appGUID,
				},
				"parameters": map[string]interface{}{"baz": "bar"},
			}))
			Expect(b.ServiceAdapter.adapter.CreateBinding().ReceivedManifest()).To(Equal(manifestForFirstDeployment))

			// logs the bind request with a request id
		})
	})

})

func withBroker(test func(*Broker)) {
	broker := NewBroker(NewBosh(), NewCloudFoundry(), NewServiceAdapter(serviceAdapterPath.Path()), brokerPath.Path())
	defer broker.Close()
	broker.ServiceAdapter.ReturnsBinding()
	broker.Start()
	test(broker)
	broker.Verify()
}

func bodyOf(response *http.Response) []byte {
	body, err := ioutil.ReadAll(response.Body)
	Expect(err).NotTo(HaveOccurred())
	return body
}
