package gpm

import (
	"encoding/json"
	"go/build"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoBuild(t *testing.T) {
	pkg, err := build.Import("github.com/dynamicgo/dy/gpm", "./test", build.IgnoreVendor)
	require.NoError(t, err)

	println(printResult(pkg))
	println(printResult(build.Default.SrcDirs()))

}

func printResult(val interface{}) string {
	buff, _ := json.MarshalIndent(val, "", "\t")

	return string(buff)
}
