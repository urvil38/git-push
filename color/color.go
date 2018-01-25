package color

import (
	"fmt"
)

type attribute int

const (
	fgblack attribute = iota + 30
	fgred
	fggreen
	fgyellow
	fgblue
	fgmagenta
	fgcyan
	fgwhite
)

const (
	reset attribute = iota
	bold
	faint
	italic
	underline
	blinkSlow
	blinkRapid
	reverseVideo
	concealed
	crossedOut
)

const escape = "\x1b"

type colorMap map[string]attribute
type effectMap map[string]attribute

var colormap = colorMap{
	"FgBlack":   fgblack,
	"FgRed":     fgred,
	"FgGreen":   fggreen,
	"FgYellow":  fgyellow,
	"FgBlue":    fgblue,
	"FgMagenta": fgmagenta,
	"FgCyan":    fgcyan,
	"FgWhite":   fgwhite,
}

var effectemap = effectMap{
	"Reset":        reset,
	"Bold":         bold,
	"Faint":        faint,
	"Italic":       italic,
	"Underline":    underline,
	"BlinkSlow":    blinkSlow,
	"BlinkRapid":   blinkRapid,
	"ReverseVideo": reverseVideo,
	"Concealed":    concealed,
	"CrossedOut":   crossedOut,
}

// Wrap returns colored ascii text .Now it only use cyan color but you can change it
func Wrap(s, color, effect string) string {
	return format(colormap[color], effectemap[effect]) + s + unformat()
}

func format(a, c attribute) string {
	return fmt.Sprintf("\x1b[%v;%vm", a, c)
}

func unformat() string {
	return fmt.Sprintf("\x1b[%vm", effectemap["Reset"])
}
