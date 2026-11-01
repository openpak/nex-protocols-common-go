package webnotificationstorage

import (
	"bytes"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/web-notification-storage"
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

func TestRegisterUser_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.registerUser(nil, fakePacket{}, 7)
	checkNotImplemented(t, message, nexErr)
}

func TestRegisterUser_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.RegisterUser = func(packet nex.PacketInterface) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.registerUser(errors.New("bad"), fakePacket{}, 7)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestRegisterUser_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.RegisterUser = func(packet nex.PacketInterface) *nex.Error {
		called = true

		return nil
	}

	message, nexErr := commonProtocol.registerUser(nil, fakePacket{}, 9)
	checkMessage(t, message, nexErr, ubi.MethodRegisterUser, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.RegisterUser = func(packet nex.PacketInterface) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.registerUser(nil, fakePacket{}, 9)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestPollNotifications_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.pollNotifications(nil, fakePacket{}, 7)
	checkMessage(t, message, nexErr, ubi.MethodPollNotifications, 7, []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestPollNotifications_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.PollNotifications = func(packet nex.PacketInterface) (types.String, types.Int32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewString(""), types.NewInt32(0), nil
	}

	message, nexErr := commonProtocol.pollNotifications(errors.New("bad"), fakePacket{}, 7)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestPollNotifications_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.PollNotifications = func(packet nex.PacketInterface) (types.String, types.Int32, *nex.Error) {
		called = true

		return types.NewString("s1"), types.NewInt32(-63), nil
	}

	message, nexErr := commonProtocol.pollNotifications(nil, fakePacket{}, 9)
	checkMessage(t, message, nexErr, ubi.MethodPollNotifications, 9, []byte{0x03, 0x00, 0x73, 0x31, 0x00, 0xc1, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.PollNotifications = func(packet nex.PacketInterface) (types.String, types.Int32, *nex.Error) {
		return types.NewString(""), types.NewInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.pollNotifications(nil, fakePacket{}, 9)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestUnregisterUser_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.unregisterUser(nil, fakePacket{}, 7)
	checkNotImplemented(t, message, nexErr)
}

func TestUnregisterUser_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.UnregisterUser = func(packet nex.PacketInterface) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.unregisterUser(errors.New("bad"), fakePacket{}, 7)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestUnregisterUser_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.UnregisterUser = func(packet nex.PacketInterface) *nex.Error {
		called = true

		return nil
	}

	message, nexErr := commonProtocol.unregisterUser(nil, fakePacket{}, 9)
	checkMessage(t, message, nexErr, ubi.MethodUnregisterUser, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.UnregisterUser = func(packet nex.PacketInterface) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.unregisterUser(nil, fakePacket{}, 9)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}
