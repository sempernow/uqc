package id_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/sempernow/uqc/kit/id"

	"github.com/gofrs/uuid"
	"lukechampine.com/blake3"
)

func TestSumSHA1(t *testing.T) {
	t.Skip()
	t.Log(id.SumSHA1([]byte("1594562637"))) //c2db288bfaeaf238a793d7a7f9c51ab50b2b9196
	t.Log(id.SumSHA1("1594562637"))         //c2db288bfaeaf238a793d7a7f9c51ab50b2b9196
	t.Log(id.SumSHA1(1594562637))           //c2db288bfaeaf238a793d7a7f9c51ab50b2b9196
	t.Log(id.SumSHA1("☩"))                  //ec7b5357c265c5832ad34aaa6eb4305b3bd83bc5
	t.Log(id.SumSHA1('☩'))                  //0000000000000000000000000000000000000000

	f, err := os.Open("kit_test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	t.Log(id.SumSHA1(f)) //0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33
}

func TestULID(t *testing.T) {
	t.Skip()
	id.DoULID()
	//t.Fatal()
}

func BenchmarkULID(b *testing.B) {
	// Run it b.N times
	//fmt.Println(id.NowULID().String())
	for n := 0; n < b.N; n++ {
		id.NowULID() //.String()
	}
} // BenchmarkULID-4   	   63837	     18485 ns/op	    9696 B/op	       6 allocs/p

func BenchmarkNewID(b *testing.B) {
	// Run it b.N times
	//fmt.Println(id.NewID()) // wbj36uter1d8i8wazuwjwqaevm
	for n := 0; n < b.N; n++ {
		id.NewID() //.String()
	}
} // BenchmarkNewID-4   	  857222	      1199 ns/op	    1648 B/op	       7 allocs/op

func BenchmarkSHA256(b *testing.B) {
	// Run it b.N times
	//s := sha256.Sum256([]byte("foobar"))
	for n := 0; n < b.N; n++ {
		id.SumSHA256([]byte("foobar"))
		// if s == sha256.Sum256([]byte("foobar")) {
		// 	//true
		// } else {
		// 	//false
		// }
	}
} // BenchmarkSHA256-4   	 2368640	       500 ns/op	     168 B/op	       4 allocs/op
//   BenchmarkSHA256-4       4286598           271 ns/op           0 B/op          0 allocs/op @ sha256 + test

func BenchmarkBlake3256(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	for n := 0; n < b.N; n++ {
		blake3.Sum256(x)
	}
} // BenchmarkBlake3256-4     7118822        162 ns/op               0 B/op          0 allocs/op
func BenchmarkSumBlake3_256(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.String()
	for n := 0; n < b.N; n++ {
		//id.SumBlake3_256(bb)
		id.SumBlake3_256(x)
	}
} // BenchmarkSumBlake3_256-4    3558349	       337 ns/op	     160 B/op	       3 allocs/op  @ bb
// BenchmarkSumBlake3_256-4      2787312	       387 ns/op	     208 B/op	       4 allocs/op  @ s

func BenchmarkSHA512(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	for n := 0; n < b.N; n++ {
		sha512.Sum512(x)
	}
} // BenchmarkSHA512-4   	 3251862	       364 ns/op	       0 B/op	       0 allocs/op

// https://pkg.go.dev/lukechampine.com/blake3
func BenchmarkBlake3512(b *testing.B) {
	// Run it b.N times
	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	for n := 0; n < b.N; n++ {
		//blake3.Sum512(bb)
		b := blake3.Sum512(x)
		hex.EncodeToString(b[:])
		//id.SumBlake3_512(bb)
	}
} // BenchmarkBlake3512-4   7205895          164 ns/op               0 B/op          0 allocs/op
//   BenchmarkBlake3512-4   2819072          417 ns/op             320 B/op          3 allocs/op @ string

func BenchmarkNewUUID(b *testing.B) {
	// Run it b.N times
	for n := 0; n < b.N; n++ {
		uuid.NewV4() //.String()
	}
} // BenchmarkNewUUID-4   	 3045514	       393 ns/op	      64 B/op	       2 allocs/o

func BenchmarkHMAC(b *testing.B) {
	// Run it b.N times

	var u uuid.UUID
	u, _ = uuid.NewV4()
	x := u.Bytes()
	u, _ = uuid.NewV4()
	y := u.Bytes()

	rid := fmt.Sprintf("%x", sha256.Sum256(x))
	rKey := fmt.Sprintf("%x", sha256.Sum256(y))
	for n := 0; n < b.N; n++ {
		hmac.New(sha256.New, []byte(rKey+rid))
	}
} // BenchmarkHMAC-4   	  399993	      2914 ns/op	    1152 B/op	      10 allocs/op
