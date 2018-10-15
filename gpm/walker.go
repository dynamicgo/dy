// Package gpm the unofficial golang package/project manager
package gpm

import (
	"go/build"
	"path/filepath"

	"github.com/dynamicgo/xerrors"

	"github.com/dynamicgo/slf4go"
)

type walkerImpl struct {
	slf4go.Logger
	packages []*build.Package
	skipped  []string
}

// Option .
type Option func(impl *walkerImpl)

// WithSkip .
func WithSkip(packages []string) Option {
	return func(impl *walkerImpl) {
		impl.skipped = packages
	}
}

// NewWalker create package walker object
func NewWalker(options ...Option) Walker {
	impl := &walkerImpl{}

	for _, opt := range options {
		opt(impl)
	}

	return impl
}

func (impl *walkerImpl) imported(key string) (*build.Package, bool) {
	for _, pkg := range impl.packages {
		if pkg.Name == key {
			return pkg, true
		}
	}

	return nil, false
}

func (impl *walkerImpl) importedDir(key string) (*build.Package, bool) {
	for _, pkg := range impl.packages {
		if pkg.ImportPath == key {
			return pkg, true
		}
	}

	return nil, false
}

func (impl *walkerImpl) skip(key string) bool {
	for _, pkg := range impl.skipped {
		if pkg == key {
			return true
		}
	}

	return false
}

func (impl *walkerImpl) Import(dir string, recursion bool) error {

	path, err := filepath.Abs(dir)

	if err != nil {
		return xerrors.Wrapf(err, "get fullpath of %s error", dir)
	}

	_, ok := impl.importedDir(path)

	if ok {
		return nil
	}

	var fifo []string

	pkg, err := impl.importDir(path)

	if err != nil {
		return xerrors.Wrapf(err, "get package from dir %s error", dir)
	}

	impl.packages = append(impl.packages, pkg)

	if !recursion {
		return nil
	}

	fifo = append(fifo, pkg.Imports...)

	for len(fifo) > 0 {
		top := fifo[0]

		fifo = fifo[1:]

		_, ok := impl.imported(top)

		if ok {
			continue
		}

		if impl.skip(top) {
			continue
		}

		pkg, err := impl.importPackage(top)

		if err != nil {
			return err
		}

		impl.packages = append(impl.packages, pkg)

		fifo = append(fifo, pkg.Imports...)

	}

	return nil
}

func (impl *walkerImpl) Packages() []*build.Package {
	return impl.packages
}

func (impl *walkerImpl) importPackage(path string) (*build.Package, error) {
	pkg, err := build.Import(path, "", build.IgnoreVendor)

	if err != nil {
		return nil, xerrors.Wrap(err, "call build.ImportDir -- failed")
	}

	return pkg, nil
}

func (impl *walkerImpl) importDir(dir string) (*build.Package, error) {
	pkg, err := build.ImportDir(dir, build.IgnoreVendor)

	if err != nil {
		return nil, xerrors.Wrap(err, "call build.ImportDir -- failed")
	}

	return pkg, nil
}
