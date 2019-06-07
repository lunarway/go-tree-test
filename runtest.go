package treetest

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

type internalTest struct {
	testCase *TestCase
	subTests []*internalSubTest
}

type internalSubTest struct {
	subTestCase SubTestCase
	t           *testing.T
	funcChan    chan func()
	doneWG      *sync.WaitGroup
}

func runTest(t *testing.T, wg *sync.WaitGroup, tc *TestCase) {
	var startedTest internalTest

	var subTests []*internalSubTest
	for j := range tc.tests {
		subTest := tc.tests[j]
		addSubTest := make(chan *internalSubTest)
		go func() {
			t.Run(tc.name+"/"+subTest.name, func(t *testing.T) {
				if tc.skip {
					addSubTest <- nil
					t.Skip("skipped for now")
					return
				}
				doneWG := &sync.WaitGroup{}
				doneWG.Add(1)
				defer doneWG.Done()
				ist := &internalSubTest{
					subTestCase: subTest,
					t:           t,
					funcChan:    make(chan func()),
					doneWG:      doneWG,
				}

				var prevSubTest *internalSubTest
				if j != 0 {
					prevSubTest = subTests[j-1]
				}
				addSubTest <- ist
				testFunc := <-ist.funcChan
				if prevSubTest != nil {
					prevSubTest.doneWG.Wait()
				}
				testFunc()
			})
		}()
		newSubTest := <-addSubTest
		subTests = append(subTests, newSubTest)
	}
	startedTest = internalTest{
		testCase: tc,
		subTests: subTests,
	}


	wg.Add(1)
	go func() {
		defer wg.Done()
		shutdownWG := &sync.WaitGroup{}
		shutdownWG.Add(1)
		defer func() {
			for _, teardown := range startedTest.testCase.teardowns {
				teardown(TeardownContext{
					T:                t,
					Directory:        startedTest.testCase.directory,
					TestCase:         startedTest.testCase,
				})
			}
		}()
		go func() {
			defer shutdownWG.Done()
			for i := range startedTest.testCase.setups {
				setup := startedTest.testCase.setups[i]
				setup(SetupContext{
					T:                t,
					Directory:        startedTest.testCase.directory,
					TestCase:         startedTest.testCase,
				})
			}

			tc := startedTest.testCase
			for index := range startedTest.subTests {
				shutdownWG.Add(1)
				subTest := startedTest.subTests[index]
				subTest.funcChan <- func() {
					defer shutdownWG.Done()
					subTest.subTestCase.test(&TestContext{
						TestCase:         tc,
						SubTestCase:      &subTest.subTestCase,
						Testing:          subTest.t,
						Require:          require.New(subTest.t),
					})
				}
			}
		}()
		shutdownWG.Wait()
	}()
}
