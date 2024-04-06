// Package main implements a Composition Function.
package main

import (
	"github.com/alecthomas/kong"
	"github.com/apple/pkl-go/pkl"

	"github.com/crossplane/function-sdk-go"

	fn "github.com/avarei/function-pkl/internal/function"
	"github.com/avarei/function-pkl/internal/pkl/reader"
)

// CLI of this Function.
type CLI struct {
	Debug bool `short:"d" help:"Emit debug logs in addition to info logs."`

	Network     string `help:"Network on which to listen for gRPC connections." default:"tcp"`
	Address     string `help:"Address at which to listen for gRPC connections." default:":9443"`
	TLSCertsDir string `help:"Directory containing server certs (tls.key, tls.crt) and the CA used to verify client certificates (ca.crt)" env:"TLS_SERVER_CERTS_DIR"`
	Insecure    bool   `help:"Run without mTLS credentials. If you supply this flag --tls-server-certs-dir will be ignored."`
}

// Run this Function.
func (c *CLI) Run() error {
	log, err := function.NewLogger(c.Debug)
	if err != nil {
		return err
	}

	evaluatorManager := pkl.NewEvaluatorManager()
	defer evaluatorManager.Close()
	defer reader.Close()

	return function.Serve(&fn.Function{Log: log, EvaluatorManager: evaluatorManager},
		function.Listen(c.Network, c.Address),
		function.MTLSCertificates(c.TLSCertsDir),
		function.Insecure(c.Insecure))
}

func main() {
	ctx := kong.Parse(&CLI{}, kong.Description("A Crossplane Composition Function."))
	ctx.FatalIfErrorf(ctx.Run())
}
