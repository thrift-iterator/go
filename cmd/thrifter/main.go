package main

import (
	"flag"
	"github.com/v2pro/wombat"
	"os"
)

func main() {
	pkgPath := flag.String("pkg", "", "the package to generate generic code for")
	flag.Parse()
	if *pkgPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	wombat.Codegen(*pkgPath)
}