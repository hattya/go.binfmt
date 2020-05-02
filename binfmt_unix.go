//
// go.binfmt :: binfmt_unix.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

// +build !plan9,!windows

package binfmt

func init() {
	Register("script", script)
}
