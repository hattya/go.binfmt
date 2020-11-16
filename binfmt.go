//
// go.binfmt :: binfmt.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package binfmt

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func Command(name string, arg ...string) *exec.Cmd {
	return CommandContext(context.Background(), name, arg...)
}

func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	name = filepath.FromSlash(name)
	args := append([]string{name}, arg...)
	f, err := os.Open(name)
	if err == nil {
		if fi, err := f.Stat(); err != nil || fi.IsDir() {
			f.Close()
			f = nil
		}
	} else {
		if path, err := exec.LookPath(name); err == nil {
			args[0] = path
			f, _ = os.Open(path)
		}
	}
	if f != nil {
		defer f.Close()
	}

L:
	for i := len(formats) - 1; i >= 0; i-- {
		var cmd *exec.Cmd
		switch command := formats[i].command.(type) {
		case func([]string) *exec.Cmd:
			cmd = command(args)
		case func(context.Context, []string) *exec.Cmd:
			cmd = command(ctx, args)
		case func(io.Reader, []string) *exec.Cmd:
			if f == nil {
				continue
			}
			if _, err := f.Seek(0, io.SeekStart); err != nil {
				break L
			}
			cmd = command(f, args)
		case func(context.Context, io.Reader, []string) *exec.Cmd:
			if f == nil {
				continue
			}
			if _, err := f.Seek(0, io.SeekStart); err != nil {
				break L
			}
			cmd = command(ctx, f, args)
		}
		if cmd != nil {
			return cmd
		}
	}
	return exec.CommandContext(ctx, name, arg...)
}

var (
	mu      sync.RWMutex
	formats []format
)

type format struct {
	name    string
	command CommandFunc
}

type CommandFunc interface{}

func Register(name string, command CommandFunc) {
	mu.Lock()
	defer mu.Unlock()

	switch command.(type) {
	case func([]string) *exec.Cmd:
	case func(context.Context, []string) *exec.Cmd:
	case func(io.Reader, []string) *exec.Cmd:
	case func(context.Context, io.Reader, []string) *exec.Cmd:
	default:
		panic("binfmt: unknown type")
	}
	formats = append(formats, format{name, command})
}
