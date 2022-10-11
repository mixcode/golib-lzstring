package lzstring_test

import (
	"fmt"

	lzstring "github.com/mixcode/golib-lzstring"
)

func ExampleCompressToBase64() {
	text := "I have some string to be compressed"
	fmt.Println(lzstring.CompressToBase64(text))
	// Output: JIAgFghgbgpiDOB7AtneAXATgSwHYHMR1EQAjOAYxQAdMZ54YATIAA==
}

func ExampleDecompressBase64() {
	base64text := "JIAgFghgbgpiDOB7AtneAXATgSwHYHMR1EQAjOAYxQAdMZ54YATIAA=="
	decompressed, err := lzstring.DecompressBase64(base64text)
	if err != nil {
		panic(err)
	}
	fmt.Println(decompressed)
	// Output: I have some string to be compressed
}
