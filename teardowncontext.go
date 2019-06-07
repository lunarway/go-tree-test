package treetest

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type TeardownContext struct {
	T                *testing.T
	Require          *require.Assertions
	TestCase         *TestCase
	Directory        directory
}

func (tc *TeardownContext) MustGetItem(name string) interface{} {
	item, ok := tc.TestCase.registeredItems[name]
	if !ok {
		var items []string
		for name := range tc.TestCase.registeredItems {
			items = append(items, name)
		}
		tc.Require.Fail(fmt.Sprintf("trying to get '%s' but it isn't registered. Available items: %s", name, strings.Join(items, ", ")))
	}
	return item
}
