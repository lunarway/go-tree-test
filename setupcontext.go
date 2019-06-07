package treetest

import (
	"fmt"
	"strings"
	"testing"
)

type SetupContext struct {
	T                *testing.T
	TestCase         *TestCase
	Directory        directory
}

func (tc *SetupContext) MustReadData(fileName string, withObject interface{}) interface{} {
	file := tc.MustFindFile(fileName)
	return file.MustReadData(withObject)
}

func (tc *SetupContext) MustFindFile(fileName string) File {
	file := tc.Directory.findUpwards(fileName)
	if file == nil {
		panic(fmt.Sprintf("failed finding file '%s' for test dir '%s'", fileName, tc.Directory.path))
	}
	return *file
}

func (tc *SetupContext) RegisterItem(name string, item interface{}) *SetupContext {
	tc.TestCase.registeredItems[name] = item
	return tc
}

func (tc *SetupContext) MustGetItem(name string) interface{} {
	item, ok := tc.TestCase.registeredItems[name]
	if !ok {
		var items []string
		for name := range tc.TestCase.registeredItems {
			items = append(items, name)
		}
		panic(fmt.Sprintf("trying to get '%s' but it isn't registered. Available items: %s", name, strings.Join(items, ", ")))
	}
	return item
}

func (tc *SetupContext) Teardown(teardown func(ctx TeardownContext)) *SetupContext {
	tc.TestCase.teardowns = append(tc.TestCase.teardowns, teardown)
	return tc
}
