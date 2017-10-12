//
// go.binfmt :: binfmt_windows.go
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
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unsafe"

	"golang.org/x/sys/windows"
)

func extension(ctx context.Context, args []string) *exec.Cmd {
	ext := filepath.Ext(args[0])
	switch ext {
	case "", ".":
		return nil
	}
	assoc, err := windows.UTF16PtrFromString(ext)
	if err != nil {
		return nil
	}
	var size uint32
	if assocQueryString(_ASSOCF_NOTRUNCATE, _ASSOCSTR_COMMAND, assoc, nil, nil, &size) != _S_FALSE {
		return nil
	}
	out := make([]uint16, size)
	if assocQueryString(_ASSOCF_NOTRUNCATE, _ASSOCSTR_COMMAND, assoc, nil, &out[0], &size) != _S_OK {
		return nil
	}

	command := evalCommand(windows.UTF16ToString(out), args)
	if len(command) < 2 {
		return nil
	}
	return exec.CommandContext(ctx, command[0], command[1:]...)
}

func evalCommand(s string, args []string) []string {
	// parse
	var command []string
	var ok bool
	if s[0] == '"' {
		command = commandFields(s)
		if 0 < len(command) {
			ok = strings.ToLower(filepath.Ext(command[0])) == ".exe"
		}
	} else {
		i := strings.Index(strings.ToLower(s), ".exe")
		if i != -1 {
			command = append([]string{s[:i+4]}, commandFields(s[i+4:])...)
			ok = true
		}
	}
	if !ok {
		return nil
	}
	// eval
	n := len(args)
	i := 0
	for j, a := range command {
		if a[0] == '%' {
			if len(a) == 2 && a[1] == '*' {
				command = append(command[:j], args[i:]...)
				break
			}
			v, err := strconv.ParseInt(a[1:], 10, 0)
			if err != nil || int(v) != i+1 || n < int(v) {
				return nil
			}
			command[j] = args[i]
			i++
		}
	}
	return command
}

func commandFields(s string) []string {
	q := false
	return strings.FieldsFunc(s, func(r rune) bool {
		if r == '"' {
			q = !q
			return true
		}
		return !q && unicode.IsSpace(r)
	})
}

var (
	shlwapi = windows.NewLazySystemDLL("shlwapi.dll")

	pAssocQueryString = shlwapi.NewProc("AssocQueryStringW")
)

// type HRESULT
type hresult int32

func (hr hresult) Error() string {
	i := uint32(hr)
	if 1 < i {
		flags := uint32(windows.FORMAT_MESSAGE_FROM_SYSTEM | windows.FORMAT_MESSAGE_ARGUMENT_ARRAY | windows.FORMAT_MESSAGE_IGNORE_INSERTS)
		b := make([]uint16, 300)
		if _, err := windows.FormatMessage(flags, 0, i, 0, b, nil); err == nil {
			return strings.TrimSpace(windows.UTF16ToString(b))
		}
	}
	return fmt.Sprintf("0x%08x", i)
}

const (
	_S_OK hresult = iota
	_S_FALSE
)

// enum ASSOCF
type assocf int32

const _ASSOCF_NOTRUNCATE assocf = 0x00000020

// enum ASSOCSTR
type assocstr int32

const _ASSOCSTR_COMMAND assocstr = 1

func assocQueryString(flags assocf, str assocstr, assoc, extra, out *uint16, size *uint32) error {
	r0, _, _ := pAssocQueryString.Call(uintptr(flags), uintptr(str), uintptr(unsafe.Pointer(assoc)), uintptr(unsafe.Pointer(extra)), uintptr(unsafe.Pointer(out)), uintptr(unsafe.Pointer(size)))
	return hresult(r0)
}
