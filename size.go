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

// Size represents a quantity of bytes. The zero value represents zero bytes.
type Size uint64

// Base-2 (IEC) units of size. Sizes are by a magnitde of 1024 apart.
const (
	Byte Size = 1 << (10 * iota)
	Kibibyte
	Mebibyte
	Gibibyte
	Tebibyte
	Pebibyte
	Exbibyte

	B   = Byte
	Kib = Kibibyte
	Mib = Mebibyte
	Gib = Gibibyte
	Tib = Tebibyte
	Pib = Pebibyte

	// Exibytes are only supported til 2^64
	Eib = Exbibyte

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

	Kb = Kilobyte
	Mb = Megabyte
	Gb = Gigabyte
	Tb = Terabyte
	Pb = Petabyte

	// Exabyte is limited by uint64
	Eb = Exabyte

	// zebibyte and above are not supported limited by uint64
)

// Bytes returns the byte size as uint64.
func (s Size) Bytes() uint64 {
	return uint64(s)
}
