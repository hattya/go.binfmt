//
// go.binfmt :: binfmt_windows.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package binfmt

func init() {
	Register("extension", extension)
	Register("script", script)
}
