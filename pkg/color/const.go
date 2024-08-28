// Description: Constants for color package
package color

// Escape character
const (
	Escape = "\x1b"
)

// Standard colors
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Text formatting options
const (
	Reset Color = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground Hi-Intensity (Bright) colors
const (
	HiBlack Color = iota + 90
	HiRed
	HiGreen
	HiYellow
	HiBlue
	HiMagenta
	HiCyan
	HiWhite
)

// Background colors
const (
	BgBlack Color = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity (Bright) colors
const (
	BrightBgBlack Color = iota + 100
	BrightBgRed
	BrightBgGreen
	BrightBgYellow
	BrightBgBlue
	BrightBgMagenta
	BrightBgCyan
	BrightBgWhite
)

// Additional colors
const (
	Orange     Color = 208
	Purple     Color = 129
	LightBlue  Color = 117
	Pink       Color = 213
	LightGreen Color = 119
)
