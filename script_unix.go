//
// go.binfmt :: script_unix.go
//
//   Copyright (c) 2014-2023 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

//go:build unix

package binfmt

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"
)

func script(ctx context.Context, r io.Reader, args []string) *exec.Cmd {
	br := bufio.NewReader(r)
	// check #!
	b := make([]byte, 2)
	br.Read(b)
	if b[0] != '#' || b[1] != '!' {
		return nil
	}

	l, err := br.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil
	}
	args = append(strings.Fields(strings.TrimSpace(l)), args...)
	if len(args) < 2 {
		return nil
	}
	return exec.CommandContext(ctx, args[0], args[1:]...)
}
