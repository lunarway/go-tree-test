package treetest

type SubTestCase struct {
	name string
	test func(*TestContext)
}
