//
// go.binfmt :: script_windows_test.go
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
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hattya/go.binfmt"
)

type scriptTest struct {
	data string
	args []string
}

func TestScript(t *testing.T) {
	dir, err := tempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	defer os.Setenv("PATH", os.Getenv("PATH"))
	os.Setenv("PATH", dir)

	script := filepath.Join(dir, "script")
	python := filepath.Join(dir, "python.exe")

	for _, tt := range []scriptTest{
		{
			data: fmt.Sprintf("#! %v\n", python),
			args: []string{script},
		},
		{
			data: fmt.Sprintf("\xee\xbb\xbf#! %v\n", python),
			args: []string{script},
		},
		{
			data: "#! /usr/bin/env python",
			args: []string{script},
		},
		{
			data: "\xee\xbb\xbf#! /usr/bin/env python",
			args: []string{script},
		},
	} {
		if err := writeFile(script, tt.data); err != nil {
			t.Fatal(err)
		}
		cmd := binfmt.Command(script)
		if g, e := cmd.Args, tt.args; !reflect.DeepEqual(g, e) {
			t.Errorf("expected %v, got %v", e, g)
		}
	}

	if err := writeFile(python, ""); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []scriptTest{
		{
			data: fmt.Sprintf("#! %v\n", python),
			args: []string{python, script},
		},
		{
			data: fmt.Sprintf("\xee\xbb\xbf#! %v\n", python),
			args: []string{python, script},
		},
		{
			data: "#! /usr/bin/env python\n",
			args: []string{"python", script},
		},
		{
			data: "\xee\xbb\xbf#! /usr/bin/env python\n",
			args: []string{"python", script},
		},
	} {
		if err := writeFile(script, tt.data); err != nil {
			t.Fatal(err)
		}
		cmd := binfmt.Command(script)
		if g, e := cmd.Args, tt.args; !reflect.DeepEqual(g, e) {
			t.Errorf("expected %v, got %v", e, g)
		}
	}
}
