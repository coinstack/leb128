// Written in 2014 by Sheran Gunasekera
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package leb128 provides methods to read and write LEB128 (Little-Endian Base 128) quantities.
package leb128

import (
	"bytes"
	"errors"
	"io"
)

// These consts are here in their entirity even though Minint64 is all that is used in this package
const (
	MaxInt8   = 1<<7 - 1
	MinInt8   = -1 << 7
	MaxInt16  = 1<<15 - 1
	MinInt16  = -1 << 15
	MaxInt32  = 1<<31 - 1
	MinInt32  = -1 << 31
	MaxInt64  = 1<<63 - 1
	MinInt64  = -1 << 63
	MaxUint8  = 1<<8 - 1
	MaxUint16 = 1<<16 - 1
	Maxuint64 = 1<<32 - 1
	MaxUint64 = 1<<64 - 1
)

// EncodeULeb128 encode's an unsigned int64 value to an unsigned LEB128 value. Returns the result in a byte slice
func EncodeULeb128(value uint64) []byte {
	remaining := value >> 7
	var buf = new(bytes.Buffer)
	for remaining != 0 {
		buf.WriteByte(byte(value&0x7f | 0x80))
		value = remaining
		remaining >>= 7
	}
	buf.WriteByte(byte(value & 0x7f))
	return buf.Bytes()
}

// DecodeULeb128 decodes an unsigned LEB128 value to an unsigned int64 value. Returns the result as a uint64
func DecodeULeb128(value []byte) uint64 {
	var result uint64
	var ctr uint
	var cur byte = 0x80
	for (cur&0x80 == 0x80) && ctr < 9 {
		cur = value[ctr] & 0xff
		result += uint64((cur & 0x7f)) << (ctr * 7)
		ctr++
	}
	return result
}

// ReadULeb128 reads and decodes an unsigned LEB128 value from a ByteReader to an unsigned int64 value. Returns the result as a uint64
func ReadULeb128(reader io.ByteReader) (uint64, error) {
	var result uint64
	var ctr uint
	var cur byte = 0x80
	var err error
	for (cur&0x80 == 0x80) && ctr < 9 {
		cur, err = reader.ReadByte()
		if err == io.EOF {
			return 0, errors.New("unexpected EOF")
		}
		result += uint64((cur & 0x7f)) << (ctr * 7)
		ctr++
	}
	return result, nil
}

// EncodeSLeb128 encode a signed int64 value to a signed LEB128 value. Returns the result in a byte slice
func EncodeSLeb128(value int64) []byte {
	var buf = new(bytes.Buffer)
	remaining := value >> 7
	hasMore := true
	var end int64
	if (value & MinInt64) == 0 {
		end = 0
	} else {
		end = -1
	}
	for hasMore {
		hasMore = (remaining != end) || ((remaining & 1) != ((value >> 6) & 1))
		var t int
		if hasMore {
			t = 0x80
		} else {
			t = 0
		}
		buf.WriteByte(byte(int((value & 0x7f)) | t))
		value = remaining
		remaining >>= 7
	}
	return buf.Bytes()
}

// DecodeSLeb128 decodes a signed LEB128 value to a signed int64 value. Returns the result as a int64
func DecodeSLeb128(value []byte) int64 {
	var result int64
	var ctr uint
	var cur byte = 0x80
	var signBits int64 = -1
	for (cur&0x80 == 0x80) && ctr < 5 {
		cur = value[ctr] & 0xff
		result += int64((cur & 0x7f)) << (ctr * 7)
		signBits <<= 7
		ctr++
	}
	if ((signBits >> 1) & result) != 0 {
		result += signBits
	}
	return result
}

// ReadSLeb128 reads and decodes a signed LEB128 value from a ByteReader to a signed int64 value. Returns the result as a int64
func ReadSLeb128(reader io.ByteReader) int64 {
	var result int64
	var ctr uint
	var cur byte = 0x80
	var signBits int64 = -1
	var err error
	for (cur&0x80 == 0x80) && ctr < 5 {
		cur, err = reader.ReadByte()
		if err != nil {
			panic(err)
		}
		result += int64((cur & 0x7f)) << (ctr * 7)
		signBits <<= 7
		ctr++
	}
	if ((signBits >> 1) & result) != 0 {
		result += signBits
	}
	return result
}
