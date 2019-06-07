package simple

import (
	"github.com/lunarway/go-tree-test"
	"testing"
)

func TestSetup(t *testing.T) {
	treetest.DefineTesting().
		Setup(func(context treetest.SetupContext) {

		}).
		Test(treetest.DefineFileTest("somefile", "somefile.json", func(file treetest.File, ctx *treetest.TestContext) {
			something := &Something{}
			file.MustReadData(something)
			ctx.Require.Equal(3, something.Value)
		})).
		Test(treetest.DefineMultiFileTest("mytest", []string{"req.json", "res.json"}, func(files treetest.Files, ctx *treetest.TestContext) {
			something1 := &Something{}
			something2 := &Something{}
			files.MustReadData("req.json", something1)
			files.MustReadData("res.json", something2)
			ctx.Require.Equal(something1.Value, something2.Value)
		})).
		RunTestsIn(t, ".")
}

type Something struct {
	Value int
}
