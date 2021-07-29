# go.binfmt

An extension library for [`os/exec`](https://golang.org/pkg/os/exec/).

[![pkg.go.dev](https://pkg.go.dev/badge/github.com/hattya/go.binfmt)](https://pkg.go.dev/github.com/hattya/go.binfmt)
[![GitHub Actions](https://github.com/hattya/go.binfmt/actions/workflows/ci.yml/badge.svg)](https://github.com/hattya/go.binfmt/actions/workflows/ci.yml)
[![Semaphore](https://semaphoreci.com/api/v1/hattya/go-binfmt/branches/master/badge.svg)](https://semaphoreci.com/hattya/go-binfmt)
[![Appveyor](https://ci.appveyor.com/api/projects/status/uhkkibn9gen71du9/branch/master?svg=true)](https://ci.appveyor.com/project/hattya/go-binfmt)
[![Codecov](https://codecov.io/gh/hattya/go.binfmt/branch/master/graph/badge.svg)](https://codecov.io/gh/hattya/go.binfmt)


## Installation

```console
$ go get -u github.com/hattya/go.binfmt
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
