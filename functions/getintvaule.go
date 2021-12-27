package functions

import (
	"errors"
	"strconv"
)

func GetIntV(i interface{}) (int, error) {
	switch v := i.(type) {
	case int:
	case int8:
	case int16:
	case int32:
	case int64:
	case float32:
	case float64:
	case uint:
	case uint8:
	case uint16:
	case uint32:
	case uint64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	}
	return 0, errors.New("unkown type")
}
