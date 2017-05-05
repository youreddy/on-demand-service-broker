// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package broker_response_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf/on-demand-service-broker/broker"
	"github.com/pivotal-cf/on-demand-service-broker/mgmtapi"
	"github.com/pivotal-cf/on-demand-service-broker/upgrader/broker_response"
)

var _ = Describe("BrokerResponse", func() {
	Context("list instances", func() {
		It("returns service instance IDs", func() {
			response := http.Response{
				StatusCode: http.StatusOK,
				Body:       asBody(listInstancesJSON("instance1-guid", "instance2-guid")),
			}

			instances, err := broker_response.ListInstancesFrom(&response)

			Expect(err).NotTo(HaveOccurred())
			Expect(instances).To(ConsistOf("instance1-guid", "instance2-guid"))
		})

		It("returns an error when the response status is not OK", func() {
			response := http.Response{
				Status:     "500 Internal Server Error",
				StatusCode: 500,
				Body:       asBody(""),
			}

			_, err := broker_response.ListInstancesFrom(&response)

			Expect(err).To(MatchError(
				ContainSubstring("HTTP response status: 500 Internal Server Error"),
			))
		})

		It("returns an error when the response body cannot be decoded", func() {
			response := http.Response{
				StatusCode: http.StatusOK,
				Body:       asBody("{ invalid json }"),
			}

			_, err := broker_response.ListInstancesFrom(&response)

			Expect(err).To(MatchError(
				ContainSubstring("invalid character"),
			))
		})
	})

	Context("last operation", func() {
		It("returns the last operation data", func() {
			expectedOperation := brokerapi.LastOperation{
				State:       brokerapi.InProgress,
				Description: "some-description",
			}
			response := http.Response{
				StatusCode: http.StatusOK,
				Body:       asBody(lastOperationJSON(expectedOperation)),
			}

			lastOperation, err := broker_response.LastOperationFrom(&response)

			Expect(err).NotTo(HaveOccurred())
			Expect(lastOperation).To(Equal(expectedOperation))
		})

		It("returns an error when the response status is not OK", func() {
			response := http.Response{
				Status:     "500 Internal Server Error",
				StatusCode: 500,
				Body:       asBody(""),
			}

			_, err := broker_response.LastOperationFrom(&response)

			Expect(err).To(MatchError(
				ContainSubstring("HTTP response status: 500 Internal Server Error"),
			))
		})

		It("returns an error when the response body cannot be decoded", func() {
			response := http.Response{
				StatusCode: http.StatusOK,
				Body:       asBody("{ invalid json }"),
			}

			_, err := broker_response.LastOperationFrom(&response)

			Expect(err).To(MatchError(ContainSubstring("invalid character")))
		})
	})

	Context("upgrade operation", func() {
		Context("when the upgrade is accepted", func() {
			It("returns the upgrade operation data", func() {
				response := http.Response{
					StatusCode: http.StatusAccepted,
					Body:       asBody(upgradeOperationJSON()),
				}

				result, err := broker_response.UpgradeOperationFrom(&response)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Data.OperationType).To(Equal(broker.OperationTypeUpgrade))
				Expect(result.Type).To(Equal(broker_response.ResultAccepted))
			})

			It("returns an error when the response body cannot be decoded", func() {
				response := http.Response{
					StatusCode: http.StatusAccepted,
					Body:       asBody("{ invalid json }"),
				}

				_, err := broker_response.UpgradeOperationFrom(&response)

				Expect(err).To(MatchError(SatisfyAll(
					ContainSubstring("cannot parse upgrade response"),
					ContainSubstring("invalid character"),
				)))
			})
		})

		Context("when the cf service instance is not found", func() {
			It("returns a not found result", func() {
				response := http.Response{
					StatusCode: http.StatusNotFound,
					Body:       asBody(""),
				}

				result, err := broker_response.UpgradeOperationFrom(&response)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Type).To(Equal(broker_response.ResultNotFound))
			})
		})

		Context("when bosh deployment for the service instance is gone", func() {
			It("returns an orphan service instance result", func() {
				response := http.Response{
					StatusCode: http.StatusGone,
					Body:       asBody(""),
				}

				result, err := broker_response.UpgradeOperationFrom(&response)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Type).To(Equal(broker_response.ResultOrphan))
			})
		})

		Context("when the service instance has an operation in progress", func() {
			It("returns an operation in progress result", func() {
				response := http.Response{
					StatusCode: http.StatusConflict,
					Body:       asBody(""),
				}

				result, err := broker_response.UpgradeOperationFrom(&response)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Type).To(Equal(broker_response.ResultOperationInProgress))
			})
		})

		Context("when the upgrade response is internal server error", func() {
			It("returns the error description", func() {
				response := http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       asBody(upgradeErrorJSON("upgrade failed")),
				}

				_, err := broker_response.UpgradeOperationFrom(&response)

				Expect(err).To(MatchError(SatisfyAll(
					ContainSubstring("unexpected status code: 500"),
					ContainSubstring("description: upgrade failed"),
				)))
			})

			It("fails to decode the response body", func() {
				response := http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       asBody("{invalid json}"),
				}

				_, err := broker_response.UpgradeOperationFrom(&response)

				Expect(err).To(MatchError(SatisfyAll(
					ContainSubstring("unexpected status code: 500"),
					ContainSubstring("cannot parse upgrade response: '{invalid json}'"),
				)))
			})
		})

		Context("when the upgrade response status code is unexpected", func() {
			It("returns an error", func() {
				response := http.Response{
					StatusCode: http.StatusTeapot,
					Body:       asBody("an unexpected error occurred"),
				}

				_, err := broker_response.UpgradeOperationFrom(&response)

				Expect(err).To(MatchError(SatisfyAll(
					ContainSubstring("unexpected status code: 418"),
					ContainSubstring("body: an unexpected error occurred"),
				)))
			})
		})
	})
})

func upgradeOperationJSON() string {
	operation := broker.OperationData{
		OperationType: broker.OperationTypeUpgrade,
	}
	content, err := json.Marshal(operation)
	Expect(err).NotTo(HaveOccurred())
	return string(content)
}

func upgradeErrorJSON(description string) string {
	errorResponse := brokerapi.ErrorResponse{
		Description: description,
	}
	content, err := json.Marshal(errorResponse)
	Expect(err).NotTo(HaveOccurred())
	return string(content)
}

func lastOperationJSON(operation brokerapi.LastOperation) string {
	content, err := json.Marshal(operation)
	Expect(err).NotTo(HaveOccurred())
	return string(content)
}

func listInstancesJSON(instanceIDs ...string) string {
	list := []mgmtapi.Instance{}
	for _, instanceID := range instanceIDs {
		list = append(list, mgmtapi.Instance{InstanceID: instanceID})
	}
	content, err := json.Marshal(list)
	Expect(err).NotTo(HaveOccurred())
	return string(content)
}

func asBody(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}