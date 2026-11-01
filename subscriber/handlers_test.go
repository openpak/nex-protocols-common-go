package subscriber

import (
	"bytes"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
	subscriber "github.com/PretendoNetwork/nex-protocols-go/v2/subscriber"
	subscriber_types "github.com/PretendoNetwork/nex-protocols-go/v2/subscriber/types"
)

func TestProtocolIDIsNotUtility(t *testing.T) {
	// * Subscriber is protocol 121. Utility is 110 (0x6E), the value an older protocols-go wrongly used here
	if subscriber.ProtocolID != 0x79 {
		t.Fatalf("Subscriber protocol ID = %#x, want 0x79", subscriber.ProtocolID)
	}
}

func TestPostContentEncodesUint64(t *testing.T) {
	commonProtocol := &CommonProtocol{
		PostContent: func(packet nex.PacketInterface, param subscriber_types.SubscriberPostContentParam) (types.UInt64, *nex.Error) {
			return 0x0102030405060708, nil
		},
	}

	response, nexError := commonProtocol.postContent(nil, rmctest.NewPacket(1), 1, subscriber_types.NewSubscriberPostContentParam())
	if nexError != nil {
		t.Fatal(nexError)
	}

	if want := []byte{8, 7, 6, 5, 4, 3, 2, 1}; !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}

	if response.ProtocolID != subscriber.ProtocolID || response.MethodID != subscriber.MethodPostContent {
		t.Fatalf("bad header: %+v", response)
	}
}

func TestGetUserStatusesEncoding(t *testing.T) {
	commonProtocol := &CommonProtocol{
		GetUserStatuses: func(packet nex.PacketInterface, pids types.List[types.PID], keys types.List[types.UInt8]) (types.List[subscriber_types.SubscriberUserStatusInfo], *nex.Error) {
			info := subscriber_types.NewSubscriberUserStatusInfo()
			info.PID = pids[0]
			info.Values = types.List[types.QBuffer]{types.NewQBuffer([]byte{1, 2})}

			return types.List[subscriber_types.SubscriberUserStatusInfo]{info}, nil
		},
	}

	response, nexError := commonProtocol.getUserStatuses(nil, rmctest.NewPacket(1), 1, types.List[types.PID]{9}, types.List[types.UInt8]{0})
	if nexError != nil {
		t.Fatal(nexError)
	}

	// * count=1, PID (4 bytes), inner count=1, qBuffer (u16 length + bytes)
	want := []byte{1, 0, 0, 0, 9, 0, 0, 0, 1, 0, 0, 0, 2, 0, 1, 2}
	if !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}
}

func TestGetContentMultiEncodesNestedEmptyLists(t *testing.T) {
	commonProtocol := &CommonProtocol{
		GetContentMulti: func(packet nex.PacketInterface, params types.List[subscriber_types.SubscriberGetContentParam]) (types.List[types.List[subscriber_types.SubscriberContent]], *nex.Error) {
			return types.List[types.List[subscriber_types.SubscriberContent]]{types.NewList[subscriber_types.SubscriberContent](), types.NewList[subscriber_types.SubscriberContent]()}, nil
		},
	}

	response, nexError := commonProtocol.getContentMulti(nil, rmctest.NewPacket(1), 1, nil)
	if nexError != nil {
		t.Fatal(nexError)
	}

	if want := []byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}; !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}
}

func TestRawMethodsPassPayloadThrough(t *testing.T) {
	commonProtocol := &CommonProtocol{
		Follow: func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error) {
			return append([]byte{0xAA}, packetPayload...), nil
		},
		Unfollow: func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error) { return nil, nil },
	}

	response, nexError := commonProtocol.follow(nil, rmctest.NewPacket(1), 1, []byte{1, 2})
	if nexError != nil || !bytes.Equal(response.Parameters, []byte{0xAA, 1, 2}) || response.MethodID != subscriber.MethodFollow {
		t.Fatalf("response=%+v err=%v", response, nexError)
	}

	response, nexError = commonProtocol.unfollow(nil, rmctest.NewPacket(1), 1, nil)
	if nexError != nil || len(response.Parameters) != 0 {
		t.Fatalf("response=%+v err=%v", response, nexError)
	}
}

func TestHelloHasNoBody(t *testing.T) {
	commonProtocol := &CommonProtocol{
		Hello: func(packet nex.PacketInterface, unknown types.String) *nex.Error { return nil },
	}

	response, nexError := commonProtocol.hello(nil, rmctest.NewPacket(1), 1, "hi")
	if nexError != nil || len(response.Parameters) != 0 {
		t.Fatalf("response=%+v err=%v", response, nexError)
	}
}
