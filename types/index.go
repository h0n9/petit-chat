package types

import "strconv"

type (
	Index  = uint64
	Length = uint64
)

func IndexToByteSlice(index Index) []byte {
	return []byte(strconv.FormatUint(index, 10))
}

func IndexFromByteSlice(data []byte) (Index, error) {
	return strconv.ParseUint(string(data), 10, 64)
}
