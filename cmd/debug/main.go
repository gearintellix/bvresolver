package main

import (
	"github.com/gearintellix/bvresolver"
)

func main() {
	bld, err := bvresolver.NewBivrostResolver("default")
	if err != nil {
		panic(err)
	}

	_ = bld
}
