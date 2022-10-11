package lzstring

import (
	b64 "encoding/base64"
	"unicode/utf16"
)

// write bits to []byte
type bitWriter struct {
	Data []byte

	offset int // current offset of Data
	bitPos int // number of written bits in Data[offset]. 0 is no data written
}

func newBitWriter() *bitWriter {
	return &bitWriter{Data: make([]byte, 1)}
}

// write low bitLength bits of value to the writer
func (w *bitWriter) writeBits(value uint64, bitLength int) {
	if bitLength == 0 {
		return
	}
	value <<= 64 - bitLength - w.bitPos // align data to MSB + w.bitPos
	for bitLength > 0 {
		k := 8 - w.bitPos // number of bits that currently available to write
		if bitLength < k {
			k = bitLength
		}
		w.Data[w.offset] = w.Data[w.offset] | byte(value>>56)
		value <<= 8
		bitLength -= k
		w.bitPos += k
		if w.bitPos == 8 {
			w.Data = append(w.Data, 0)
			w.offset++
			w.bitPos = 0
		}
	}
}

// write bitLength bits, in reversed order
func writeBitsInv(w *bitWriter, value uint64, bitLength int) {
	var v uint64
	for i := 0; i < bitLength; i++ {
		v = v<<1 | value&1
		value >>= 1
	}
	w.writeBits(v, bitLength)
}

// compress a string as lzstring, then encode with base64.
func CompressToBase64(s string) (base64 string) {
	return b64.StdEncoding.EncodeToString(Compress(s))
}

// compress a string as raw lzstring bytes.
func Compress(s string) (compressed []byte) {
	return compressUTF16(utf16.Encode([]rune(s)))
}

func compressUTF16(u16string []uint16) (compressed []byte) {

	dict := make(map[string]int)
	dictToCreate := make(map[string]bool)
	makeKey := func(u []uint16) string {
		return string(utf16.Decode(u))
	}

	dictSize := 3
	indexBits := 2
	enlargeIn := 1 << (indexBits - 1)

	bw := newBitWriter()

	var word []uint16

	for _, c := range u16string {
		charKey := makeKey([]uint16{c})
		if _, ok := dict[charKey]; !ok {
			dict[charKey] = dictSize
			dictSize++
			dictToCreate[charKey] = true
		}

		word_c := append(word, c)
		wcKey := makeKey(word_c)
		if _, ok := dict[wcKey]; ok {
			word = word_c
		} else {
			wKey := makeKey(word)
			if _, ok := dictToCreate[wKey]; ok {
				// a new dictionary
				if len(word) > 0 && word[0] < 256 {
					// 8-bit ascii character
					writeBitsInv(bw, 0, indexBits)
					writeBitsInv(bw, uint64(word[0]), 8)
				} else {
					// UTF16 charcter
					writeBitsInv(bw, 1, indexBits)
					writeBitsInv(bw, uint64(word[0]), 16)
				}
				enlargeIn--
				if enlargeIn == 0 {
					enlargeIn = 1 << indexBits //contextEnLargeIn = int(math.Pow(2, float64(contextNumBits)))
					indexBits++
				}
				delete(dictToCreate, wKey)
			} else {
				// write existing dictionary index
				writeBitsInv(bw, uint64(dict[wKey]), indexBits)
			}
			enlargeIn--
			if enlargeIn == 0 {
				enlargeIn = 1 << indexBits //contextEnLargeIn = int(math.Pow(2, float64(contextNumBits)))
				indexBits++
			}
			dict[wcKey] = dictSize
			dictSize++
			word = []uint16{c}
		}
	}

	// write the last word
	if len(word) != 0 {
		wKey := makeKey(word)
		if _, ok := dictToCreate[wKey]; ok {
			if word[0] < 256 {
				// write a 8-bit ascii character
				writeBitsInv(bw, 0, indexBits)
				writeBitsInv(bw, uint64(word[0]), 8)
			} else {
				// write a UTF16 charcter
				writeBitsInv(bw, 1, indexBits)
				writeBitsInv(bw, uint64(word[0]), 16)
			}
			enlargeIn--
			if enlargeIn == 0 {
				enlargeIn = 1 << indexBits //contextEnLargeIn = int(math.Pow(2, float64(contextNumBits)))
				indexBits++
			}
			delete(dictToCreate, wKey)
		} else {
			// write existing dictionary index
			writeBitsInv(bw, uint64(dict[wKey]), indexBits)
		}
		enlargeIn--
		if enlargeIn == 0 {
			//enlargeIn = 1 << indexBits //contextEnLargeIn = int(math.Pow(2, float64(contextNumBits)))
			indexBits++
		}
	}

	writeBitsInv(bw, 2, indexBits)
	return bw.Data
}
