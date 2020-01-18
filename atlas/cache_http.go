// +build !noHttpCache

package atlas

// The point of this file is to load and register the http cache backend.
// the http cache can be excluded during the build with the `noHttpCache` build flag
// for example from the cmd/tegola direcotry:
//
// go build -tags 'noHttpCache'
import (
	_ "github.com/go-spatial/tegola/cache/http"
)
