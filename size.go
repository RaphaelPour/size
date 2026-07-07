package size

type Size uint64

const (
	Byte Size = 1
	B         = Byte

	// base-2
	Kibibyte = Byte * 1024
	Kib      = Kibibyte
	Mebibyte = Kibibyte * 1024
	Mib      = Mebibyte
	Gibibyte = Mebibyte * 1024
	Gib      = Gibibyte
	Tebibyte = Gibibyte * 1024
	Tib      = Tebibyte
	Pebibyte = Tebibyte * 1024
	Pib      = Pebibyte

	// Exibytes are only supported til 2^64
	Exbibyte = Pebibyte * 1024
	Eib      = Exbibyte

	// zebibyte is 2^70 and larger as uint64 2^64

	// base-10
	Kilobyte = Byte * 1000
	Kb       = Kilobyte
	Megabyte = Kilobyte * 1000
	Mb       = Megabyte
	Gigabyte = Megabyte * 1000
	Gb       = Gigabyte
	Terabyte = Gigabyte * 1000
	Tb       = Terabyte
	Petabyte = Terabyte * 1000
	Pb       = Petabyte

	// base-10 Exabyte fits entirely into uint64 as it only needs ~2*60
	Exabyte = Petabyte * 1000
	Eb      = Exabyte
)

// Bytes returns the byte size as uint64.
func (s Size) Bytes() uint64 {
	return uint64(s)
}
