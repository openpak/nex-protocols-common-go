package persistentstore

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
)

// TestHooksAbsentAnswerNotImplemented checks that every method without a hook answers Core::NotImplemented
func TestHooksAbsentAnswerNotImplemented(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	packet := rmctest.NewPacket(1)

	calls := map[string]func() (*nex.RMCMessage, *nex.Error){
		"FindByGroup": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.findByGroup(nil, packet, 1, *new(types.UInt32))
		},
		"InsertItem": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.insertItem(nil, packet, 1, *new(types.UInt32), *new(types.String), *new(types.Buffer), *new(types.Bool))
		},
		"RemoveItem": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.removeItem(nil, packet, 1, *new(types.UInt32), *new(types.String))
		},
		"GetItem": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getItem(nil, packet, 1, *new(types.UInt32), *new(types.String))
		},
		"InsertCustomItem": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.insertCustomItem(nil, packet, 1, *new(types.UInt32), *new(types.String), *new(types.DataHolder), *new(types.Bool))
		},
		"GetCustomItem": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getCustomItem(nil, packet, 1, *new(types.UInt32), *new(types.String))
		},
		"FindItemsBySQLQuery": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.findItemsBySQLQuery(nil, packet, 1, *new(types.UInt32), *new(types.String), *new(types.String))
		},
	}

	if len(calls) != 7 {
		t.Fatalf("expected 7 methods, got %d", len(calls))
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
		"FindByGroup": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.findByGroup(errBadRequest, packet, 1, *new(types.UInt32))
		},
		"InsertItem": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.insertItem(errBadRequest, packet, 1, *new(types.UInt32), *new(types.String), *new(types.Buffer), *new(types.Bool))
		},
		"RemoveItem": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.removeItem(errBadRequest, packet, 1, *new(types.UInt32), *new(types.String))
		},
		"GetItem": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getItem(errBadRequest, packet, 1, *new(types.UInt32), *new(types.String))
		},
		"InsertCustomItem": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.insertCustomItem(errBadRequest, packet, 1, *new(types.UInt32), *new(types.String), *new(types.DataHolder), *new(types.Bool))
		},
		"GetCustomItem": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.getCustomItem(errBadRequest, packet, 1, *new(types.UInt32), *new(types.String))
		},
		"FindItemsBySQLQuery": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.findItemsBySQLQuery(errBadRequest, packet, 1, *new(types.UInt32), *new(types.String), *new(types.String))
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
