package utils

import (
	"os"
)

func WriteFile(writer *os.File, data []byte) {
	data = append(data, newLine...)
	writer.Write(data)
}
