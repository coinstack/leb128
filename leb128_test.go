// Written in 2014 by Sheran Gunasekera
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leb128

import (
	"bytes"
	// "encoding/hex"
	"testing"
)

var testEncode = map[uint64][]byte{
	0:     []byte{0x00},
	1:     []byte{0x01},
	2:     []byte{0x02},
	127:   []byte{0x7f},
	128:   []byte{0x80, 0x01},
	129:   []byte{0x81, 0x01},
	130:   []byte{0x82, 0x01},
	12857: []byte{0xB9, 0x64},
	16256: []byte{0x80, 0x7f},
}

var testDecode = map[int64][]byte{
	0:    []byte{0x00},
	1:    []byte{0x01},
	2:    []byte{0x02},
	127:  []byte{0xFF, 0x00},
	128:  []byte{0x80, 0x01},
	129:  []byte{0x81, 0x01},
	-1:   []byte{0x7f},
	-2:   []byte{0x7e},
	-127: []byte{0x81, 0x7f},
	-128: []byte{0x80, 0x7f},
	-129: []byte{0xFF, 0x7e},
}

var testEncodeLength = map[uint64]int{
	0:     1,
	1:     1,
	2:     1,
	127:   1,
	128:   2,
	129:   2,
	130:   2,
	12857: 2,
	16256: 2,
}

func TestEncodeDecode(t *testing.T) {
	var before1 uint64 = 922337203685477600
	var before2 uint64 = 23123214212

	temp := []byte{}
	temp = append(temp, EncodeULeb128(before1)...)
	temp = append(temp, EncodeULeb128(before2)...)

	tempBuffer := bytes.NewBuffer(temp)
	after1, _ := ReadULeb128(tempBuffer)
	after2, _ := ReadULeb128(tempBuffer)

	if before1 != after1 {
		t.Errorf("Wanted %v, got %v", before1, after1)
	}

	if before2 != after2 {
		t.Errorf("Wanted %v, got %v", before2, after2)
	}

	_, eof := ReadULeb128(tempBuffer)
	if eof == nil {
		t.Error("unexpected EOF not detected")
	}
}

func TestDecodeULeb128(t *testing.T) {
	for k, v := range testEncode {
		res := DecodeULeb128(v)
		if res != k {
			t.Errorf("Wanted %d, got %d", k, res)
		}
	}
}

// func TestDecodeSLeb128(t *testing.T) {
// 	for k, v := range testDecode {
// 		res := DecodeSLeb128(v)
// 		if res != k {
// 			t.Errorf("Wanted %d, got %d", k, res)
// 		}
// 	}
// }

func TestEnecodeULeb128(t *testing.T) {
	for k, v := range testEncode {
		res := EncodeULeb128(k)
		if bytes.Compare(res, v) != 0 {
			t.Errorf("Wanted %d, got %d", v, res)
		}
	}
}

// func TestEnecodeSLeb128(t *testing.T) {
// 	for k, v := range testDecode {
// 		res := EncodeSLeb128(k)
// 		if bytes.Compare(res, v) != 0 {
// 			t.Errorf("Wanted %d, got %d", v, res)
// 		}
// 	}
// }
