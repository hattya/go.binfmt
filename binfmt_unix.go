//
// go.binfmt :: binfmt_unix.go
//
//   Copyright (c) 2014-2022 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

//go:build !plan9 && !windows
// +build !plan9,!windows

package binfmt

func init() {
	Register("script", script)
}
