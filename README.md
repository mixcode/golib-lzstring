
# golib-lzstring

This package provides encoding and decoding of lzstring data.

lzstring is a custom javascript string compression format. See following URL for
details on lzstring.

https://pieroxy.net/blog/pages/lz-string/index.html

# Example

```go
	text := "I have some string to be compressed"
	fmt.Println(lzstring.CompressToBase64(text))
	// Output: JIAgFghgbgpiDOB7AtneAXATgSwHYHMR1EQAjOAYxQAdMZ54YATIAA==


	base64text := "JIAgFghgbgpiDOB7AtneAXATgSwHYHMR1EQAjOAYxQAdMZ54YATIAA=="
	decompressed, err := lzstring.DecompressBase64(base64text)
	if err != nil {
		panic(err)
	}
	fmt.Println(decompressed)
	// Output: I have some string to be compressed
```

