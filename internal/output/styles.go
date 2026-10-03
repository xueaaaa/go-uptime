package output

import (
	"github.com/pterm/pterm"
	"github.com/xueaaaa/go-uptime/internal/site/model"
)

var Header = pterm.DefaultHeader.
	WithTextStyle(pterm.NewStyle(pterm.FgLightGreen)).
	WithBackgroundStyle(pterm.NewStyle()).
	WithFullWidth()

var Table = pterm.DefaultTable.
	WithHasHeader().
	WithBoxed()

func PingColor(ping int) string {
	switch {
	case ping < 100:
		return pterm.LightGreen(ping)
	case ping <= 300:
		return pterm.LightYellow(ping)
	default:
		return pterm.LightRed(ping)
	}
}

func StatusColor(status model.Status) string {
	switch {
	case status == model.Unknown:
		return pterm.Gray("Unknown")
	case status == model.Unavailable:
		return pterm.LightRed("Unavailable")
	case status == model.Available:
		return pterm.LightGreen("Available")
	}

	return ""
}
