package on_demand_service_broker_test

import (
	"fmt"

	brokerConfig "github.com/pivotal-cf/on-demand-service-broker/config"
	sdk "github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"

	"net/http"

	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
)

var _ = Describe("Catalog", func() {
	Context("without optional fields", func() {
		BeforeEach(func() {
			serviceCatalogConfig := defaultServiceCatalogConfig()
			serviceCatalogConfig.DashboardClient = nil
			conf := brokerConfig.Config{
				Broker: brokerConfig.Broker{
					Port: serverPort, Username: brokerUsername, Password: brokerPassword,
				},
				ServiceCatalog: serviceCatalogConfig,
			}

			StartServer(conf)
		})

		It("returns catalog metadata", func() {
			response := doCatalogRequest()

			By("returning the correct HTTP status")
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			By("returning the correct catalog response")
			catalog := make(map[string][]brokerapi.Service)
			Expect(json.NewDecoder(response.Body).Decode(&catalog)).To(Succeed())
			Expect(catalog).To(Equal(map[string][]brokerapi.Service{
				"services": {
					{
						ID:            serviceID,
						Name:          serviceName,
						Description:   serviceDescription,
						Bindable:      serviceBindable,
						PlanUpdatable: servicePlanUpdatable,
						Metadata: &brokerapi.ServiceMetadata{
							DisplayName:         serviceMetadataDisplayName,
							ImageUrl:            serviceMetadataImageURL,
							LongDescription:     serviceMetaDataLongDescription,
							ProviderDisplayName: serviceMetaDataProviderDisplayName,
							DocumentationUrl:    serviceMetaDataDocumentationURL,
							SupportUrl:          serviceMetaDataSupportURL,
							Shareable:           &trueVar,
						},
						DashboardClient: nil,
						Tags:            serviceTags,
						Plans: []brokerapi.ServicePlan{
							{
								ID:          dedicatedPlanID,
								Name:        dedicatedPlanName,
								Description: dedicatedPlanDescription,
								Free:        &trueVar,
								Bindable:    &trueVar,
								Metadata: &brokerapi.ServicePlanMetadata{
									Bullets:     dedicatedPlanBullets,
									DisplayName: dedicatedPlanDisplayName,
									Costs: []brokerapi.ServicePlanCost{
										{
											Unit:   dedicatedPlanCostUnit,
											Amount: dedicatedPlanCostAmount,
										},
									},
								},
							},
							{
								ID:          highMemoryPlanID,
								Name:        highMemoryPlanName,
								Description: highMemoryPlanDescription,
								Metadata: &brokerapi.ServicePlanMetadata{
									Bullets:     highMemoryPlanBullets,
									DisplayName: highMemoryPlanDisplayName,
								},
							},
						},
					},
				},
			}))
		})
	})

	Context("with optional fields", func() {
		BeforeEach(func() {
			serviceCatalogConfig := defaultServiceCatalogConfig()
			serviceCatalogConfig.Requires = []string{"syslog_drain", "route_forwarding"}
			conf := brokerConfig.Config{
				Broker: brokerConfig.Broker{
					Port: serverPort, Username: brokerUsername, Password: brokerPassword,
				},
				ServiceCatalog: serviceCatalogConfig,
			}

			StartServer(conf)
		})

		It("returns catalog metadata", func() {
			response := doCatalogRequest()

			By("returning the correct HTTP status")
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			By("returning the correct catalog response")
			catalog := make(map[string][]brokerapi.Service)
			Expect(json.NewDecoder(response.Body).Decode(&catalog)).To(Succeed())
			Expect(catalog).To(Equal(map[string][]brokerapi.Service{
				"services": {
					{
						ID:            serviceID,
						Name:          serviceName,
						Description:   serviceDescription,
						Bindable:      serviceBindable,
						PlanUpdatable: servicePlanUpdatable,
						Metadata: &brokerapi.ServiceMetadata{
							DisplayName:         serviceMetadataDisplayName,
							ImageUrl:            serviceMetadataImageURL,
							LongDescription:     serviceMetaDataLongDescription,
							ProviderDisplayName: serviceMetaDataProviderDisplayName,
							DocumentationUrl:    serviceMetaDataDocumentationURL,
							SupportUrl:          serviceMetaDataSupportURL,
							Shareable:           &trueVar,
						},
						DashboardClient: &brokerapi.ServiceDashboardClient{
							ID:          "client-id-1",
							Secret:      "secret-1",
							RedirectURI: "https://dashboard.url",
						},
						Requires: []brokerapi.RequiredPermission{"syslog_drain", "route_forwarding"},
						Tags:     serviceTags,
						Plans: []brokerapi.ServicePlan{
							{
								ID:          dedicatedPlanID,
								Name:        dedicatedPlanName,
								Description: dedicatedPlanDescription,
								Free:        &trueVar,
								Bindable:    &trueVar,
								Metadata: &brokerapi.ServicePlanMetadata{
									Bullets:     dedicatedPlanBullets,
									DisplayName: dedicatedPlanDisplayName,
									Costs: []brokerapi.ServicePlanCost{
										{
											Unit:   dedicatedPlanCostUnit,
											Amount: dedicatedPlanCostAmount,
										},
									},
								},
							},
							{
								ID:          highMemoryPlanID,
								Name:        highMemoryPlanName,
								Description: highMemoryPlanDescription,
								Metadata: &brokerapi.ServicePlanMetadata{
									Bullets:     highMemoryPlanBullets,
									DisplayName: highMemoryPlanDisplayName,
								},
							},
						},
					},
				},
			}))
		})
	})
})

func doCatalogRequest() *http.Response {

	catalogReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/v2/catalog", serverURL), nil)
	Expect(err).ToNot(HaveOccurred())
	catalogReq.SetBasicAuth(brokerUsername, brokerPassword)

	catalogResponse, err := http.DefaultClient.Do(catalogReq)
	Expect(err).ToNot(HaveOccurred())

	return catalogResponse
}

func defaultServiceCatalogConfig() brokerConfig.ServiceOffering {
	return brokerConfig.ServiceOffering{
		ID:            serviceID,
		Name:          serviceName,
		Description:   serviceDescription,
		Bindable:      serviceBindable,
		PlanUpdatable: servicePlanUpdatable,
		Metadata: brokerConfig.ServiceMetadata{
			DisplayName:         serviceMetadataDisplayName,
			ImageURL:            serviceMetadataImageURL,
			LongDescription:     serviceMetaDataLongDescription,
			ProviderDisplayName: serviceMetaDataProviderDisplayName,
			DocumentationURL:    serviceMetaDataDocumentationURL,
			SupportURL:          serviceMetaDataSupportURL,
			Shareable:           serviceMetaDataShareable,
		},
		DashboardClient: &brokerConfig.DashboardClient{
			ID:          "client-id-1",
			Secret:      "secret-1",
			RedirectUri: "https://dashboard.url",
		},
		Tags:             serviceTags,
		GlobalProperties: sdk.Properties{"global_property": "global_value"},
		GlobalQuotas:     brokerConfig.Quotas{},
		Plans: []brokerConfig.Plan{
			{
				Name:        dedicatedPlanName,
				ID:          dedicatedPlanID,
				Description: dedicatedPlanDescription,
				Free:        &trueVar,
				Bindable:    &trueVar,
				Update:      dedicatedPlanUpdateBlock,
				Metadata: brokerConfig.PlanMetadata{
					DisplayName: dedicatedPlanDisplayName,
					Bullets:     dedicatedPlanBullets,
					Costs: []brokerConfig.PlanCost{
						{
							Amount: dedicatedPlanCostAmount,
							Unit:   dedicatedPlanCostUnit,
						},
					},
				},
				Quotas: brokerConfig.Quotas{
					ServiceInstanceLimit: &dedicatedPlanQuota,
				},
				Properties: sdk.Properties{
					"type": "dedicated-plan-property",
				},
				InstanceGroups: []sdk.InstanceGroup{
					{
						Name:               "instance-group-name",
						VMType:             dedicatedPlanVMType,
						VMExtensions:       dedicatedPlanVMExtensions,
						PersistentDiskType: dedicatedPlanDisk,
						Instances:          dedicatedPlanInstances,
						Networks:           dedicatedPlanNetworks,
						AZs:                dedicatedPlanAZs,
					},
					{
						Name:               "instance-group-errand",
						Lifecycle:          "errand",
						VMType:             dedicatedPlanVMType,
						PersistentDiskType: dedicatedPlanDisk,
						Instances:          dedicatedPlanInstances,
						Networks:           dedicatedPlanNetworks,
						AZs:                dedicatedPlanAZs,
					},
				},
			},
			{
				Name:        highMemoryPlanName,
				ID:          highMemoryPlanID,
				Description: highMemoryPlanDescription,
				Metadata: brokerConfig.PlanMetadata{
					DisplayName: highMemoryPlanDisplayName,
					Bullets:     highMemoryPlanBullets,
				},
				Properties: sdk.Properties{
					"type":            "high-memory-plan-property",
					"global_property": "overrides_global_value",
				},
				InstanceGroups: []sdk.InstanceGroup{
					{
						Name:         "instance-group-name",
						VMType:       highMemoryPlanVMType,
						VMExtensions: highMemoryPlanVMExtensions,
						Instances:    highMemoryPlanInstances,
						Networks:     highMemoryPlanNetworks,
						AZs:          highMemoryPlanAZs,
					},
				},
			},
		},
	}
}
