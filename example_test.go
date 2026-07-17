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
