package persistentstore

import (
	"bytes"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
	persistent_store "github.com/PretendoNetwork/nex-protocols-go/v2/persistent-store"
)

func TestFindByGroupEncoding(t *testing.T) {
	commonProtocol := &CommonProtocol{
		FindByGroup: func(packet nex.PacketInterface, uiGroup types.UInt32) (types.List[types.String], *nex.Error) {
			if uiGroup != 9 {
				t.Errorf("uiGroup = %d", uiGroup)
			}

			return types.List[types.String]{"a", "b"}, nil
		},
	}

	response, nexError := commonProtocol.findByGroup(nil, rmctest.NewPacket(1), 11, 9)
	if nexError != nil {
		t.Fatal(nexError)
	}

	// * List<String>: u32 count, then u16 length (including the NUL) + bytes + NUL for each
	want := []byte{2, 0, 0, 0, 2, 0, 'a', 0, 2, 0, 'b', 0}
	if !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}

	if response.ProtocolID != persistent_store.ProtocolID || response.MethodID != persistent_store.MethodFindByGroup || response.CallID != 11 {
		t.Fatalf("bad header: %+v", response)
	}
}

func TestGetItemEncoding(t *testing.T) {
	commonProtocol := &CommonProtocol{
		GetItem: func(packet nex.PacketInterface, uiGroup types.UInt32, strTag types.String) (types.Buffer, types.Bool, *nex.Error) {
			return types.NewBuffer([]byte{1, 2}), true, nil
		},
	}

	response, nexError := commonProtocol.getItem(nil, rmctest.NewPacket(1), 1, 1, "t")
	if nexError != nil {
		t.Fatal(nexError)
	}

	// * Buffer (u32 length + bytes) then the Bool result
	want := []byte{2, 0, 0, 0, 1, 2, 1}
	if !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}
}

func TestInsertCustomItemHasNoBody(t *testing.T) {
	called := false
	commonProtocol := &CommonProtocol{
		InsertCustomItem: func(packet nex.PacketInterface, uiGroup types.UInt32, strTag types.String, hData types.DataHolder, bReplace types.Bool) *nex.Error {
			called = uiGroup == 3 && strTag == "k" && bool(bReplace)
			return nil
		},
	}

	response, nexError := commonProtocol.insertCustomItem(nil, rmctest.NewPacket(1), 1, 3, "k", types.NewDataHolder(), true)
	if nexError != nil || !called || len(response.Parameters) != 0 {
		t.Fatalf("called=%v err=%v response=%+v", called, nexError, response)
	}
}

func TestHookErrorIsReturned(t *testing.T) {
	commonProtocol := &CommonProtocol{
		RemoveItem: func(packet nex.PacketInterface, uiGroup types.UInt32, strTag types.String) (types.Bool, *nex.Error) {
			return false, nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
		},
	}

	response, nexError := commonProtocol.removeItem(nil, rmctest.NewPacket(1), 1, 1, "t")
	if response != nil || !rmctest.HasCode(nexError, nex.ResultCodes.Core.Unknown) {
		t.Fatalf("expected the hook error, got %v", nexError)
	}
}
