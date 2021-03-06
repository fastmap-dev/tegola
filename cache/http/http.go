package http

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-spatial/tegola"
	"github.com/go-spatial/tegola/cache"
	"github.com/go-spatial/tegola/dict"
)

var (
	ErrMissingUrl = errors.New("filecache: missing required param 'url'")
)

const CacheType = "http"

const (
	ConfigKeyUrl     = "url"
	ConfigKeyMaxZoom = "max_zoom"
)

func init() {
	cache.Register(CacheType, New)
}

// New instantiates a Cache. The config expects the following params:
//
// 	url (string): a path to where the cache will be written
// 	max_zoom (int): max zoom to use the cache. beyond this zoom cache Set() calls will be ignored
//
func New(config dict.Dicter) (cache.Interface, error) {
	var err error

	// new httpcache
	fc := Cache{}

	defaultMaxZoom := uint(tegola.MaxZ)
	fc.MaxZoom, err = config.Uint(ConfigKeyMaxZoom, &defaultMaxZoom)
	if err != nil {
		return nil, err
	}

	fc.URL, err = config.String(ConfigKeyUrl, nil)
	if err != nil {
		return nil, ErrMissingUrl
	}

	if fc.URL == "" {
		return nil, ErrMissingUrl
	}

	return &fc, nil
}

// Cache ...
type Cache struct {
	URL string
	// MaxZoom determines the max zoom the cache to persist. Beyond this
	// zoom, cache Set() calls will be ignored. This is useful if the cache
	// should not be leveraged for higher zooms when data changes often.
	MaxZoom uint
}

//Get reads a z,x,y entry from the cache and returns the contents
// if there is a hit. the second argument denotes a hit or miss
// so the consumer does not need to sniff errors for cache read misses
func (fc *Cache) Get(key *cache.Key) ([]byte, bool, error) {
	requestURL := fmt.Sprintf("%s/%d/%d/%d", fc.URL, key.Z, key.X, key.Y)

	// Get the data
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, false, nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, true, nil
}

// Set data to cache
func (fc *Cache) Set(key *cache.Key, val []byte) error {
	strURL := fmt.Sprintf("%s/%d/%d/%d", fc.URL, key.Z, key.X, key.Y)
	resp, err := http.Post(strURL, "binary/octet-stream", bytes.NewReader(val))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// Purge request
func (fc *Cache) Purge(key *cache.Key) error {
	var err error
	requestURL := fmt.Sprintf("%s/%d/%d/%d", fc.URL, key.Z, key.X, key.Y)
	req, err := http.NewRequest("DELETE", requestURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
