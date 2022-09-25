package id

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	//uuip "github.com/pborman/uuid"
	"github.com/sempernow/uqc/kit/timestamp"
	"github.com/sempernow/uqc/kit/types"

	"github.com/gofrs/uuid"

	"github.com/oklog/ulid"
	"lukechampine.com/blake3"
)

// UUIDv5 returns the static ID per namespace (DNS, OID, URL or x500)
// and name, both case-insensitive.
func UUIDv5(space, name string) (string, error) {

	if name == "" {
		return "", errors.New("missing NAME")
	}

	ns := uuid.NamespaceDNS
	switch strings.ToUpper(space) {
	case "DNS":
		ns = uuid.NamespaceDNS
	case "OID":
		ns = uuid.NamespaceOID
	case "URL":
		ns = uuid.NamespaceURL
	case "X500":
		ns = uuid.NamespaceX500
	default:
		return "", errors.New(
			"missing or bad NAMESPACE (dns, oid, url or x500)",
		)
	}
	return uuid.NewV5(ns, strings.ToLower(name)).String(), nil
}

// SumBlake3_256 of  string, int, ... to hex string
func SumBlake3_256(i interface{}) string {
	xx := [32]byte{}
	switch x := i.(type) {
	case []byte:
		xx = blake3.Sum256(x)
	case string:
		xx = blake3.Sum256([]byte(x))
	case int:
		xx = blake3.Sum256([]byte(strconv.Itoa(x)))
	case int64:
		xx = blake3.Sum256(types.Int64ToBytes(x))
	case uint64:
		xx = blake3.Sum256([]byte(types.Uint64ToString(x)))
	default:
		return types.BytesToString(xx[:])
	}
	return hex.EncodeToString(xx[:])
}

// SumBlake3_512 of  string, int, ... to hex string
func SumBlake3_512(i interface{}) string {
	xx := [64]byte{}
	switch x := i.(type) {
	case []byte:
		xx = blake3.Sum512(x)
	case string:
		xx = blake3.Sum512([]byte(x))
	case int:
		xx = blake3.Sum512([]byte(strconv.Itoa(x)))
	case int64:
		xx = blake3.Sum512(types.Int64ToBytes(x))
	case uint64:
		xx = blake3.Sum512([]byte(types.Uint64ToString(x)))
	default:
		return types.BytesToString(xx[:])
	}
	return hex.EncodeToString(xx[:])
}

// SumSHA1 of *os.File, string, int, ... to hex string
func SumSHA1(i interface{}) string {
	xx := [20]byte{}
	switch x := i.(type) {
	case []byte:
		xx = sha1.Sum(x)
	case string:
		xx = sha1.Sum([]byte(x))
	case int:
		xx = sha1.Sum([]byte(strconv.Itoa(x)))
	case int64:
		xx = sha1.Sum(types.Int64ToBytes(x))
	case uint64:
		xx = sha1.Sum([]byte(types.Uint64ToString(x)))
	case *os.File:
		h := sha1.New()
		if _, err := io.Copy(h, x); err != nil {
			return types.BytesToString(xx[:])
		}
		return fmt.Sprintf("%x", h.Sum(nil))
	default:
		return types.BytesToString(xx[:])
	}
	return hex.EncodeToString(xx[:])
}

// SumSHA256 of *os.File, string, int, ... to hex string
func SumSHA256(i interface{}) string {
	xx := [32]byte{}
	switch x := i.(type) {
	case []byte:
		xx = sha256.Sum256(x)
	case string:
		xx = sha256.Sum256([]byte(x))
	case int:
		xx = sha256.Sum256([]byte(strconv.Itoa(x)))
	case int64:
		xx = sha256.Sum256(types.Int64ToBytes(x))
	case uint64:
		xx = sha256.Sum256([]byte(types.Uint64ToString(x)))
	case *os.File:
		h := sha256.New()
		if _, err := io.Copy(h, x); err != nil {
			return types.BytesToString(xx[:])
		}
		return fmt.Sprintf("%x", h.Sum(nil))
	default:
		return types.BytesToString(xx[:])
	}
	return hex.EncodeToString(xx[:])
}

// NewID is a base32 encoded UUID v4 GUID sans padding.
// It is a 26 character alpha-num [a-z0-9] string.
// https://sourcegraph.com/github.com/mattermost/mattermost-server/-/blob/model/utils.go#L201
func NewID() string {
	// Define the alphabet
	var encoding = base32.NewEncoding("s6dwh85iz12qmkueap4cgfbr7nv3xytj")
	var b bytes.Buffer
	encoder := base32.NewEncoder(encoding, &b)
	x := uuid.Must(uuid.NewV4())
	encoder.Write(x.Bytes())
	encoder.Close()
	b.Truncate(26)
	return b.String()
}

// // Etag : NNN-NNN.part1.part2.part3....
// func Etag(cfg Config, parts ...interface{}) string {
// 	etag := cfg.ver // SERVER_VER + "-" + APP_VER
// 	for _, part := range parts {
// 		etag += fmt.Sprintf(".%v", part)
// 	}
// 	return etag
// }

// UUID is 36 runes; 32+ 4 hyphens
// xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Format
// UUID v4 a.k.a. "Random"

// NowULID generates a ULID having a timestamp of `tt[0]` else `time.Now()`.
func NowULID(tt ...time.Time) ulid.ULID { // 26 runes
	//t := time.Unix(0, 0) // 000000000006AFVGQT5ZYC0GEK
	var t time.Time
	if len(tt) == 0 {
		t = time.Now()
	} else {
		t = tt[0]
	} // Nanosecond precision
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

// DoULID ...
func DoULID() {
	id := NowULID()
	debug(id.String())
	parse, _ := ulid.Parse("01EC5ZTYF6C7CMX079XS81D1G9")
	debug(parse)
	debug(string(id.Entropy()))
	debug(id.Time()) // Msec
	debug(timestamp.NowEpochMsec())
	debug(ulid.Time(id.Time()))
	debug(time.Now())
	debug(NowULID(time.Now()))
}
func debug(msgs interface{}) {
	fmt.Println(msgs)
}
