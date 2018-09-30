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
	return &walkerImpl{
		Logger: slf4go.Get("package-walker"),
	}
}

func (impl *walkerImpl) imported(abspath string) (*build.Package, bool) {

	for _, pkg := range impl.packages {
		if pkg.ImportPath == abspath {
			return pkg, true
		}
	}

	return nil, false
}

func (impl *walkerImpl) importedDir(dir string) (*build.Package, bool) {

	for _, pkg := range impl.packages {
		if pkg.Dir == dir {
			return pkg, true
		}
	}

	return nil, false
}

func (impl *walkerImpl) Import(dir string, recursion bool) error {

	path, err := filepath.Abs(dir)

	if err != nil {
		return xerrors.Wrapf(err, "get fullpath of %s error", dir)
	}

	pkg, ok := impl.importedDir(path)

	if ok {
		impl.DebugF("get imported package: %s", pkg.ImportPath)
		return nil
	}

	pkg, err = impl.importDir(dir)

	if err != nil {
		return xerrors.Wrapf(err, "get package from dir %s error", dir)
	}

	impl.packages = append(impl.packages, pkg)

	impl.DebugF("find package: %s", pkg.ImportPath)

	if !recursion {
		return nil
	}

	var fifo []string

	fifo = append(fifo, pkg.Imports...)

	for len(fifo) > 0 {

		top := fifo[0]

		fifo = fifo[1:]

		_, ok := impl.imported(top)

		if ok {
			continue
		}

		pkg, err := impl.importPackage(top)

		if err != nil {
			impl.DebugF("get imported package: %s", pkg.ImportPath)
			return err
		}

		impl.packages = append(impl.packages, pkg)

		impl.DebugF("find package: %s", pkg.ImportPath)
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
