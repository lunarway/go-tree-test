package treetest

type multiFileTest struct {
	name          string
	requiredFiles []string
	testFunc      func(files Files, ctx *TestContext)
}

func (st *multiFileTest) GetName() string {
	return st.name
}

func (st *multiFileTest) Run(files Files, tc *TestContext) {
	st.testFunc(files, tc)
}

func (st *multiFileTest) GetRequiredFiles() []string {
	return st.requiredFiles
}

func DefineMultiFileTest(name string, requiredFiles []string, testFunc func(files Files, ctx *TestContext)) Test {
	return &multiFileTest{
		name:          name,
		requiredFiles: requiredFiles,
		testFunc:      testFunc,
	}
}

var _ Test = &multiFileTest{}
