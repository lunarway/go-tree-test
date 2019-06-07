package treetest

import (
	"sync"
	"testing"
)

type SpecTest interface {
	Test(name string, test func(*TestContext)) SpecTest
	GetName() string
	Run(t *testing.T, wg *sync.WaitGroup)
	Setup(f func(ctx SetupContext)) SpecTest
}
