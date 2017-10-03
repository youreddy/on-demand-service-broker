// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package boshdirector_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/on-demand-service-broker/mockhttp/mockbosh"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
)

var _ = Describe("get releases", func() {
	It("succeeds", func() {
		expectedReleases := []bosh.Release{
			{
				Name:    "release",
				Version: "1.0",
			},
		}

		director.VerifyAndMock(
			mockbosh.Releases().RespondsOKWithJSON(expectedReleases),
		)

		actualReleases, err := c.GetReleases(logger)
		Expect(err).NotTo(HaveOccurred())
		Expect(actualReleases).To(Equal(expectedReleases))
	})

	It("returns an error when an unexpected status is returned", func() {
		director.VerifyAndMock(
			mockbosh.Releases().RespondsInternalServerErrorWith("error"),
		)

		_, err := c.GetReleases(logger)
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("expected status 200, was 500. Response Body: error"))
	})

	It("returns an error when the bosh director cannot be reached", func() {
		director.Close()

		_, err := c.GetReleases(logger)
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("error reaching bosh director")))
	})
})
