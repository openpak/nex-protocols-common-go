package localization

import (
	"bytes"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/localization"
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

func TestGetLocaleCode_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getLocaleCode(nil, fakePacket{}, 7)
	checkNotImplemented(t, message, nexErr)
}

func TestGetLocaleCode_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetLocaleCode = func(packet nex.PacketInterface) (types.String, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewString(""), nil
	}

	message, nexErr := commonProtocol.getLocaleCode(errors.New("bad"), fakePacket{}, 7)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetLocaleCode_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetLocaleCode = func(packet nex.PacketInterface) (types.String, *nex.Error) {
		called = true

		return types.NewString("s1"), nil
	}

	message, nexErr := commonProtocol.getLocaleCode(nil, fakePacket{}, 9)
	checkMessage(t, message, nexErr, ubi.MethodGetLocaleCode, 9, []byte{0x03, 0x00, 0x73, 0x31, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetLocaleCode = func(packet nex.PacketInterface) (types.String, *nex.Error) {
		return types.NewString(""), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getLocaleCode(nil, fakePacket{}, 9)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSetLocaleCode_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.setLocaleCode(nil, fakePacket{}, 7, types.NewString("s1"))
	checkNotImplemented(t, message, nexErr)
}

func TestSetLocaleCode_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SetLocaleCode = func(packet nex.PacketInterface, _ types.String) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.setLocaleCode(errors.New("bad"), fakePacket{}, 7, types.NewString(""))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSetLocaleCode_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SetLocaleCode = func(packet nex.PacketInterface, got0 types.String) *nex.Error {
		called = true

		if !(got0.Equals(types.NewString("s1"))) {
			t.Errorf("localCode differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.setLocaleCode(nil, fakePacket{}, 9, types.NewString("s1"))
	checkMessage(t, message, nexErr, ubi.MethodSetLocaleCode, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SetLocaleCode = func(packet nex.PacketInterface, _ types.String) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.setLocaleCode(nil, fakePacket{}, 9, types.NewString("s1"))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}
