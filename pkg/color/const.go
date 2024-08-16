// Description: Constants for color package
package color

const (
	// Escape character
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

// Constants for text formatting
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

// Constants for Foreground Hi-Intensity colors
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

// Constants for bold colors
const (
	BoldBlack Color = 90 + iota
	BoldRed
	BoldGreen
	BoldYellow
	BoldBlue
	BoldMagenta
	BoldCyan
	BoldWhite
)

// Constants for background colors
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

// Constants for bright colors
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

// Constants for additional colors
const (
	Orange     Color = 208
	Purple     Color = 129
	LightBlue  Color = 117
	Pink       Color = 213
	LightGreen Color = 119
)
