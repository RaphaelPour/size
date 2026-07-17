// Package size provides a Size type for representing data sizes in bytes,
// analogous to how the standard library's time.Duration represents durations.
//
// Size supports both base-2 (IEC) units such as Kibibyte, Mebibyte, and
// Gibibyte, and base-10 (SI) units such as Kilobyte, Megabyte, and
// Gigabyte. Short aliases like Kib, Mib, Gib and Kb, Mb, Gb are also
// provided for convenience.
//
// Because the underlying type is uint64, values above roughly 16 Exbibytes
// (2^64 bytes) cannot be represented; Zebibyte and larger units are
// therefore intentionally omitted.
//
// Basic usage:
//
//	s := 5 * size.Mebibyte
//	fmt.Println(s.Bytes()) // 5242880
package size

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Size represents a quantity of bytes. The zero value represents zero bytes.
type Size uint64

// Base-2 (IEC) units of size. Sizes are by a magnitude of 1024 apart.
const (
	Byte Size = 1 << (10 * iota)
	Kibibyte
	Mebibyte
	Gibibyte
	Tebibyte
	Pebibyte
	Exbibyte

	B   = Byte
	KiB = Kibibyte
	MiB = Mebibyte
	GiB = Gibibyte
	TiB = Tebibyte
	PiB = Pebibyte

	// Exbibytes are only supported til 2^64
	EiB = Exbibyte

	// zebibyte and higher is not supported,
	// zebibite is 2^70 and larger as uint64 2^64
)

// Base-10 (SI) units of size. Sizes are by a magnitude of 1000 apart.
const (
	Kilobyte Size = 1000
	Megabyte Size = Kilobyte * 1000
	Gigabyte Size = Megabyte * 1000
	Terabyte Size = Gigabyte * 1000
	Petabyte Size = Terabyte * 1000
	Exabyte  Size = Petabyte * 1000

	KB = Kilobyte
	MB = Megabyte
	GB = Gigabyte
	TB = Terabyte
	PB = Petabyte

	// Exabyte is limited by uint64
	EB = Exabyte

	// zebibyte and above are not supported limited by uint64
)

type unit struct {
	base   Size
	suffix string
}

type Unit uint8

func (u Unit) String() string {
	return u.internalUnit().suffix
}

func (u Unit) internalUnit() unit {
	unt := unitByte
	if found, ok := units[u]; ok {
		unt = found
	}
	return unt
}

const (
	UnitByte Unit = iota
	UnitKiB
	UnitMiB
	UnitGiB
	UnitTiB
	UnitPiB
	UnitEiB

	UnitKB
	UnitMB
	UnitGB
	UnitTB
	UnitPB
	UnitEB
	unitMax
)

var (
	unitByte = unit{base: B, suffix: "B"}

	unitKiB = unit{base: KiB, suffix: "KiB"}
	unitMiB = unit{base: MiB, suffix: "MiB"}
	unitGiB = unit{base: GiB, suffix: "GiB"}
	unitTiB = unit{base: TiB, suffix: "TiB"}
	unitPiB = unit{base: PiB, suffix: "PiB"}
	unitEiB = unit{base: EiB, suffix: "EiB"}

	unitsIEC = []unit{
		unitEiB,
		unitPiB,
		unitTiB,
		unitGiB,
		unitMiB,
		unitKiB,
		unitByte,
	}

	unitKB = unit{base: KB, suffix: "kB"}
	unitMB = unit{base: MB, suffix: "MB"}
	unitGB = unit{base: GB, suffix: "GB"}
	unitTB = unit{base: TB, suffix: "TB"}
	unitPB = unit{base: PB, suffix: "PB"}
	unitEB = unit{base: EB, suffix: "EB"}

	unitsSI = []unit{
		unitEB,
		unitPB,
		unitTB,
		unitGB,
		unitMB,
		unitKB,
		unitByte,
	}

	units = map[Unit]unit{
		UnitEiB:  unitEiB,
		UnitPiB:  unitPiB,
		UnitTiB:  unitTiB,
		UnitGiB:  unitGiB,
		UnitMiB:  unitMiB,
		UnitKiB:  unitKiB,
		UnitByte: unitByte,

		UnitEB: unitEB,
		UnitPB: unitPB,
		UnitTB: unitTB,
		UnitGB: unitGB,
		UnitMB: unitMB,
		UnitKB: unitKB,
	}
)

var (
	ErrUnknownUnit = errors.New("unknown unit")
)

// Bytes returns the byte size as uint64.
func (s Size) Bytes() uint64 {
	return uint64(s)
}

// String returns a formated size string with unit suffix using opinionated format options:
// - IEC (2-based)
// - cut empty fractions (42.00 -> 42, 12.30 -> 12.3, 3.14 -> 3.14)
// - prevision: two decimals (max)
// Use Format, FormatIEC and FormatSI for full customization.
func (s Size) String() string {
	// format with opinionated settings:
	// - IEC (2-based)
	// - cut empty fraction
	// - prevision: two decimals (max)
	return s.format(s.fit(unitsIEC), WithCutEmptyFraction(), WithPrecision(2))
}

func (s Size) fit(units []unit) unit {
	for _, uc := range units {
		if s >= uc.base {
			return uc
		}
	}
	return unitByte
}

func (s Size) format(u unit, opts ...FormatOption) string {
	// use default format options as starting point
	options := defaultFormatOptions

	for _, opt := range opts {
		opt(&options)
	}
	str := strconv.FormatFloat(float64(s)/float64(u.base), 'f', options.precision, 64)

	if options.cutEmptyFraction && strings.Contains(str, ".") {
		// Trims: convert 42.00 to 42
		str = strings.TrimRight(str, "0")
		str = strings.TrimRight(str, ".")
	}
	// append size suffix
	return str + u.suffix
}

var (
	defaultFormatOptions = formatOptions{
		cutEmptyFraction: false,
		precision:        2,
	}
)

type formatOptions struct {
	cutEmptyFraction bool
	precision        int
}

type FormatOption func(fOpt *formatOptions)

func WithCutEmptyFraction() FormatOption {
	return func(fOpt *formatOptions) {
		fOpt.cutEmptyFraction = true
	}
}

func WithPrecision(precision int) FormatOption {
	return func(fOpt *formatOptions) {
		fOpt.precision = precision
	}
}

func (s Size) Format(u Unit, opts ...FormatOption) string {
	return s.format(u.internalUnit(), opts...)
}

func (s Size) FormatIEC(opts ...FormatOption) string {
	return s.format(s.fit(unitsIEC), opts...)
}

func (s Size) FormatSI(opts ...FormatOption) string {
	return s.format(s.fit(unitsSI), opts...)
}

// MarshalText implements encoding.TextMarshaler. It emits the exact byte count
// with a "B" suffix (for example "5242880B"). Losless encoding and reversible
// via UnmarshalText
func (s Size) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(s), 10) + "B"), nil
}

// UnmarshalText implements encoding.TextUnmarshaler, parsing values such as
// "5242880B", "5MiB", "5 MB", or "0.04TiB".
func (s *Size) UnmarshalText(text []byte) error {
	raw := strings.TrimSpace(string(text))

	// longest suffix matching so MB wins over B
	var match unit
	found := false
	for _, u := range units {
		if strings.HasSuffix(raw, u.suffix) && (!found || len(u.suffix) > len(match.suffix)) {
			match = u
			found = true
		}
	}

	if !found {
		return ErrUnknownUnit
	}

	num := strings.TrimSpace(strings.TrimSuffix(raw, match.suffix))

	// assume bytes first to avoid losing precision due to float64
	if v, err := strconv.ParseUint(num, 10, 64); err == nil {
		*s = Size(v) * match.base
		return nil
	}

	val, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return fmt.Errorf("parsing size value failed: %w", err)
	}

	*s = Size(val * float64(match.base))

	return nil
}
