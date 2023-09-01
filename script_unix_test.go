//
// go.binfmt :: script_unix_test.go
//
//   Copyright (c) 2014-2023 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

//go:build unix

package binfmt_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hattya/go.binfmt"
)

func TestScript(t *testing.T) {
	dir := t.TempDir()
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
		if err := file(script, tt.data); err != nil {
			t.Fatal(err)
		}
		cmd := binfmt.Command(script)
		if g, e := cmd.Args, tt.args; !reflect.DeepEqual(g, e) {
			t.Errorf("expected %v, got %v", e, g)
		}
	}
}
