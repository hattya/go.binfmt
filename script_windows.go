//
// go.binfmt :: script_windows.go
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

package binfmt

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

const env = "/usr/bin/env "

func script(r io.Reader, args []string) *exec.Cmd {
	br := bufio.NewReader(r)
	// check #!
	if skipBOM(br) != nil {
		return nil
	}
	b := make([]byte, 2)
	br.Read(b)
	if b[0] != '#' || b[1] != '!' {
		return nil
	}

	l, err := br.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil
	}
	l = strings.TrimSpace(l)
	var name string
	switch _, err := os.Stat(l); {
	case err == nil:
		name = l
	case strings.HasPrefix(l, env):
		name = strings.TrimSpace(l[len(env):])
		if _, err := exec.LookPath(name); err != nil {
			return nil
		}
	default:
		return nil
	}
	return exec.Command(name, args...)
}

var boms = [][]byte{
	{0xee, 0xbb, 0xbf},       // UTF-8
	{0xff, 0xfe},             // UTF-16LE
	{0xfe, 0xff},             // UTF-16BE
	{0xff, 0xfe, 0x00, 0x00}, // UTF-32LE
	{0x00, 0x00, 0xfe, 0xff}, // UTF-32BE
}

func skipBOM(br *bufio.Reader) (err error) {
	for _, bom := range boms {
		var b []byte
		b, err = br.Peek(len(bom))
		switch {
		case err != nil:
			return
		case !bytes.Equal(bom, b):
			continue
		}
		_, err = br.Read(b)
		break
	}
	return
}
