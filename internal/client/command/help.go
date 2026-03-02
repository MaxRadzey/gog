package command

import "context"

const helpText = `GOG — password manager client.

Commands:

  register <login> <password>     — register a user
  login <login> <password>        — sign in
  logout                          — sign out

  secret-list                     — list secrets (id, type, created)
  secret-get <id>                 — show secret by id
  secret-create <type> [key=val ...]  — create secret (e.g. secret-create login login=u password=p name=site)
  secret-update <id> key=val [...]    — update secret
  secret-delete <id>              — delete secret

  version                         — show version and build date
  help                            — this help
`

// HelpCommand выводит справку по командам.
type HelpCommand struct{}

// NewHelpCommand создаёт команду help.
func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

// Execute возвращает текст справки (helpText).
func (c *HelpCommand) Execute(ctx context.Context, args []string) (string, error) {
	return helpText, nil
}
