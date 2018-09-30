// Package gpm the unofficial golang package/project manager
package gpm

import "go/build"

// Walker The src packages walker
type Walker interface {
	Import(dir string, recursion bool) error
	Packages() []*build.Package
}
