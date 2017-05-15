// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_new_tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	brokerBasePort = 37890
	brokerUsername = "broker-username"
	brokerPassword = "a-very-strong-password"

	spaceGUID        = "space-guid"
	organizationGUID = "organizationGuid"
)

type Broker struct {
	port             int
	brokerBinaryPath string
	session          *gexec.Session
	tempDirPath      string
}

func NewBroker(brokerBinaryPath string) *Broker {
	node := GinkgoParallelNode()
	tempDirPath, err := ioutil.TempDir("", fmt.Sprintf("broker-integration-tests-%d", node))
	Expect(err).ToNot(HaveOccurred())
	return &Broker{
		port:             brokerBasePort + node,
		tempDirPath:      tempDirPath,
		brokerBinaryPath: brokerBinaryPath,
	}
}

func (b *Broker) Start(configuration *config.Config) {
	params := []string{"-configFilePath", b.configurationFile(configuration)}
	session, err := gexec.Start(exec.Command(b.brokerBinaryPath, params...), GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gbytes.Say("listening on"))

	b.session = session
}

func (b *Broker) configurationFile(configuration *config.Config) string {
	testConfigFilePath := filepath.Join(b.tempDirPath, "broker.yml")

	configContents, err := yaml.Marshal(configuration)
	Expect(err).ToNot(HaveOccurred())
	Expect(ioutil.WriteFile(testConfigFilePath, configContents, 0644)).To(Succeed())
	return testConfigFilePath
}

func (b *Broker) Configuration() config.Broker {
	return config.Broker{
		Port:     b.port,
		Username: brokerUsername,
		Password: brokerPassword,
	}
}

func (b *Broker) Close() {
	if b.session != nil {
		b.session.Kill()
	}
	Expect(os.RemoveAll(b.tempDirPath)).To(Succeed())
}

func (b *Broker) HasLogged(expectedString string) {
	Eventually(b.session).Should(gbytes.Say(expectedString))
}

func (b *Broker) CreateBindingRequest(serviceInstanceID ServiceInstanceID) *http.Request {
	reqJson := fmt.Sprintf(`{
		"plan_id" : "%s",
		"service_id":  "%s",
		"app_guid": "%s",
		"bind_resource": { "app_guid": "%s"},
		"parameters": {"baz": "bar"}
	}`,
		basePlanID, theServiceID, appGUIDfromCF, appGUIDfromCF,
	)

	bindingReq, err := http.NewRequest("PUT",
		fmt.Sprintf("http://localhost:%d/v2/service_instances/%s/service_bindings/%s", b.port, serviceInstanceID, bindingGUIDfromCF),
		bytes.NewReader([]byte(reqJson)))
	Expect(err).ToNot(HaveOccurred())
	return withBasicAuth(bindingReq)
}

func (b *Broker) UpdateServiceInstanceRequest(serviceInstanceID ServiceInstanceID) *http.Request {
	reqJson := fmt.Sprintf(`{
		"plan_id" : "%s",
		"service_id":  "%s",
		"parameters": {"baz": "bar"},
		"previous_values": {
		  "organization_id": "%s",
			"service_id":      "%s",
			"plan_id":         "%s",
			"space_id":        "%s"
		}
	}`,
		basePlanID, theServiceID, organizationGUID, theServiceID, basePlanID, spaceGUID,
	)

	updateReq, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("http://localhost:%d/v2/service_instances/%s?accepts_incomplete=true", b.port, serviceInstanceID),
		bytes.NewReader([]byte(reqJson)),
	)
	Expect(err).NotTo(HaveOccurred())
	return withBasicAuth(updateReq)
}

func withBasicAuth(req *http.Request) *http.Request {
	req.SetBasicAuth(brokerUsername, brokerPassword)
	return req
}
