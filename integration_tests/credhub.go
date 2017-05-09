// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_tests

import (
	"fmt"

	"github.com/pivotal-cf/on-demand-service-broker/config"
	"github.com/pivotal-cf/on-demand-service-broker/mockcredhub"
	"github.com/pivotal-cf/on-demand-service-broker/mockhttp"
	"github.com/pivotal-cf/on-demand-service-broker/mockuaa"
)

const (
	credhubClientID     = "credhubAdminUsername"
	credhubClientSecret = "credhubAdminPassword"
)

type Credhub interface {
	Configuration() *config.Credhub
	Verify()
	Close()
}

type MockCredhub struct {
	Credhub    *mockhttp.Server
	credhubUaa *mockuaa.UserCredentialsServer
}

func NewCredhub() *MockCredhub {
	return &MockCredhub{
		Credhub:    mockcredhub.New(),
		credhubUaa: mockuaa.NewUserCredentialsServer("credhub", "", credhubClientID, credhubClientSecret, "Credhub token"),
	}
}

func (c *MockCredhub) WillReceiveCredentials(serviceInstanceID ServiceInstanceID) {
	credhubBindingServiceId := fmt.Sprintf("%s/%s", serviceInstanceID, bindingGUIDfromCF)

	c.Credhub.VerifyAndMock(
		mockcredhub.GetInfo().RespondsWithUAAURL(c.credhubUaa.URL),
		mockcredhub.PutCredential(credhubBindingServiceId).WithPassword(BindingCredentials).RespondsWithPasswordData(BindingCredentials),
	)
}

func (c *MockCredhub) Configuration() *config.Credhub {
	return &config.Credhub{
		APIURL: c.Credhub.URL,
		ID:     credhubClientID,
		Secret: credhubClientSecret,
	}
}

func (c *MockCredhub) Verify() {
	c.Credhub.VerifyMocks()
}

func (c *MockCredhub) Close() {
	c.Credhub.Close()
	c.credhubUaa.Close()
}

type noop struct{}

func (n *noop) Configuration() *config.Credhub { return nil }
func (n *noop) Verify()                        {}
func (n *noop) Close()                         {}

var NoCredhub = new(noop)
