// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_new_tests

import (
	"fmt"
	"math/rand"

	"github.com/pivotal-cf/on-demand-service-broker/broker"
	"github.com/pivotal-cf/on-demand-service-broker/config"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

const (
	bindingGUIDfromCF = "Gjklh45ljkhn"

	theServiceID  = "the-service-id"
	basePlanID    = "base-plan-id"
	appGUIDfromCF = "app-guid-from-cf"
)

type BrokerEnvironment struct {
	Broker            *Broker
	Bosh              *Bosh
	CF                *CloudFoundry
	ServiceAdapter    *ServiceAdapter
	Credhub           Credhub
	serviceInstanceID ServiceInstanceID
}

func NewBrokerEnvironment(bosh *Bosh, cf *CloudFoundry, serviceAdapter *ServiceAdapter, credhub Credhub, brokerBinaryPath string) *BrokerEnvironment {
	return &BrokerEnvironment{
		Broker:            NewBroker(brokerBinaryPath),
		Bosh:              bosh,
		CF:                cf,
		ServiceAdapter:    serviceAdapter,
		Credhub:           credhub,
		serviceInstanceID: AServiceInstanceID(),
	}
}

func (be *BrokerEnvironment) Start() {
	be.CF.RespondsToInitialChecks()
	be.Bosh.RespondsToInitialChecks()
	be.Broker.Start(be.Configuration())
}

func (be *BrokerEnvironment) Configuration() *config.Config {
	return &config.Config{
		Broker:            be.Broker.Configuration(),
		Bosh:              be.Bosh.Configuration(),
		CF:                be.CF.Configuration(),
		ServiceAdapter:    be.ServiceAdapter.Configuration(),
		Credhub:           be.Credhub.Configuration(),
		ServiceCatalog:    theServiceOffering(),
		ServiceDeployment: theServiceDeployment(),
	}
}

func (be *BrokerEnvironment) Verify() {
	be.Bosh.Verify()
	be.CF.Verify()
	be.Credhub.Verify()
}

func (be *BrokerEnvironment) Close() {
	be.Broker.Close()
	be.CF.Close()
	be.Bosh.Close()
	be.Credhub.Close()
}

func (be *BrokerEnvironment) DeploymentName() string {
	return broker.DeploymentNameFrom(string(be.serviceInstanceID))
}

type ServiceInstanceID string

func AServiceInstanceID() ServiceInstanceID {
	return ServiceInstanceID(fmt.Sprintf("service-instance-ID-%d", rand.Int()))
}

func theServiceOffering() config.ServiceOffering {
	return config.ServiceOffering{

		Plans: []config.Plan{
			{
				ID: basePlanID,
				InstanceGroups: []serviceadapter.InstanceGroup{
					{
						VMType: "the-vm-type",
						Name: "the-instance-group",
						Instances: 1,
						AZs: []string{ "the-az" },
					},
				},
			},
		},
	}
}

func theServiceDeployment() config.ServiceDeployment {
	return config.ServiceDeployment{
		Stemcell: serviceadapter.Stemcell{
			OS: "ubuntu-trusty",
			Version: "10.0.01",
		},
	}
}
