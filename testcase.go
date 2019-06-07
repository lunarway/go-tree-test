package treetest

import (
	"sync"
	"testing"
	"time"
)

type TestCase struct {
	name            string
	skip            bool
	tests           []SubTestCase
	startAt         *time.Time
	setups          []func(ctx SetupContext)
	teardowns       []func(ctx TeardownContext)
	registeredItems map[string]interface{}
	directory       directory

}

func (tc *TestCase) Setup(f func(ctx SetupContext)) SpecTest {
	tc.setups = append(tc.setups, f)
	return tc
}

func (tc *TestCase) Test(name string, f func(*TestContext)) SpecTest {
	tc.tests = append(tc.tests, SubTestCase{
		name: name,
		test: f,
	})
	return tc
}

func (tc *TestCase) GetName() string {
	return tc.name
}

func (tc *TestCase) Run(t *testing.T, wg *sync.WaitGroup) {
	runTest(t, wg, tc)
}

func (tc *TestCase) Skip() SpecTest {
	tc.skip = true
	return tc
}

func DefineTest(name string, directory directory) SpecTest {
	return &TestCase{
		name:            name,
		directory:       directory,
		registeredItems: make(map[string]interface{}),
	}
}
