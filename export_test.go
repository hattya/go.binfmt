//
// go.binfmt :: export_test.go
//
//   Copyright (c) 2017-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package binfmt

func Save() func() {
	save := formats[:]
	return func() { formats = save }
}
