// Package size provides a Size type for representing data sizes in bytes,
// analogous to how the standard library's [time.Duration] represents durations.
//
// Size supports both base-2 (IEC) units such as Kibibyte, Mebibyte, and
// Gibibyte, and base-10 (SI) units such as Kilobyte, Megabyte, and
// Gigabyte. Short aliases like KiB, MiB, GiB and KB, MB, GB are also
// provided for convenience.
//
// Because the underlying type is uint64, values above roughly 16 Exbibytes
// (2^64 bytes) cannot be represented; Zebibyte and larger units are
// therefore intentionally omitted.
//
// Use Bytes for the raw byte count or the per-unit float accessors (Gibibytes,
// Megabytes, ...) for a scaled value, String/Format/FormatIEC/FormatSI to render
// a size as text, Parse to read one back, and MarshalText/UnmarshalText for text
// (un)marshaling.
//
// Size also implements [flag.Value] (via Set) for use as a command-line flag and
// [slog.LogValuer] (via LogValue) for structured logging. Because it implements
// [encoding.TextMarshaler] and [encoding.TextUnmarshaler], it round-trips through
// [encoding/json] and other text-based codecs (as the string form, e.g.
// "5242880B") with no additional code.
//
// Basic usage:
//
//	s := 5 * size.Mebibyte
//	fmt.Println(s.Bytes()) // 5242880
package size

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
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

	// Zebibyte and higher are not supported: a Zebibyte is 2^70, which
	// exceeds the 2^64 range of the underlying uint64.
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

	// Zettabyte and above are not supported: they exceed the 2^64 range of
	// the underlying uint64.
)

// The not exported unit contains the actual unit information and
// is guarded by Unit
type unit struct {
	base   Size
	suffix string
}

// Unit selects an explicit unit for Size.Format. Valid values are the Unit*
// constants; an unrecognized Unit falls back to bytes.
type Unit uint8

// String returns the unit's suffix, for example "MiB" or "kB".
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

// Unit constants for use with Size.Format, covering the base-2 (IEC) and
// base-10 (SI) units.
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

// ErrUnknownUnit is returned by UnmarshalText when the input has no recognized
// unit suffix.
var ErrUnknownUnit = errors.New("unknown unit")

// Bytes returns the byte size as uint64.
func (s Size) Bytes() uint64 {
	return uint64(s)
}

// The following accessors return the size scaled to a unit as a float64,
// analogous to [time.Duration.Hours] and [time.Duration.Minutes]. They may lose
// precision for very large values; use Bytes for the exact count.

// Kibibytes returns the size in base-2 kibibytes (1024 bytes).
func (s Size) Kibibytes() float64 { return float64(s) / float64(Kibibyte) }

// Mebibytes returns the size in base-2 mebibytes (1024 kibibytes).
func (s Size) Mebibytes() float64 { return float64(s) / float64(Mebibyte) }

// Gibibytes returns the size in base-2 gibibytes (1024 mebibytes).
func (s Size) Gibibytes() float64 { return float64(s) / float64(Gibibyte) }

// Tebibytes returns the size in base-2 tebibytes (1024 gibibytes).
func (s Size) Tebibytes() float64 { return float64(s) / float64(Tebibyte) }

// Pebibytes returns the size in base-2 pebibytes (1024 tebibytes).
func (s Size) Pebibytes() float64 { return float64(s) / float64(Pebibyte) }

// Exbibytes returns the size in base-2 exbibytes (1024 pebibytes).
func (s Size) Exbibytes() float64 { return float64(s) / float64(Exbibyte) }

// Kilobytes returns the size in base-10 kilobytes (1000 bytes).
func (s Size) Kilobytes() float64 { return float64(s) / float64(Kilobyte) }

// Megabytes returns the size in base-10 megabytes (1000 kilobytes).
func (s Size) Megabytes() float64 { return float64(s) / float64(Megabyte) }

// Gigabytes returns the size in base-10 gigabytes (1000 megabytes).
func (s Size) Gigabytes() float64 { return float64(s) / float64(Gigabyte) }

// Terabytes returns the size in base-10 terabytes (1000 gigabytes).
func (s Size) Terabytes() float64 { return float64(s) / float64(Terabyte) }

// Petabytes returns the size in base-10 petabytes (1000 terabytes).
func (s Size) Petabytes() float64 { return float64(s) / float64(Petabyte) }

// Exabytes returns the size in base-10 exabytes (1000 petabytes).
func (s Size) Exabytes() float64 { return float64(s) / float64(Exabyte) }

// String returns a formatted size string with unit suffix using opinionated
// format options:
//
//   - IEC (2-based)
//   - cut empty fractions (42.00 -> 42, 12.30 -> 12.3, 3.14 -> 3.14)
//   - precision: two decimals (max)
//
// Use Format, FormatIEC and FormatSI for full customization.
func (s Size) String() string {
	// format with opinionated settings:
	// - IEC (2-based)
	// - cut empty fraction
	// - precision: two decimals (max)
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
	// append size suffix, optionally separated by a space
	if options.space {
		return str + " " + u.suffix
	}
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
	space            bool
}

// FormatOption configures how a Size is rendered as text. See
// WithCutEmptyFraction and WithPrecision.
type FormatOption func(fOpt *formatOptions)

// WithCutEmptyFraction trims trailing zeros and a trailing decimal point from
// the formatted value (for example 42.00 -> 42 and 12.30 -> 12.3).
func WithCutEmptyFraction() FormatOption {
	return func(fOpt *formatOptions) {
		fOpt.cutEmptyFraction = true
	}
}

// WithPrecision sets the number of decimal places in the formatted value. The
// default is 2. Combine it with WithCutEmptyFraction to treat this as an upper
// bound rather than an exact width. A negative precision selects the smallest
// number of digits that represents the value exactly.
func WithPrecision(precision int) FormatOption {
	return func(fOpt *formatOptions) {
		fOpt.precision = precision
	}
}

// WithSpace inserts a single space between the value and its unit suffix (for
// example "5 MiB" instead of "5MiB"). Parse accepts both forms.
func WithSpace() FormatOption {
	return func(fOpt *formatOptions) {
		fOpt.space = true
	}
}

// Format renders the size in the given Unit, applying any FormatOption.
func (s Size) Format(u Unit, opts ...FormatOption) string {
	return s.format(u.internalUnit(), opts...)
}

// FormatIEC renders the size using the largest base-2 (IEC) unit that keeps the
// value at or above 1, applying any FormatOption.
func (s Size) FormatIEC(opts ...FormatOption) string {
	return s.format(s.fit(unitsIEC), opts...)
}

// FormatSI renders the size using the largest base-10 (SI) unit that keeps the
// value at or above 1, applying any FormatOption.
func (s Size) FormatSI(opts ...FormatOption) string {
	return s.format(s.fit(unitsSI), opts...)
}

// MarshalText implements [encoding.TextMarshaler]. It emits the exact byte count
// with a "B" suffix (for example "5242880B"). The encoding is lossless and is
// reversible via UnmarshalText.
func (s Size) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(s), 10) + "B"), nil
}

// Parse parses a size string such as "5242880B", "5MiB", "5 MB", "5kB", or
// "0.04TiB" into a Size. It tolerates surrounding whitespace and accepts
// fractional values. Negative, infinite, and NaN values are rejected because a
// Size is a non-negative byte count.
//
// Parsing is case-insensitive except for the ambiguous kilo suffix, where the
// case of the "k" selects the base, following common convention:
//
//   - "kB" (lowercase k) is the SI kilobyte and scales by 1000
//   - "KB" (uppercase K) is the JEDEC kilobyte and scales by 1024
//
// Every other unit is unambiguous because its base-2 form carries an "i"
// ("MiB" vs "MB"), so "gb"/"GB" and "mib"/"MiB" each parse regardless of case.
// Parse is the inverse of the Format methods and underlies UnmarshalText.
func Parse(text string) (Size, error) {
	trimmed := strings.TrimSpace(text)
	raw := strings.ToLower(trimmed)

	// Longest suffix matching case-insensitive so MB wins over B
	var match unit
	found := false
	for _, u := range units {
		if strings.HasSuffix(raw, strings.ToLower(u.suffix)) && (!found || len(u.suffix) > len(match.suffix)) {
			match = u
			found = true
		}
	}

	if !found {
		return 0, ErrUnknownUnit
	}

	base := match.base

	// Disambiguate the kilo suffix by the case of its "k": lowercase "kB" is
	// the SI kilobyte (1000), uppercase "KB" is the JEDEC kilobyte (1024). All
	// larger units carry an "i" in their base-2 form, so they never collide.
	if match == unitKB {
		if k := len(trimmed) - len(match.suffix); k >= 0 && trimmed[k] == 'K' {
			base = Kibibyte
		}
	}

	// Isolate the actual number from the suffix, so we can parse it in the next step.
	num := strings.TrimSpace(strings.TrimSuffix(raw, strings.ToLower(match.suffix)))

	// assume bytes first to avoid losing precision due to float64
	if v, err := strconv.ParseUint(num, 10, 64); err == nil {
		return Size(v) * base, nil
	}

	val, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing size value failed: %w", err)
	}

	// Sizes are non-negative and finite; reject values that would otherwise
	// wrap or convert to an implementation-defined uint64.
	if val < 0 || math.IsInf(val, 0) || math.IsNaN(val) {
		return 0, fmt.Errorf("invalid size value %q", num)
	}

	return Size(val * float64(base)), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler], parsing values such as
// "5242880B", "5MiB", "5 MB", or "0.04TiB". See Parse for the accepted syntax.
func (s *Size) UnmarshalText(text []byte) error {
	v, err := Parse(string(text))
	if err != nil {
		return err
	}
	*s = v
	return nil
}

// Set implements [flag.Value], parsing a size string such as "5MiB" so a Size can
// be used directly as a command-line flag. See Parse for the accepted syntax.
func (s *Size) Set(text string) error {
	v, err := Parse(text)
	if err != nil {
		return err
	}
	*s = v
	return nil
}

// LogValue implements [slog.LogValuer], rendering the size via String so it
// appears as a human-readable string (for example "42GiB") in structured logs.
func (s Size) LogValue() slog.Value {
	return slog.StringValue(s.String())
}
