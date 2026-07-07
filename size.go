package size

type Size uint64

const (
	Byte     Size = 1
	B             = Byte
	Kibibyte      = Byte * 1024
	Kib           = Kibibyte
	Mebibyte      = Kibibyte * 1024
	Mib           = Mebibyte
	Gibibyte      = Mebibyte * 1024
	Gib           = Gibibyte
	Tebibyte      = Gibibyte * 1024
	Tib           = Tebibyte
	Pebibyte      = Tebibyte * 1024
	Pib           = Pebibyte

	// Exibytes are only supported til 2^64
	Exbibyte = Pebibyte * 1024
	Eib      = Exbibyte

	// zebibyte is 2^70 and larger as uint64 2^64
)

func (s Size) Bytes() uint64 {
	return uint64(s)
}
