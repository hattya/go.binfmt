go.binfmt
=========

An extension library for `os/exec`_.

.. image:: https://semaphoreci.com/api/v1/hattya/go-binfmt/branches/master/badge.svg
   :target: https://semaphoreci.com/hattya/go-binfmt

.. image:: https://ci.appveyor.com/api/projects/status/uhkkibn9gen71du9/branch/master?svg=true
   :target: https://ci.appveyor.com/project/hattya/go-binfmt

.. _os/exec: https://golang.org/pkg/os/exec/


Installation
------------

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
   	cmd.Stdin = os.Stdin
   	cmd.Stdout = os.Stdout
   	cmd.Stderr = os.Stderr
   	if err := cmd.Run(); err != nil {
   		fmt.Fprintln(os.Stderr, err)
   		os.Exit(1)
   	}
   }


License
-------

go.binfmt is distributed under the terms of the MIT License
