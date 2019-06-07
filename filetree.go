package treetest

import (
	"os"
	"strings"
	"sync"
	"testing"
)

type TestRunner func(context TestRunnerContext, run func(TestRunnerContext))
type TestSetup func(SetupContext)
type Test interface {
	Run(Files, *TestContext)
	GetRequiredFiles() []string
	GetName() string
}

type TreeTestSpec interface {
	Setup(TestSetup) TreeTestSpec
	Test(Test) TreeTestSpec
	RunTestsIn(t *testing.T, path string)
}

type TestingDefinitionBuilder interface {
	InSync() TestingDefinitionBuilder
	UsingRunner(TestRunner) TestingDefinitionBuilder
	Setup(TestSetup) TestingDefinitionBuilder
	Test(Test) TestingDefinitionBuilder
	RunTestsIn(t *testing.T, path string)
}

type testingDefinitionBuilder struct {
	setups     []TestSetup
	tests      []Test
	inSync     bool
	testRunner TestRunner
}

func (b *testingDefinitionBuilder) Setup(setup TestSetup) TestingDefinitionBuilder {
	b.setups = append(b.setups, setup)
	return b
}

func (b *testingDefinitionBuilder) UsingRunner(testRunner TestRunner) TestingDefinitionBuilder {
	b.testRunner = testRunner
	return b
}

func (b *testingDefinitionBuilder) InSync() TestingDefinitionBuilder {
	b.inSync = true
	return b
}

func (b *testingDefinitionBuilder) Test(test Test) TestingDefinitionBuilder {
	b.tests = append(b.tests, test)
	return b
}

func (b *testingDefinitionBuilder) RunTestsIn(t *testing.T, path string) {
	onlyRun := os.Getenv("ONLY_RUN")

	testDirectory, err := getFiles(path)
	if err != nil {
		panic(err)
	}

	dirs := testDirectory.getAllDirectories()
	for i := range dirs {
		dir := dirs[i]
		var testSpec SpecTest
		for i := range b.tests {
			test := b.tests[i]
			foundAll, foundFiles := dir.findAllFilesUpwards(test.GetRequiredFiles()...)
			if !foundAll {
				continue // Test are missing files, so skip the test
			}
			if testSpec == nil {
				testSpec = dir.getTestSpec()
				for i := range b.setups {
					testSpec.Setup(func(setupContext SetupContext) {
						b.setups[i](setupContext)
					})
				}
			}
			testSpec.Test(test.GetName(), func(context *TestContext) {
				test.Run(foundFiles, context)
			})
		}
	}

	var runner = func(context TestRunnerContext, run func(TestRunnerContext)) {
		run(context)
	}
	if b.testRunner != nil {
		runner = b.testRunner
	}

	allSpecs := testDirectory.getAllDefinedTestSpecs()
	if b.inSync {
		for ts := range allSpecs {
			wg := &sync.WaitGroup{}
			spec := allSpecs[ts]
			if onlyRun == "" || strings.HasPrefix(spec.GetName(), onlyRun) {
				runner(TestRunnerContext{
					t:    t,
					wg:   wg,
					spec: spec.(*TestCase), // TODO fix this possibly bad type cast
				}, func(context TestRunnerContext) {
					context.GetSpec().Run(context.GetT(), context.GetWaitGroup())
				})
			}
			wg.Wait()
		}
	} else {
		wg := &sync.WaitGroup{}
		for ts := range allSpecs {
			spec := allSpecs[ts]

			if onlyRun == "" || strings.HasPrefix(spec.GetName(), onlyRun) {
				runner(TestRunnerContext{
					t:    t,
					wg:   wg,
					spec: spec.(*TestCase), // TODO fix this possibly bad type cast
				}, func(context TestRunnerContext) {
					context.GetSpec().Run(context.GetT(), context.GetWaitGroup())
				})
			}
		}
		wg.Wait()
	}
}

func DefineTesting() TestingDefinitionBuilder {
	return &testingDefinitionBuilder{}
}
