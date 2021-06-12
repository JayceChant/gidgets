// Here's a wrapper for official encoding/base64,
// adding sort consistency and other useful features.
package base64

import (
	"encoding/base64"
	"encoding/binary"
)

const (
	// URL-safe encode chars ordered by ASCII.
	// With this set of encoding-characters, the character range is the same as the default Base64URL encoding,
	// except it's sorting consistent.
	// In other words, the ordering is the same before and after encoding.
	urlSafeAsc = "-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
)

var (
	stdPadding = base64.NewEncoding(urlSafeAsc)
	noPadding  = base64.NewEncoding(urlSafeAsc).WithPadding(base64.NoPadding)
)

// Encode encodes the source bytes into base64 string,
// with URL safe characters in ascending order.
func Encode(src []byte) string {
	return stdPadding.EncodeToString(src)
}

// Decode decodes the bytes represented by the string,
// is the inverse of Encode.
func Decode(s string) ([]byte, error) {
	return stdPadding.DecodeString(s)
}

// EncodeNoPadding encodes the source bytes into base64 string,
// with URL safe characters in ascending order and no padding.
func EncodeNoPadding(src []byte) string {
	return noPadding.EncodeToString(src)
}

// DecodeNoPadding decodes the bytes represented by the string,
// is the inverse of EncodeNoPadding.
func DecodeNoPadding(s string) ([]byte, error) {
	return noPadding.DecodeString(s)
}

// EncodeUint encodes the unsigned int into base64 string.
func EncodeUint(n uint64, littleEndian bool) string {
	return stdPadding.EncodeToString(uintToBytes(n, littleEndian))
}

// DecodeUint return the unsigned int representing by the string,
// is the inverse of EncodeUint.
func DecodeUint(s string, littleEndian bool) (uint64, error) {
	buf, err := stdPadding.DecodeString(s)
	return bytesToUint(buf, littleEndian), err
}

// EncodeUintNoPadding encodes the unsigned int into base64 string, no padding.
func EncodeUintNoPadding(n uint64, littleEndian bool) string {
	return noPadding.EncodeToString(uintToBytes(n, littleEndian))
}

// DecodeUintNoPadding return the unsigned int representing by the string,
// is the inverse of EncodeUintNoPadding.
func DecodeUintNoPadding(s string, littleEndian bool) (uint64, error) {
	buf, err := noPadding.DecodeString(s)
	return bytesToUint(buf, littleEndian), err
}

func uintToBytes(n uint64, littleEndian bool) []byte {
	buf := make([]byte, 8)
	if littleEndian {
		binary.LittleEndian.PutUint64(buf, n)
	} else {
		binary.BigEndian.PutUint64(buf, n)
	}
	return buf
}

func bytesToUint(b []byte, littleEndian bool) uint64 {
	if littleEndian {
		return binary.LittleEndian.Uint64(b)
	}
	return binary.BigEndian.Uint64(b)
}
