// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package runner

import (
	"fmt"
	"log"

	"time"

	"github.com/pivotal-cf/on-demand-service-broker/cf"
	"github.com/pivotal-cf/on-demand-service-broker/config"
	"github.com/pivotal-cf/on-demand-service-broker/service"
	"github.com/cloudfoundry/bosh-deployment-resource/tools"
)


//go:generate counterfeiter -o fakes/fake_sleeper.go . Sleeper
type Sleeper interface {
	Sleep(d time.Duration)
}

type Config struct {
	BrokerAPI           BrokerAPI           `yaml:"broker_api"`
	ServiceInstancesAPI ServiceInstancesAPI `yaml:"service_instances_api"`
	PollingInterval     int                 `yaml:"polling_interval"`
	AttemptInterval     int                 `yaml:"attempt_interval"`
	AttemptLimit        int                 `yaml:"attempt_limit"`
}

type ServiceCatalog struct {
	ID string `yaml:"id"`
}

type Deleter struct {
	logger               *log.Logger
	pollingInitialOffset time.Duration
	pollingInterval      time.Duration
	cfClient             CloudFoundryClient
	sleeper              Sleeper
}

func NewBuilder(
	conf config.UpgradeAllInstanceErrandConfig,
	logger *log.Logger,
) (*Builder, error) {

	brokerServices, err := brokerServices(conf, logger)
	if err != nil {
		return nil, err
	}

	instanceLister, err := serviceInstanceLister(conf, logger)
	if err != nil {
		return nil, err
	}

	pollingInterval, err := pollingInterval(conf)
	if err != nil {
		return nil, err
	}

	attemptInterval, err := attemptInterval(conf)
	if err != nil {
		return nil, err
	}

	attemptLimit, err := attemptLimit(conf)
	if err != nil {
		return nil, err
	}

	listener := NewLoggingListener(logger)

	b := &Builder{
		BrokerServices:        brokerServices,
		ServiceInstanceLister: instanceLister,
		PollingInterval:       pollingInterval,
		AttemptInterval:       attemptInterval,
		AttemptLimit:          attemptLimit,
		Listener:              listener,
		Sleeper:               &tools.RealSleeper{},
	}

	return b, nil
}

func (d *Deleter) InspectAllServiceInstances(serviceUniqueID string) error {
	d.logger.Printf("Deleter Configuration: polling_intial_offset: %v, polling_interval: %v.", d.pollingInitialOffset.Seconds(), d.pollingInterval.Seconds())
	serviceInstances, err := d.cfClient.GetInstancesOfServiceOffering(serviceUniqueID, d.logger)
	if err != nil {
		return err
	}

	if len(serviceInstances) == 0 {
		d.logger.Println("No service instances found.")
		return nil
	}

	for _, instance := range serviceInstances {
		err = d.deleteBindings(instance.GUID)
		if err != nil {
			return err
		}

		err = d.deleteServiceKeys(instance.GUID)
		if err != nil {
			return err
		}

		err = d.deleteServiceInstance(instance.GUID)
		if err != nil {
			return err
		}

		d.logger.Printf("Waiting for service instance %s to be deleted", instance.GUID)

		err = d.pollInstanceDeleteStatus(instance.GUID)
		if err != nil {
			return err
		}
	}

	serviceInstances, err = d.cfClient.GetInstancesOfServiceOffering(serviceUniqueID, d.logger)
	if err != nil {
		return err
	}

	if len(serviceInstances) != 0 {
		return fmt.Errorf("expected 0 instances for service offering with unique ID: %s. Got %d instance(s).", serviceUniqueID, len(serviceInstances))
	}

	return nil
}

