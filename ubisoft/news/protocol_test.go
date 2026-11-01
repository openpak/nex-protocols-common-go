package news

import (
	"bytes"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/news"
	news_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/news/types"
)

type fakeEndpoint struct{ nex.EndpointInterface }

func (fakeEndpoint) LibraryVersions() *nex.LibraryVersions {
	versions := nex.NewLibraryVersions()
	versions.Main = nex.NewLibraryVersion(3, 0, 0)

	return versions
}

func (fakeEndpoint) ByteStreamSettings() *nex.ByteStreamSettings {
	return nex.NewByteStreamSettings()
}

type fakeConnection struct{ nex.ConnectionInterface }

func (fakeConnection) Endpoint() nex.EndpointInterface { return fakeEndpoint{} }

type fakePacket struct{ nex.PacketInterface }

func (fakePacket) Sender() nex.ConnectionInterface { return fakeConnection{} }

func checkMessage(t *testing.T, message *nex.RMCMessage, nexErr *nex.Error, methodID uint32, callID uint32, want []byte) {
	t.Helper()

	if nexErr != nil {
		t.Fatalf("unexpected error: %s", nexErr.Error())
	}

	if message == nil {
		t.Fatal("nil message")
	}

	if message.ProtocolID != ubi.ProtocolID || message.MethodID != methodID || message.CallID != callID || !message.IsSuccess {
		t.Fatalf("bad message header: protocol %d method %d call %d success %v", message.ProtocolID, message.MethodID, message.CallID, message.IsSuccess)
	}

	if !bytes.Equal(message.Parameters, want) {
		t.Fatalf("response bytes differ\n got %x\nwant %x", message.Parameters, want)
	}
}

func checkNotImplemented(t *testing.T, message *nex.RMCMessage, nexErr *nex.Error) {
	t.Helper()

	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "").ResultCode {
		t.Fatalf("want Core::NotImplemented, got message %v error %v", message, nexErr)
	}
}

func TestGetChannels_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getChannels(nil, fakePacket{}, 7, types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)})
	checkMessage(t, message, nexErr, ubi.MethodGetChannels, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetChannels_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetChannels = func(packet nex.PacketInterface, _ types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[news_types.NewsChannel](), nil
	}

	message, nexErr := commonProtocol.getChannels(errors.New("bad"), fakePacket{}, 7, types.NewResultRange())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetChannels_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetChannels = func(packet nex.PacketInterface, got0 types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		called = true

		if !(got0 == types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)}) {
			t.Errorf("resultRange differs")
		}

		return types.List[news_types.NewsChannel]{func() news_types.NewsChannel {
			v := news_types.NewNewsChannel()
			v.ID = types.NewUInt32(4041)
			v.OwnerPID = types.NewUInt32(5050)
			v.Name = types.NewString("s6")
			v.Description = types.NewString("s7")
			v.CreationTime = types.NewDateTime(8000101)
			v.ExpirationTime = types.NewDateTime(9000104)
			v.Type = types.NewString("s10")
			v.Locale = types.NewString("s11")
			v.Subscribable = types.NewBool(false)
			return v
		}(), func() news_types.NewsChannel {
			v := news_types.NewNewsChannel()
			v.ID = types.NewUInt32(14131)
			v.OwnerPID = types.NewUInt32(15140)
			v.Name = types.NewString("s16")
			v.Description = types.NewString("s17")
			v.CreationTime = types.NewDateTime(18000131)
			v.ExpirationTime = types.NewDateTime(19000134)
			v.Type = types.NewString("s20")
			v.Locale = types.NewString("s21")
			v.Subscribable = types.NewBool(false)
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getChannels(nil, fakePacket{}, 9, types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)})
	checkMessage(t, message, nexErr, ubi.MethodGetChannels, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xc9, 0x0f, 0x00, 0x00, 0xba, 0x13, 0x00, 0x00, 0x03, 0x00, 0x73, 0x36, 0x00, 0x03, 0x00, 0x73, 0x37, 0x00, 0x65, 0x12, 0x7a, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa8, 0x54, 0x89, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x30, 0x00, 0x04, 0x00, 0x73, 0x31, 0x31, 0x00, 0x00, 0x33, 0x37, 0x00, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x36, 0x00, 0x04, 0x00, 0x73, 0x31, 0x37, 0x00, 0x03, 0xa9, 0x12, 0x01, 0x00, 0x00, 0x00, 0x00, 0x46, 0xeb, 0x21, 0x01, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x30, 0x00, 0x04, 0x00, 0x73, 0x32, 0x31, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetChannels = func(packet nex.PacketInterface, _ types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		return types.NewList[news_types.NewsChannel](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getChannels(nil, fakePacket{}, 9, types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetChannelsByTypes_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getChannelsByTypes(nil, fakePacket{}, 7, types.List[types.String]{types.NewString("s2"), types.NewString("s3")}, types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)})
	checkMessage(t, message, nexErr, ubi.MethodGetChannelsByTypes, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetChannelsByTypes_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetChannelsByTypes = func(packet nex.PacketInterface, _ types.List[types.String], _ types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[news_types.NewsChannel](), nil
	}

	message, nexErr := commonProtocol.getChannelsByTypes(errors.New("bad"), fakePacket{}, 7, types.NewList[types.String](), types.NewResultRange())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetChannelsByTypes_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetChannelsByTypes = func(packet nex.PacketInterface, got0 types.List[types.String], got1 types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[types.String]{types.NewString("s2"), types.NewString("s3")})) {
			t.Errorf("newsChannelTypes differs")
		}

		if !(got1 == types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)}) {
			t.Errorf("resultRange differs")
		}

		return types.List[news_types.NewsChannel]{func() news_types.NewsChannel {
			v := news_types.NewNewsChannel()
			v.ID = types.NewUInt32(7068)
			v.OwnerPID = types.NewUInt32(8077)
			v.Name = types.NewString("s9")
			v.Description = types.NewString("s10")
			v.CreationTime = types.NewDateTime(11000110)
			v.ExpirationTime = types.NewDateTime(12000113)
			v.Type = types.NewString("s13")
			v.Locale = types.NewString("s14")
			v.Subscribable = types.NewBool(true)
			return v
		}(), func() news_types.NewsChannel {
			v := news_types.NewNewsChannel()
			v.ID = types.NewUInt32(17158)
			v.OwnerPID = types.NewUInt32(18167)
			v.Name = types.NewString("s19")
			v.Description = types.NewString("s20")
			v.CreationTime = types.NewDateTime(21000140)
			v.ExpirationTime = types.NewDateTime(22000143)
			v.Type = types.NewString("s23")
			v.Locale = types.NewString("s24")
			v.Subscribable = types.NewBool(true)
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getChannelsByTypes(nil, fakePacket{}, 9, types.List[types.String]{types.NewString("s2"), types.NewString("s3")}, types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)})
	checkMessage(t, message, nexErr, ubi.MethodGetChannelsByTypes, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x8d, 0x1f, 0x00, 0x00, 0x03, 0x00, 0x73, 0x39, 0x00, 0x04, 0x00, 0x73, 0x31, 0x30, 0x00, 0x2e, 0xd9, 0xa7, 0x00, 0x00, 0x00, 0x00, 0x00, 0x71, 0x1b, 0xb7, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x33, 0x00, 0x04, 0x00, 0x73, 0x31, 0x34, 0x00, 0x01, 0x06, 0x43, 0x00, 0x00, 0xf7, 0x46, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x39, 0x00, 0x04, 0x00, 0x73, 0x32, 0x30, 0x00, 0xcc, 0x6f, 0x40, 0x01, 0x00, 0x00, 0x00, 0x00, 0x0f, 0xb2, 0x4f, 0x01, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x33, 0x00, 0x04, 0x00, 0x73, 0x32, 0x34, 0x00, 0x01})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetChannelsByTypes = func(packet nex.PacketInterface, _ types.List[types.String], _ types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		return types.NewList[news_types.NewsChannel](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getChannelsByTypes(nil, fakePacket{}, 9, types.List[types.String]{types.NewString("s2"), types.NewString("s3")}, types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetSubscribableChannels_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getSubscribableChannels(nil, fakePacket{}, 7, types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)})
	checkMessage(t, message, nexErr, ubi.MethodGetSubscribableChannels, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetSubscribableChannels_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetSubscribableChannels = func(packet nex.PacketInterface, _ types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[news_types.NewsChannel](), nil
	}

	message, nexErr := commonProtocol.getSubscribableChannels(errors.New("bad"), fakePacket{}, 7, types.NewResultRange())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetSubscribableChannels_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetSubscribableChannels = func(packet nex.PacketInterface, got0 types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		called = true

		if !(got0 == types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)}) {
			t.Errorf("resultRange differs")
		}

		return types.List[news_types.NewsChannel]{func() news_types.NewsChannel {
			v := news_types.NewNewsChannel()
			v.ID = types.NewUInt32(4041)
			v.OwnerPID = types.NewUInt32(5050)
			v.Name = types.NewString("s6")
			v.Description = types.NewString("s7")
			v.CreationTime = types.NewDateTime(8000101)
			v.ExpirationTime = types.NewDateTime(9000104)
			v.Type = types.NewString("s10")
			v.Locale = types.NewString("s11")
			v.Subscribable = types.NewBool(false)
			return v
		}(), func() news_types.NewsChannel {
			v := news_types.NewNewsChannel()
			v.ID = types.NewUInt32(14131)
			v.OwnerPID = types.NewUInt32(15140)
			v.Name = types.NewString("s16")
			v.Description = types.NewString("s17")
			v.CreationTime = types.NewDateTime(18000131)
			v.ExpirationTime = types.NewDateTime(19000134)
			v.Type = types.NewString("s20")
			v.Locale = types.NewString("s21")
			v.Subscribable = types.NewBool(false)
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getSubscribableChannels(nil, fakePacket{}, 9, types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)})
	checkMessage(t, message, nexErr, ubi.MethodGetSubscribableChannels, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xc9, 0x0f, 0x00, 0x00, 0xba, 0x13, 0x00, 0x00, 0x03, 0x00, 0x73, 0x36, 0x00, 0x03, 0x00, 0x73, 0x37, 0x00, 0x65, 0x12, 0x7a, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa8, 0x54, 0x89, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x30, 0x00, 0x04, 0x00, 0x73, 0x31, 0x31, 0x00, 0x00, 0x33, 0x37, 0x00, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x36, 0x00, 0x04, 0x00, 0x73, 0x31, 0x37, 0x00, 0x03, 0xa9, 0x12, 0x01, 0x00, 0x00, 0x00, 0x00, 0x46, 0xeb, 0x21, 0x01, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x30, 0x00, 0x04, 0x00, 0x73, 0x32, 0x31, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetSubscribableChannels = func(packet nex.PacketInterface, _ types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		return types.NewList[news_types.NewsChannel](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getSubscribableChannels(nil, fakePacket{}, 9, types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetChannelsByIDs_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getChannelsByIDs(nil, fakePacket{}, 7, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	checkMessage(t, message, nexErr, ubi.MethodGetChannelsByIDs, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetChannelsByIDs_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetChannelsByIDs = func(packet nex.PacketInterface, _ types.List[types.UInt32]) (types.List[news_types.NewsChannel], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[news_types.NewsChannel](), nil
	}

	message, nexErr := commonProtocol.getChannelsByIDs(errors.New("bad"), fakePacket{}, 7, types.NewList[types.UInt32]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetChannelsByIDs_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetChannelsByIDs = func(packet nex.PacketInterface, got0 types.List[types.UInt32]) (types.List[news_types.NewsChannel], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})) {
			t.Errorf("newsChannelIDs differs")
		}

		return types.List[news_types.NewsChannel]{func() news_types.NewsChannel {
			v := news_types.NewNewsChannel()
			v.ID = types.NewUInt32(6059)
			v.OwnerPID = types.NewUInt32(7068)
			v.Name = types.NewString("s8")
			v.Description = types.NewString("s9")
			v.CreationTime = types.NewDateTime(10000107)
			v.ExpirationTime = types.NewDateTime(11000110)
			v.Type = types.NewString("s12")
			v.Locale = types.NewString("s13")
			v.Subscribable = types.NewBool(false)
			return v
		}(), func() news_types.NewsChannel {
			v := news_types.NewNewsChannel()
			v.ID = types.NewUInt32(16149)
			v.OwnerPID = types.NewUInt32(17158)
			v.Name = types.NewString("s18")
			v.Description = types.NewString("s19")
			v.CreationTime = types.NewDateTime(20000137)
			v.ExpirationTime = types.NewDateTime(21000140)
			v.Type = types.NewString("s22")
			v.Locale = types.NewString("s23")
			v.Subscribable = types.NewBool(false)
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getChannelsByIDs(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	checkMessage(t, message, nexErr, ubi.MethodGetChannelsByIDs, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xab, 0x17, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x03, 0x00, 0x73, 0x38, 0x00, 0x03, 0x00, 0x73, 0x39, 0x00, 0xeb, 0x96, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0xd9, 0xa7, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x32, 0x00, 0x04, 0x00, 0x73, 0x31, 0x33, 0x00, 0x00, 0x15, 0x3f, 0x00, 0x00, 0x06, 0x43, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x38, 0x00, 0x04, 0x00, 0x73, 0x31, 0x39, 0x00, 0x89, 0x2d, 0x31, 0x01, 0x00, 0x00, 0x00, 0x00, 0xcc, 0x6f, 0x40, 0x01, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x32, 0x00, 0x04, 0x00, 0x73, 0x32, 0x33, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetChannelsByIDs = func(packet nex.PacketInterface, _ types.List[types.UInt32]) (types.List[news_types.NewsChannel], *nex.Error) {
		return types.NewList[news_types.NewsChannel](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getChannelsByIDs(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetSubscribedChannels_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getSubscribedChannels(nil, fakePacket{}, 7, types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)})
	checkMessage(t, message, nexErr, ubi.MethodGetSubscribedChannels, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetSubscribedChannels_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetSubscribedChannels = func(packet nex.PacketInterface, _ types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[news_types.NewsChannel](), nil
	}

	message, nexErr := commonProtocol.getSubscribedChannels(errors.New("bad"), fakePacket{}, 7, types.NewResultRange())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetSubscribedChannels_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetSubscribedChannels = func(packet nex.PacketInterface, got0 types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		called = true

		if !(got0 == types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)}) {
			t.Errorf("resultRange differs")
		}

		return types.List[news_types.NewsChannel]{func() news_types.NewsChannel {
			v := news_types.NewNewsChannel()
			v.ID = types.NewUInt32(4041)
			v.OwnerPID = types.NewUInt32(5050)
			v.Name = types.NewString("s6")
			v.Description = types.NewString("s7")
			v.CreationTime = types.NewDateTime(8000101)
			v.ExpirationTime = types.NewDateTime(9000104)
			v.Type = types.NewString("s10")
			v.Locale = types.NewString("s11")
			v.Subscribable = types.NewBool(false)
			return v
		}(), func() news_types.NewsChannel {
			v := news_types.NewNewsChannel()
			v.ID = types.NewUInt32(14131)
			v.OwnerPID = types.NewUInt32(15140)
			v.Name = types.NewString("s16")
			v.Description = types.NewString("s17")
			v.CreationTime = types.NewDateTime(18000131)
			v.ExpirationTime = types.NewDateTime(19000134)
			v.Type = types.NewString("s20")
			v.Locale = types.NewString("s21")
			v.Subscribable = types.NewBool(false)
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getSubscribedChannels(nil, fakePacket{}, 9, types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)})
	checkMessage(t, message, nexErr, ubi.MethodGetSubscribedChannels, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xc9, 0x0f, 0x00, 0x00, 0xba, 0x13, 0x00, 0x00, 0x03, 0x00, 0x73, 0x36, 0x00, 0x03, 0x00, 0x73, 0x37, 0x00, 0x65, 0x12, 0x7a, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa8, 0x54, 0x89, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x30, 0x00, 0x04, 0x00, 0x73, 0x31, 0x31, 0x00, 0x00, 0x33, 0x37, 0x00, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x36, 0x00, 0x04, 0x00, 0x73, 0x31, 0x37, 0x00, 0x03, 0xa9, 0x12, 0x01, 0x00, 0x00, 0x00, 0x00, 0x46, 0xeb, 0x21, 0x01, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x30, 0x00, 0x04, 0x00, 0x73, 0x32, 0x31, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetSubscribedChannels = func(packet nex.PacketInterface, _ types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error) {
		return types.NewList[news_types.NewsChannel](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getSubscribedChannels(nil, fakePacket{}, 9, types.ResultRange{Offset: types.NewUInt32(3), Length: types.NewUInt32(13)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSubscribeChannel_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.subscribeChannel(nil, fakePacket{}, 7, types.NewUInt32(1014))
	checkNotImplemented(t, message, nexErr)
}

func TestSubscribeChannel_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SubscribeChannel = func(packet nex.PacketInterface, _ types.UInt32) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.subscribeChannel(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSubscribeChannel_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SubscribeChannel = func(packet nex.PacketInterface, got0 types.UInt32) *nex.Error {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("newsChannelID differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.subscribeChannel(nil, fakePacket{}, 9, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodSubscribeChannel, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SubscribeChannel = func(packet nex.PacketInterface, _ types.UInt32) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.subscribeChannel(nil, fakePacket{}, 9, types.NewUInt32(1014))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestUnsubscribeChannel_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.unsubscribeChannel(nil, fakePacket{}, 7, types.NewUInt32(1014))
	checkNotImplemented(t, message, nexErr)
}

func TestUnsubscribeChannel_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.UnsubscribeChannel = func(packet nex.PacketInterface, _ types.UInt32) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.unsubscribeChannel(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestUnsubscribeChannel_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.UnsubscribeChannel = func(packet nex.PacketInterface, got0 types.UInt32) *nex.Error {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("newsChannelID differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.unsubscribeChannel(nil, fakePacket{}, 9, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodUnsubscribeChannel, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.UnsubscribeChannel = func(packet nex.PacketInterface, _ types.UInt32) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.unsubscribeChannel(nil, fakePacket{}, 9, types.NewUInt32(1014))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetNewsHeaders_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getNewsHeaders(nil, fakePacket{}, 7, func() news_types.NewsRecipient {
		v := news_types.NewNewsRecipient()
		v.RecipientID = types.NewUInt32(2023)
		v.RecipientType = types.NewUInt32(3032)
		return v
	}(), types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)})
	checkMessage(t, message, nexErr, ubi.MethodGetNewsHeaders, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetNewsHeaders_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetNewsHeaders = func(packet nex.PacketInterface, _ news_types.NewsRecipient, _ types.ResultRange) (types.List[news_types.NewsHeader], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[news_types.NewsHeader](), nil
	}

	message, nexErr := commonProtocol.getNewsHeaders(errors.New("bad"), fakePacket{}, 7, news_types.NewNewsRecipient(), types.NewResultRange())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetNewsHeaders_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetNewsHeaders = func(packet nex.PacketInterface, got0 news_types.NewsRecipient, got1 types.ResultRange) (types.List[news_types.NewsHeader], *nex.Error) {
		called = true

		if !(got0.Equals(func() news_types.NewsRecipient {
			v := news_types.NewNewsRecipient()
			v.RecipientID = types.NewUInt32(2023)
			v.RecipientType = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("recipient differs")
		}

		if !(got1 == types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)}) {
			t.Errorf("range differs")
		}

		return types.List[news_types.NewsHeader]{func() news_types.NewsHeader {
			v := news_types.NewNewsHeader()
			v.ID = types.NewUInt32(7068)
			v.RecipientID = types.NewUInt32(8077)
			v.RecipientType = types.NewUInt32(9086)
			v.PublisherPID = types.NewUInt32(10095)
			v.PublisherName = types.NewString("s11")
			v.PublicationTime = types.NewDateTime(12000113)
			v.DisplayTime = types.NewDateTime(13000116)
			v.ExpirationTime = types.NewDateTime(14000119)
			v.Title = types.NewString("s15")
			v.Link = types.NewString("s16")
			return v
		}(), func() news_types.NewsHeader {
			v := news_types.NewNewsHeader()
			v.ID = types.NewUInt32(18167)
			v.RecipientID = types.NewUInt32(19176)
			v.RecipientType = types.NewUInt32(20185)
			v.PublisherPID = types.NewUInt32(21194)
			v.PublisherName = types.NewString("s22")
			v.PublicationTime = types.NewDateTime(23000146)
			v.DisplayTime = types.NewDateTime(24000149)
			v.ExpirationTime = types.NewDateTime(25000152)
			v.Title = types.NewString("s26")
			v.Link = types.NewString("s27")
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getNewsHeaders(nil, fakePacket{}, 9, func() news_types.NewsRecipient {
		v := news_types.NewNewsRecipient()
		v.RecipientID = types.NewUInt32(2023)
		v.RecipientType = types.NewUInt32(3032)
		return v
	}(), types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)})
	checkMessage(t, message, nexErr, ubi.MethodGetNewsHeaders, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x8d, 0x1f, 0x00, 0x00, 0x7e, 0x23, 0x00, 0x00, 0x6f, 0x27, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x31, 0x00, 0x71, 0x1b, 0xb7, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb4, 0x5d, 0xc6, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf7, 0x9f, 0xd5, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x35, 0x00, 0x04, 0x00, 0x73, 0x31, 0x36, 0x00, 0xf7, 0x46, 0x00, 0x00, 0xe8, 0x4a, 0x00, 0x00, 0xd9, 0x4e, 0x00, 0x00, 0xca, 0x52, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x32, 0x00, 0x52, 0xf4, 0x5e, 0x01, 0x00, 0x00, 0x00, 0x00, 0x95, 0x36, 0x6e, 0x01, 0x00, 0x00, 0x00, 0x00, 0xd8, 0x78, 0x7d, 0x01, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x36, 0x00, 0x04, 0x00, 0x73, 0x32, 0x37, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetNewsHeaders = func(packet nex.PacketInterface, _ news_types.NewsRecipient, _ types.ResultRange) (types.List[news_types.NewsHeader], *nex.Error) {
		return types.NewList[news_types.NewsHeader](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getNewsHeaders(nil, fakePacket{}, 9, func() news_types.NewsRecipient {
		v := news_types.NewNewsRecipient()
		v.RecipientID = types.NewUInt32(2023)
		v.RecipientType = types.NewUInt32(3032)
		return v
	}(), types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetNewsMessages_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getNewsMessages(nil, fakePacket{}, 7, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	checkMessage(t, message, nexErr, ubi.MethodGetNewsMessages, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetNewsMessages_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetNewsMessages = func(packet nex.PacketInterface, _ types.List[types.UInt32]) (types.List[news_types.NewsMessage], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[news_types.NewsMessage](), nil
	}

	message, nexErr := commonProtocol.getNewsMessages(errors.New("bad"), fakePacket{}, 7, types.NewList[types.UInt32]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetNewsMessages_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetNewsMessages = func(packet nex.PacketInterface, got0 types.List[types.UInt32]) (types.List[news_types.NewsMessage], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})) {
			t.Errorf("newsMessageIDs differs")
		}

		return types.List[news_types.NewsMessage]{func() news_types.NewsMessage {
			v := news_types.NewNewsMessage()
			v.Body = types.NewString("s6")
			return v
		}(), func() news_types.NewsMessage {
			v := news_types.NewNewsMessage()
			v.Body = types.NewString("s8")
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getNewsMessages(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	checkMessage(t, message, nexErr, ubi.MethodGetNewsMessages, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x73, 0x36, 0x00, 0x03, 0x00, 0x73, 0x38, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetNewsMessages = func(packet nex.PacketInterface, _ types.List[types.UInt32]) (types.List[news_types.NewsMessage], *nex.Error) {
		return types.NewList[news_types.NewsMessage](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getNewsMessages(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetNumberOfNews_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getNumberOfNews(nil, fakePacket{}, 7, func() news_types.NewsRecipient {
		v := news_types.NewNewsRecipient()
		v.RecipientID = types.NewUInt32(2023)
		v.RecipientType = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetNumberOfNews, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetNumberOfNews_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetNumberOfNews = func(packet nex.PacketInterface, _ news_types.NewsRecipient) (types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.getNumberOfNews(errors.New("bad"), fakePacket{}, 7, news_types.NewNewsRecipient())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetNumberOfNews_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetNumberOfNews = func(packet nex.PacketInterface, got0 news_types.NewsRecipient) (types.UInt32, *nex.Error) {
		called = true

		if !(got0.Equals(func() news_types.NewsRecipient {
			v := news_types.NewNewsRecipient()
			v.RecipientID = types.NewUInt32(2023)
			v.RecipientType = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("recipient differs")
		}

		return types.NewUInt32(4041), nil
	}

	message, nexErr := commonProtocol.getNumberOfNews(nil, fakePacket{}, 9, func() news_types.NewsRecipient {
		v := news_types.NewNewsRecipient()
		v.RecipientID = types.NewUInt32(2023)
		v.RecipientType = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetNumberOfNews, 9, []byte{0xc9, 0x0f, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetNumberOfNews = func(packet nex.PacketInterface, _ news_types.NewsRecipient) (types.UInt32, *nex.Error) {
		return types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getNumberOfNews(nil, fakePacket{}, 9, func() news_types.NewsRecipient {
		v := news_types.NewNewsRecipient()
		v.RecipientID = types.NewUInt32(2023)
		v.RecipientType = types.NewUInt32(3032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetChannelByType_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getChannelByType(nil, fakePacket{}, 7, types.NewString("s1"))
	checkNotImplemented(t, message, nexErr)
}

func TestGetChannelByType_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetChannelByType = func(packet nex.PacketInterface, _ types.String) (news_types.NewsChannel, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return news_types.NewNewsChannel(), nil
	}

	message, nexErr := commonProtocol.getChannelByType(errors.New("bad"), fakePacket{}, 7, types.NewString(""))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetChannelByType_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetChannelByType = func(packet nex.PacketInterface, got0 types.String) (news_types.NewsChannel, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewString("s1"))) {
			t.Errorf("newsChannelType differs")
		}

		return func() news_types.NewsChannel {
			v := news_types.NewNewsChannel()
			v.ID = types.NewUInt32(3032)
			v.OwnerPID = types.NewUInt32(4041)
			v.Name = types.NewString("s5")
			v.Description = types.NewString("s6")
			v.CreationTime = types.NewDateTime(7000098)
			v.ExpirationTime = types.NewDateTime(8000101)
			v.Type = types.NewString("s9")
			v.Locale = types.NewString("s10")
			v.Subscribable = types.NewBool(true)
			return v
		}(), nil
	}

	message, nexErr := commonProtocol.getChannelByType(nil, fakePacket{}, 9, types.NewString("s1"))
	checkMessage(t, message, nexErr, ubi.MethodGetChannelByType, 9, []byte{0xd8, 0x0b, 0x00, 0x00, 0xc9, 0x0f, 0x00, 0x00, 0x03, 0x00, 0x73, 0x35, 0x00, 0x03, 0x00, 0x73, 0x36, 0x00, 0x22, 0xd0, 0x6a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x65, 0x12, 0x7a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x73, 0x39, 0x00, 0x04, 0x00, 0x73, 0x31, 0x30, 0x00, 0x01})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetChannelByType = func(packet nex.PacketInterface, _ types.String) (news_types.NewsChannel, *nex.Error) {
		return news_types.NewNewsChannel(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getChannelByType(nil, fakePacket{}, 9, types.NewString("s1"))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetNewsHeadersByType_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getNewsHeadersByType(nil, fakePacket{}, 7, types.NewString("s1"), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	checkMessage(t, message, nexErr, ubi.MethodGetNewsHeadersByType, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetNewsHeadersByType_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetNewsHeadersByType = func(packet nex.PacketInterface, _ types.String, _ types.ResultRange) (types.List[news_types.NewsHeader], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[news_types.NewsHeader](), nil
	}

	message, nexErr := commonProtocol.getNewsHeadersByType(errors.New("bad"), fakePacket{}, 7, types.NewString(""), types.NewResultRange())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetNewsHeadersByType_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetNewsHeadersByType = func(packet nex.PacketInterface, got0 types.String, got1 types.ResultRange) (types.List[news_types.NewsHeader], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewString("s1"))) {
			t.Errorf("newsChannelType differs")
		}

		if !(got1 == types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)}) {
			t.Errorf("range differs")
		}

		return types.List[news_types.NewsHeader]{func() news_types.NewsHeader {
			v := news_types.NewNewsHeader()
			v.ID = types.NewUInt32(5050)
			v.RecipientID = types.NewUInt32(6059)
			v.RecipientType = types.NewUInt32(7068)
			v.PublisherPID = types.NewUInt32(8077)
			v.PublisherName = types.NewString("s9")
			v.PublicationTime = types.NewDateTime(10000107)
			v.DisplayTime = types.NewDateTime(11000110)
			v.ExpirationTime = types.NewDateTime(12000113)
			v.Title = types.NewString("s13")
			v.Link = types.NewString("s14")
			return v
		}(), func() news_types.NewsHeader {
			v := news_types.NewNewsHeader()
			v.ID = types.NewUInt32(16149)
			v.RecipientID = types.NewUInt32(17158)
			v.RecipientType = types.NewUInt32(18167)
			v.PublisherPID = types.NewUInt32(19176)
			v.PublisherName = types.NewString("s20")
			v.PublicationTime = types.NewDateTime(21000140)
			v.DisplayTime = types.NewDateTime(22000143)
			v.ExpirationTime = types.NewDateTime(23000146)
			v.Title = types.NewString("s24")
			v.Link = types.NewString("s25")
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getNewsHeadersByType(nil, fakePacket{}, 9, types.NewString("s1"), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	checkMessage(t, message, nexErr, ubi.MethodGetNewsHeadersByType, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xba, 0x13, 0x00, 0x00, 0xab, 0x17, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x8d, 0x1f, 0x00, 0x00, 0x03, 0x00, 0x73, 0x39, 0x00, 0xeb, 0x96, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0xd9, 0xa7, 0x00, 0x00, 0x00, 0x00, 0x00, 0x71, 0x1b, 0xb7, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x33, 0x00, 0x04, 0x00, 0x73, 0x31, 0x34, 0x00, 0x15, 0x3f, 0x00, 0x00, 0x06, 0x43, 0x00, 0x00, 0xf7, 0x46, 0x00, 0x00, 0xe8, 0x4a, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x30, 0x00, 0xcc, 0x6f, 0x40, 0x01, 0x00, 0x00, 0x00, 0x00, 0x0f, 0xb2, 0x4f, 0x01, 0x00, 0x00, 0x00, 0x00, 0x52, 0xf4, 0x5e, 0x01, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x34, 0x00, 0x04, 0x00, 0x73, 0x32, 0x35, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetNewsHeadersByType = func(packet nex.PacketInterface, _ types.String, _ types.ResultRange) (types.List[news_types.NewsHeader], *nex.Error) {
		return types.NewList[news_types.NewsHeader](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getNewsHeadersByType(nil, fakePacket{}, 9, types.NewString("s1"), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetNewsMessagesByType_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getNewsMessagesByType(nil, fakePacket{}, 7, types.NewString("s1"), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	checkMessage(t, message, nexErr, ubi.MethodGetNewsMessagesByType, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetNewsMessagesByType_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetNewsMessagesByType = func(packet nex.PacketInterface, _ types.String, _ types.ResultRange) (types.List[news_types.NewsMessage], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[news_types.NewsMessage](), nil
	}

	message, nexErr := commonProtocol.getNewsMessagesByType(errors.New("bad"), fakePacket{}, 7, types.NewString(""), types.NewResultRange())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetNewsMessagesByType_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetNewsMessagesByType = func(packet nex.PacketInterface, got0 types.String, got1 types.ResultRange) (types.List[news_types.NewsMessage], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewString("s1"))) {
			t.Errorf("newsChannelType differs")
		}

		if !(got1 == types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)}) {
			t.Errorf("range differs")
		}

		return types.List[news_types.NewsMessage]{func() news_types.NewsMessage {
			v := news_types.NewNewsMessage()
			v.Body = types.NewString("s5")
			return v
		}(), func() news_types.NewsMessage {
			v := news_types.NewNewsMessage()
			v.Body = types.NewString("s7")
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getNewsMessagesByType(nil, fakePacket{}, 9, types.NewString("s1"), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	checkMessage(t, message, nexErr, ubi.MethodGetNewsMessagesByType, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x73, 0x35, 0x00, 0x03, 0x00, 0x73, 0x37, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetNewsMessagesByType = func(packet nex.PacketInterface, _ types.String, _ types.ResultRange) (types.List[news_types.NewsMessage], *nex.Error) {
		return types.NewList[news_types.NewsMessage](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getNewsMessagesByType(nil, fakePacket{}, 9, types.NewString("s1"), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}
