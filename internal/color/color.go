package color

import (
	"fmt"

	"github.com/muesli/termenv"
)

var p = termenv.ColorProfile()

func FgBlue(s string) string {
	styled := termenv.String(s).Foreground(p.Color("#2F99EE"))
	return fmt.Sprintf("%s", styled)
}

func FgGreen(s string) string {
	styled := termenv.String(s).Foreground(p.Color("#9ece6a"))
	return fmt.Sprintf("%s", styled)
}

func FgYellow(s string) string {
	styled := termenv.String(s).Foreground(p.Color("#e0af68"))
	return fmt.Sprintf("%s", styled)
}

func FgGray(s string) string {
	styled := termenv.String(s).Foreground(p.Color("#777777"))
	return fmt.Sprintf("%s", styled)
}
