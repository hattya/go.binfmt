//
// go.binfmt :: export_windows_test.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package binfmt

var EvalCommand = evalCommand

func NewHRESULT(hr uint32) error {
	return hresult(hr)
}
