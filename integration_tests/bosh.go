// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_tests

import (
	"math/rand"

	"github.com/pivotal-cf/on-demand-service-broker/boshclient"
	"github.com/pivotal-cf/on-demand-service-broker/config"
	"github.com/pivotal-cf/on-demand-service-broker/mockbosh"
	"github.com/pivotal-cf/on-demand-service-broker/mockhttp"
	"github.com/pivotal-cf/on-demand-service-broker/mockuaa"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
)

const (
	boshClientID      = "bosh-client-id"
	boshClientSecret  = "boshClientSecret"
	boshVMDescription = `{"IPs" : ["ip.from.bosh"], "job_name": "some-instance-group"}`
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

func (b *Bosh) RespondsToInitialChecks() {
	b.Director.VerifyAndMock(mockbosh.Info().RespondsWithSufficientVersionForLifecycleErrands())
}

func (b *Bosh) Verify() {
	b.Director.VerifyMocks()
}

func (b *Bosh) Close() {
	b.UAA.Close()
	b.Director.Close()
}

func (b *Bosh) HasManifestFor(deploymentName string) {
	b.Director.AppendMocks(
		mockbosh.GetDeployment(deploymentName).RespondsWithManifest(&bosh.BoshManifest{Name: deploymentName}),
	)
}

func (b *Bosh) HasVMsFor(deploymentName string) {
	taskID := rand.Int()

	b.Director.AppendMocks(
		mockbosh.VMsForDeployment(deploymentName).RedirectsToTask(taskID),
		mockbosh.Task(taskID).RespondsWithTaskContainingState(boshclient.BoshTaskDone),
		mockbosh.TaskOutput(taskID).RespondsWithBody(boshVMDescription),
	)
}

func (b *Bosh) HasNoVMsFor(deploymentName string) {
	b.Director.AppendMocks(
		mockbosh.VMsForDeployment(deploymentName).RespondsNotFoundWith(""),
	)
}

func (b *Bosh) DeploysWithoutContextId(deploymentName string, taskID int) {
	b.Director.AppendMocks(
		mockbosh.Deploy().WithManifest(bosh.BoshManifest{Name: deploymentName}).WithoutContextID().RedirectsToTask(taskID),
	)
}

func (b *Bosh) HasNoTasksFor(deploymentName string) {
	b.Director.AppendMocks(
		mockbosh.Tasks(deploymentName).RespondsWithNoTasks(),
	)
}
