package treetest

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

type TestRunnerContext struct {
	t    *testing.T
	wg   *sync.WaitGroup
	spec *TestCase
}

func (ctx *TestRunnerContext) GetT() *testing.T {
	return ctx.t
}

func (ctx *TestRunnerContext) GetSpec() SpecTest {
	return ctx.spec
}

func (ctx *TestRunnerContext) GetWaitGroup() *sync.WaitGroup {
	return ctx.wg
}

func (ctx *TestRunnerContext) RegisterItem(name string, item interface{}) *TestRunnerContext {
	ctx.spec.registeredItems[name] = item
	return ctx
}

func (tc *TestRunnerContext) MustGetItem(name string) interface{} {
	item, ok := tc.spec.registeredItems[name]
	if !ok {
		var items []string
		for name := range tc.spec.registeredItems {
			items = append(items, name)
		}
		panic(fmt.Sprintf("trying to get '%s' but it isn't registered. Available items: %s", name, strings.Join(items, ", ")))
	}
	return item
}