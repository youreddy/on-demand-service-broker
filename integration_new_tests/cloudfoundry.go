// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_new_tests

import (
	"github.com/pivotal-cf/on-demand-service-broker/config"
	"github.com/pivotal-cf/on-demand-service-broker/mockhttp"
	"github.com/pivotal-cf/on-demand-service-broker/mockhttp/mockcfapi"
	"github.com/pivotal-cf/on-demand-service-broker/mockuaa"
)

const (
	cfUaaClientID     = "cfAdminUsername"
	cfUaaClientSecret = "cfAdminPassword"
)

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

func (cf *CloudFoundry) RespondsToInitialChecks() {
	cf.API.VerifyAndMock(
		mockcfapi.GetInfo().RespondsWithSufficientAPIVersion(),
		mockcfapi.ListServiceOfferings().RespondsWithNoServiceOfferings(),
	)
}

func (cf *CloudFoundry) Verify() {
	cf.API.VerifyMocks()
}

func (cf *CloudFoundry) Close() {
	cf.API.Close()
	cf.UAA.Close()
}
