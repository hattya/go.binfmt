//
// go.binfmt :: script_unix_test.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

// +build !plan9,!windows

package binfmt_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hattya/go.binfmt"
)

func TestScript(t *testing.T) {
	dir, err := tempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	script := filepath.Join(dir, "script")
	sh := filepath.Join(dir, "sh")

	for _, tt := range []struct {
		data string
		args []string
	}{
		{
			data: fmt.Sprintf("#! %v\n", sh),
			args: []string{sh, script},
		},
		{
			data: "#!\n",
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
}
