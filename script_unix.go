//
// go.binfmt :: script_unix.go
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

// +build !plan9,!windows

package binfmt

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

func script(r io.Reader, args []string) *exec.Cmd {
	br := bufio.NewReader(r)
	// check #!
	b := make([]byte, 2)
	br.Read(b)
	if b[0] != '#' || b[1] != '!' {
		return nil
	}

	l, err := br.ReadString('\n')
	switch err {
	case nil, io.EOF:
	default:
		return nil
	}
	args = append(strings.Fields(strings.TrimSpace(l)), args...)
	if len(args) < 2 {
		return nil
	}
	return exec.Command(args[0], args[1:]...)
}
