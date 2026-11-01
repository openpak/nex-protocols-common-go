package friends

import (
	"bytes"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
	friends "github.com/PretendoNetwork/nex-protocols-go/v2/friends"
	friends_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends/types"
)

func TestAddFriendEncodesBool(t *testing.T) {
	commonProtocol := &CommonProtocol{
		AddFriend: func(packet nex.PacketInterface, uiPlayer types.UInt32, uiDetails types.UInt32, strMessage types.String) (types.Bool, *nex.Error) {
			return uiPlayer == 5, nil
		},
	}

	response, nexError := commonProtocol.addFriend(nil, rmctest.NewPacket(1), 6, 5, 0, "hi")
	if nexError != nil {
		t.Fatal(nexError)
	}

	if want := []byte{1}; !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}

	if response.ProtocolID != friends.ProtocolID || response.MethodID != friends.MethodAddFriend || response.CallID != 6 {
		t.Fatalf("bad header: %+v", response)
	}
}

func TestGetListEncoding(t *testing.T) {
	commonProtocol := &CommonProtocol{
		GetList: func(packet nex.PacketInterface, byRelationship types.UInt8, bReversed types.Bool) (types.List[types.UInt32], *nex.Error) {
			return types.List[types.UInt32]{5, 6}, nil
		},
	}

	response, nexError := commonProtocol.getList(nil, rmctest.NewPacket(1), 1, 3, false)
	if nexError != nil {
		t.Fatal(nexError)
	}

	if want := []byte{2, 0, 0, 0, 5, 0, 0, 0, 6, 0, 0, 0}; !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}
}

func TestAddFriendWithDetailsEncodesRelationshipData(t *testing.T) {
	commonProtocol := &CommonProtocol{
		AddFriendWithDetails: func(packet nex.PacketInterface, uiPlayer types.UInt32, uiDetails types.UInt32, strMessage types.String) (friends_types.RelationshipData, *nex.Error) {
			data := friends_types.NewRelationshipData()
			data.PID = 1
			data.StrName = "n"
			data.ByRelationship = 2
			data.UIDetails = 3
			data.ByStatus = 4

			return data, nil
		},
	}

	response, nexError := commonProtocol.addFriendWithDetails(nil, rmctest.NewPacket(1), 1, 1, 3, "")
	if nexError != nil {
		t.Fatal(nexError)
	}

	// * m_pid u32, m_strName String, m_byRelationship u8, m_uiDetails u32, m_byStatus u8
	want := []byte{1, 0, 0, 0, 2, 0, 'n', 0, 2, 3, 0, 0, 0, 4}
	if !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}
}

func TestGetRelationshipsEncodesTotalThenList(t *testing.T) {
	commonProtocol := &CommonProtocol{
		GetRelationships: func(packet nex.PacketInterface, resultRange types.ResultRange) (types.UInt32, types.List[friends_types.RelationshipData], *nex.Error) {
			return 3, types.NewList[friends_types.RelationshipData](), nil
		},
	}

	response, nexError := commonProtocol.getRelationships(nil, rmctest.NewPacket(1), 1, types.NewResultRange())
	if nexError != nil {
		t.Fatal(nexError)
	}

	if want := []byte{3, 0, 0, 0, 0, 0, 0, 0}; !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}
}
