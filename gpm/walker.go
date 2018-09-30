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
}

// NewWalker create package walker object
func NewWalker() Walker {
	return &walkerImpl{}
}

func (impl *walkerImpl) imported(key string) (*build.Package, bool) {
	return nil, false
}

func (impl *walkerImpl) importedDir(key string) (*build.Package, bool) {
	return nil, false
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

	pkg, err := impl.importDir(dir)

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
