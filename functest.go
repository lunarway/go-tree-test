package treetest

type funcTest struct {
	name     string
	testFunc func(ctx *TestContext)
}

func (st *funcTest) GetName() string {
	return st.name
}

func (st *funcTest) Run(files Files, tc *TestContext) {
	st.testFunc(tc)
}

func (st *funcTest) GetRequiredFiles() []string {
	return []string{}
}

func DefineFuncTest(name string, testFunc func(ctx *TestContext)) Test {
	return &funcTest{
		name:     name,
		testFunc: testFunc,
	}
}

var _ Test = &funcTest{}
