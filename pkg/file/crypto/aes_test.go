package crypto

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
	"fmt"
	"log"
)

func ExampleAesEncrypt_Encrypt() {
	aesEnc := AesEncrypt{"1234334"}
	arrEncrypt, err := aesEnc.Encrypt([]byte("abcdef"))
	if err != nil {
		log.Println(arrEncrypt)
		return
	}
	strMsg, err := aesEnc.Decrypt(arrEncrypt)
	if err != nil {
		log.Println(arrEncrypt)
		return
	}
	fmt.Println(string(strMsg))

	// Output: abcdef
}

func ExampleAesEncrypt_Decrypt() {
	aesEnc := AesEncrypt{"1234334"}
	arrEncrypt, err := aesEnc.Encrypt([]byte("abcdef"))
	if err != nil {
		log.Println(arrEncrypt)
		return
	}
	aesDec := AesEncrypt{"1234335"}
	strMsg, err := aesDec.Decrypt(arrEncrypt)
	if err != nil {
		fmt.Printf("error password")
		return
	}
	fmt.Println(string(strMsg))

	// Output: error password
}
