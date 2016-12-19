package neurolog

import (
	"fmt"
	"strconv"
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
