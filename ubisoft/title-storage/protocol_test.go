package titlestorage

import (
	"bytes"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/title-storage"
	title_storage_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/title-storage/types"
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

func TestBrowseContents_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.browseContents(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewUInt8(15), types.NewUInt32(3032), types.NewUInt8(29))
	checkMessage(t, message, nexErr, ubi.MethodBrowseContents, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestBrowseContents_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.BrowseContents = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt8, _ types.UInt32, _ types.UInt8) (types.List[title_storage_types.TitleContent], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[title_storage_types.TitleContent](), nil
	}

	message, nexErr := commonProtocol.browseContents(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewUInt8(0), types.NewUInt32(0), types.NewUInt8(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestBrowseContents_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.BrowseContents = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.UInt8, got2 types.UInt32, got3 types.UInt8) (types.List[title_storage_types.TitleContent], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("flags differs")
		}

		if !(got1.Equals(types.NewUInt8(15))) {
			t.Errorf("flagsMaskType differs")
		}

		if !(got2.Equals(types.NewUInt32(3032))) {
			t.Errorf("extendedFlags differs")
		}

		if !(got3.Equals(types.NewUInt8(29))) {
			t.Errorf("extendedFlagsMaskType differs")
		}

		return types.List[title_storage_types.TitleContent]{func() title_storage_types.TitleContent {
			v := title_storage_types.NewTitleContent()
			v.Id = types.NewUInt32(7068)
			v.Name = types.NewString("s8")
			v.Description = types.NewString("s9")
			v.Size = types.NewInt64(-40000000071)
			v.ContentType = types.NewString("s11")
			v.Flags = types.NewUInt32(12113)
			v.ExtendedFlags = types.NewUInt32(13122)
			return v
		}(), func() title_storage_types.TitleContent {
			v := title_storage_types.NewTitleContent()
			v.Id = types.NewUInt32(15140)
			v.Name = types.NewString("s16")
			v.Description = types.NewString("s17")
			v.Size = types.NewInt64(-72000000127)
			v.ContentType = types.NewString("s19")
			v.Flags = types.NewUInt32(20185)
			v.ExtendedFlags = types.NewUInt32(21194)
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.browseContents(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt8(15), types.NewUInt32(3032), types.NewUInt8(29))
	checkMessage(t, message, nexErr, ubi.MethodBrowseContents, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x03, 0x00, 0x73, 0x38, 0x00, 0x03, 0x00, 0x73, 0x39, 0x00, 0xb9, 0x6f, 0xd0, 0xaf, 0xf6, 0xff, 0xff, 0xff, 0x04, 0x00, 0x73, 0x31, 0x31, 0x00, 0x51, 0x2f, 0x00, 0x00, 0x42, 0x33, 0x00, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x36, 0x00, 0x04, 0x00, 0x73, 0x31, 0x37, 0x00, 0x81, 0x2f, 0x77, 0x3c, 0xef, 0xff, 0xff, 0xff, 0x04, 0x00, 0x73, 0x31, 0x39, 0x00, 0xd9, 0x4e, 0x00, 0x00, 0xca, 0x52, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.BrowseContents = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt8, _ types.UInt32, _ types.UInt8) (types.List[title_storage_types.TitleContent], *nex.Error) {
		return types.NewList[title_storage_types.TitleContent](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.browseContents(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt8(15), types.NewUInt32(3032), types.NewUInt8(29))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestBrowseContentsByType_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.browseContentsByType(nil, fakePacket{}, 7, types.NewString("s1"), types.NewUInt32(2023), types.NewUInt8(22), types.NewUInt32(4041), types.NewUInt8(36))
	checkMessage(t, message, nexErr, ubi.MethodBrowseContentsByType, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestBrowseContentsByType_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.BrowseContentsByType = func(packet nex.PacketInterface, _ types.String, _ types.UInt32, _ types.UInt8, _ types.UInt32, _ types.UInt8) (types.List[types.DataHolder], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[types.DataHolder](), nil
	}

	message, nexErr := commonProtocol.browseContentsByType(errors.New("bad"), fakePacket{}, 7, types.NewString(""), types.NewUInt32(0), types.NewUInt8(0), types.NewUInt32(0), types.NewUInt8(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestBrowseContentsByType_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.BrowseContentsByType = func(packet nex.PacketInterface, got0 types.String, got1 types.UInt32, got2 types.UInt8, got3 types.UInt32, got4 types.UInt8) (types.List[types.DataHolder], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewString("s1"))) {
			t.Errorf("contentType differs")
		}

		if !(got1.Equals(types.NewUInt32(2023))) {
			t.Errorf("flags differs")
		}

		if !(got2.Equals(types.NewUInt8(22))) {
			t.Errorf("flagsMaskType differs")
		}

		if !(got3.Equals(types.NewUInt32(4041))) {
			t.Errorf("extendedFlags differs")
		}

		if !(got4.Equals(types.NewUInt8(36))) {
			t.Errorf("extendedFlagsMaskType differs")
		}

		return types.List[types.DataHolder]{types.DataHolder{Object: types.NewData()}, types.DataHolder{Object: types.NewData()}}, nil
	}

	message, nexErr := commonProtocol.browseContentsByType(nil, fakePacket{}, 9, types.NewString("s1"), types.NewUInt32(2023), types.NewUInt8(22), types.NewUInt32(4041), types.NewUInt8(36))
	checkMessage(t, message, nexErr, ubi.MethodBrowseContentsByType, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x05, 0x00, 0x44, 0x61, 0x74, 0x61, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x44, 0x61, 0x74, 0x61, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.BrowseContentsByType = func(packet nex.PacketInterface, _ types.String, _ types.UInt32, _ types.UInt8, _ types.UInt32, _ types.UInt8) (types.List[types.DataHolder], *nex.Error) {
		return types.NewList[types.DataHolder](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.browseContentsByType(nil, fakePacket{}, 9, types.NewString("s1"), types.NewUInt32(2023), types.NewUInt8(22), types.NewUInt32(4041), types.NewUInt8(36))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetContentDetails_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getContentDetails(nil, fakePacket{}, 7, types.NewUInt32(1014))
	checkNotImplemented(t, message, nexErr)
}

func TestGetContentDetails_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetContentDetails = func(packet nex.PacketInterface, _ types.UInt32) (types.DataHolder, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewDataHolder(), nil
	}

	message, nexErr := commonProtocol.getContentDetails(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetContentDetails_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetContentDetails = func(packet nex.PacketInterface, got0 types.UInt32) (types.DataHolder, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("contentID differs")
		}

		return types.DataHolder{Object: types.NewData()}, nil
	}

	message, nexErr := commonProtocol.getContentDetails(nil, fakePacket{}, 9, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodGetContentDetails, 9, []byte{0x05, 0x00, 0x44, 0x61, 0x74, 0x61, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetContentDetails = func(packet nex.PacketInterface, _ types.UInt32) (types.DataHolder, *nex.Error) {
		return types.NewDataHolder(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getContentDetails(nil, fakePacket{}, 9, types.NewUInt32(1014))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestDownloadContent_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.downloadContent(nil, fakePacket{}, 7, types.NewUInt32(1014))
	checkNotImplemented(t, message, nexErr)
}

func TestDownloadContent_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.DownloadContent = func(packet nex.PacketInterface, _ types.UInt32) (types.String, types.String, types.String, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewString(""), types.NewString(""), types.NewString(""), nil
	}

	message, nexErr := commonProtocol.downloadContent(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestDownloadContent_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.DownloadContent = func(packet nex.PacketInterface, got0 types.UInt32) (types.String, types.String, types.String, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("contentID differs")
		}

		return types.NewString("s2"), types.NewString("s3"), types.NewString("s4"), nil
	}

	message, nexErr := commonProtocol.downloadContent(nil, fakePacket{}, 9, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodDownloadContent, 9, []byte{0x03, 0x00, 0x73, 0x32, 0x00, 0x03, 0x00, 0x73, 0x33, 0x00, 0x03, 0x00, 0x73, 0x34, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.DownloadContent = func(packet nex.PacketInterface, _ types.UInt32) (types.String, types.String, types.String, *nex.Error) {
		return types.NewString(""), types.NewString(""), types.NewString(""), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.downloadContent(nil, fakePacket{}, 9, types.NewUInt32(1014))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}
