package service_test

import (
	"errors"

	"github.com/pivotal-cf/on-demand-service-broker/service"

	"io/ioutil"
	"net/http"
	"strings"

	"fmt"

	"net/url"

	"crypto/x509"

	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fakes2 "github.com/pivotal-cf/on-demand-service-broker/authorizationheader/fakes"
	"github.com/pivotal-cf/on-demand-service-broker/loggerfactory"
	"github.com/pivotal-cf/on-demand-service-broker/service/fakes"
)

var _ = Describe("ServiceInstanceLister", func() {
	var (
		client            *fakes.FakeDoer
		authHeaderBuilder *fakes2.FakeAuthHeaderBuilder
		logger            *log.Logger
	)

	BeforeEach(func() {
		client = new(fakes.FakeDoer)
		authHeaderBuilder = new(fakes2.FakeAuthHeaderBuilder)
		loggerFactory := loggerfactory.New(os.Stdout, "service-instance-lister-test", loggerfactory.Flags)
		logger = loggerFactory.New()
	})

	It("lists service instances", func() {
		client.DoReturns(response(http.StatusOK, `[{"service_instance_id": "foo", "plan_id": "plan"}, {"service_instance_id": "bar", "plan_id": "another-plan"}]`), nil)
		serviceInstanceLister := service.NewInstanceLister(client, authHeaderBuilder, "", false, logger)
		instances, err := serviceInstanceLister.Instances()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(instances)).To(Equal(2))
		Expect(instances[0]).To(Equal(service.Instance{GUID: "foo", PlanUniqueID: "plan"}))
	})

	It("returns an error when the request fails", func() {
		client.DoReturns(nil, errors.New("connection error"))
		serviceInstanceLister := service.NewInstanceLister(client, authHeaderBuilder, "", false, logger)
		_, err := serviceInstanceLister.Instances()
		Expect(err).To(MatchError("connection error"))
	})

	It("returns an error when the broker response is unrecognised", func() {
		client.DoReturns(response(http.StatusOK, `{"not": "a list"}`), nil)
		serviceInstanceLister := service.NewInstanceLister(client, authHeaderBuilder, "", false, logger)
		_, err := serviceInstanceLister.Instances()
		Expect(err).To(HaveOccurred())
	})

	It("returns an error when the HTTP status is not OK", func() {
		client.DoReturns(response(http.StatusInternalServerError, ``), nil)
		serviceInstanceLister := service.NewInstanceLister(client, authHeaderBuilder, "", false, logger)
		_, err := serviceInstanceLister.Instances()
		Expect(err).To(MatchError(fmt.Sprintf(
			"HTTP response status: %d %s",
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)))
	})

	It("returns a service instance API error when the HTTP status is not OK and service API is configured", func() {
		client.DoReturns(response(http.StatusInternalServerError, ``), nil)
		serviceInstanceLister := service.NewInstanceLister(client, authHeaderBuilder, "", true, logger)
		_, err := serviceInstanceLister.Instances()
		Expect(err).To(MatchError(fmt.Sprintf(
			"error communicating with service_instances_api (%s): HTTP response status: %d %s",
			"http://example.org/some-path",
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)))
	})

	It("returns SSL validation error when service instance API request fails due to unknown authority", func() {
		expectedURL := "https://example.org/service-instances"
		expectedError := &url.Error{
			URL: expectedURL,
			Err: x509.UnknownAuthorityError{},
		}
		client.DoReturns(nil, expectedError)
		serviceInstanceLister := service.NewInstanceLister(client, authHeaderBuilder, expectedURL, true, logger)
		_, err := serviceInstanceLister.Instances()
		Expect(err).To(MatchError(fmt.Sprintf(
			"SSL validation error for `service_instances_api.url`: %s. Please configure a `service_instances_api.root_ca_cert` and use a valid SSL certificate",
			expectedURL,
		)))
	})

	It("returns the expected error when service instance API request fails due to generic certificate error", func() {
		expectedURL := "https://example.org/service-instances"
		expectedError := &url.Error{
			URL: expectedURL,
			Err: x509.CertificateInvalidError{},
		}
		client.DoReturns(nil, expectedError)
		serviceInstanceLister := service.NewInstanceLister(client, authHeaderBuilder, expectedURL, true, logger)
		_, err := serviceInstanceLister.Instances()
		Expect(err).To(MatchError(Equal(fmt.Sprintf(
			"error communicating with service_instances_api (%s): %s",
			expectedURL,
			expectedError.Error(),
		))))
	})

	It("returns the expected error when service instance API request fails due to a url error with no Err", func() {
		expectedURL := "https://example.org/service-instances"
		expectedError := &url.Error{
			URL: expectedURL,
		}
		client.DoReturns(nil, expectedError)
		serviceInstanceLister := service.NewInstanceLister(client, authHeaderBuilder, expectedURL, true, logger)
		_, err := serviceInstanceLister.Instances()
		Expect(err).To(Equal(expectedError))
	})

	It("passes expected request to authHeaderBuilder", func() {
		expectedURL := "https://example.org/service-instances"
		client.DoReturns(response(http.StatusOK, `[]`), nil)
		serviceInstanceLister := service.NewInstanceLister(client, authHeaderBuilder, expectedURL, false, logger)
		serviceInstanceLister.Instances()
		requestForAuth, _ := authHeaderBuilder.AddAuthHeaderArgsForCall(0)
		expectedRequest, err := http.NewRequest(
			http.MethodGet,
			expectedURL,
			nil,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(requestForAuth).To(Equal(expectedRequest))
	})

	It("passes expected request to httpClient", func() {
		expectedURL := "https://example.org/service-instances"
		client.DoReturns(response(http.StatusOK, `[]`), nil)
		serviceInstanceLister := service.NewInstanceLister(client, authHeaderBuilder, expectedURL, false, logger)
		serviceInstanceLister.Instances()
		requestForClient := client.DoArgsForCall(0)
		expectedRequest, err := http.NewRequest(
			http.MethodGet,
			expectedURL,
			nil,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(requestForClient).To(Equal(expectedRequest))
	})

	It("passes logger to authHeaderBuilder", func() {
		expectedURL := "https://example.org/service-instances"
		client.DoReturns(response(http.StatusOK, `[]`), nil)
		serviceInstanceLister := service.NewInstanceLister(client, authHeaderBuilder, expectedURL, false, logger)
		serviceInstanceLister.Instances()
		_, actualLogger := authHeaderBuilder.AddAuthHeaderArgsForCall(0)
		Expect(actualLogger).To(BeIdenticalTo(logger))
	})
})

func response(statusCode int, body string) *http.Response {
	parsedUrl, err := url.Parse("http://example.org/some-path")
	Expect(err).NotTo(HaveOccurred())
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Body:       ioutil.NopCloser(strings.NewReader(body)),
		Request: &http.Request{
			URL: parsedUrl,
		},
	}
}
