// Package client provides an http client of uqrate services.
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// GhostPrint(..) prints args per format, e.g., "want: %s\nhave: %s\n", to os.Stderr.
func GhostPrint(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// SetCache writes data to key file in Env.Cache folder.
func (env *Env) SetCache(key, data string) error {
	if key == "" {
		GhostPrint("\nWARN @ SetCache : missing key\n")
		return errors.New("missing key")
	}
	if data == "" {
		GhostPrint("\nWARN @ SetCache : NO DATA\n")
		return errors.New("no data")
	}
	err := os.WriteFile(
		filepath.Join(env.Cache, key),
		[]byte(data),
		0664,
	)
	if err != nil {
		GhostPrint("\nERR @ SetCache : writing file: %v\n", err)
		return err
	}
	fi, err := os.Stat(filepath.Join(env.Cache, key))
	if err == nil {
		GhostPrint("\nINFO @ SetCache : size: %v @ file: '%s'\n", fi.Size(), key)
	}
	return nil
}

// GetCache reads key file of Env.Cache folder.
func (env *Env) GetCache(key string) []byte {
	if key == "" {
		GhostPrint("\nERR @ GetCache : missing key\n")
		return []byte{}
	}
	bb, err := os.ReadFile(filepath.Join(env.Cache, key))
	if err != nil {
		GhostPrint("\nINFO @ GetCache : miss @ %s\n", key)
		return []byte{}
	}
	//GhostPrint("\n@ %s_CACHE : %s\n", env.NS, filepath.Join(env.Cache, key))
	return bb
}

// GetCacheJSON reads key file of Env.Cache folder into struct of pointer.
func (env *Env) GetCacheJSON(key string, ptr interface{}) {
	if key == "" {
		GhostPrint("\nERR @ GetCacheJSON : missing key\n")
	}
	bb, err := os.ReadFile(filepath.Join(env.Cache, key))
	if err != nil {
		GhostPrint("\nERR @ GetCacheJSON : reading file: %s\n", err)
	}
	if err := json.Unmarshal(bb, &ptr); err != nil {
		GhostPrint("\nERR @ GetCacheJSON : decoding JSON: %s\n", err)

	}
}
