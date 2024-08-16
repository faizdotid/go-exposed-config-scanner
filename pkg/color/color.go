package color

import "strconv"

// Coloring function
func Coloring(data string, colors ...Color) string {
	return string(Color(0).AnsiFormat(data, colors...))
}

// append the color code to the buffer
func (c Color) parseColor(buffer []byte, color Color) []byte {
	// directly append the color code to avoid unnecessary allocations
	buffer = append(buffer, Escape...)
	buffer = append(buffer, '[')
	buffer = append(buffer, strconv.Itoa(int(color))...)
	buffer = append(buffer, 'm')
	return buffer
}

// AnsiFormat function
func (c Color) AnsiFormat(data string, colors ...Color) string {
	// if c != 0 {
	// 	colors = append([]Color{c}, colors...)
	// }
	// calculate the total length of the buffer
	totalLen := len(data) + len(colors)*5 + 4 // each color is "\x1b[XXm" (5 bytes), reset is "\x1b[0m" (4 bytes)

	// create a buffer
	buffer := make([]byte, 0, totalLen)

	// append the color codes
	for _, color := range colors {
		buffer = color.parseColor(buffer, color)
	}
	// append data
	buffer = append(buffer, data...)

	// append reset code
	buffer = append(buffer, Escape...)
	buffer = append(buffer, '[', byte(Reset), 'm')
	return string(buffer)
}
