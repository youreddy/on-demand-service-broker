package cf_test

import (
	"github.com/pivotal-cf/on-demand-service-broker/cf"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/on-demand-service-broker/mockhttp"
	"github.com/pivotal-cf/on-demand-service-broker/mockhttp/mockcfapi"
	"fmt"
	"github.com/onsi/gomega/gbytes"
	"log"
	"io"
)

var _ = FDescribe("UnauthenticatedClient", func() {
	It("GetAuthenticationEndpoint returns authentication url", func() {
		var server *mockhttp.Server
		server = mockcfapi.New()
		logBuffer := gbytes.NewBuffer()
		testLogger := log.New(io.MultiWriter(logBuffer, GinkgoWriter), "my-app", log.LstdFlags)
		expectedAuthURL := "non-auth-url"

		server.VerifyAndMock(
			mockcfapi.GetInfo().RespondsOKWith(
				fmt.Sprintf("{\"authorization_endpoint\": \"%s\"}", expectedAuthURL),
			),
		)

		unauthCF, err := cf.NewUnauthenticated(server.URL, nil, true)
		Expect(err).NotTo(HaveOccurred())


		authURL, err := unauthCF.GetAuthURL(testLogger)
		Expect(err).NotTo(HaveOccurred())
		Expect(authURL).To(Equal(expectedAuthURL))
		Expect(logBuffer).To(gbytes.Say(fmt.Sprintf("GET %s/v2/info", server.URL)))
	})

	It("GetAuthenticationEndpoint returns useful error message", func() {
		var server *mockhttp.Server
		server = mockcfapi.New()
		logBuffer := gbytes.NewBuffer()
		testLogger := log.New(io.MultiWriter(logBuffer, GinkgoWriter), "my-app", log.LstdFlags)

		server.VerifyAndMock(
			mockcfapi.GetInfo().RespondsOKWith("{}"),
		)

		unauthCF, err := cf.NewUnauthenticated(server.URL, nil, true)
		Expect(err).NotTo(HaveOccurred())


		authURL, err := unauthCF.GetAuthURL(testLogger)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("Non-valid auth url."))
		Expect(authURL).To(BeEmpty())
	})


	It("GetAuthenticationEndpoint deals with invalid JSON", func() {
		var server *mockhttp.Server
		server = mockcfapi.New()
		logBuffer := gbytes.NewBuffer()
		testLogger := log.New(io.MultiWriter(logBuffer, GinkgoWriter), "my-app", log.LstdFlags)

		server.VerifyAndMock(
			mockcfapi.GetInfo().RespondsOKWith("{heads :::[}"),
		)

		unauthCF, err := cf.NewUnauthenticated(server.URL, nil, true)
		Expect(err).NotTo(HaveOccurred())


		authURL, err := unauthCF.GetAuthURL(testLogger)
		Expect(err).To(HaveOccurred())
		Expect(err).To(BeAssignableToTypeOf(cf.InvalidResponseError{}))
		Expect(err.Error()).To(ContainSubstring("Invalid response body:"))
		Expect(authURL).To(BeEmpty())
	})


})
