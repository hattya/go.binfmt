//
// go.binfmt :: binfmt_test.go
//
//   Copyright (c) 2014-2025 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package binfmt_test

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hattya/go.binfmt"
)

func TestCommand(t *testing.T) {
	restore := binfmt.Save()
	defer restore()

	v := 0
	funcs := []any{
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
		binfmt.Register(string(rune('a'+i)), fn)
	}

	script := filepath.Join(t.TempDir(), "script")
	if err := file(script, ""); err != nil {
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

func file(name, data string) error {
	return os.WriteFile(name, []byte(data), 0o666)
}
