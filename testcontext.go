package treetest

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

type TestContext struct {
	Testing          *testing.T
	Require          *require.Assertions
	TestCase         *TestCase
	SubTestCase      *SubTestCase
}

func (tc *TestContext) ReadString(filePath string) string {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		tc.Require.Fail(fmt.Sprintf("failed reading File: %s", filePath))
	}
	return string(bytes)
}

func (tc *TestContext) MustReadData(withObject interface{}, filePath string) interface{} {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		tc.Require.Fail(fmt.Sprintf("failed reading File: %s", filePath))
	}

	err = json.Unmarshal(bytes, withObject)
	if err != nil {
		tc.Require.Fail(fmt.Sprintf("failed unmarshalling %T for '%s'. Error:\n%s", withObject, filePath, err))
	}
	return withObject
}

func (tc *TestContext) MustGetItem(name string) interface{} {
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
