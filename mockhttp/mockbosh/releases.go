// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package mockbosh

import "github.com/pivotal-cf/on-demand-service-broker/mockhttp"

type releasesMock struct {
	*mockhttp.Handler
}

func Releases() *releasesMock {
	return &releasesMock{
		mockhttp.NewMockedHttpRequest("GET", "/releases"),
	}
}

func (m *releasesMock) RespondsWithRequiredReleases() *mockhttp.Handler {
	return m.RespondsOKWith(`[
	  {
		"name": "required-bosh-release",
		"release_versions": [
		  {
			"version": "1.0",
			"commit_hash": "4c36884a",
			"uncommitted_changes": false,
			"currently_deployed": false,
			"job_names": [  ]
		  }
		]
	  }
	]`)
}

func (m *releasesMock) RespondsWithMissingReleases() *mockhttp.Handler {
	return m.RespondsOKWith(`[
	  {
		"name": "not-required-bosh-release",
		"release_versions": [
		  {
			"version": "1.0",
			"commit_hash": "4c36884a",
			"uncommitted_changes": false,
			"currently_deployed": false,
			"job_names": [  ]
		  }
		]
	  }
	]`)
}
