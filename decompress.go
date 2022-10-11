// This package provides encoding and decoding of lzstring data.
//
// lzstring is a custom javascript string compression format.
// See following URL for details on lzstring.
//
// https://pieroxy.net/blog/pages/lz-string/index.html
package lzstring

import (
	"encoding/base64"
	"errors"
	"io"
	"unicode/utf16"
)

var (
	ErrNotDecodable = errors.New("input not decodable")
	ErrEmptyInput   = errors.New("empty input")
)

// decompress raw lzstring bytes.
func DecompressBytes(data []byte) (decompressed string, err error) {
	br := &byteBitReader{Bytes: data}
	res, err := decompress(br)
	if err != nil {
		return
	}
	return string(utf16.Decode(res)), nil
}

// decompress base64 encoded lzstring.
func DecompressBase64(src string) (decompressed string, err error) {
	if src == "" {
		err = ErrEmptyInput
		return
	}
	data, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return
	}
	return DecompressBytes(data)
}

// a bit reader. returns n bits from the stream.
type bitReader interface {
	getBits(n int) (uint64, error) // Read first n bits. The read bits must be right-aligned. i.e. the last read bit should be the LSB.
}

// read the bits and inverse the bit order
func getBitsInv(br bitReader, n int) (data uint64, err error) {
	c, err := br.getBits(n)
	if err != nil {
		return
	}
	for i := 0; i < n; i++ {
		data = data<<1 | c&1
		c >>= 1
	}
	return
}

// []byte based bit reader
type byteBitReader struct {
	Bytes []byte // source

	currentByte byte // current reading byte
	offset      int  // offset of next byte
	bitsLeft    int  // bits left in the buffer, from the LSB. if bitpos==0, then read a new byte
}

func (r *byteBitReader) getBits(n int) (data uint64, err error) {
	if r.offset == len(r.Bytes) && r.bitsLeft == 0 {
		err = io.EOF
		return
	}
	for n > 0 {
		if r.bitsLeft == 0 {
			if r.offset >= len(r.Bytes) {
				err = io.ErrUnexpectedEOF
				return
			}
			r.currentByte = r.Bytes[r.offset]
			r.offset++
			r.bitsLeft = 8
		}
		k := n
		if k > r.bitsLeft {
			k = r.bitsLeft
		}
		data = (data << k) | uint64(r.currentByte>>(8-k))
		r.currentByte <<= k
		r.bitsLeft -= k
		n -= k
	}
	return
}

// the main lzstring decomporessor
func decompress(br bitReader) (decompressed []uint16, err error) {

	// the dictionary
	dict := make([][]uint16, 0)
	dictSize := 0               // number of entires in the dictionary
	indexBits := 2              // the current bit size of dictionary index
	enlargeIn := 1 << indexBits // reverse counter to grow the dictionary index bit size when dictionary gets large

	addWordToDictionary := func(word []uint16) { // add a new word to the dictionary
		wc := make([]uint16, len(word))
		copy(wc, word) // copy data to prevent appending new characters to existing dictionary
		dict = append(dict, wc)
		dictSize++
		enlargeIn--
		if enlargeIn == 0 {
			// directory size gets big, so increase the index bit size
			enlargeIn = 1 << indexBits
			indexBits++
		}
	}

	// prepare the dictionary for index 0, 1, 2
	for i := 0; i < 3; i++ {
		addWordToDictionary([]uint16{uint16(i)})
	}

	// read the first rune
	bits, err := getBitsInv(br, indexBits)
	if err != nil {
		return
	}
	switch bits {
	case 0:
		bits, err = getBitsInv(br, 8)
		if err != nil {
			return
		}
	case 1:
		bits, err = getBitsInv(br, 16)
		if err != nil {
			return
		}
	case 2:
		// empty string
		return nil, nil
	}
	c := uint16(bits)
	word := []uint16{c} // last processed character array
	addWordToDictionary(word)

	output := make([]uint16, 0) // decoded output in UTF16 data
	output = append(output, c)

	for {
		bits, err = getBitsInv(br, indexBits)
		if err != nil {
			return
		}
		n := int(bits)
		switch n {
		case 0:
			// read a new 8-bit ASCII char and append it to the dictionary
			bits, err = getBitsInv(br, 8)
			if err != nil {
				return
			}
			n = dictSize // save the last dictionary index for later use
			addWordToDictionary([]uint16{uint16(bits)})
		case 1:
			// read a new UTF16 char and append it to the dictionary
			bits, err = getBitsInv(br, 16)
			if err != nil {
				return
			}
			n = dictSize // save the last dictionary index for later use
			addWordToDictionary([]uint16{uint16(bits)})
		case 2:
			// end of stream
			return output, nil
		}

		// append the found entry to the output
		var entry []uint16
		if n < dictSize {
			// entry found in the dictionary
			entry = dict[n]
		} else if n == dictSize {
			// make a new dictionary with the working memory
			entry = append(word, word[0])
		} else {
			return nil, ErrNotDecodable
		}
		output = append(output, entry...)

		// add the last word to the dictionary
		addWordToDictionary(append(word, entry[0]))
		word = entry
	}
}
