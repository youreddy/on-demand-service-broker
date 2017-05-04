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
	"github.com/pivotal-cf/on-demand-service-broker/mockbosh"
	"github.com/pivotal-cf/on-demand-service-broker/mockcfapi"
	"github.com/pivotal-cf/on-demand-service-broker/mockuaa"
	"gopkg.in/yaml.v2"
)

const (
	brokerPort        = 37890
	boshClientID      = "bosh-client-id"
	boshClientSecret  = "boshClientSecret"
	cfUaaClientID     = "cfAdminUsername"
	cfUaaClientSecret = "cfAdminPassword"
)

var (
	brokerPath         string
	serviceAdapterPath string
	brokerSession      *gexec.Session
)

var _ = Describe("binding service instances", func() {
	BeforeSuite(func() {
		var err error
		brokerPath, err = gexec.Build("github.com/pivotal-cf/on-demand-service-broker/cmd/on-demand-service-broker")
		Expect(err).NotTo(HaveOccurred())

		serviceAdapterPath, err = gexec.Build("github.com/pivotal-cf/on-demand-service-broker/old_integration_tests/mock/adapter")
		Expect(err).NotTo(HaveOccurred())

	})

	BeforeEach(func() {
		boshDirector := mockbosh.New()
		boshUAA := mockuaa.NewClientCredentialsServer(boshClientID, boshClientSecret, "bosh uaa token")
		cfAPI := mockcfapi.New()
		cfUAA := mockuaa.NewClientCredentialsServer(cfUaaClientID, cfUaaClientSecret, "CF UAA token")

		brokerSession = startBroker(brokerPath, boshDirector.URL, boshUAA.URL, cfAPI.URL, cfUAA.URL, serviceAdapterPath)
	})

	AfterEach(func() {
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

func startBroker(brokerPath, boshURL, boshUaaURL, cfURL, cfUAAURL, serviceAdapterPath string) *gexec.Session {
	configContents, err := yaml.Marshal(brokerConfig(boshURL, boshUaaURL, cfURL, cfUAAURL, serviceAdapterPath))
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

func brokerConfig(boshURL, boshUaaURL, cfURL, cfUAAURL, serviceAdapterPath string) config.Config {
	return config.Config{
		Broker: config.Broker{
			Port:          brokerPort,
			Username:      "boshUsername",
			Password:      "boshPassword",
			StartUpBanner: false,
		},
		Bosh: config.Bosh{
			URL: boshURL,
			Authentication: config.BOSHAuthentication{
				UAA: config.BOSHUAAAuthentication{UAAURL: boshUaaURL, ID: boshClientID, Secret: boshClientSecret},
			},
		},
		CF: config.CF{
			URL: cfURL,
			Authentication: config.UAAAuthentication{
				URL: cfUAAURL,
				ClientCredentials: config.ClientCredentials{
					ID:     cfUaaClientID,
					Secret: cfUaaClientSecret,
				},
			},
		},
		ServiceAdapter: config.ServiceAdapter{
			Path: serviceAdapterPath,
		},
	}
}
