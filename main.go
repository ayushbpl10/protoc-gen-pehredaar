package main

import (
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func main() {
	pgs.Init().
		RegisterModule(&RightsGen{ModuleBase: pgs.ModuleBase{}}).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
