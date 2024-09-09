package eth

import (
	"fmt"
	"strconv"
)

func ParseHex(hexStr string) int64 {
	intValue, err := strconv.ParseInt(hexStr[2:], 16, 64)
	if err != nil {
		// Handle the error
		fmt.Println("Error:", err)
		return -1
	}

	return intValue
}
