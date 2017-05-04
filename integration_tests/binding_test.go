// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_tests

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf/on-demand-service-broker/config"
	"gopkg.in/yaml.v2"
)

const (
	brokerPort = 37890
)

var (
	brokerSession *gexec.Session
)

var _ = Describe("binding service instances", func() {
	BeforeSuite(func() {
		brokerPath, err := gexec.Build("github.com/pivotal-cf/on-demand-service-broker/cmd/on-demand-service-broker")
		Expect(err).NotTo(HaveOccurred())

		brokerSession = startBroker(brokerPath)
	})

	AfterSuite(func() {
		if brokerSession != nil {
			brokerSession.Kill()
		}
	})

	It("binds a service to an application instance", func() {
		// director, with authentication
		// CF, with authentication
		// service adapter
		// create broker

		// request a new binding from service to application application

		// responds with Created and the binding details
		// logs the bind request with an ID
	})

})

func startBroker(brokerPath string) *gexec.Session {
	configContents, err := yaml.Marshal(brokerConfig)
	Expect(err).ToNot(HaveOccurred())

	tempDirPath, err := ioutil.TempDir("", fmt.Sprintf("broker-integration-tests-%d", GinkgoParallelNode()))
	Expect(err).ToNot(HaveOccurred())

	testConfigFilePath := filepath.Join(tempDirPath, "broker.yml")
	Expect(ioutil.WriteFile(testConfigFilePath, configContents, 0644)).To(Succeed())

	params := []string{"-configFilePath", testConfigFilePath}

	cmd := exec.Command(brokerPath, params...)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gbytes.Say("listening on"))

	return session
}

var brokerConfig = config.Config{
	Broker: config.Broker{
		Port:          brokerPort,
		Username:      "boshUsername",
		Password:      "boshPassword",
		StartUpBanner: false,
	},
}
