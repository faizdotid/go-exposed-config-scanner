package utils

import (
	"os"
	"strconv"
	"strings"
)

func WriteFile(writer *os.File, data []byte) {
	data = append(data, newLine...)
	writer.Write(data)
}

// generic function to convert any type to interface{}
func ExplodeString[T string | int](data string) []T {
	var result []T
	parts := strings.Split(data, ",")
	for i := 0; i < len(parts); i++ {
		var value T
		switch any(value).(type) {
		case string:
			strValue := parts[i]
			if strings.HasSuffix(strValue, "\\") && i+1 < len(parts) {
				strValue = strValue[:len(strValue)-1] + "," + parts[i+1]
				i++ // Skip the next part as we've already included it
			}
			value = any(strValue).(T)
		case int:
			if intValue, err := strconv.Atoi(parts[i]); err == nil {
				value = any(intValue).(T)
			}
		}
		result = append(result, value)
	}
	return result
}
