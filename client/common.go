// Package client provides an http client of uqrate services.
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// GhostPrint(..) prints args per format, e.g., "want: %s\nhave: %s\n", to os.Stderr.
func GhostPrint(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// SetCache writes data to key file in Env.Cache folder.
func (env *Env) SetCache(key, data string) error {
	if key == "" || data == "" {
		msg := "@ env.SetCache(..) : NO DATA"
		GhostPrint("\n%s\n  key: '%s'\n  data: '%s'\n", msg, key, data)
		return errors.New(key)
	}
	err := ioutil.WriteFile(
		filepath.Join(env.Cache, key),
		[]byte(data),
		0664,
	)
	if err != nil {
		GhostPrint("\nERR @ env.SetCache(..) : writing file: %v\n", err)
		return err
	}
	fi, err := os.Stat(filepath.Join(env.Cache, key))
	if err == nil {
		GhostPrint("\nsize: %vKB @ file: '%s'\n", (fi.Size() / 1024), key)
	}
	return nil
}

// GetCache reads key file of Env.Cache folder.
func (env *Env) GetCache(key string) []byte {
	if key == "" {
		GhostPrint("\nERR @ env.GetCache(..) : missing key parameter\n")
		GhostPrint("\nkey: '%s'\n", key)
		return []byte{}
	}
	bb, err := ioutil.ReadFile(filepath.Join(env.Cache, key))
	if err != nil {
		GhostPrint("\nMISS @ %s\n", key)
		return []byte{}
	}
	//GhostPrint("\n@ %s_CACHE : %s\n", env.NS, filepath.Join(env.Cache, key))
	return bb
}

// GetCacheJSON reads key file of Env.Cache folder into struct of pointer.
func (env *Env) GetCacheJSON(key string, ptr interface{}) {
	if key == "" {
		GhostPrint("\nERR @ env.GetCacheJSON(..) : missing key parameter\n")
	}
	bb, err := ioutil.ReadFile(filepath.Join(env.Cache, key))
	if err != nil {
		GhostPrint("\nERR @ env.GetCacheJSON(..) : reading file: %s\n", err)
	}
	if err := json.Unmarshal(bb, &ptr); err != nil {
		GhostPrint("\nERR @ env.GetCacheJSON(..) : decoding JSON: %s\n", err)

	}
}
