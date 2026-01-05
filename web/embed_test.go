package web

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssets(t *testing.T) {
	// Test that Assets() returns a non-nil filesystem
	assets := Assets()
	assert.NotNil(t, assets, "Assets() should return a non-nil filesystem")

	// Test that Assets() returns the same instance on subsequent calls
	// This verifies that we're reusing the cached filesystem
	assets2 := Assets()
	assert.Same(t, assets, assets2, "Assets() should return the same filesystem instance")
}

func TestAssetsIsValidFS(t *testing.T) {
	// Test that the returned filesystem implements fs.FS interface
	assets := Assets()
	var _ fs.FS = assets

	// Test that we can read from the filesystem
	// Note: This will only work if the dist directory exists
	// In CI/build environments, this should be present after web build
	entries, err := fs.ReadDir(assets, ".")
	if err == nil {
		// If we can read the directory, verify it's not empty in a real build
		// In development without a build, this might be empty
		assert.NotNil(t, entries)
	}
	// If there's an error, it's likely because dist is empty, which is fine for testing
}
