//
// go.binfmt :: script_windows.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package binfmt

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

const env = "/usr/bin/env "

func script(ctx context.Context, r io.Reader, args []string) *exec.Cmd {
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
	return exec.CommandContext(ctx, name, args...)
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
