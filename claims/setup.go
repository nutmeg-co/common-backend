package claims

import (
	"github.com/urfave/cli/v2"
)

var (
	ClaimsFactory Factory
)

func SetupSecret(cli *cli.Context, secret string) error {
	ClaimsFactory.Secret = []byte(secret)
	return nil
}
