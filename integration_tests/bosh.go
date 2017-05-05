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
	boshClientID       = "bosh-client-id"
	boshClientSecret   = "boshClientSecret"
	boshDeploymentName = "service-instance_some-service-instance-ID"
	boshVMDescription  = `{"IPs" : ["ip.from.bosh"], "job_name": "some-instance-group"}`
)

var manifestForFirstDeployment = bosh.BoshManifest{
	Name:           boshDeploymentName,
	Releases:       []bosh.Release{},
	Stemcells:      []bosh.Stemcell{},
	InstanceGroups: []bosh.InstanceGroup{},
}

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

func (b *Bosh) ReturnsDeployment() {
	taskID := rand.Int()
	b.Director.VerifyAndMock(
		mockbosh.VMsForDeployment(boshDeploymentName).RedirectsToTask(taskID),
		mockbosh.Task(taskID).RespondsWithTaskContainingState(boshclient.BoshTaskDone),
		mockbosh.TaskOutput(taskID).RespondsWithBody(boshVMDescription),
		mockbosh.GetDeployment(boshDeploymentName).RespondsWithManifest(manifestForFirstDeployment),
	)
}
