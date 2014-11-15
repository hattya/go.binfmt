//
// go.binfmt :: binfmt_test.go
//
//   Copyright (c) 2014 Akinori Hattori <hattya@gmail.com>
//
//   Permission is hereby granted, free of charge, to any person
//   obtaining a copy of this software and associated documentation files
//   (the "Software"), to deal in the Software without restriction,
//   including without limitation the rights to use, copy, modify, merge,
//   publish, distribute, sublicense, and/or sell copies of the Software,
//   and to permit persons to whom the Software is furnished to do so,
//   subject to the following conditions:
//
//   The above copyright notice and this permission notice shall be
//   included in all copies or substantial portions of the Software.
//
//   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//   EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
//   MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//   NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
//   BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
//   ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
//   CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//   SOFTWARE.
//

package binfmt_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hattya/go.binfmt"
)

func TestCommand(t *testing.T) {
	cmd := binfmt.Command("go", "version")
	if err := testArgs(cmd.Args, []string{"go", "version"}); err != nil {
		t.Error(err)
	}

	cmd = binfmt.Command(".")
	if err := testArgs(cmd.Args, []string{"."}); err != nil {
		t.Error(err)
	}
}

func testArgs(g, e []string) error {
	if len(g) == len(e) {
		i := 0
		for ; i < len(e) && g[i] == e[i]; i++ {
		}
		if i == len(e) {
			return nil
		}
	}
	return fmt.Errorf("expected %#v, got %#v", e, g)
}

func write(name, data string) error {
	return ioutil.WriteFile(name, []byte(data), 0666)
}

func tempDir() (string, error) {
	return ioutil.TempDir("", "go.binfmt.test")
}
