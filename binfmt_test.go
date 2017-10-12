//
// go.binfmt :: binfmt_test.go
//
//   Copyright (c) 2014-2017 Akinori Hattori <hattya@gmail.com>
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
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hattya/go.binfmt"
)

func TestCommand(t *testing.T) {
	dir, err := tempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	restore := binfmt.Save()
	defer restore()

	v := 0
	funcs := []interface{}{
		func([]string) *exec.Cmd {
			v |= 1 << 0
			return nil
		},
		func(context.Context, []string) *exec.Cmd {
			v |= 1 << 1
			return nil
		},
		func(io.Reader, []string) *exec.Cmd {
			v |= 1 << 2
			return nil
		},
		func(context.Context, io.Reader, []string) *exec.Cmd {
			v |= 1 << 3
			return nil
		},
	}
	for i, fn := range funcs {
		binfmt.Register(string('a'+i), fn)
	}

	script := filepath.Join(dir, "script")
	if err := writeFile(script, ""); err != nil {
		t.Fatal(err)
	}

	for _, args := range [][]string{
		{"go", "version"},
		{"."},
		{script},
	} {
		cmd := binfmt.Command(args[0], args[1:]...)
		if g, e := cmd.Args, args; !reflect.DeepEqual(g, e) {
			t.Errorf("expected %v, got %v", e, g)
		}
	}

	for i, fn := range funcs {
		if v&(1<<uint(i)) == 0 {
			t.Errorf("%T is not called", fn)
		}
	}
}

func TestRegister(t *testing.T) {
	restore := binfmt.Save()
	defer restore()

	binfmt.Register("a", func([]string) *exec.Cmd { return nil })
	binfmt.Register("b", func(context.Context, []string) *exec.Cmd { return nil })
	binfmt.Register("c", func(io.Reader, []string) *exec.Cmd { return nil })
	binfmt.Register("d", func(context.Context, io.Reader, []string) *exec.Cmd { return nil })

	func() {
		defer func() {
			if recover() == nil {
				t.Error("expected panic")
			}
		}()

		binfmt.Register("painc", func() {})
	}()
}

func writeFile(name, data string) error {
	return ioutil.WriteFile(name, []byte(data), 0666)
}

func tempDir() (string, error) {
	return ioutil.TempDir("", "go.binfmt.test")
}
