// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf/on-demand-service-broker/config"
	"github.com/pivotal-cf/on-demand-service-broker/mockbosh"
	"github.com/pivotal-cf/on-demand-service-broker/mockcfapi"
	"github.com/pivotal-cf/on-demand-service-broker/mockhttp"
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

type Bosh struct {
	Director *mockhttp.Server
	UAA      *mockuaa.ClientCredentialsServer
}

func NewBosh() *Bosh {
	return &Bosh{
		Director: mockbosh.New(),
		UAA:      mockuaa.NewClientCredentialsServer(boshClientID, boshClientSecret, "bosh uaa token"),
	}
}

func (b *Bosh) Configuration() config.Bosh {
	return config.Bosh{
		URL: b.Director.URL,
		Authentication: config.BOSHAuthentication{
			UAA: config.BOSHUAAAuthentication{UAAURL: b.UAA.URL, ID: boshClientID, Secret: boshClientSecret},
		},
	}
}

type CloudFoundry struct {
	API *mockhttp.Server
	UAA *mockuaa.ClientCredentialsServer
}

func NewCloudFoundry() *CloudFoundry {
	return &CloudFoundry{
		API: mockcfapi.New(),
		UAA: mockuaa.NewClientCredentialsServer(cfUaaClientID, cfUaaClientSecret, "CF UAA token"),
	}
}

func (cf *CloudFoundry) Configuration() config.CF {
	return config.CF{
		URL: cf.API.URL,
		Authentication: config.UAAAuthentication{
			URL: cf.UAA.URL,
			ClientCredentials: config.ClientCredentials{
				ID:     cfUaaClientID,
				Secret: cfUaaClientSecret,
			},
		},
	}
}

type ServiceAdapter struct {
	Path string
}

func NewServiceAdapter(path string) *ServiceAdapter {
	return &ServiceAdapter{
		Path: path,
	}
}

func (sa *ServiceAdapter) Configuration() config.ServiceAdapter {
	return config.ServiceAdapter{
		Path: sa.Path,
	}
}

type Broker struct {
	Bosh           *Bosh
	CF             *CloudFoundry
	ServiceAdapter *ServiceAdapter
	Path           string
	tempDirPath    string
	Session        *gexec.Session
}

func NewBroker(bosh *Bosh, cf *CloudFoundry, serviceAdapter *ServiceAdapter, path string) *Broker {
	tempDirPath, err := ioutil.TempDir("", fmt.Sprintf("broker-integration-tests-%d", GinkgoParallelNode()))
	Expect(err).ToNot(HaveOccurred())

	return &Broker{
		Bosh:           bosh,
		CF:             cf,
		ServiceAdapter: serviceAdapter,
		Path:           path,
		tempDirPath:    tempDirPath,
	}
}

func (b *Broker) Start() {
	params := []string{"-configFilePath", b.writeConfigurationToFile()}

	session, err := gexec.Start(exec.Command(brokerPath, params...), GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gbytes.Say("listening on"))

	b.Session = session
}

func (b *Broker) writeConfigurationToFile() string {
	testConfigFilePath := filepath.Join(b.tempDirPath, "broker.yml")

	configContents, err := yaml.Marshal(b.configuration())
	Expect(err).ToNot(HaveOccurred())
	Expect(ioutil.WriteFile(testConfigFilePath, configContents, 0644)).To(Succeed())
	return testConfigFilePath
}

func (b *Broker) configuration() config.Config {
	return config.Config{
		Broker: config.Broker{
			Port:          brokerPort,
			Username:      "boshUsername",
			Password:      "boshPassword",
			StartUpBanner: false,
		},
		Bosh:           b.Bosh.Configuration(),
		CF:             b.CF.Configuration(),
		ServiceAdapter: b.ServiceAdapter.Configuration(),
	}
}

func (b *Broker) Cleanup() {
	if b.Session != nil {
		b.Session.Kill()
	}
	Expect(os.RemoveAll(b.tempDirPath)).To(Succeed())
}

var (
	brokerPath         string
	serviceAdapterPath string
)

var _ = Describe("binding service instances", func() {
	BeforeSuite(func() {
		var err error
		brokerPath, err = gexec.Build("github.com/pivotal-cf/on-demand-service-broker/cmd/on-demand-service-broker")
		Expect(err).NotTo(HaveOccurred())

		serviceAdapterPath, err = gexec.Build("github.com/pivotal-cf/on-demand-service-broker/old_integration_tests/mock/adapter")
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

func withBroker(body func(*Broker)) {
	bosh := NewBosh()
	cloudFoundry := NewCloudFoundry()
	serviceAdapter := NewServiceAdapter(serviceAdapterPath)

	theBroker := NewBroker(bosh, cloudFoundry, serviceAdapter, brokerPath)
	defer theBroker.Cleanup()
	theBroker.Start()
	body(theBroker)
}
