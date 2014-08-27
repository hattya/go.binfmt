go.binfmt
=========

An extension library for `os/exec`_.

.. _os/exec: https://golang.org/pkg/os/exec/


Install
-------

.. code:: console

   $ go get -u github.com/hattya/go.binfmt


Usage
-----

.. code:: go

   package main

   import (
   	"fmt"
   	"os"

   	"github.com/hattya/go.binfmt"
   )

   func main() {
   	cmd := binfmt.Command("rst2html.py", "README.rst", "README.html")
   	if err := cmd.Run(); err != nil {
   		fmt.Fprintln(os.Stderr, err)
   		os.Exit(1)
   	}
   }


License
-------

go.binfmt is distributed under the terms of the MIT License
