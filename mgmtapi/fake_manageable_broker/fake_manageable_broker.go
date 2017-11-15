// Code generated by counterfeiter. DO NOT EDIT.
package fake_manageable_broker

import (
	"context"
	"log"
	"sync"

	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-cf/on-demand-service-broker/broker"
	"github.com/pivotal-cf/on-demand-service-broker/cf"
	"github.com/pivotal-cf/on-demand-service-broker/mgmtapi"
	"github.com/pivotal-cf/on-demand-service-broker/service"
)

type FakeManageableBroker struct {
	InstancesStub        func(logger *log.Logger) ([]service.Instance, error)
	instancesMutex       sync.RWMutex
	instancesArgsForCall []struct {
		logger *log.Logger
	}
	instancesReturns struct {
		result1 []service.Instance
		result2 error
	}
	instancesReturnsOnCall map[int]struct {
		result1 []service.Instance
		result2 error
	}
	OrphanDeploymentsStub        func(logger *log.Logger) ([]string, error)
	orphanDeploymentsMutex       sync.RWMutex
	orphanDeploymentsArgsForCall []struct {
		logger *log.Logger
	}
	orphanDeploymentsReturns struct {
		result1 []string
		result2 error
	}
	orphanDeploymentsReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	UpgradeStub        func(ctx context.Context, instanceID string, updateDetails brokerapi.UpdateDetails, logger *log.Logger) (broker.OperationData, error)
	upgradeMutex       sync.RWMutex
	upgradeArgsForCall []struct {
		ctx           context.Context
		instanceID    string
		updateDetails brokerapi.UpdateDetails
		logger        *log.Logger
	}
	upgradeReturns struct {
		result1 broker.OperationData
		result2 error
	}
	upgradeReturnsOnCall map[int]struct {
		result1 broker.OperationData
		result2 error
	}
	CountInstancesOfPlansStub        func(logger *log.Logger) (map[cf.ServicePlan]int, error)
	countInstancesOfPlansMutex       sync.RWMutex
	countInstancesOfPlansArgsForCall []struct {
		logger *log.Logger
	}
	countInstancesOfPlansReturns struct {
		result1 map[cf.ServicePlan]int
		result2 error
	}
	countInstancesOfPlansReturnsOnCall map[int]struct {
		result1 map[cf.ServicePlan]int
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeManageableBroker) Instances(logger *log.Logger) ([]service.Instance, error) {
	fake.instancesMutex.Lock()
	ret, specificReturn := fake.instancesReturnsOnCall[len(fake.instancesArgsForCall)]
	fake.instancesArgsForCall = append(fake.instancesArgsForCall, struct {
		logger *log.Logger
	}{logger})
	fake.recordInvocation("Instances", []interface{}{logger})
	fake.instancesMutex.Unlock()
	if fake.InstancesStub != nil {
		return fake.InstancesStub(logger)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.instancesReturns.result1, fake.instancesReturns.result2
}

func (fake *FakeManageableBroker) InstancesCallCount() int {
	fake.instancesMutex.RLock()
	defer fake.instancesMutex.RUnlock()
	return len(fake.instancesArgsForCall)
}

func (fake *FakeManageableBroker) InstancesArgsForCall(i int) *log.Logger {
	fake.instancesMutex.RLock()
	defer fake.instancesMutex.RUnlock()
	return fake.instancesArgsForCall[i].logger
}

func (fake *FakeManageableBroker) InstancesReturns(result1 []service.Instance, result2 error) {
	fake.InstancesStub = nil
	fake.instancesReturns = struct {
		result1 []service.Instance
		result2 error
	}{result1, result2}
}

func (fake *FakeManageableBroker) InstancesReturnsOnCall(i int, result1 []service.Instance, result2 error) {
	fake.InstancesStub = nil
	if fake.instancesReturnsOnCall == nil {
		fake.instancesReturnsOnCall = make(map[int]struct {
			result1 []service.Instance
			result2 error
		})
	}
	fake.instancesReturnsOnCall[i] = struct {
		result1 []service.Instance
		result2 error
	}{result1, result2}
}

func (fake *FakeManageableBroker) OrphanDeployments(logger *log.Logger) ([]string, error) {
	fake.orphanDeploymentsMutex.Lock()
	ret, specificReturn := fake.orphanDeploymentsReturnsOnCall[len(fake.orphanDeploymentsArgsForCall)]
	fake.orphanDeploymentsArgsForCall = append(fake.orphanDeploymentsArgsForCall, struct {
		logger *log.Logger
	}{logger})
	fake.recordInvocation("OrphanDeployments", []interface{}{logger})
	fake.orphanDeploymentsMutex.Unlock()
	if fake.OrphanDeploymentsStub != nil {
		return fake.OrphanDeploymentsStub(logger)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.orphanDeploymentsReturns.result1, fake.orphanDeploymentsReturns.result2
}

func (fake *FakeManageableBroker) OrphanDeploymentsCallCount() int {
	fake.orphanDeploymentsMutex.RLock()
	defer fake.orphanDeploymentsMutex.RUnlock()
	return len(fake.orphanDeploymentsArgsForCall)
}

func (fake *FakeManageableBroker) OrphanDeploymentsArgsForCall(i int) *log.Logger {
	fake.orphanDeploymentsMutex.RLock()
	defer fake.orphanDeploymentsMutex.RUnlock()
	return fake.orphanDeploymentsArgsForCall[i].logger
}

func (fake *FakeManageableBroker) OrphanDeploymentsReturns(result1 []string, result2 error) {
	fake.OrphanDeploymentsStub = nil
	fake.orphanDeploymentsReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeManageableBroker) OrphanDeploymentsReturnsOnCall(i int, result1 []string, result2 error) {
	fake.OrphanDeploymentsStub = nil
	if fake.orphanDeploymentsReturnsOnCall == nil {
		fake.orphanDeploymentsReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.orphanDeploymentsReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeManageableBroker) Upgrade(ctx context.Context, instanceID string, updateDetails brokerapi.UpdateDetails, logger *log.Logger) (broker.OperationData, error) {
	fake.upgradeMutex.Lock()
	ret, specificReturn := fake.upgradeReturnsOnCall[len(fake.upgradeArgsForCall)]
	fake.upgradeArgsForCall = append(fake.upgradeArgsForCall, struct {
		ctx           context.Context
		instanceID    string
		updateDetails brokerapi.UpdateDetails
		logger        *log.Logger
	}{ctx, instanceID, updateDetails, logger})
	fake.recordInvocation("Upgrade", []interface{}{ctx, instanceID, updateDetails, logger})
	fake.upgradeMutex.Unlock()
	if fake.UpgradeStub != nil {
		return fake.UpgradeStub(ctx, instanceID, updateDetails, logger)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.upgradeReturns.result1, fake.upgradeReturns.result2
}

func (fake *FakeManageableBroker) UpgradeCallCount() int {
	fake.upgradeMutex.RLock()
	defer fake.upgradeMutex.RUnlock()
	return len(fake.upgradeArgsForCall)
}

func (fake *FakeManageableBroker) UpgradeArgsForCall(i int) (context.Context, string, brokerapi.UpdateDetails, *log.Logger) {
	fake.upgradeMutex.RLock()
	defer fake.upgradeMutex.RUnlock()
	return fake.upgradeArgsForCall[i].ctx, fake.upgradeArgsForCall[i].instanceID, fake.upgradeArgsForCall[i].updateDetails, fake.upgradeArgsForCall[i].logger
}

func (fake *FakeManageableBroker) UpgradeReturns(result1 broker.OperationData, result2 error) {
	fake.UpgradeStub = nil
	fake.upgradeReturns = struct {
		result1 broker.OperationData
		result2 error
	}{result1, result2}
}

func (fake *FakeManageableBroker) UpgradeReturnsOnCall(i int, result1 broker.OperationData, result2 error) {
	fake.UpgradeStub = nil
	if fake.upgradeReturnsOnCall == nil {
		fake.upgradeReturnsOnCall = make(map[int]struct {
			result1 broker.OperationData
			result2 error
		})
	}
	fake.upgradeReturnsOnCall[i] = struct {
		result1 broker.OperationData
		result2 error
	}{result1, result2}
}

func (fake *FakeManageableBroker) CountInstancesOfPlans(logger *log.Logger) (map[cf.ServicePlan]int, error) {
	fake.countInstancesOfPlansMutex.Lock()
	ret, specificReturn := fake.countInstancesOfPlansReturnsOnCall[len(fake.countInstancesOfPlansArgsForCall)]
	fake.countInstancesOfPlansArgsForCall = append(fake.countInstancesOfPlansArgsForCall, struct {
		logger *log.Logger
	}{logger})
	fake.recordInvocation("CountInstancesOfPlans", []interface{}{logger})
	fake.countInstancesOfPlansMutex.Unlock()
	if fake.CountInstancesOfPlansStub != nil {
		return fake.CountInstancesOfPlansStub(logger)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.countInstancesOfPlansReturns.result1, fake.countInstancesOfPlansReturns.result2
}

func (fake *FakeManageableBroker) CountInstancesOfPlansCallCount() int {
	fake.countInstancesOfPlansMutex.RLock()
	defer fake.countInstancesOfPlansMutex.RUnlock()
	return len(fake.countInstancesOfPlansArgsForCall)
}

func (fake *FakeManageableBroker) CountInstancesOfPlansArgsForCall(i int) *log.Logger {
	fake.countInstancesOfPlansMutex.RLock()
	defer fake.countInstancesOfPlansMutex.RUnlock()
	return fake.countInstancesOfPlansArgsForCall[i].logger
}

func (fake *FakeManageableBroker) CountInstancesOfPlansReturns(result1 map[cf.ServicePlan]int, result2 error) {
	fake.CountInstancesOfPlansStub = nil
	fake.countInstancesOfPlansReturns = struct {
		result1 map[cf.ServicePlan]int
		result2 error
	}{result1, result2}
}

func (fake *FakeManageableBroker) CountInstancesOfPlansReturnsOnCall(i int, result1 map[cf.ServicePlan]int, result2 error) {
	fake.CountInstancesOfPlansStub = nil
	if fake.countInstancesOfPlansReturnsOnCall == nil {
		fake.countInstancesOfPlansReturnsOnCall = make(map[int]struct {
			result1 map[cf.ServicePlan]int
			result2 error
		})
	}
	fake.countInstancesOfPlansReturnsOnCall[i] = struct {
		result1 map[cf.ServicePlan]int
		result2 error
	}{result1, result2}
}

func (fake *FakeManageableBroker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.instancesMutex.RLock()
	defer fake.instancesMutex.RUnlock()
	fake.orphanDeploymentsMutex.RLock()
	defer fake.orphanDeploymentsMutex.RUnlock()
	fake.upgradeMutex.RLock()
	defer fake.upgradeMutex.RUnlock()
	fake.countInstancesOfPlansMutex.RLock()
	defer fake.countInstancesOfPlansMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeManageableBroker) recordInvocation(key string, args []interface{}) {
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

var _ mgmtapi.ManageableBroker = new(FakeManageableBroker)
