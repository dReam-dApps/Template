package template

import "github.com/fyne-io/terminal"

// This a Fyne cli terminal that can be used use with a wallets or other dev tools

var cli *terminal.Terminal

// Starts a Fyne terminal app in Template
func startTerminal() *terminal.Terminal {
	cli = terminal.New()
	go func() {
		_ = cli.RunLocalShell()
	}()

	return cli
}

// Exit running terminal
func exitTerminal() {
	if cli != nil {
		cli.Exit()
	}
}
