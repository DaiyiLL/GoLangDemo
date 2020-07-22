package main

import (
	"fmt"
	"strconv"
)



func LogError(con string, err error)  {
	fmt.Println(con, "err: ", err)
}

func DYToString(obj interface{}) string  {
	if obj == nil {
		return ""
	}
	if v, ok := obj.(string); ok {
		return v
	}
	return fmt.Sprintf("%v", obj)
}

func DYToInt64(obj interface{}) int64  {
	if obj == nil {
		return 0
	}

	if v, ok := obj.(int64); ok {
		return v
	}

	switch value := obj.(type) {
	case int:
		return int64(value)

	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return value

	case bool:
		if value {
			return 1
		} else {
			return 0
		}
	case uint:
		return int64(value)
	case uint8:
		return int64(value)
	case uint16:
		return int64(value)
	case uint32:
		return int64(value)
	case uint64:
		return int64(value)
	case float32:
		return int64(value)
	case float64:
		return int64(value)

	case string:
		{
			num, err := strconv.ParseInt(DYToString(obj), 0, 64)
			if err != nil {
				return 0
			}
			return num
		}
	default:
		return 0
	}
}

//func DYMapToText(data map[string]interface{}, keys ...string)  {
//	tempData
//}