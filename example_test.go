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
	fmt.Println(42 * size.Gib)
	fmt.Println(42 * size.Gb)
	fmt.Println(1024 * size.Kib)
	fmt.Println(1024 * size.Kb)

	// Output:
	// 42Gib
	// 39.12Gib
	// 1Mib
	// 1000Kib
}
