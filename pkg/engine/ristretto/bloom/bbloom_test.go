package bloom

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/bhojpur/cache/pkg/hack"
)

var (
	wordlist1 [][]byte
	n         = uint64(1 << 16)
	bf        *Bloom
)

func TestMain(m *testing.M) {
	wordlist1 = make([][]byte, n)
	for i := range wordlist1 {
		b := make([]byte, 32)
		_, _ = rand.Read(b)
		wordlist1[i] = b
	}
	fmt.Println("\n###############\nbbloom_test.go")
	fmt.Print("Benchmarks relate to 2**16 OP. --> output/65536 op/ns\n###############\n\n")

	os.Exit(m.Run())
}

func TestM_NumberOfWrongs(t *testing.T) {
	bf = NewBloomFilter(n*10, 7)

	cnt := 0
	for i := range wordlist1 {
		hash := hack.RuntimeMemhash(wordlist1[i], 0)
		if !bf.AddIfNotHas(hash) {
			cnt++
		}
	}
	fmt.Printf("Bloomfilter New(7* 2**16, 7) (-> size=%v bit): \n            Check for 'false positives': %v wrong positive 'Has' results on 2**16 entries => %v %%\n", len(bf.bitset)<<6, cnt, float64(cnt)/float64(n))

}

func BenchmarkM_New(b *testing.B) {
	for r := 0; r < b.N; r++ {
		_ = NewBloomFilter(n*10, 7)
	}
}

func BenchmarkM_Clear(b *testing.B) {
	bf = NewBloomFilter(n*10, 7)
	for i := range wordlist1 {
		hash := hack.RuntimeMemhash(wordlist1[i], 0)
		bf.Add(hash)
	}
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		bf.Clear()
	}
}

func BenchmarkM_Add(b *testing.B) {
	bf = NewBloomFilter(n*10, 7)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := range wordlist1 {
			hash := hack.RuntimeMemhash(wordlist1[i], 0)
			bf.Add(hash)
		}
	}

}

func BenchmarkM_Has(b *testing.B) {
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := range wordlist1 {
			hash := hack.RuntimeMemhash(wordlist1[i], 0)
			bf.Has(hash)
		}
	}
}
