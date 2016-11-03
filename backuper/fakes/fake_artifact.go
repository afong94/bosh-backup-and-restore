// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/pivotal-cf/pcf-backup-and-restore/backuper"
)

type FakeArtifact struct {
	CreateFileStub        func(string) error
	createFileMutex       sync.RWMutex
	createFileArgsForCall []struct {
		arg1 string
	}
	createFileReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeArtifact) CreateFile(arg1 string) error {
	fake.createFileMutex.Lock()
	fake.createFileArgsForCall = append(fake.createFileArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("CreateFile", []interface{}{arg1})
	fake.createFileMutex.Unlock()
	if fake.CreateFileStub != nil {
		return fake.CreateFileStub(arg1)
	} else {
		return fake.createFileReturns.result1
	}
}

func (fake *FakeArtifact) CreateFileCallCount() int {
	fake.createFileMutex.RLock()
	defer fake.createFileMutex.RUnlock()
	return len(fake.createFileArgsForCall)
}

func (fake *FakeArtifact) CreateFileArgsForCall(i int) string {
	fake.createFileMutex.RLock()
	defer fake.createFileMutex.RUnlock()
	return fake.createFileArgsForCall[i].arg1
}

func (fake *FakeArtifact) CreateFileReturns(result1 error) {
	fake.CreateFileStub = nil
	fake.createFileReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeArtifact) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createFileMutex.RLock()
	defer fake.createFileMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeArtifact) recordInvocation(key string, args []interface{}) {
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

var _ backuper.Artifact = new(FakeArtifact)