package checks

import (
	"context"
	"net"
	"time"
)

const (
	dnssecControlDomain = "cloudflare.com"
	dnssecTestDomain    = "dnssec-failed.org"
)

type DNSSECResult struct {
	ControlDomain       string
	ControlResolved     bool
	ControlError        string
	TestDomain          string
	TestResolved        bool
	TestError           string
	ProtectionEffective bool
}

func CheckDNSSEC(ctx context.Context) DNSSECResult {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	controlResolved, controlErr := lookupHost(ctx, dnssecControlDomain)
	testResolved, testErr := lookupHost(ctx, dnssecTestDomain)

	return DNSSECResult{
		ControlDomain:       dnssecControlDomain,
		ControlResolved:     controlResolved,
		ControlError:        errorString(controlErr),
		TestDomain:          dnssecTestDomain,
		TestResolved:        testResolved,
		TestError:           errorString(testErr),
		ProtectionEffective: controlResolved && !testResolved,
	}
}

func lookupHost(ctx context.Context, hostname string) (bool, error) {
	_, err := net.DefaultResolver.LookupHost(ctx, hostname)
	if err != nil {
		return false, err
	}

	return true, nil
}

func errorString(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
