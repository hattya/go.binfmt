//
// go.binfmt :: binfmt.go
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

package binfmt

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

func Command(name string, arg ...string) *exec.Cmd {
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
	for i := len(formats) - 1; 0 <= i; i-- {
		var cmd *exec.Cmd
		switch command := formats[i].command.(type) {
		case func([]string) *exec.Cmd:
			cmd = command(args)
		case func(io.Reader, []string) *exec.Cmd:
			if f == nil {
				continue
			}
			if _, err := f.Seek(0, os.SEEK_SET); err != nil {
				break L
			}
			cmd = command(f, args)
		}
		if cmd != nil {
			return cmd
		}
	}
	return exec.Command(name, arg...)
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
	case func(io.Reader, []string) *exec.Cmd:
	default:
		panic("binfmt: unknown type")
	}
	formats = append(formats, format{name, command})
}
