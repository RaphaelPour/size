package size_test

import (
	"fmt"

	"github.com/raphaelpour/size"
)

func ExampleSize() {
	blockSize := 4 * size.Mebibyte
	fmt.Printf("%d", blockSize.Bytes())

	// Output:
	// 4194304
}
