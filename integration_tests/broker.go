// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_tests

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
	brokerPort         = 37890
	brokerUsername     = "broker-username"
	brokerPassword     = "a-very-strong-password"
	aServiceInstanceID = "some-service-instance-ID"
	bindingId          = "Gjklh45ljkhn"

	bindingPlanID    = "plan-guid-from-cc"
	bindingServiceID = "service-guid-from-cc"
	appGUID          = "app-guid-from-cc"
)

//TODO should be a TestContext?

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
	b.CF.RespondsToInitialChecks()
	b.Bosh.RespondsToInitialChecks()

	params := []string{"-configFilePath", b.configurationFile()}
	session, err := gexec.Start(exec.Command(b.Path, params...), GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gbytes.Say("listening on"))

	b.Session = session
}

func (b *Broker) configurationFile() string {
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
			Username:      brokerUsername,
			Password:      brokerPassword,
			StartUpBanner: false,
		},
		Bosh:           b.Bosh.Configuration(),
		CF:             b.CF.Configuration(),
		ServiceAdapter: b.ServiceAdapter.Configuration(),
	}
}

func (b *Broker) Verify() {
	b.Bosh.Verify()
	b.CF.Verify()
}

func (b *Broker) Close() {
	if b.Session != nil {
		b.Session.Kill()
	}
	b.CF.Close()
	b.Bosh.Close()
	Expect(os.RemoveAll(b.tempDirPath)).To(Succeed())
}

func (b *Broker) CreationRequest() *http.Request {
	reqJson := fmt.Sprintf(`{
		"plan_id" : "%s",
		"service_id":  "%s",
		"app_guid": "%s",
		"bind_resource": { "app_guid": "%s"},
		"parameters": {"baz": "bar"}
	}`,
		bindingPlanID, bindingServiceID, appGUID, appGUID,
	)

	bindingReq, err := http.NewRequest("PUT",
		fmt.Sprintf("http://localhost:%d/v2/service_instances/%s/service_bindings/%s", brokerPort, aServiceInstanceID, bindingId),
		bytes.NewReader([]byte(reqJson)))
	Expect(err).ToNot(HaveOccurred())
	return withBasicAuth(bindingReq)
}

func withBasicAuth(req *http.Request) *http.Request {
	req.SetBasicAuth(brokerUsername, brokerPassword)
	return req
}
