package types

import "strconv"

type Count = uint64

func CountToByteSlice(count Count) []byte {
	return []byte(strconv.FormatUint(count, 10))
}

func CountFromByteSlice(data []byte) (Count, error) {
	return strconv.ParseUint(string(data), 10, 64)
}
