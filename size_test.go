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
		{name: "bytes", input: 1234 * size.Byte, expected: 1234},
		{name: "kibibytes", input: 1234 * size.Kibibyte, expected: 1234 * 1024},
		{name: "mebibytes", input: 1234 * size.Mebibyte, expected: 1234 * 1024 * 1024},
		{name: "gebibytes", input: 1234 * size.Gibibyte, expected: 1234 * 1024 * 1024 * 1024},
		{name: "tebibytes", input: 1234 * size.Tebibyte, expected: 1234 * 1024 * 1024 * 1024 * 1024},
		{name: "pebibytes", input: 1234 * size.Pebibyte, expected: 1234 * 1024 * 1024 * 1024 * 1024 * 1024},
		{name: "ebibytes", input: 12 * size.Exbibyte, expected: 12 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(t, testCase.expected, testCase.input.Bytes())
		})
	}
}
