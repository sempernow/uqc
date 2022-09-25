package meta

import (
	crnd "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/sempernow/uqc/kit/types"

	"github.com/pkg/errors"
)

// https://golang.org/ref/spec#Handling_panics
func ProtectFromPanicAt(g func()) {
	defer func() {
		recover()
	}()
	g()
}

// ----------------------------------------------------------------------------
// DEBUG

type PrintDebug func(...interface{})

// NewDebug returns a fmt.Println(...interface{}) that prints only when in debug mode (on=true).
func NewDebug(on bool) PrintDebug {
	return func(msgs ...interface{}) {
		if !on {
			return
		}
		fmt.Println(msgs)
	}
}

// Unique returns input slice stripped of redundant elements
func Unique(ss []string) []string {
	sort.Strings(ss)
	j := 0
	for i := 1; i < len(ss); i++ {
		if ss[j] == ss[i] {
			continue
		}
		j++
		ss[j] = ss[i]
	}
	return ss[:j+1]
}

// ----------------------------------------------------------------------------
// FILESYSTEM

// IsDir ...
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := fi.Mode()
	return mode.IsDir()
}

// IsRegFile ...
func IsRegFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := fi.Mode()
	return mode.IsRegular()
}

// StripQuotes from a string
// https://stackoverflow.com/questions/44222554/how-to-remove-quotes-from-around-a-string-in-golang#44222606
func StripQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go

// RandAlphaNum ... 5x FASTER; 92ns @ 16 chars; 160ns @ 32 chars
func RandAlphaNum(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// Nonce of size bytes.
func Nonce(size int) ([]byte, error) {
	bb := make([]byte, size)
	_, err := io.ReadFull(crnd.Reader, bb)
	if err != nil {
		return []byte{}, errors.Wrap(err, "Nonce : source of randomness unavailable")
	}
	return bb, nil
}

// UniqueStrings removes duplicates.
func UniqueStrings(ss []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, s := range ss {
		if _, val := keys[s]; !val {
			keys[s] = true
			list = append(list, s)
		}
	}
	return list
} // https://www.golangprograms.com/remove-duplicate-values-from-slice.html

// XOR ...
func XOR(input, key string) (output string) {
	var str strings.Builder

	for i := 0; i < len(input); i++ {
		str.WriteString(string(input[i] ^ key[i%len(key)]))
	}
	return str.String()
}

// EncodeSecret per XOR with nonce bisected for use as both key and padding.
func EncodeSecret(s, nonce string) (string, error) {
	n := nonce[0:(len(nonce)/2 - len(nonce)%2)]
	p := nonce[len(n):]
	//delta := (len(n) - len(s) - 1)
	delta := (len(n) - len(s))
	xx := (len(n) - (delta / 2) - (delta % 2))
	pad := p[0:(delta / 2)]
	//bb, err := XORbytes([]byte(s+"|"+pad), []byte(n[0:xx]))
	bb, err := XORbytes([]byte(s+pad), []byte(n[0:xx]))
	if err != nil {
		return "", err
	}
	return types.BytesToString(bb), nil
}

// DecodeSecret per XOR with its nonce key/pad.
func DecodeSecret(secret, nonce string) (string, error) {
	bb, err := XORbytes([]byte(secret), []byte(nonce[0:len(secret)]))
	if err != nil {
		return "", err
	}
	//return strings.Split(BytesToString(bb), "|")[0], nil
	return types.BytesToString(bb), nil
}

// XORbytes performs XOR (encrypt/decrypt) of two equal-sized byte slices.
// https://sourcegraph.com/github.com/hashicorp/vault/-/blob/helper/xor/xor.go#L8
func XORbytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("Slice LENGTHs differ : [%d]a vs. [%d]b", len(a), len(b))
	}
	buf := make([]byte, len(a))
	for i := range a {
		buf[i] = a[i] ^ b[i]
	}
	return buf, nil
}

// XORstrings performs XOR (encrypt/decrypt) on input against a nonce.
// https://kylewbanks.com/blog/xor-encryption-using-go
func XORstrings(input, nonce string) (output string) {
	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ nonce[i%len(nonce)])
	}
	return output
}

// Reverse a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
