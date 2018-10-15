package gpm // import "github.com/dynamicgo/dy/gpm"

import (
	"encoding/json"
	"go/build"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoBuild(t *testing.T) {
	pkg, err := build.ImportDir("./", build.IgnoreVendor|build.ImportComment)
	require.NoError(t, err)

	println(printResult(pkg))
	println(printResult(build.Default.SrcDirs()))

	println(pkg.Dir)

}

func TestWalker(t *testing.T) {

	walker := NewWalker(WithSkip([]string{"C"}))

	err := walker.Import("/Users/yayanyang/Workspace/src/cjoy.tech/cjpool", true)

	require.NoError(t, err)

	packages := walker.Packages()

	for _, pkg := range packages {
		println("find package: ", pkg.Name, pkg.ImportPath)

		// println(printResult(pkg.Imports))
	}

}

func TestPath(t *testing.T) {
	path, _ := filepath.Abs("./")

	println(path)
}

func printResult(val interface{}) string {
	buff, _ := json.MarshalIndent(val, "", "\t")

	return string(buff)
}
