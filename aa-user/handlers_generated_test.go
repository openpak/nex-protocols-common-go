package aauser

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
	aauser_types "github.com/PretendoNetwork/nex-protocols-go/v2/aa-user/types"
)

// TestHooksAbsentAnswerNotImplemented checks that every method without a hook answers Core::NotImplemented
func TestHooksAbsentAnswerNotImplemented(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	packet := rmctest.NewPacket(1)

	calls := map[string]func() (*nex.RMCMessage, *nex.Error){
		"RegisterApplication": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.registerApplication(nil, packet, 1, *new(types.UInt64))
		},
		"UnregisterApplication": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.unregisterApplication(nil, packet, 1, *new(types.UInt64))
		},
		"SetApplicationInfo": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.setApplicationInfo(nil, packet, 1, *new(types.List[aauser_types.ApplicationInfo]))
		},
		"GetApplicationInfo": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.getApplicationInfo(nil, packet, 1) },
	}

	if len(calls) != 4 {
		t.Fatalf("expected 4 methods, got %d", len(calls))
	}

	for name, call := range calls {
		response, nexError := call()
		if response != nil || !rmctest.HasCode(nexError, nex.ResultCodes.Core.NotImplemented) {
			t.Errorf("%s: expected Core::NotImplemented, got %v", name, nexError)
		}
	}
}

// TestDecodeErrorsAreInvalidArgument checks that a request that failed to decode never reaches a hook
func TestDecodeErrorsAreInvalidArgument(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	packet := rmctest.NewPacket(1)

	calls := map[string]func() (*nex.RMCMessage, *nex.Error){
		"RegisterApplication": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.registerApplication(errBadRequest, packet, 1, *new(types.UInt64))
		},
		"UnregisterApplication": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.unregisterApplication(errBadRequest, packet, 1, *new(types.UInt64))
		},
		"SetApplicationInfo": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.setApplicationInfo(errBadRequest, packet, 1, *new(types.List[aauser_types.ApplicationInfo]))
		},
		"GetApplicationInfo": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getApplicationInfo(errBadRequest, packet, 1)
		},
	}

	for name, call := range calls {
		response, nexError := call()
		if response != nil || !rmctest.HasCode(nexError, nex.ResultCodes.Core.InvalidArgument) {
			t.Errorf("%s: expected Core::InvalidArgument, got %v", name, nexError)
		}
	}
}

var errBadRequest = errors.New("bad request")
