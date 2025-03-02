# go.binfmt

An extension library for [`os/exec`](https://pkg.go.dev/os/exec).

[![pkg.go.dev](https://pkg.go.dev/badge/github.com/hattya/go.binfmt)](https://pkg.go.dev/github.com/hattya/go.binfmt)
[![GitHub Actions](https://github.com/hattya/go.binfmt/actions/workflows/ci.yml/badge.svg)](https://github.com/hattya/go.binfmt/actions/workflows/ci.yml)
[![Appveyor](https://ci.appveyor.com/api/projects/status/uhkkibn9gen71du9/branch/master?svg=true)](https://ci.appveyor.com/project/hattya/go-binfmt)
[![Codecov](https://codecov.io/gh/hattya/go.binfmt/branch/master/graph/badge.svg)](https://codecov.io/gh/hattya/go.binfmt)


## Installation

```console
$ go get github.com/hattya/go.binfmt
```


## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/hattya/go.binfmt"
)

func main() {
	cmd := binfmt.Command("rst2html5.py", "README.rst", "README.html")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```


## License

go.binfmt is distributed under the terms of the MIT License.
