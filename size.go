package size

type Size uint64

const (
	Byte Size = 1
	B         = Byte
)

// base-2
const (
	Bytes Size = 1 << (10 * iota)
	Kibibyte
	Mebibyte
	Gibibyte
	Tebibyte
	Pebibyte
	Exbibyte

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

// base-10
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
