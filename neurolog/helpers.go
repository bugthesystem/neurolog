package neurolog

import (
	"fmt"
	"strconv"
	"encoding/binary"
	"math"
)

func toString(item interface{}) string {
	if key, ok := item.(string); ok {
		return key
	} else {
		if integer, ok := item.(int64); ok {
			return strconv.FormatInt(integer, 10)
		} else {
			return fmt.Sprintf("%v", item)
		}
	}
}

func Float64frombytes(bytes []uint8) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}