package types_test

import (
	"strconv"
	"testing"

	"github.com/sempernow/uqc/kit/types"
)

// ******************************************************
// go test -count=1 -benchmem -bench=EncodeSecret ./kit
// ******************************************************

// ================================================
// ToString functions

var i uint64 = 999999999999

func BenchmarkFormatUint(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		strconv.FormatUint(i, 10)
	}
}

func BenchmarkUint64ToString(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		types.Uint64ToString(i)
	}
}

func BenchmarkToString(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		types.ToString(i)
	} //... 3x SLOWER than other 3 "ToString" funcs.
}

func BenchmarkIntToString(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		types.IntToString(int(999999999999))
	}
}

func BenchmarkInt64ToString(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		types.Int64ToString(int64(999999999999))
	}
}

// goos: windows
// goarch: amd64

// BenchmarkFormatUint-4       	21818617	        53.3 ns/op	      16 B/op	       1 allocs/op
// BenchmarkUint64ToString-4   	23076434	        52.7 ns/op	      16 B/op	       1 allocs/op
// BenchmarkToString-4         	 8275861	       145 ns/op	      24 B/op	       2 allocs/op
// BenchmarkIntToString-4      	23530380	        51.0 ns/op	      16 B/op	       1 allocs/op
