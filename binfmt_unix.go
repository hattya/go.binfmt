//
// go.binfmt :: binfmt_unix.go
//
//   Copyright (c) 2014-2023 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

//go:build unix

package binfmt

func init() {
	Register("script", script)
}
