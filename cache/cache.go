package cache

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

const DefaultTTL = 5 * time.Minute

type CacheKey string

func (k CacheKey) String() string {
	return string(k)
}

func (k CacheKey) WithArg(arg any) CacheKey {
	return CacheKey(fmt.Sprintf("%s_%v", k, arg))
}

const (
	KeyCloudServers     CacheKey = "cloud_servers"
	KeyVpsImages        CacheKey = "vps_images"
	KeyVpsProducts      CacheKey = "vps_products"
	KeyFlavours         CacheKey = "flavours"
	KeyAttachedNetworks CacheKey = "attached_networks"
	KeyVirtualNetworks  CacheKey = "virtual_networks"
)

const NoCache = CacheKey("")

type cacheEntry[T any] struct {
	Timestamp time.Time `json:"ts"`
	Data      T         `json:"data"`
}

func cacheFilePath(key CacheKey) (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "oh")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, string(key)+".json"), nil
}

// loadEntry reads and unmarshals a cacheEntry[T] from disk.
func loadEntry[T any](path string) (*cacheEntry[T], error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var e cacheEntry[T]
	if err := json.Unmarshal(b, &e); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &e, nil
}

// saveEntry marshals and writes a cacheEntry[T] to disk.
func saveEntry[T any](path string, e *cacheEntry[T]) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// Call wraps a fetch func with on-disk caching
func Call[T any](key CacheKey, ttl time.Duration, fetch func() (T, error)) (T, error) {
	// Honor the global --no-cache flag or explicit NoCache key
	if viper.GetBool("no-cache") || key == NoCache {
		return fetch()
	}

	var zero T

	path, err := cacheFilePath(key)
	if err != nil {
		return zero, err
	}

	if e, err := loadEntry[T](path); err == nil {
		if time.Since(e.Timestamp) < ttl {
			return e.Data, nil
		}
	}

	fresh, err := fetch()
	if err != nil {
		return zero, err
	}

	entry := &cacheEntry[T]{Timestamp: time.Now(), Data: fresh}
	_ = saveEntry(path, entry) // best-effort

	return fresh, nil
}

// Store explicitly overwrites the cache
func Store[T any](key CacheKey, data T) error {
	path, err := cacheFilePath(key)
	if err != nil {
		return err
	}
	entry := &cacheEntry[T]{Timestamp: time.Now(), Data: data}
	return saveEntry(path, entry)
}
