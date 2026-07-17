package size_test

import (
	"fmt"

	"github.com/raphaelpour/size"
)

func ExampleSize_Bytes() {
	blockSize := 4 * size.Mebibyte
	fmt.Printf("%d", blockSize.Bytes())

	// Output:
	// 4194304
}

func ExampleSize_String() {
	fmt.Println(42 * size.GiB)
	fmt.Println(42 * size.GB)
	fmt.Println(1024 * size.KiB)
	fmt.Println(1024 * size.KB)

	// Output:
	// 42GiB
	// 39.12GiB
	// 1MiB
	// 1000KiB
}

func ExampleSize_Format() {
	s := 42 * size.GiB
	fmt.Println(s.Format(size.UnitTiB))
	fmt.Println(s.Format(size.UnitGiB, size.WithCutEmptyFraction()))
	fmt.Println(s.Format(size.UnitMiB, size.WithCutEmptyFraction()))
	fmt.Println(s.Format(size.UnitKiB))
	fmt.Println(s.Format(size.UnitByte))

	// Output:
	// 0.04TiB
	// 42GiB
	// 43008MiB
	// 44040192.00KiB
	// 45097156608.00B
}

func ExampleSize_FormatIEC() {
	s := 42 * size.GiB
	fmt.Println(s.FormatIEC())
	fmt.Println(s.FormatIEC(size.WithCutEmptyFraction()))

	// Output:
	// 42.00GiB
	// 42GiB
}

func ExampleSize_FormatSI() {
	s := 42 * size.GiB
	fmt.Println(s.FormatSI())
	fmt.Println(s.FormatSI(size.WithPrecision(4)))

	// Output:
	// 45.10GB
	// 45.0972GB
}

func ExampleParse() {
	s, _ := size.Parse("5MiB")
	fmt.Println(s.Bytes())

	// Parsing is case-insensitive; kB/MB are base-10, KiB/MiB base-2.
	kb, _ := size.Parse("1kb")
	kib, _ := size.Parse("1KiB")
	fmt.Println(kb.Bytes(), kib.Bytes())

	// Output:
	// 5242880
	// 1000 1024
}

func ExampleSize_MarshalText() {
	s := 5 * size.MiB

	text, _ := s.MarshalText()
	fmt.Println(string(text))

	var back size.Size
	_ = back.UnmarshalText(text)
	fmt.Println(back == s)

	// Output:
	// 5242880B
	// true
}
