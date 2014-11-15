//
// go.binfmt :: binfmt_windows_test.go
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
	"os"
	"path/filepath"
	"testing"

	"github.com/hattya/go.binfmt"
)

var (
	SystemDrive  = os.Getenv("SystemDrive")
	SystemRoot   = os.Getenv("SystemRoot")
	ProgramFiles = os.Getenv("ProgramFiles")
)

func TestExtension(t *testing.T) {
	cmd := binfmt.Command("a.txt")
	if err := testArgs(cmd.Args, []string{filepath.Join(SystemRoot, "system32", "NOTEPAD.EXE"), "a.txt"}); err != nil {
		t.Error(err)
	}

	cmd = binfmt.Command("a.bat", "1", "2")
	if err := testArgs(cmd.Args, []string{"a.bat", "1", "2"}); err != nil {
		t.Error(err)
	}
}

func TestEvalCommand(t *testing.T) {
	rst2html := []string{"rst2html.py", "README.rst", "README.html"}

	python := filepath.Join(SystemDrive, "PythonXY", "python.exe")
	args := binfmt.EvalCommand(fmt.Sprintf(`"%s" "%%1" %%*`, python), rst2html)
	if err := testArgs(args, append([]string{python}, rst2html...)); err != nil {
		t.Error(err)
	}
	python = filepath.Join(ProgramFiles, "PythonXY", "python.exe")
	args = binfmt.EvalCommand(fmt.Sprintf(`"%s" "%%1" %%*`, python), rst2html)
	if err := testArgs(args, append([]string{python}, rst2html...)); err != nil {
		t.Error(err)
	}

	python = filepath.Join(SystemDrive, "PythonXY", "python")
	args = binfmt.EvalCommand(fmt.Sprintf(`"%s" "%%1" %%*`, python), rst2html)
	if err := testArgs(args, []string{}); err != nil {
		t.Error(err)
	}
	python = filepath.Join(ProgramFiles, "PythonXY", "python")
	args = binfmt.EvalCommand(fmt.Sprintf(`"%s" "%%1" %%*`, python), rst2html)
	if err := testArgs(args, []string{}); err != nil {
		t.Error(err)
	}

	notepad := filepath.Join(SystemRoot, "system32", "NOTEPAD.EXE")
	args = binfmt.EvalCommand(fmt.Sprintf(`%s %%a`, notepad), []string{})
	if err := testArgs(args, []string{}); err != nil {
		t.Error(err)
	}
	args = binfmt.EvalCommand(fmt.Sprintf(`%s %%1`, notepad), []string{})
	if err := testArgs(args, []string{}); err != nil {
		t.Error(err)
	}
	args = binfmt.EvalCommand(fmt.Sprintf(`%s %%2`, notepad), []string{})
	if err := testArgs(args, []string{}); err != nil {
		t.Error(err)
	}
}

func TestHRESULT(t *testing.T) {
	hr := binfmt.NewHRESULT(0)
	if g, e := hr.Error(), "0x00000000"; g != e {
		t.Errorf("expected %v, got %q", e, g)
	}

	hr = binfmt.NewHRESULT(1)
	if g, e := hr.Error(), "0x00000001"; g != e {
		t.Errorf("expected %v, got %q", e, g)
	}

	hr = binfmt.NewHRESULT(0x80070483)
	if g, e := hr.Error(), "0x80070483"; g == e {
		t.Errorf("expected error message, got %q", g)
	}
}
