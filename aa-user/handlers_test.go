package aauser

import (
	"bytes"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
	aa_user "github.com/PretendoNetwork/nex-protocols-go/v2/aa-user"
	aauser_types "github.com/PretendoNetwork/nex-protocols-go/v2/aa-user/types"
)

func TestProtocolID(t *testing.T) {
	if aa_user.ProtocolID != 0x7B {
		t.Fatalf("AA User protocol ID = %#x", aa_user.ProtocolID)
	}
}

func TestGetApplicationInfoEncoding(t *testing.T) {
	commonProtocol := &CommonProtocol{
		GetApplicationInfo: func(packet nex.PacketInterface) (types.List[aauser_types.ApplicationInfo], *nex.Error) {
			info := aauser_types.NewApplicationInfo()
			info.TitleID = 0x0102030405060708
			info.TitleVersion = 0x0A0B

			return types.List[aauser_types.ApplicationInfo]{info}, nil
		},
	}

	response, nexError := commonProtocol.getApplicationInfo(nil, rmctest.NewPacket(1), 1)
	if nexError != nil {
		t.Fatal(nexError)
	}

	// * count=1, Uint64 title id, Uint16 title version, all little endian
	want := []byte{1, 0, 0, 0, 8, 7, 6, 5, 4, 3, 2, 1, 0x0B, 0x0A}
	if !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}
}

func TestRegisterAndUnregisterHaveNoBody(t *testing.T) {
	var registered, unregistered types.UInt64

	commonProtocol := &CommonProtocol{
		RegisterApplication:   func(packet nex.PacketInterface, titleID types.UInt64) *nex.Error { registered = titleID; return nil },
		UnregisterApplication: func(packet nex.PacketInterface, titleID types.UInt64) *nex.Error { unregistered = titleID; return nil },
	}

	for _, response := range []*nex.RMCMessage{
		must(t)(commonProtocol.registerApplication(nil, rmctest.NewPacket(1), 1, 10)),
		must(t)(commonProtocol.unregisterApplication(nil, rmctest.NewPacket(1), 1, 11)),
	} {
		if len(response.Parameters) != 0 {
			t.Fatalf("unexpected body %x", response.Parameters)
		}
	}

	if registered != 10 || unregistered != 11 {
		t.Fatalf("registered=%d unregistered=%d", registered, unregistered)
	}
}

func must(t *testing.T) func(*nex.RMCMessage, *nex.Error) *nex.RMCMessage {
	return func(response *nex.RMCMessage, nexError *nex.Error) *nex.RMCMessage {
		if nexError != nil {
			t.Fatal(nexError)
		}

		return response
	}
}
