package screening

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
	screening_types "github.com/PretendoNetwork/nex-protocols-go/v2/screening/types"
)

// TestHooksAbsentAnswerNotImplemented checks that every method without a hook answers Core::NotImplemented
func TestHooksAbsentAnswerNotImplemented(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	packet := rmctest.NewPacket(1)

	calls := map[string]func() (*nex.RMCMessage, *nex.Error){
		"ReportDataStoreContent": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.reportDataStoreContent(nil, packet, 1, *new(screening_types.ScreeningDataStoreContentParam), *new(screening_types.ScreeningUGCViolationParam))
		},
		"ReportUser": func() (*nex.RMCMessage, *nex.Error) { return commonProtocol.reportUser(nil, packet, 1, *new([]byte)) },
	}

	if len(calls) != 2 {
		t.Fatalf("expected 2 methods, got %d", len(calls))
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
		"ReportDataStoreContent": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.reportDataStoreContent(errBadRequest, packet, 1, *new(screening_types.ScreeningDataStoreContentParam), *new(screening_types.ScreeningUGCViolationParam))
		},
		"ReportUser": func() (*nex.RMCMessage, *nex.Error) {
			return commonProtocol.reportUser(errBadRequest, packet, 1, *new([]byte))
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
