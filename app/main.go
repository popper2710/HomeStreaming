package main

import (
	"github.com/my-repo/home_streaming/src"
)

func main() {
	r := src.Router()
	if err := r.Run(":3000"); err != nil {
		panic(err)
	}
}
