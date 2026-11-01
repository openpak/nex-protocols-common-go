package fakenet

import (
	"bytes"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
)

// RequireError fails the test unless err carries the given result code.
func RequireError(t *testing.T, msg *nex.RMCMessage, err *nex.Error, code uint32) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error 0x%08X, got success", nex.NewError(code, "").ResultCode)
	}

	if msg != nil {
		t.Fatalf("expected no RMC message alongside an error")
	}

	want := nex.NewError(code, "").ResultCode
	if err.ResultCode != want {
		t.Fatalf("result code 0x%08X, want 0x%08X (%s)", err.ResultCode, want, err.Message)
	}
}

// RequireBody fails the test unless msg is a success response for the given
// protocol, method and call ID carrying exactly the wanted body bytes.
func RequireBody(t *testing.T, msg *nex.RMCMessage, err *nex.Error, protocolID uint16, methodID uint32, callID uint32, want []byte) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if msg == nil || !msg.IsSuccess {
		t.Fatalf("expected a success response, got %#v", msg)
	}

	if msg.ProtocolID != protocolID || msg.MethodID != methodID || msg.CallID != callID {
		t.Fatalf("header proto=0x%X method=%d call=%d, want proto=0x%X method=%d call=%d", msg.ProtocolID, msg.MethodID, msg.CallID, protocolID, methodID, callID)
	}

	if !bytes.Equal(msg.Parameters, want) {
		t.Fatalf("body\n got % X\nwant % X", msg.Parameters, want)
	}
}
