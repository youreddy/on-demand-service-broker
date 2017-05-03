// This file was generated by counterfeiter
package fakes

import (
	"log"
	"sync"

	"github.com/pivotal-cf/on-demand-service-broker/broker"
	"github.com/pivotal-cf/on-demand-service-broker/cloud_foundry_client"
)

type FakeCloudFoundryClient struct {
	GetAPIVersionStub        func(logger *log.Logger) (string, error)
	getAPIVersionMutex       sync.RWMutex
	getAPIVersionArgsForCall []struct {
		logger *log.Logger
	}
	getAPIVersionReturns struct {
		result1 string
		result2 error
	}
	getAPIVersionReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	CountInstancesOfPlanStub        func(serviceOfferingID, planID string, logger *log.Logger) (int, error)
	countInstancesOfPlanMutex       sync.RWMutex
	countInstancesOfPlanArgsForCall []struct {
		serviceOfferingID string
		planID            string
		logger            *log.Logger
	}
	countInstancesOfPlanReturns struct {
		result1 int
		result2 error
	}
	countInstancesOfPlanReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	CountInstancesOfServiceOfferingStub        func(serviceOfferingID string, logger *log.Logger) (instanceCountByPlanID map[string]int, err error)
	countInstancesOfServiceOfferingMutex       sync.RWMutex
	countInstancesOfServiceOfferingArgsForCall []struct {
		serviceOfferingID string
		logger            *log.Logger
	}
	countInstancesOfServiceOfferingReturns struct {
		result1 map[string]int
		result2 error
	}
	countInstancesOfServiceOfferingReturnsOnCall map[int]struct {
		result1 map[string]int
		result2 error
	}
	GetInstanceStateStub        func(serviceInstanceGUID string, logger *log.Logger) (cloud_foundry_client.InstanceState, error)
	getInstanceStateMutex       sync.RWMutex
	getInstanceStateArgsForCall []struct {
		serviceInstanceGUID string
		logger              *log.Logger
	}
	getInstanceStateReturns struct {
		result1 cloud_foundry_client.InstanceState
		result2 error
	}
	getInstanceStateReturnsOnCall map[int]struct {
		result1 cloud_foundry_client.InstanceState
		result2 error
	}
	GetInstancesOfServiceOfferingStub        func(serviceOfferingID string, logger *log.Logger) ([]string, error)
	getInstancesOfServiceOfferingMutex       sync.RWMutex
	getInstancesOfServiceOfferingArgsForCall []struct {
		serviceOfferingID string
		logger            *log.Logger
	}
	getInstancesOfServiceOfferingReturns struct {
		result1 []string
		result2 error
	}
	getInstancesOfServiceOfferingReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCloudFoundryClient) GetAPIVersion(logger *log.Logger) (string, error) {
	fake.getAPIVersionMutex.Lock()
	ret, specificReturn := fake.getAPIVersionReturnsOnCall[len(fake.getAPIVersionArgsForCall)]
	fake.getAPIVersionArgsForCall = append(fake.getAPIVersionArgsForCall, struct {
		logger *log.Logger
	}{logger})
	fake.recordInvocation("GetAPIVersion", []interface{}{logger})
	fake.getAPIVersionMutex.Unlock()
	if fake.GetAPIVersionStub != nil {
		return fake.GetAPIVersionStub(logger)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getAPIVersionReturns.result1, fake.getAPIVersionReturns.result2
}

func (fake *FakeCloudFoundryClient) GetAPIVersionCallCount() int {
	fake.getAPIVersionMutex.RLock()
	defer fake.getAPIVersionMutex.RUnlock()
	return len(fake.getAPIVersionArgsForCall)
}

func (fake *FakeCloudFoundryClient) GetAPIVersionArgsForCall(i int) *log.Logger {
	fake.getAPIVersionMutex.RLock()
	defer fake.getAPIVersionMutex.RUnlock()
	return fake.getAPIVersionArgsForCall[i].logger
}

func (fake *FakeCloudFoundryClient) GetAPIVersionReturns(result1 string, result2 error) {
	fake.GetAPIVersionStub = nil
	fake.getAPIVersionReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeCloudFoundryClient) GetAPIVersionReturnsOnCall(i int, result1 string, result2 error) {
	fake.GetAPIVersionStub = nil
	if fake.getAPIVersionReturnsOnCall == nil {
		fake.getAPIVersionReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getAPIVersionReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeCloudFoundryClient) CountInstancesOfPlan(serviceOfferingID string, planID string, logger *log.Logger) (int, error) {
	fake.countInstancesOfPlanMutex.Lock()
	ret, specificReturn := fake.countInstancesOfPlanReturnsOnCall[len(fake.countInstancesOfPlanArgsForCall)]
	fake.countInstancesOfPlanArgsForCall = append(fake.countInstancesOfPlanArgsForCall, struct {
		serviceOfferingID string
		planID            string
		logger            *log.Logger
	}{serviceOfferingID, planID, logger})
	fake.recordInvocation("CountInstancesOfPlan", []interface{}{serviceOfferingID, planID, logger})
	fake.countInstancesOfPlanMutex.Unlock()
	if fake.CountInstancesOfPlanStub != nil {
		return fake.CountInstancesOfPlanStub(serviceOfferingID, planID, logger)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.countInstancesOfPlanReturns.result1, fake.countInstancesOfPlanReturns.result2
}

func (fake *FakeCloudFoundryClient) CountInstancesOfPlanCallCount() int {
	fake.countInstancesOfPlanMutex.RLock()
	defer fake.countInstancesOfPlanMutex.RUnlock()
	return len(fake.countInstancesOfPlanArgsForCall)
}

func (fake *FakeCloudFoundryClient) CountInstancesOfPlanArgsForCall(i int) (string, string, *log.Logger) {
	fake.countInstancesOfPlanMutex.RLock()
	defer fake.countInstancesOfPlanMutex.RUnlock()
	return fake.countInstancesOfPlanArgsForCall[i].serviceOfferingID, fake.countInstancesOfPlanArgsForCall[i].planID, fake.countInstancesOfPlanArgsForCall[i].logger
}

func (fake *FakeCloudFoundryClient) CountInstancesOfPlanReturns(result1 int, result2 error) {
	fake.CountInstancesOfPlanStub = nil
	fake.countInstancesOfPlanReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeCloudFoundryClient) CountInstancesOfPlanReturnsOnCall(i int, result1 int, result2 error) {
	fake.CountInstancesOfPlanStub = nil
	if fake.countInstancesOfPlanReturnsOnCall == nil {
		fake.countInstancesOfPlanReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.countInstancesOfPlanReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeCloudFoundryClient) CountInstancesOfServiceOffering(serviceOfferingID string, logger *log.Logger) (instanceCountByPlanID map[string]int, err error) {
	fake.countInstancesOfServiceOfferingMutex.Lock()
	ret, specificReturn := fake.countInstancesOfServiceOfferingReturnsOnCall[len(fake.countInstancesOfServiceOfferingArgsForCall)]
	fake.countInstancesOfServiceOfferingArgsForCall = append(fake.countInstancesOfServiceOfferingArgsForCall, struct {
		serviceOfferingID string
		logger            *log.Logger
	}{serviceOfferingID, logger})
	fake.recordInvocation("CountInstancesOfServiceOffering", []interface{}{serviceOfferingID, logger})
	fake.countInstancesOfServiceOfferingMutex.Unlock()
	if fake.CountInstancesOfServiceOfferingStub != nil {
		return fake.CountInstancesOfServiceOfferingStub(serviceOfferingID, logger)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.countInstancesOfServiceOfferingReturns.result1, fake.countInstancesOfServiceOfferingReturns.result2
}

func (fake *FakeCloudFoundryClient) CountInstancesOfServiceOfferingCallCount() int {
	fake.countInstancesOfServiceOfferingMutex.RLock()
	defer fake.countInstancesOfServiceOfferingMutex.RUnlock()
	return len(fake.countInstancesOfServiceOfferingArgsForCall)
}

func (fake *FakeCloudFoundryClient) CountInstancesOfServiceOfferingArgsForCall(i int) (string, *log.Logger) {
	fake.countInstancesOfServiceOfferingMutex.RLock()
	defer fake.countInstancesOfServiceOfferingMutex.RUnlock()
	return fake.countInstancesOfServiceOfferingArgsForCall[i].serviceOfferingID, fake.countInstancesOfServiceOfferingArgsForCall[i].logger
}

func (fake *FakeCloudFoundryClient) CountInstancesOfServiceOfferingReturns(result1 map[string]int, result2 error) {
	fake.CountInstancesOfServiceOfferingStub = nil
	fake.countInstancesOfServiceOfferingReturns = struct {
		result1 map[string]int
		result2 error
	}{result1, result2}
}

func (fake *FakeCloudFoundryClient) CountInstancesOfServiceOfferingReturnsOnCall(i int, result1 map[string]int, result2 error) {
	fake.CountInstancesOfServiceOfferingStub = nil
	if fake.countInstancesOfServiceOfferingReturnsOnCall == nil {
		fake.countInstancesOfServiceOfferingReturnsOnCall = make(map[int]struct {
			result1 map[string]int
			result2 error
		})
	}
	fake.countInstancesOfServiceOfferingReturnsOnCall[i] = struct {
		result1 map[string]int
		result2 error
	}{result1, result2}
}

func (fake *FakeCloudFoundryClient) GetInstanceState(serviceInstanceGUID string, logger *log.Logger) (cloud_foundry_client.InstanceState, error) {
	fake.getInstanceStateMutex.Lock()
	ret, specificReturn := fake.getInstanceStateReturnsOnCall[len(fake.getInstanceStateArgsForCall)]
	fake.getInstanceStateArgsForCall = append(fake.getInstanceStateArgsForCall, struct {
		serviceInstanceGUID string
		logger              *log.Logger
	}{serviceInstanceGUID, logger})
	fake.recordInvocation("GetInstanceState", []interface{}{serviceInstanceGUID, logger})
	fake.getInstanceStateMutex.Unlock()
	if fake.GetInstanceStateStub != nil {
		return fake.GetInstanceStateStub(serviceInstanceGUID, logger)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getInstanceStateReturns.result1, fake.getInstanceStateReturns.result2
}

func (fake *FakeCloudFoundryClient) GetInstanceStateCallCount() int {
	fake.getInstanceStateMutex.RLock()
	defer fake.getInstanceStateMutex.RUnlock()
	return len(fake.getInstanceStateArgsForCall)
}

func (fake *FakeCloudFoundryClient) GetInstanceStateArgsForCall(i int) (string, *log.Logger) {
	fake.getInstanceStateMutex.RLock()
	defer fake.getInstanceStateMutex.RUnlock()
	return fake.getInstanceStateArgsForCall[i].serviceInstanceGUID, fake.getInstanceStateArgsForCall[i].logger
}

func (fake *FakeCloudFoundryClient) GetInstanceStateReturns(result1 cloud_foundry_client.InstanceState, result2 error) {
	fake.GetInstanceStateStub = nil
	fake.getInstanceStateReturns = struct {
		result1 cloud_foundry_client.InstanceState
		result2 error
	}{result1, result2}
}

func (fake *FakeCloudFoundryClient) GetInstanceStateReturnsOnCall(i int, result1 cloud_foundry_client.InstanceState, result2 error) {
	fake.GetInstanceStateStub = nil
	if fake.getInstanceStateReturnsOnCall == nil {
		fake.getInstanceStateReturnsOnCall = make(map[int]struct {
			result1 cloud_foundry_client.InstanceState
			result2 error
		})
	}
	fake.getInstanceStateReturnsOnCall[i] = struct {
		result1 cloud_foundry_client.InstanceState
		result2 error
	}{result1, result2}
}

func (fake *FakeCloudFoundryClient) GetInstancesOfServiceOffering(serviceOfferingID string, logger *log.Logger) ([]string, error) {
	fake.getInstancesOfServiceOfferingMutex.Lock()
	ret, specificReturn := fake.getInstancesOfServiceOfferingReturnsOnCall[len(fake.getInstancesOfServiceOfferingArgsForCall)]
	fake.getInstancesOfServiceOfferingArgsForCall = append(fake.getInstancesOfServiceOfferingArgsForCall, struct {
		serviceOfferingID string
		logger            *log.Logger
	}{serviceOfferingID, logger})
	fake.recordInvocation("GetInstancesOfServiceOffering", []interface{}{serviceOfferingID, logger})
	fake.getInstancesOfServiceOfferingMutex.Unlock()
	if fake.GetInstancesOfServiceOfferingStub != nil {
		return fake.GetInstancesOfServiceOfferingStub(serviceOfferingID, logger)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getInstancesOfServiceOfferingReturns.result1, fake.getInstancesOfServiceOfferingReturns.result2
}

func (fake *FakeCloudFoundryClient) GetInstancesOfServiceOfferingCallCount() int {
	fake.getInstancesOfServiceOfferingMutex.RLock()
	defer fake.getInstancesOfServiceOfferingMutex.RUnlock()
	return len(fake.getInstancesOfServiceOfferingArgsForCall)
}

func (fake *FakeCloudFoundryClient) GetInstancesOfServiceOfferingArgsForCall(i int) (string, *log.Logger) {
	fake.getInstancesOfServiceOfferingMutex.RLock()
	defer fake.getInstancesOfServiceOfferingMutex.RUnlock()
	return fake.getInstancesOfServiceOfferingArgsForCall[i].serviceOfferingID, fake.getInstancesOfServiceOfferingArgsForCall[i].logger
}

func (fake *FakeCloudFoundryClient) GetInstancesOfServiceOfferingReturns(result1 []string, result2 error) {
	fake.GetInstancesOfServiceOfferingStub = nil
	fake.getInstancesOfServiceOfferingReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeCloudFoundryClient) GetInstancesOfServiceOfferingReturnsOnCall(i int, result1 []string, result2 error) {
	fake.GetInstancesOfServiceOfferingStub = nil
	if fake.getInstancesOfServiceOfferingReturnsOnCall == nil {
		fake.getInstancesOfServiceOfferingReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.getInstancesOfServiceOfferingReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeCloudFoundryClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getAPIVersionMutex.RLock()
	defer fake.getAPIVersionMutex.RUnlock()
	fake.countInstancesOfPlanMutex.RLock()
	defer fake.countInstancesOfPlanMutex.RUnlock()
	fake.countInstancesOfServiceOfferingMutex.RLock()
	defer fake.countInstancesOfServiceOfferingMutex.RUnlock()
	fake.getInstanceStateMutex.RLock()
	defer fake.getInstanceStateMutex.RUnlock()
	fake.getInstancesOfServiceOfferingMutex.RLock()
	defer fake.getInstancesOfServiceOfferingMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeCloudFoundryClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ broker.CloudFoundryClient = new(FakeCloudFoundryClient)
