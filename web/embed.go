package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var web embed.FS

// dist is the sub-filesystem computed once at initialization time
var dist fs.FS

func init() {
	var err error
	dist, err = fs.Sub(web, "dist")
	if err != nil {
		panic(err)
	}
}

func Assets() fs.FS {
	return dist
}
