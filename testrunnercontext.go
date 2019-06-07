package treetest

import (
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
