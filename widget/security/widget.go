// +build !windows

package security

import (
	"fmt"
	"strings"

	"github.com/senorprogrammer/wtf/wtf"
)

type Config struct {
	wtf.WidgetConfig `yaml:",inline"`
}

type Widget struct {
	*wtf.TextWidget
	config *Config
	logger wtf.Logger
}

func (widget *Widget) Name() string {
	return "security"
}

func (widget *Widget) Init(configure wtf.UnmarshalFunc, context *wtf.AppContext) error {
	context.Logger.Debug("Security: init")

	widget.config = &Config{}
	if err := configure(widget.config); err != nil {
		return err
	}

	widget.TextWidget = wtf.NewTextWidget(context.App, "Security", widget.config.WidgetConfig, false)

	widget.logger = context.Logger

	return nil
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	data := NewSecurityData()
	data.Fetch()

	widget.TextView.SetText(widget.contentFrom(data))
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) contentFrom(data *SecurityData) string {
	str := " [red]WiFi[white]\n"
	str = str + fmt.Sprintf(" %8s: %s\n", "Network", data.WifiName)
	str = str + fmt.Sprintf(" %8s: %s\n", "Crypto", data.WifiEncryption)
	str = str + "\n"
	str = str + " [red]Firewall[white]        [red]DNS[white]\n"
	str = str + fmt.Sprintf(" %8s: [%s]%-3s[white]   %-16s\n", "Enabled", widget.labelColor(data.FirewallEnabled), data.FirewallEnabled, data.DnsAt(0))
	str = str + fmt.Sprintf(" %8s: [%s]%-3s[white]   %-16s\n", "Stealth", widget.labelColor(data.FirewallStealth), data.FirewallStealth, data.DnsAt(1))
	str = str + "\n"
	str = str + " [red]Users[white]\n"
	str = str + fmt.Sprintf(" %s", strings.Join(data.LoggedInUsers, ", "))

	return str
}

func (widget *Widget) labelColor(label string) string {
	switch label {
	case "on":
		return "green"
	case "off":
		return "red"
	default:
		return "white"
	}
}
