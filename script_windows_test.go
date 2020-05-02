//
// go.binfmt :: script_windows_test.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
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
