package tools

import (
	"strconv"
)

// Convert string to specify type.
type StrTo string

func (f StrTo) Exist() bool {
	return string(f) != string(rune(0x1E))
}

func (f StrTo) Uint8() (uint8, error) {
	v, err := strconv.ParseUint(f.String(), 10, 8)
	return uint8(v), err
}

func (f StrTo) Int() (int, error) {
	v, err := strconv.ParseInt(f.String(), 10, 0)
	return int(v), err
}

func (f StrTo) Int64() (int64, error) {
	v, err := strconv.ParseInt(f.String(), 10, 64)
	return int64(v), err
}

func (f StrTo) Float64() (float64, error) {
	v, err := strconv.ParseFloat(f.String(), 64)
	return float64(v), err
}

func (f StrTo) MustUint8() uint8 {
	v, _ := f.Uint8()
	return v
}

func (f StrTo) MustInt() int {
	v, _ := f.Int()
	return v
}

func (f StrTo) MustInt64() int64 {
	v, _ := f.Int64()
	return v
}

func (f StrTo) MustFloat64() float64 {
	v, _ := f.Float64()
	return v
}

func (f StrTo) String() string {
	if f.Exist() {
		return string(f)
	}
	return ""
}
