package credhubbroker_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
	apifakes "github.com/pivotal-cf/on-demand-service-broker/apiserver/fakes"
	"github.com/pivotal-cf/on-demand-service-broker/credhubbroker"
	credfakes "github.com/pivotal-cf/on-demand-service-broker/credhubbroker/fakes"
)

var _ = Describe("CredHub broker", func() {
	Describe("binding", func() {
		It("stores key-value credentials in CredHub", func() {
			fakeBroker := new(apifakes.FakeCombinedBrokers)
			fakeCredStore := new(credfakes.FakeCredentialStore)

			creds := map[string]interface{}{
				"foo": 42,
				"bar": "IPA",
			}

			bindingResponse := brokerapi.Binding{
				Credentials: creds,
			}
			fakeBroker.BindReturns(bindingResponse, nil)

			credhubBroker := credhubbroker.New(fakeBroker, fakeCredStore)

			ctx := new(context.Context)
			instanceID := "ohai"
			bindingID := "rofl"
			details := brokerapi.BindDetails{
				ServiceID: "big-hybrid-cloud-of-things",
			}

			credhubBroker.Bind(*ctx, instanceID, bindingID, details)

			key, receivedCreds := fakeCredStore.SetArgsForCall(0)
			Expect(key).To(Equal("/c/big-hybrid-cloud-of-things/ohai/rofl/credentials"))
			Expect(receivedCreds).To(Equal(creds))
		})

		It("stores string credentials in CredHub", func() {
			fakeBroker := new(apifakes.FakeCombinedBrokers)
			fakeCredStore := new(credfakes.FakeCredentialStore)

			creds := "justAString"

			bindingResponse := brokerapi.Binding{
				Credentials: creds,
			}
			fakeBroker.BindReturns(bindingResponse, nil)

			credhubBroker := credhubbroker.New(fakeBroker, fakeCredStore)

			ctx := new(context.Context)
			instanceID := "ohai"
			bindingID := "rofl"
			details := brokerapi.BindDetails{
				ServiceID: "big-hybrid-cloud-of-things",
			}

			credhubBroker.Bind(*ctx, instanceID, bindingID, details)

			key, receivedCreds := fakeCredStore.SetArgsForCall(0)
			Expect(key).To(Equal("/c/big-hybrid-cloud-of-things/ohai/rofl/credentials"))
			Expect(receivedCreds).To(Equal(creds))
		})

		It("passes the return value through from the wrapped broker", func() {
			fakeBroker := new(apifakes.FakeCombinedBrokers)
			fakeCredStore := new(credfakes.FakeCredentialStore)

			expectedBindingResponse := brokerapi.Binding{
				Credentials: "anything",
			}
			fakeBroker.BindReturns(expectedBindingResponse, nil)

			credhubBroker := credhubbroker.New(fakeBroker, fakeCredStore)

			ctx := new(context.Context)
			instanceID := "ohai"
			bindingID := "rofl"
			details := brokerapi.BindDetails{}

			Expect(credhubBroker.Bind(*ctx, instanceID, bindingID, details)).To(Equal(expectedBindingResponse))
		})

		XContext("when cannot store credentials in credhub", func() {
			It("calls unbind on the wrapped broker", func() {

			})

			It("returns an error", func() {

			})
		})

	})
})
