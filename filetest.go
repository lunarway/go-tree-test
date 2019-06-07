package treetest

type fileTest struct {
	testName     string
	requiredFile string
	testFunc     func(file File, ctx *TestContext)
}

func (st *fileTest) GetName() string {
	return st.testName
}

func (st *fileTest) Run(files Files, tc *TestContext) {
	st.testFunc(files[st.requiredFile], tc)
}

func (st *fileTest) GetRequiredFiles() []string {
	return []string{st.requiredFile}
}

func DefineFileTest(testName string, requiredFile string, testFunc func(file File, ctx *TestContext)) Test {
	return &fileTest{
		testName:     testName,
		requiredFile: requiredFile,
		testFunc:     testFunc,
	}
}

var _ Test = &fileTest{}
