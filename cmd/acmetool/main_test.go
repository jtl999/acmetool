package main

import (
	"fmt"
	"github.com/hlandau/acme/acmeapi"
	"github.com/hlandau/acme/interaction"
	"github.com/hlandau/acme/responder"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

type interceptor struct {
}

func (i *interceptor) Prompt(c *interaction.Challenge) (*interaction.Response, error) {
	switch c.UniqueID {
	case "acmetool-quickstart-choose-server":
		return &interaction.Response{Value: "url"}, nil
	case "acmetool-quickstart-enter-directory-url":
		return &interaction.Response{Value: "https://127.0.0.1:4000/directory"}, nil
	case "acmetool-quickstart-choose-method":
		return &interaction.Response{Value: "redirector"}, nil
	case "acme-enter-email":
		return &interaction.Response{Value: "nobody@example.com"}, nil
	case "acmetool-quickstart-complete":
		return &interaction.Response{}, nil
	case "acmetool-quickstart-install-cronjob", "acmetool-quickstart-install-haproxy-script", "acmetool-quickstart-install-redirector-systemd":
		return &interaction.Response{Cancelled: true}, nil
	default:
		if strings.HasPrefix(c.UniqueID, "acme-agreement:") {
			return &interaction.Response{}, nil
		}

		return nil, fmt.Errorf("unsupported challenge for interceptor: %v", c)
	}
}

func (i *interceptor) Status(info *interaction.StatusInfo) (interaction.StatusSink, error) {
	return nil, fmt.Errorf("status not supported")
}

func TestCLI(t *testing.T) {
	log.Warnf("This test requires a configured Boulder instance listening at http://127.0.0.1:4000/ and the ability to successfully complete challenges. You must change the Boulder configuration to use ports 80 and 5443. Also change the rate limits per certificate name. Consider ensuring that the user you run these tests as can write to /var/run/acme/acme-challenge and that that directory is served on port 80 /.well-known/acme-challenge/")

	acmeapi.TestingNoTLS = true

	interaction.Interceptor = &interceptor{}

	tmpDir, err := ioutil.TempDir("", "acmetool-test")
	if err != nil {
		panic(err)
	}

	*stateFlag = filepath.Join(tmpDir, "state")
	*hooksFlag = filepath.Join(tmpDir, "hooks")

	responder.InternalTLSSNIPort = 5443
	quickstart()

	*wantArg = []string{"dom1.acmetool-test.devever.net", "dom2.acmetool-test.devever.net"}

	want()
	reconcile()
}

// © 2015 Hugo Landau <hlandau@devever.net>    MIT License
