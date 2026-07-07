package size_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/raphaelpour/size"
)

func TestSize(t *testing.T) {
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
		{name: "kib", input: 1234 * size.Kib, expected: 1234 * 1024},
		{name: "mib", input: 1234 * size.Mib, expected: 1234 * 1024 * 1024},
		{name: "gib", input: 1234 * size.Gib, expected: 1234 * 1024 * 1024 * 1024},
		{name: "tib", input: 1234 * size.Tib, expected: 1234 * 1024 * 1024 * 1024 * 1024},
		{name: "pib", input: 1234 * size.Pib, expected: 1234 * 1024 * 1024 * 1024 * 1024 * 1024},
		{name: "eib", input: 12 * size.Eib, expected: 12 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024},
		{name: "kilobytes", input: 1234 * size.Kilobyte, expected: 1234 * 1000},
		{name: "megabytes", input: 1234 * size.Megabyte, expected: 1234 * 1000 * 1000},
		{name: "gigabytes", input: 1234 * size.Gigabyte, expected: 1234 * 1000 * 1000 * 1000},
		{name: "terabytes", input: 1234 * size.Terabyte, expected: 1234 * 1000 * 1000 * 1000 * 1000},
		{name: "petabytes", input: 1234 * size.Petabyte, expected: 1234 * 1000 * 1000 * 1000 * 1000 * 1000},
		{name: "exabytes", input: 12 * size.Exabyte, expected: 12 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000},
		{name: "kb", input: 1234 * size.Kb, expected: 1234 * 1000},
		{name: "mb", input: 1234 * size.Mb, expected: 1234 * 1000 * 1000},
		{name: "gb", input: 1234 * size.Gb, expected: 1234 * 1000 * 1000 * 1000},
		{name: "tb", input: 1234 * size.Tb, expected: 1234 * 1000 * 1000 * 1000 * 1000},
		{name: "pb", input: 1234 * size.Pb, expected: 1234 * 1000 * 1000 * 1000 * 1000 * 1000},
		{name: "eb", input: 12 * size.Eb, expected: 12 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(t, testCase.expected, testCase.input.Bytes())
		})
	}
}
