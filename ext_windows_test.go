//
// go.binfmt :: binfmt_windows_test.go
//
//   Copyright (c) 2014-2025 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
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

var (
	SystemDrive  = os.Getenv("SystemDrive")
	SystemRoot   = os.Getenv("SystemRoot")
	ProgramFiles = os.Getenv("ProgramFiles")
)

func TestExtension(t *testing.T) {
	for _, tt := range []struct {
		in, out []string
	}{
		{
			in:  []string{"a.txt"},
			out: []string{filepath.Join(SystemRoot, "system32", "NOTEPAD.EXE"), "a.txt"},
		},
		{
			in:  []string{"a.bat", "1", "2"},
			out: []string{"a.bat", "1", "2"},
		},
	} {
		cmd := binfmt.Command(tt.in[0], tt.in[1:]...)
		if g, e := cmd.Args, tt.out; !reflect.DeepEqual(g, e) {
			t.Errorf("expected %v, got %v", e, g)
		}
	}
}

func TestEvalCommand(t *testing.T) {
	python1 := filepath.Join(SystemDrive, "PythonXY", "python.exe")
	python2 := filepath.Join(ProgramFiles, "PythonXY", "python.exe")
	rst2html := []string{"rst2html.py", "README.rst", "README.html"}
	for _, tt := range []struct {
		python string
		args   []string
	}{
		{
			python: python1,
			args:   append([]string{python1}, rst2html...),
		},
		{
			python: python2,
			args:   append([]string{python2}, rst2html...),
		},
		{
			python: python1[:len(python1)-4],
		},
		{
			python: python2[:len(python2)-4],
		},
	} {
		args := binfmt.EvalCommand(fmt.Sprintf(`"%s" "%%1" %%*`, tt.python), rst2html)
		if g, e := args, tt.args; !reflect.DeepEqual(g, e) {
			t.Errorf("expected %v, got %v", e, g)
		}
	}

	notepad := filepath.Join(SystemRoot, "system32", "NOTEPAD.EXE")
	for _, s := range []string{
		`%s %%a`,
		`%s %%1`,
		`%s %%2`,
	} {
		args := binfmt.EvalCommand(fmt.Sprintf(s, notepad), nil)
		if g, e := args, []string(nil); !reflect.DeepEqual(g, e) {
			t.Errorf("expected %v, got %v", e, g)
		}
	}
}

func TestHRESULT(t *testing.T) {
	for i := range uint32(2) {
		hr := binfmt.NewHRESULT(i)
		if g, e := hr.Error(), fmt.Sprintf("0x%08x", i); g != e {
			t.Errorf("expected %v, got %q", e, g)
		}
	}

	hr := binfmt.NewHRESULT(0x80070483)
	if g, e := hr.Error(), fmt.Sprintf("0x%08x", hr); g == e {
		t.Errorf("expected error message, got %q", g)
	}
}
