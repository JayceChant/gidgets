package base64

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

var strCases = []string{
	"test",
	"test string",
	"a longer sentence.",
	`1234567890-=qwertyuiop[]\asdfghjkl;'zxcvbnm,./`,
	"6akCZrAL9jcdzLw85Rq0",
	"Sqwhh4DgoKJCQh5pm8rR",
	`p8G5ZHz2tyOwNxcuyMCO
	Hj0fGWegp8nnYWkgS6uc
	S41pxEOI0HkaHqwFvnsg
	mPvpwv2zd4cTJAHfBZNa
	VCQ8910chMsBR0q0EYVW
	Oz8jY675qhuKQBS1dGEX
	HAcN5KaDzV9pXGxMmDGB
	0uCO71VFIQTafVJTUf28
	y814ZyUU1B8BUHS2aJEU
	Ooerhf85QF19z7UXCYrN`,
}

func TestEncodeAndDecode(t *testing.T) {
	for _, s := range strCases {
		t.Run(s, func(t *testing.T) {
			buf, err := Decode(Encode([]byte(s)))
			if err != nil {
				t.Error(err)
				return
			}
			if got := string(buf); got != s {
				t.Errorf("after encode and decode got %v, want %v", got, s)
			}
		})
	}
}

func TestEncodeAndDecodeNoPadding(t *testing.T) {
	for _, s := range strCases {
		t.Run(s, func(t *testing.T) {
			buf, err := DecodeNoPadding(EncodeNoPadding([]byte(s)))
			if err != nil {
				t.Error(err)
				return
			}
			if got := string(buf); got != s {
				t.Errorf("after encode and decode (no padding) got %v, want %v", got, s)
			}
		})
	}
}

var sortCases = []string{}

func genSortCases() {
	if len(sortCases) > 0 {
		return
	}
	for i := 0; i < 1000; i++ {
		sortCases = append(sortCases, fmt.Sprintf("%03d", i))
	}
}

func TestSortingConsistencyCasesStrength(t *testing.T) {
	genSortCases()
	sortBefore := make([]string, len(sortCases))
	copy(sortBefore, sortCases)
	sort.Strings(sortBefore)

	sortAfter := make([]string, len(sortCases))
	for i := range sortCases {
		sortAfter[i] = base64.URLEncoding.EncodeToString([]byte(sortCases[i]))
	}
	sort.Strings(sortAfter)
	for i := range sortAfter {
		buf, err := base64.URLEncoding.DecodeString(sortAfter[i])
		if err != nil {
			t.Error(err)
			return
		}
		sortAfter[i] = string(buf)
	}
	if reflect.DeepEqual(sortBefore, sortAfter) {
		t.Error("sorting consistency test cases are too week")
	}
}

func TestSortingConsistency(t *testing.T) {
	genSortCases()
	sortBefore := make([]string, len(sortCases))
	copy(sortBefore, sortCases)
	sort.Strings(sortBefore)

	sortAfter := make([]string, len(sortCases))
	for i := range sortCases {
		sortAfter[i] = Encode([]byte(sortCases[i]))
	}
	sort.Strings(sortAfter)
	for i := range sortAfter {
		buf, err := Decode(sortAfter[i])
		if err != nil {
			t.Error(err)
			return
		}
		sortAfter[i] = string(buf)
	}
	if !reflect.DeepEqual(sortBefore, sortAfter) {
		t.Errorf("sorting consistency is break after encoded, got %v, want %v", sortAfter, sortBefore)
	}
}

func TestSortingConsistencyNoPadding(t *testing.T) {
	genSortCases()
	sortBefore := make([]string, len(sortCases))
	copy(sortBefore, sortCases)
	sort.Strings(sortBefore)

	sortAfter := make([]string, len(sortCases))
	for i := range sortCases {
		sortAfter[i] = EncodeNoPadding([]byte(sortCases[i]))
	}
	sort.Strings(sortAfter)
	for i := range sortAfter {
		buf, err := DecodeNoPadding(sortAfter[i])
		if err != nil {
			t.Error(err)
			return
		}
		sortAfter[i] = string(buf)
	}
	if !reflect.DeepEqual(sortBefore, sortAfter) {
		t.Errorf("sorting consistency is break after encoded, got %v, want %v", sortAfter, sortBefore)
	}
}

var intCases = []uint64{
	0,
	1,
	2,
	1234567890,
	42,
	1<<64 - 1,
}

func TestUintAndBytesConversion(t *testing.T) {
	for e := 0; e < 2; e++ {
		littleEndian := e == 1
		for _, i := range intCases {
			t.Run(strconv.FormatUint(i, 10), func(t *testing.T) {
				if got := bytesToUint(uintToBytes(i, littleEndian), littleEndian); got != i {
					t.Errorf("convert to bytes and convert back (littleEndian = %v), got %v, want %v", littleEndian, got, i)
				}
			})
		}
	}
}

func TestEncodeAndDecodeUint(t *testing.T) {
	for e := 0; e < 2; e++ {
		littleEndian := e == 1
		for _, i := range intCases {
			t.Run(strconv.FormatUint(i, 10), func(t *testing.T) {
				got, err := DecodeUint(EncodeUint(i, littleEndian), littleEndian)
				if err != nil {
					t.Error(err)
					return
				}
				if got != i {
					t.Errorf("after encode and decode uint (littleEndian = %v) got %v, want %v", littleEndian, got, i)
				}
			})
		}
	}
}

func TestEncodeAndDecodeUintNoPadding(t *testing.T) {
	for e := 0; e < 2; e++ {
		littleEndian := e == 1
		for _, i := range intCases {
			t.Run(strconv.FormatUint(i, 10), func(t *testing.T) {
				got, err := DecodeUintNoPadding(EncodeUintNoPadding(i, littleEndian), littleEndian)
				if err != nil {
					t.Error(err)
					return
				}
				if got != i {
					t.Errorf("after encode and decode uint (no padding, littleEndian = %v) got %v, want %v", littleEndian, got, i)
				}
			})
		}
	}
}
