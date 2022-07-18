// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/controllers/controllers/workloads"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CFClient struct {
	CreateStub        func(context.Context, client.Object, ...client.CreateOption) error
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		arg1 context.Context
		arg2 client.Object
		arg3 []client.CreateOption
	}
	createReturns struct {
		result1 error
	}
	createReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStub        func(context.Context, client.Object, ...client.DeleteOption) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 context.Context
		arg2 client.Object
		arg3 []client.DeleteOption
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	GetStub        func(context.Context, types.NamespacedName, client.Object) error
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 context.Context
		arg2 types.NamespacedName
		arg3 client.Object
	}
	getReturns struct {
		result1 error
	}
	getReturnsOnCall map[int]struct {
		result1 error
	}
	ListStub        func(context.Context, client.ObjectList, ...client.ListOption) error
	listMutex       sync.RWMutex
	listArgsForCall []struct {
		arg1 context.Context
		arg2 client.ObjectList
		arg3 []client.ListOption
	}
	listReturns struct {
		result1 error
	}
	listReturnsOnCall map[int]struct {
		result1 error
	}
	PatchStub        func(context.Context, client.Object, client.Patch, ...client.PatchOption) error
	patchMutex       sync.RWMutex
	patchArgsForCall []struct {
		arg1 context.Context
		arg2 client.Object
		arg3 client.Patch
		arg4 []client.PatchOption
	}
	patchReturns struct {
		result1 error
	}
	patchReturnsOnCall map[int]struct {
		result1 error
	}
	StatusStub        func() client.StatusWriter
	statusMutex       sync.RWMutex
	statusArgsForCall []struct {
	}
	statusReturns struct {
		result1 client.StatusWriter
	}
	statusReturnsOnCall map[int]struct {
		result1 client.StatusWriter
	}
	UpdateStub        func(context.Context, client.Object, ...client.UpdateOption) error
	updateMutex       sync.RWMutex
	updateArgsForCall []struct {
		arg1 context.Context
		arg2 client.Object
		arg3 []client.UpdateOption
	}
	updateReturns struct {
		result1 error
	}
	updateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CFClient) Create(arg1 context.Context, arg2 client.Object, arg3 ...client.CreateOption) error {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		arg1 context.Context
		arg2 client.Object
		arg3 []client.CreateOption
	}{arg1, arg2, arg3})
	stub := fake.CreateStub
	fakeReturns := fake.createReturns
	fake.recordInvocation("Create", []interface{}{arg1, arg2, arg3})
	fake.createMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CFClient) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *CFClient) CreateCalls(stub func(context.Context, client.Object, ...client.CreateOption) error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = stub
}

func (fake *CFClient) CreateArgsForCall(i int) (context.Context, client.Object, []client.CreateOption) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	argsForCall := fake.createArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFClient) CreateReturns(result1 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) CreateReturnsOnCall(i int, result1 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) Delete(arg1 context.Context, arg2 client.Object, arg3 ...client.DeleteOption) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 context.Context
		arg2 client.Object
		arg3 []client.DeleteOption
	}{arg1, arg2, arg3})
	stub := fake.DeleteStub
	fakeReturns := fake.deleteReturns
	fake.recordInvocation("Delete", []interface{}{arg1, arg2, arg3})
	fake.deleteMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CFClient) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *CFClient) DeleteCalls(stub func(context.Context, client.Object, ...client.DeleteOption) error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *CFClient) DeleteArgsForCall(i int) (context.Context, client.Object, []client.DeleteOption) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFClient) DeleteReturns(result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) DeleteReturnsOnCall(i int, result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) Get(arg1 context.Context, arg2 types.NamespacedName, arg3 client.Object) error {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 context.Context
		arg2 types.NamespacedName
		arg3 client.Object
	}{arg1, arg2, arg3})
	stub := fake.GetStub
	fakeReturns := fake.getReturns
	fake.recordInvocation("Get", []interface{}{arg1, arg2, arg3})
	fake.getMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CFClient) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *CFClient) GetCalls(stub func(context.Context, types.NamespacedName, client.Object) error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *CFClient) GetArgsForCall(i int) (context.Context, types.NamespacedName, client.Object) {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFClient) GetReturns(result1 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) GetReturnsOnCall(i int, result1 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) List(arg1 context.Context, arg2 client.ObjectList, arg3 ...client.ListOption) error {
	fake.listMutex.Lock()
	ret, specificReturn := fake.listReturnsOnCall[len(fake.listArgsForCall)]
	fake.listArgsForCall = append(fake.listArgsForCall, struct {
		arg1 context.Context
		arg2 client.ObjectList
		arg3 []client.ListOption
	}{arg1, arg2, arg3})
	stub := fake.ListStub
	fakeReturns := fake.listReturns
	fake.recordInvocation("List", []interface{}{arg1, arg2, arg3})
	fake.listMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CFClient) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *CFClient) ListCalls(stub func(context.Context, client.ObjectList, ...client.ListOption) error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = stub
}

func (fake *CFClient) ListArgsForCall(i int) (context.Context, client.ObjectList, []client.ListOption) {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	argsForCall := fake.listArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFClient) ListReturns(result1 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) ListReturnsOnCall(i int, result1 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	if fake.listReturnsOnCall == nil {
		fake.listReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.listReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) Patch(arg1 context.Context, arg2 client.Object, arg3 client.Patch, arg4 ...client.PatchOption) error {
	fake.patchMutex.Lock()
	ret, specificReturn := fake.patchReturnsOnCall[len(fake.patchArgsForCall)]
	fake.patchArgsForCall = append(fake.patchArgsForCall, struct {
		arg1 context.Context
		arg2 client.Object
		arg3 client.Patch
		arg4 []client.PatchOption
	}{arg1, arg2, arg3, arg4})
	stub := fake.PatchStub
	fakeReturns := fake.patchReturns
	fake.recordInvocation("Patch", []interface{}{arg1, arg2, arg3, arg4})
	fake.patchMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CFClient) PatchCallCount() int {
	fake.patchMutex.RLock()
	defer fake.patchMutex.RUnlock()
	return len(fake.patchArgsForCall)
}

func (fake *CFClient) PatchCalls(stub func(context.Context, client.Object, client.Patch, ...client.PatchOption) error) {
	fake.patchMutex.Lock()
	defer fake.patchMutex.Unlock()
	fake.PatchStub = stub
}

func (fake *CFClient) PatchArgsForCall(i int) (context.Context, client.Object, client.Patch, []client.PatchOption) {
	fake.patchMutex.RLock()
	defer fake.patchMutex.RUnlock()
	argsForCall := fake.patchArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *CFClient) PatchReturns(result1 error) {
	fake.patchMutex.Lock()
	defer fake.patchMutex.Unlock()
	fake.PatchStub = nil
	fake.patchReturns = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) PatchReturnsOnCall(i int, result1 error) {
	fake.patchMutex.Lock()
	defer fake.patchMutex.Unlock()
	fake.PatchStub = nil
	if fake.patchReturnsOnCall == nil {
		fake.patchReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.patchReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) Status() client.StatusWriter {
	fake.statusMutex.Lock()
	ret, specificReturn := fake.statusReturnsOnCall[len(fake.statusArgsForCall)]
	fake.statusArgsForCall = append(fake.statusArgsForCall, struct {
	}{})
	stub := fake.StatusStub
	fakeReturns := fake.statusReturns
	fake.recordInvocation("Status", []interface{}{})
	fake.statusMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CFClient) StatusCallCount() int {
	fake.statusMutex.RLock()
	defer fake.statusMutex.RUnlock()
	return len(fake.statusArgsForCall)
}

func (fake *CFClient) StatusCalls(stub func() client.StatusWriter) {
	fake.statusMutex.Lock()
	defer fake.statusMutex.Unlock()
	fake.StatusStub = stub
}

func (fake *CFClient) StatusReturns(result1 client.StatusWriter) {
	fake.statusMutex.Lock()
	defer fake.statusMutex.Unlock()
	fake.StatusStub = nil
	fake.statusReturns = struct {
		result1 client.StatusWriter
	}{result1}
}

func (fake *CFClient) StatusReturnsOnCall(i int, result1 client.StatusWriter) {
	fake.statusMutex.Lock()
	defer fake.statusMutex.Unlock()
	fake.StatusStub = nil
	if fake.statusReturnsOnCall == nil {
		fake.statusReturnsOnCall = make(map[int]struct {
			result1 client.StatusWriter
		})
	}
	fake.statusReturnsOnCall[i] = struct {
		result1 client.StatusWriter
	}{result1}
}

func (fake *CFClient) Update(arg1 context.Context, arg2 client.Object, arg3 ...client.UpdateOption) error {
	fake.updateMutex.Lock()
	ret, specificReturn := fake.updateReturnsOnCall[len(fake.updateArgsForCall)]
	fake.updateArgsForCall = append(fake.updateArgsForCall, struct {
		arg1 context.Context
		arg2 client.Object
		arg3 []client.UpdateOption
	}{arg1, arg2, arg3})
	stub := fake.UpdateStub
	fakeReturns := fake.updateReturns
	fake.recordInvocation("Update", []interface{}{arg1, arg2, arg3})
	fake.updateMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CFClient) UpdateCallCount() int {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return len(fake.updateArgsForCall)
}

func (fake *CFClient) UpdateCalls(stub func(context.Context, client.Object, ...client.UpdateOption) error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = stub
}

func (fake *CFClient) UpdateArgsForCall(i int) (context.Context, client.Object, []client.UpdateOption) {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	argsForCall := fake.updateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFClient) UpdateReturns(result1 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	fake.updateReturns = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) UpdateReturnsOnCall(i int, result1 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	if fake.updateReturnsOnCall == nil {
		fake.updateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CFClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	fake.patchMutex.RLock()
	defer fake.patchMutex.RUnlock()
	fake.statusMutex.RLock()
	defer fake.statusMutex.RUnlock()
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CFClient) recordInvocation(key string, args []interface{}) {
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

var _ workloads.CFClient = new(CFClient)
