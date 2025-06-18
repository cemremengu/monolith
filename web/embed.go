package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var web embed.FS

func Assets() fs.FS {
	dist, err := fs.Sub(web, "dist")
	if err != nil {
		panic(err)
	}

	return dist
}
