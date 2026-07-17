package size_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/raphaelpour/size"
)

func TestSize_Bytes(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		input    size.Size
		expected uint64
	}{
		{name: "kibibyte constant", input: size.Kibibyte, expected: 1024},
		{name: "mebibyte constant", input: size.Mebibyte, expected: 1024 * 1024},
		{name: "bytes", input: 1234 * size.Byte, expected: 1234},
		{name: "kibibytes", input: 1234 * size.Kibibyte, expected: 1234 * 1024},
		{name: "mebibytes", input: 1234 * size.Mebibyte, expected: 1234 * 1024 * 1024},
		{name: "gebibytes", input: 1234 * size.Gibibyte, expected: 1234 * 1024 * 1024 * 1024},
		{name: "tebibytes", input: 1234 * size.Tebibyte, expected: 1234 * 1024 * 1024 * 1024 * 1024},
		{name: "pebibytes", input: 1234 * size.Pebibyte, expected: 1234 * 1024 * 1024 * 1024 * 1024 * 1024},
		{name: "exbibytes", input: 12 * size.Exbibyte, expected: 12 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024},
		{name: "b", input: 1234 * size.B, expected: 1234},
		{name: "kib", input: 1234 * size.KiB, expected: 1234 * 1024},
		{name: "mib", input: 1234 * size.MiB, expected: 1234 * 1024 * 1024},
		{name: "gib", input: 1234 * size.GiB, expected: 1234 * 1024 * 1024 * 1024},
		{name: "tib", input: 1234 * size.TiB, expected: 1234 * 1024 * 1024 * 1024 * 1024},
		{name: "pib", input: 1234 * size.PiB, expected: 1234 * 1024 * 1024 * 1024 * 1024 * 1024},
		{name: "eib", input: 12 * size.EiB, expected: 12 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024},
		{name: "kilobytes", input: 1234 * size.Kilobyte, expected: 1234 * 1000},
		{name: "megabytes", input: 1234 * size.Megabyte, expected: 1234 * 1000 * 1000},
		{name: "gigabytes", input: 1234 * size.Gigabyte, expected: 1234 * 1000 * 1000 * 1000},
		{name: "terabytes", input: 1234 * size.Terabyte, expected: 1234 * 1000 * 1000 * 1000 * 1000},
		{name: "petabytes", input: 1234 * size.Petabyte, expected: 1234 * 1000 * 1000 * 1000 * 1000 * 1000},
		{name: "exabytes", input: 12 * size.Exabyte, expected: 12 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000},
		{name: "kb", input: 1234 * size.KB, expected: 1234 * 1000},
		{name: "mb", input: 1234 * size.MB, expected: 1234 * 1000 * 1000},
		{name: "gb", input: 1234 * size.GB, expected: 1234 * 1000 * 1000 * 1000},
		{name: "tb", input: 1234 * size.TB, expected: 1234 * 1000 * 1000 * 1000 * 1000},
		{name: "pb", input: 1234 * size.PB, expected: 1234 * 1000 * 1000 * 1000 * 1000 * 1000},
		{name: "eb", input: 12 * size.EB, expected: 12 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(t, testCase.expected, testCase.input.Bytes())
		})
	}
}

func TestSize_Format(t *testing.T) {
	sz := 10 * size.MiB
	require.Equal(t, sz.String(), sz.Format(size.UnitMiB, size.WithCutEmptyFraction()))
}

func TestSize_FormatWithSpace(t *testing.T) {
	sz := 5 * size.MiB
	require.Equal(t, "5 MiB", sz.Format(size.UnitMiB, size.WithCutEmptyFraction(), size.WithSpace()))
	require.Equal(t, "5.00 MiB", sz.FormatIEC(size.WithSpace()))
}

func TestSize_TextRoundTrip(t *testing.T) {
	for _, testCase := range []struct {
		name  string
		input size.Size
	}{
		{name: "zero", input: 0},
		{name: "bytes", input: 1234 * size.B},
		{name: "kib", input: 5 * size.KiB},
		{name: "mib", input: 5 * size.MiB},
		{name: "gib", input: 5 * size.GiB},
		{name: "tib", input: 5 * size.TiB},
		{name: "pib", input: 5 * size.PiB},
		{name: "eib", input: 5 * size.EiB},
		{name: "kb", input: 5 * size.KB},
		{name: "mb", input: 5 * size.MB},
		{name: "gb", input: 5 * size.GB},
		{name: "tb", input: 5 * size.TB},
		{name: "pb", input: 5 * size.PB},
		{name: "eb", input: 5 * size.EB},
		{name: "large above 2^53", input: 1 << 60},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			text, err := testCase.input.MarshalText()
			require.NoError(t, err)

			var got size.Size
			require.NoError(t, got.UnmarshalText(text))
			require.Equal(t, testCase.input, got)
		})
	}
}

func TestParse(t *testing.T) {
	frac := 0.04
	for _, testCase := range []struct {
		name     string
		input    string
		expected size.Size
		wantErr  bool
	}{
		{name: "bytes", input: "5B", expected: 5 * size.B},
		{name: "MB suffix", input: "5MB", expected: 5 * size.MB},
		{name: "kB suffix", input: "5kB", expected: 5 * size.KB},
		{name: "KiB suffix", input: "5KiB", expected: 5 * size.KiB},
		{name: "MiB suffix", input: "5MiB", expected: 5 * size.MiB},
		{name: "lowercase kb is SI kilobyte", input: "5kb", expected: 5 * size.KB},
		{name: "uppercase KB is JEDEC kibibyte", input: "5KB", expected: 5 * size.KiB},
		{name: "mixed Kb is JEDEC kibibyte", input: "5Kb", expected: 5 * size.KiB},
		{name: "leading space then KB", input: "  5 KB", expected: 5 * size.KiB},
		{name: "lowercase gb", input: "5gb", expected: 5 * size.GB},
		{name: "uppercase mib", input: "5MIB", expected: 5 * size.MiB},
		{name: "lowercase byte", input: "5b", expected: 5 * size.B},
		{name: "leading and trailing space", input: "  5 MiB  ", expected: 5 * size.MiB},
		{name: "fractional", input: "0.04TiB", expected: size.Size(frac * float64(size.TiB))},
		{name: "missing suffix", input: "5", wantErr: true},
		{name: "non-numeric value", input: "abcMiB", wantErr: true},
		{name: "negative value", input: "-5MiB", wantErr: true},
		{name: "infinite value", input: "infMiB", wantErr: true},
		{name: "nan value", input: "nanMiB", wantErr: true},
		{name: "with space", input: "5 MiB", expected: 5 * size.MiB},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			// Repeat to catch any residual map-iteration nondeterminism.
			for range 50 {
				got, err := size.Parse(testCase.input)
				if testCase.wantErr {
					require.Error(t, err)
					continue
				}
				require.NoError(t, err)
				require.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestSize_UnitAccessors(t *testing.T) {
	require.Equal(t, float64(1), (1 * size.Kibibyte).Kibibytes())
	require.Equal(t, float64(1), (1 * size.Mebibyte).Mebibytes())
	require.Equal(t, float64(1), (1 * size.Gibibyte).Gibibytes())
	require.Equal(t, float64(1), (1 * size.Tebibyte).Tebibytes())
	require.Equal(t, float64(1), (1 * size.Pebibyte).Pebibytes())
	require.Equal(t, float64(1), (1 * size.Exbibyte).Exbibytes())

	require.Equal(t, float64(1), (1 * size.Kilobyte).Kilobytes())
	require.Equal(t, float64(1), (1 * size.Megabyte).Megabytes())
	require.Equal(t, float64(1), (1 * size.Gigabyte).Gigabytes())
	require.Equal(t, float64(1), (1 * size.Terabyte).Terabytes())
	require.Equal(t, float64(1), (1 * size.Petabyte).Petabytes())
	require.Equal(t, float64(1), (1 * size.Exabyte).Exabytes())

	require.Equal(t, 0.5, (512 * size.Kibibyte).Mebibytes())
	require.Equal(t, float64(1024), size.Mebibyte.Kibibytes())
}

func TestSize_UnmarshalText(t *testing.T) {
	// Non-constant so the fractional float can be converted to Size, matching
	// exactly what UnmarshalText computes internally.
	frac := 0.04
	for _, testCase := range []struct {
		name     string
		input    string
		expected size.Size
		wantErr  bool
	}{
		{name: "bytes", input: "5B", expected: 5 * size.B},
		{name: "ambiguous MB suffix", input: "5MB", expected: 5 * size.MB},
		{name: "ambiguous kB suffix", input: "5kB", expected: 5 * size.KB},
		{name: "ambiguous KiB suffix", input: "5KiB", expected: 5 * size.KiB},
		{name: "ambiguous MiB suffix", input: "5MiB", expected: 5 * size.MiB},
		{name: "leading and trailing space", input: "  5 MiB  ", expected: 5 * size.MiB},
		{name: "fractional", input: "0.04TiB", expected: size.Size(frac * float64(size.TiB))},
		{name: "missing suffix", input: "5", wantErr: true},
		{name: "non-numeric value", input: "abcMiB", wantErr: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var got size.Size
			err := got.UnmarshalText([]byte(testCase.input))
			if testCase.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, testCase.expected, got)
		})
	}
}
