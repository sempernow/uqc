package meta_test

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/sempernow/uqc/kit/id"
	"github.com/sempernow/uqc/kit/meta"
	"github.com/sempernow/uqc/kit/types"

	"github.com/gofrs/uuid"
)

func TestOBA(t *testing.T) {
	t.Skip()
	var u uuid.UUID
	u, _ = uuid.NewV4()
	_na := sha512.Sum512(u.Bytes())
	u, _ = uuid.NewV4()
	_nb := sha512.Sum512(u.Bytes())

	na := hex.EncodeToString(_na[:])
	nb := hex.EncodeToString(_nb[:])

	email := "foo@bar.com"
	pass := "1234"

	email = meta.XOR(email, meta.Reverse(na))
	pass = meta.XOR(pass, meta.XOR(na, meta.Reverse(nb)))

	t.Log("email", hex.EncodeToString([]byte(email)))
	t.Log("pass", hex.EncodeToString([]byte(pass)))

	email = meta.XOR(email, meta.Reverse(na))
	pass = meta.XOR(pass, meta.XOR(na, meta.Reverse(nb)))

	t.Log(email, pass)

	//t.Fail()
}

// ******************************************************
// go test -count=1 -benchmem -bench=EncodeSecret ./kit
// ******************************************************

func BenchmarkEncodeSecret(b *testing.B) {
	// Run it b.N times
	secret := id.SumSHA256("a")
	nonce := id.SumSHA256("b")
	//x, _ := meta.DecodeSecret(secret[:], nonce)

	bb, _ := meta.XORbytes([]byte(secret), []byte(nonce[0:len(secret)]))
	fmt.Println("hex:", hex.EncodeToString(bb))

	for n := 0; n < b.N; n++ {
		bb, _ := meta.XORbytes([]byte(secret), []byte(nonce[0:len(secret)]))
		hex.EncodeToString(bb)
	}
} // BenchmarkEncodeSecret-4   	 5881406	       194 ns/op	     192 B/op	       3 allocs/op

func BenchmarkOBA(b *testing.B) {
	var u uuid.UUID
	u, _ = uuid.NewV4()
	_na := sha512.Sum512(u.Bytes())
	u, _ = uuid.NewV4()
	_nb := sha512.Sum512(u.Bytes())

	na := hex.EncodeToString(_na[:])
	nb := hex.EncodeToString(_nb[:])

	email := "foo@bar.com"
	pass := "1234"

	email = meta.XOR(email, meta.Reverse(na))
	pass = meta.XOR(pass, meta.XOR(na, meta.Reverse(nb)))

	for n := 0; n < b.N; n++ {
		email = meta.XOR(email, meta.Reverse(na))
		pass = meta.XOR(pass, meta.XOR(na, meta.Reverse(nb)))
	}

} // BenchmarkOBA-4   	   83332	     14496 ns/op	    1592 B/op	      12 allocs/op

func BenchmarkStringContains(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		strings.Contains("2q43o46u34pi 58303 4foo86bar l35325l4j", "foo86bar")
	}
} // BenchmarkStringContains-4   	100000000	        10.7 ns/op	       0 B/op	       0 allocs/op

func BenchmarkRandAlphaNum(b *testing.B) {
	len := 32 // 92ns @ 16, 160ns @ 32
	//fmt.Println(meta.RandAlphaNum(len)) // gy2RCNho42n21Wh7
	for i := 0; i < b.N; i++ {
		meta.RandAlphaNum(len)
	}
} // BenchmarkRandAlphaNum-4         12554822                91.4 ns/op            16 B/op          1 allocs/op

func BenchmarkXORbytes(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	u, _ = uuid.NewV4()
	y := u.Bytes()
	for n := 0; n < b.N; n++ {
		bb, _ := meta.XORbytes(x, y)
		hex.EncodeToString(bb)
	}
} //BenchmarkXORbytes-4   	18176034	        63.4 ns/op	      48 B/op	       1 allocs/op
// BenchmarkXORbytes-4   	 5940584	       219 ns/op	     208 B/op	       3 allocs/op

func BenchmarkXORbytesOnStrings(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	u, _ = uuid.NewV4()
	y := u.Bytes()
	bb, _ := meta.XORbytes(x, y)
	fmt.Println(hex.EncodeToString(bb))

	for n := 0; n < b.N; n++ {
		//meta.XORbytes([]byte(x), []byte(y))
		bb, _ := meta.XORbytes([]byte(x), []byte(y))
		hex.EncodeToString(bb)
	}
} // BenchmarkXORbytesOnStrings-4   	 3680936	       313 ns/op	     304 B/op	       5 allocs/op

func BenchmarkXOR(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.String()
	u, _ = uuid.NewV4()
	y := u.String()
	for n := 0; n < b.N; n++ {
		meta.XOR(x, y)
	}
} // go test -bench=XOR -count=1 ./kit
// BenchmarkXOR-4   	 1300089	       919 ns/op	     120 B/op	       4 allocs/op

func BenchmarkXORstrings(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.String()
	u, _ = uuid.NewV4()
	y := u.String()
	for n := 0; n < b.N; n++ {
		meta.XORstrings(x, y)
	}
} //BenchmarkXORstrings-4   	  461542	      2383 ns/op	     864 B/op	      35 allocs/op

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
