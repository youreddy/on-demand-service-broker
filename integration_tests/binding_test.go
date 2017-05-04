// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	brokerPath         string
	serviceAdapterPath string
)

var _ = Describe("binding service instances", func() {
	BeforeSuite(func() {
		var err error
		brokerPath, err = gexec.Build("github.com/pivotal-cf/on-demand-service-broker/cmd/on-demand-service-broker")
		Expect(err).NotTo(HaveOccurred())

		serviceAdapterPath, err = gexec.Build("github.com/pivotal-cf/on-demand-service-broker/integration_tests/mock/adapter")
		Expect(err).NotTo(HaveOccurred())

	})

	It("binds a service to an application instance", func() {
		withBroker(func(b *Broker) {
			Expect(true).NotTo(BeTrue(), "in the test")
			// request a new binding from service to application application

			// responds with Created and the binding details
			// logs the bind request with an ID
		})
	})

})

func withBroker(test func(*Broker)) {
	broker := NewBroker(NewBosh(), NewCloudFoundry(), NewServiceAdapter(serviceAdapterPath), brokerPath)
	defer broker.Close()
	broker.Start()
	test(broker)
	broker.Verify()
}
