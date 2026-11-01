package gamesession

import (
	"bytes"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/game-session"
	game_session_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/game-session/types"
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

func TestCreateSession_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.createSession(nil, fakePacket{}, 7, func() game_session_types.GameSession {
		v := game_session_types.NewGameSession()
		v.TypeID = types.NewUInt32(2023)
		v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(5050)
			v.Value = types.NewUInt32(6059)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(8077)
			v.Value = types.NewUInt32(9086)
			return v
		}()}
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestCreateSession_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.CreateSession = func(packet nex.PacketInterface, _ game_session_types.GameSession) (game_session_types.GameSessionKey, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return game_session_types.NewGameSessionKey(), nil
	}

	message, nexErr := commonProtocol.createSession(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSession())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestCreateSession_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.CreateSession = func(packet nex.PacketInterface, got0 game_session_types.GameSession) (game_session_types.GameSessionKey, *nex.Error) {
		called = true

		if !(got0.Equals(func() game_session_types.GameSession {
			v := game_session_types.NewGameSession()
			v.TypeID = types.NewUInt32(2023)
			v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(5050)
				v.Value = types.NewUInt32(6059)
				return v
			}(), func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(8077)
				v.Value = types.NewUInt32(9086)
				return v
			}()}
			return v
		}())) {
			t.Errorf("gameSession differs")
		}

		return func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(11104)
			v.SessionID = types.NewUInt32(12113)
			return v
		}(), nil
	}

	message, nexErr := commonProtocol.createSession(nil, fakePacket{}, 9, func() game_session_types.GameSession {
		v := game_session_types.NewGameSession()
		v.TypeID = types.NewUInt32(2023)
		v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(5050)
			v.Value = types.NewUInt32(6059)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(8077)
			v.Value = types.NewUInt32(9086)
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodCreateSession, 9, []byte{0x60, 0x2b, 0x00, 0x00, 0x51, 0x2f, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.CreateSession = func(packet nex.PacketInterface, _ game_session_types.GameSession) (game_session_types.GameSessionKey, *nex.Error) {
		return game_session_types.NewGameSessionKey(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.createSession(nil, fakePacket{}, 9, func() game_session_types.GameSession {
		v := game_session_types.NewGameSession()
		v.TypeID = types.NewUInt32(2023)
		v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(5050)
			v.Value = types.NewUInt32(6059)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(8077)
			v.Value = types.NewUInt32(9086)
			return v
		}()}
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestUpdateSession_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.updateSession(nil, fakePacket{}, 7, func() game_session_types.GameSessionUpdate {
		v := game_session_types.NewGameSessionUpdate()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(7068)
			v.Value = types.NewUInt32(8077)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.NewUInt32(11104)
			return v
		}()}
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestUpdateSession_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.UpdateSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionUpdate) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.updateSession(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionUpdate())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestUpdateSession_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.UpdateSession = func(packet nex.PacketInterface, got0 game_session_types.GameSessionUpdate) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionUpdate {
			v := game_session_types.NewGameSessionUpdate()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(3032)
				v.SessionID = types.NewUInt32(4041)
				return v
			}()
			v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(7068)
				v.Value = types.NewUInt32(8077)
				return v
			}(), func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(10095)
				v.Value = types.NewUInt32(11104)
				return v
			}()}
			return v
		}())) {
			t.Errorf("gameSessionUpdate differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.updateSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionUpdate {
		v := game_session_types.NewGameSessionUpdate()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(7068)
			v.Value = types.NewUInt32(8077)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.NewUInt32(11104)
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodUpdateSession, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.UpdateSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionUpdate) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.updateSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionUpdate {
		v := game_session_types.NewGameSessionUpdate()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(7068)
			v.Value = types.NewUInt32(8077)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.NewUInt32(11104)
			return v
		}()}
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestDeleteSession_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.deleteSession(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestDeleteSession_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.DeleteSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.deleteSession(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestDeleteSession_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.DeleteSession = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.deleteSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodDeleteSession, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.DeleteSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.deleteSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestMigrateSession_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.migrateSession(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestMigrateSession_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.MigrateSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) (game_session_types.GameSessionKey, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return game_session_types.NewGameSessionKey(), nil
	}

	message, nexErr := commonProtocol.migrateSession(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestMigrateSession_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.MigrateSession = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey) (game_session_types.GameSessionKey, *nex.Error) {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		return func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(5050)
			v.SessionID = types.NewUInt32(6059)
			return v
		}(), nil
	}

	message, nexErr := commonProtocol.migrateSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodMigrateSession, 9, []byte{0xba, 0x13, 0x00, 0x00, 0xab, 0x17, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.MigrateSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) (game_session_types.GameSessionKey, *nex.Error) {
		return game_session_types.NewGameSessionKey(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.migrateSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestLeaveSession_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.leaveSession(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestLeaveSession_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.LeaveSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.leaveSession(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestLeaveSession_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.LeaveSession = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.leaveSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodLeaveSession, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.LeaveSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.leaveSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetSession_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getSession(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestGetSession_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) (game_session_types.GameSessionSearchResult, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return game_session_types.NewGameSessionSearchResult(), nil
	}

	message, nexErr := commonProtocol.getSession(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetSession_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetSession = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey) (game_session_types.GameSessionSearchResult, *nex.Error) {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		return func() game_session_types.GameSessionSearchResult {
			v := game_session_types.NewGameSessionSearchResult()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(6059)
				v.SessionID = types.NewUInt32(7068)
				return v
			}()
			v.HostPID = types.NewUInt32(8077)
			v.HostURLs = types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.10"), types.NewStationURL("prudp:/address=10.0.0.11")}
			v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(14131)
				v.Value = types.NewUInt32(15140)
				return v
			}(), func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(17158)
				v.Value = types.NewUInt32(18167)
				return v
			}()}
			return v
		}(), nil
	}

	message, nexErr := commonProtocol.getSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetSession, 9, []byte{0xab, 0x17, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x8d, 0x1f, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x30, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x31, 0x00, 0x02, 0x00, 0x00, 0x00, 0x33, 0x37, 0x00, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x06, 0x43, 0x00, 0x00, 0xf7, 0x46, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) (game_session_types.GameSessionSearchResult, *nex.Error) {
		return game_session_types.NewGameSessionSearchResult(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSearchSessions_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.searchSessions(nil, fakePacket{}, 7, func() game_session_types.GameSessionQuery {
		v := game_session_types.NewGameSessionQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.Parameters = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(6059)
			v.Value = types.NewUInt32(7068)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(9086)
			v.Value = types.NewUInt32(10095)
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchSessions, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestSearchSessions_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SearchSessions = func(packet nex.PacketInterface, _ game_session_types.GameSessionQuery) (types.List[game_session_types.GameSessionSearchResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[game_session_types.GameSessionSearchResult](), nil
	}

	message, nexErr := commonProtocol.searchSessions(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionQuery())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSearchSessions_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SearchSessions = func(packet nex.PacketInterface, got0 game_session_types.GameSessionQuery) (types.List[game_session_types.GameSessionSearchResult], *nex.Error) {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionQuery {
			v := game_session_types.NewGameSessionQuery()
			v.TypeID = types.NewUInt32(2023)
			v.QueryID = types.NewUInt32(3032)
			v.Parameters = types.List[game_session_types.Property]{func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(6059)
				v.Value = types.NewUInt32(7068)
				return v
			}(), func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(9086)
				v.Value = types.NewUInt32(10095)
				return v
			}()}
			return v
		}())) {
			t.Errorf("gameSessionQuery differs")
		}

		return types.List[game_session_types.GameSessionSearchResult]{func() game_session_types.GameSessionSearchResult {
			v := game_session_types.NewGameSessionSearchResult()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(14131)
				v.SessionID = types.NewUInt32(15140)
				return v
			}()
			v.HostPID = types.NewUInt32(16149)
			v.HostURLs = types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.18"), types.NewStationURL("prudp:/address=10.0.0.19")}
			v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(22203)
				v.Value = types.NewUInt32(23212)
				return v
			}(), func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(25230)
				v.Value = types.NewUInt32(26239)
				return v
			}()}
			return v
		}(), func() game_session_types.GameSessionSearchResult {
			v := game_session_types.NewGameSessionSearchResult()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(29266)
				v.SessionID = types.NewUInt32(30275)
				return v
			}()
			v.HostPID = types.NewUInt32(31284)
			v.HostURLs = types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.33"), types.NewStationURL("prudp:/address=10.0.0.34")}
			v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(37338)
				v.Value = types.NewUInt32(38347)
				return v
			}(), func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(40365)
				v.Value = types.NewUInt32(41374)
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.searchSessions(nil, fakePacket{}, 9, func() game_session_types.GameSessionQuery {
		v := game_session_types.NewGameSessionQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.Parameters = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(6059)
			v.Value = types.NewUInt32(7068)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(9086)
			v.Value = types.NewUInt32(10095)
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchSessions, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x33, 0x37, 0x00, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x15, 0x3f, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x38, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x39, 0x00, 0x02, 0x00, 0x00, 0x00, 0xbb, 0x56, 0x00, 0x00, 0xac, 0x5a, 0x00, 0x00, 0x8e, 0x62, 0x00, 0x00, 0x7f, 0x66, 0x00, 0x00, 0x52, 0x72, 0x00, 0x00, 0x43, 0x76, 0x00, 0x00, 0x34, 0x7a, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x33, 0x33, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x33, 0x34, 0x00, 0x02, 0x00, 0x00, 0x00, 0xda, 0x91, 0x00, 0x00, 0xcb, 0x95, 0x00, 0x00, 0xad, 0x9d, 0x00, 0x00, 0x9e, 0xa1, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SearchSessions = func(packet nex.PacketInterface, _ game_session_types.GameSessionQuery) (types.List[game_session_types.GameSessionSearchResult], *nex.Error) {
		return types.NewList[game_session_types.GameSessionSearchResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.searchSessions(nil, fakePacket{}, 9, func() game_session_types.GameSessionQuery {
		v := game_session_types.NewGameSessionQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.Parameters = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(6059)
			v.Value = types.NewUInt32(7068)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(9086)
			v.Value = types.NewUInt32(10095)
			return v
		}()}
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestAddParticipants_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.addParticipants(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)}, types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)})
	checkNotImplemented(t, message, nexErr)
}

func TestAddParticipants_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.AddParticipants = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey, _ types.List[types.UInt32], _ types.List[types.UInt32]) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.addParticipants(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey(), types.NewList[types.UInt32](), types.NewList[types.UInt32]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestAddParticipants_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.AddParticipants = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey, got1 types.List[types.UInt32], got2 types.List[types.UInt32]) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		if !(got1.Equals(types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})) {
			t.Errorf("publicParticipantIDs differs")
		}

		if !(got2.Equals(types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)})) {
			t.Errorf("privateParticipantIDs differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.addParticipants(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)}, types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)})
	checkMessage(t, message, nexErr, ubi.MethodAddParticipants, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.AddParticipants = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey, _ types.List[types.UInt32], _ types.List[types.UInt32]) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.addParticipants(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)}, types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestRemoveParticipants_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.removeParticipants(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})
	checkNotImplemented(t, message, nexErr)
}

func TestRemoveParticipants_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.RemoveParticipants = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey, _ types.List[types.UInt32]) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.removeParticipants(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey(), types.NewList[types.UInt32]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestRemoveParticipants_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.RemoveParticipants = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey, got1 types.List[types.UInt32]) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		if !(got1.Equals(types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})) {
			t.Errorf("participantIDs differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.removeParticipants(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})
	checkMessage(t, message, nexErr, ubi.MethodRemoveParticipants, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.RemoveParticipants = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey, _ types.List[types.UInt32]) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.removeParticipants(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetParticipantCount_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getParticipantCount(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetParticipantCount, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetParticipantCount_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetParticipantCount = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) (types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.getParticipantCount(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetParticipantCount_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetParticipantCount = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey) (types.UInt32, *nex.Error) {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		return types.NewUInt32(4041), nil
	}

	message, nexErr := commonProtocol.getParticipantCount(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetParticipantCount, 9, []byte{0xc9, 0x0f, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetParticipantCount = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) (types.UInt32, *nex.Error) {
		return types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getParticipantCount(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetParticipants_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getParticipants(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)})
	checkMessage(t, message, nexErr, ubi.MethodGetParticipants, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetParticipants_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetParticipants = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey, _ types.ResultRange) (types.List[game_session_types.GameSessionParticipant], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[game_session_types.GameSessionParticipant](), nil
	}

	message, nexErr := commonProtocol.getParticipants(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey(), types.NewResultRange())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetParticipants_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetParticipants = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey, got1 types.ResultRange) (types.List[game_session_types.GameSessionParticipant], *nex.Error) {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		if !(got1 == types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)}) {
			t.Errorf("resultRange differs")
		}

		return types.List[game_session_types.GameSessionParticipant]{func() game_session_types.GameSessionParticipant {
			v := game_session_types.NewGameSessionParticipant()
			v.PID = types.NewUInt32(7068)
			v.Name = types.NewString("s8")
			v.StationURLs = types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.10"), types.NewStationURL("prudp:/address=10.0.0.11")}
			return v
		}(), func() game_session_types.GameSessionParticipant {
			v := game_session_types.NewGameSessionParticipant()
			v.PID = types.NewUInt32(13122)
			v.Name = types.NewString("s14")
			v.StationURLs = types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.16"), types.NewStationURL("prudp:/address=10.0.0.17")}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getParticipants(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)})
	checkMessage(t, message, nexErr, ubi.MethodGetParticipants, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x03, 0x00, 0x73, 0x38, 0x00, 0x02, 0x00, 0x00, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x30, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x31, 0x00, 0x42, 0x33, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x34, 0x00, 0x02, 0x00, 0x00, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x36, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x37, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetParticipants = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey, _ types.ResultRange) (types.List[game_session_types.GameSessionParticipant], *nex.Error) {
		return types.NewList[game_session_types.GameSessionParticipant](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getParticipants(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSendInvitation_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.sendInvitation(nil, fakePacket{}, 7, func() game_session_types.GameSessionInvitation {
		v := game_session_types.NewGameSessionInvitation()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.RecipientPIDs = types.List[types.UInt32]{types.NewUInt32(6059), types.NewUInt32(7068)}
		v.Message = types.NewString("s8")
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestSendInvitation_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SendInvitation = func(packet nex.PacketInterface, _ game_session_types.GameSessionInvitation) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.sendInvitation(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionInvitation())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSendInvitation_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SendInvitation = func(packet nex.PacketInterface, got0 game_session_types.GameSessionInvitation) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionInvitation {
			v := game_session_types.NewGameSessionInvitation()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(3032)
				v.SessionID = types.NewUInt32(4041)
				return v
			}()
			v.RecipientPIDs = types.List[types.UInt32]{types.NewUInt32(6059), types.NewUInt32(7068)}
			v.Message = types.NewString("s8")
			return v
		}())) {
			t.Errorf("invitation differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.sendInvitation(nil, fakePacket{}, 9, func() game_session_types.GameSessionInvitation {
		v := game_session_types.NewGameSessionInvitation()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.RecipientPIDs = types.List[types.UInt32]{types.NewUInt32(6059), types.NewUInt32(7068)}
		v.Message = types.NewString("s8")
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSendInvitation, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SendInvitation = func(packet nex.PacketInterface, _ game_session_types.GameSessionInvitation) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.sendInvitation(nil, fakePacket{}, 9, func() game_session_types.GameSessionInvitation {
		v := game_session_types.NewGameSessionInvitation()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.RecipientPIDs = types.List[types.UInt32]{types.NewUInt32(6059), types.NewUInt32(7068)}
		v.Message = types.NewString("s8")
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetInvitationReceivedCount_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getInvitationReceivedCount(nil, fakePacket{}, 7, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodGetInvitationReceivedCount, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetInvitationReceivedCount_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetInvitationReceivedCount = func(packet nex.PacketInterface, _ types.UInt32) (types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.getInvitationReceivedCount(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetInvitationReceivedCount_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetInvitationReceivedCount = func(packet nex.PacketInterface, got0 types.UInt32) (types.UInt32, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("gameSessionTypeID differs")
		}

		return types.NewUInt32(2023), nil
	}

	message, nexErr := commonProtocol.getInvitationReceivedCount(nil, fakePacket{}, 9, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodGetInvitationReceivedCount, 9, []byte{0xe7, 0x07, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetInvitationReceivedCount = func(packet nex.PacketInterface, _ types.UInt32) (types.UInt32, *nex.Error) {
		return types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getInvitationReceivedCount(nil, fakePacket{}, 9, types.NewUInt32(1014))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetInvitationsReceived_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getInvitationsReceived(nil, fakePacket{}, 7, types.NewUInt32(1014), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	checkMessage(t, message, nexErr, ubi.MethodGetInvitationsReceived, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetInvitationsReceived_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetInvitationsReceived = func(packet nex.PacketInterface, _ types.UInt32, _ types.ResultRange) (types.List[game_session_types.GameSessionInvitationReceived], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[game_session_types.GameSessionInvitationReceived](), nil
	}

	message, nexErr := commonProtocol.getInvitationsReceived(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewResultRange())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetInvitationsReceived_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetInvitationsReceived = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.ResultRange) (types.List[game_session_types.GameSessionInvitationReceived], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("gameSessionTypeID differs")
		}

		if !(got1 == types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)}) {
			t.Errorf("resultRange differs")
		}

		return types.List[game_session_types.GameSessionInvitationReceived]{func() game_session_types.GameSessionInvitationReceived {
			v := game_session_types.NewGameSessionInvitationReceived()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(6059)
				v.SessionID = types.NewUInt32(7068)
				return v
			}()
			v.SenderPID = types.NewUInt32(8077)
			v.Message = types.NewString("s9")
			v.CreationTime = types.NewDateTime(10000107)
			return v
		}(), func() game_session_types.GameSessionInvitationReceived {
			v := game_session_types.NewGameSessionInvitationReceived()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(13122)
				v.SessionID = types.NewUInt32(14131)
				return v
			}()
			v.SenderPID = types.NewUInt32(15140)
			v.Message = types.NewString("s16")
			v.CreationTime = types.NewDateTime(17000128)
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getInvitationsReceived(nil, fakePacket{}, 9, types.NewUInt32(1014), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	checkMessage(t, message, nexErr, ubi.MethodGetInvitationsReceived, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xab, 0x17, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x8d, 0x1f, 0x00, 0x00, 0x03, 0x00, 0x73, 0x39, 0x00, 0xeb, 0x96, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x42, 0x33, 0x00, 0x00, 0x33, 0x37, 0x00, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x36, 0x00, 0xc0, 0x66, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetInvitationsReceived = func(packet nex.PacketInterface, _ types.UInt32, _ types.ResultRange) (types.List[game_session_types.GameSessionInvitationReceived], *nex.Error) {
		return types.NewList[game_session_types.GameSessionInvitationReceived](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getInvitationsReceived(nil, fakePacket{}, 9, types.NewUInt32(1014), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetInvitationSentCount_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getInvitationSentCount(nil, fakePacket{}, 7, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodGetInvitationSentCount, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetInvitationSentCount_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetInvitationSentCount = func(packet nex.PacketInterface, _ types.UInt32) (types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.getInvitationSentCount(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetInvitationSentCount_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetInvitationSentCount = func(packet nex.PacketInterface, got0 types.UInt32) (types.UInt32, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("gameSessionTypeID differs")
		}

		return types.NewUInt32(2023), nil
	}

	message, nexErr := commonProtocol.getInvitationSentCount(nil, fakePacket{}, 9, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodGetInvitationSentCount, 9, []byte{0xe7, 0x07, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetInvitationSentCount = func(packet nex.PacketInterface, _ types.UInt32) (types.UInt32, *nex.Error) {
		return types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getInvitationSentCount(nil, fakePacket{}, 9, types.NewUInt32(1014))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetInvitationsSent_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getInvitationsSent(nil, fakePacket{}, 7, types.NewUInt32(1014), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	checkMessage(t, message, nexErr, ubi.MethodGetInvitationsSent, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetInvitationsSent_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetInvitationsSent = func(packet nex.PacketInterface, _ types.UInt32, _ types.ResultRange) (types.List[game_session_types.GameSessionInvitationSent], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[game_session_types.GameSessionInvitationSent](), nil
	}

	message, nexErr := commonProtocol.getInvitationsSent(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewResultRange())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetInvitationsSent_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetInvitationsSent = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.ResultRange) (types.List[game_session_types.GameSessionInvitationSent], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("gameSessionTypeID differs")
		}

		if !(got1 == types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)}) {
			t.Errorf("resultRange differs")
		}

		return types.List[game_session_types.GameSessionInvitationSent]{func() game_session_types.GameSessionInvitationSent {
			v := game_session_types.NewGameSessionInvitationSent()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(6059)
				v.SessionID = types.NewUInt32(7068)
				return v
			}()
			v.RecipientPID = types.NewUInt32(8077)
			v.Message = types.NewString("s9")
			v.CreationTime = types.NewDateTime(10000107)
			return v
		}(), func() game_session_types.GameSessionInvitationSent {
			v := game_session_types.NewGameSessionInvitationSent()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(13122)
				v.SessionID = types.NewUInt32(14131)
				return v
			}()
			v.RecipientPID = types.NewUInt32(15140)
			v.Message = types.NewString("s16")
			v.CreationTime = types.NewDateTime(17000128)
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getInvitationsSent(nil, fakePacket{}, 9, types.NewUInt32(1014), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	checkMessage(t, message, nexErr, ubi.MethodGetInvitationsSent, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xab, 0x17, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x8d, 0x1f, 0x00, 0x00, 0x03, 0x00, 0x73, 0x39, 0x00, 0xeb, 0x96, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x42, 0x33, 0x00, 0x00, 0x33, 0x37, 0x00, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x36, 0x00, 0xc0, 0x66, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetInvitationsSent = func(packet nex.PacketInterface, _ types.UInt32, _ types.ResultRange) (types.List[game_session_types.GameSessionInvitationSent], *nex.Error) {
		return types.NewList[game_session_types.GameSessionInvitationSent](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getInvitationsSent(nil, fakePacket{}, 9, types.NewUInt32(1014), types.ResultRange{Offset: types.NewUInt32(6), Length: types.NewUInt32(16)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestAcceptInvitation_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.acceptInvitation(nil, fakePacket{}, 7, func() game_session_types.GameSessionInvitationReceived {
		v := game_session_types.NewGameSessionInvitationReceived()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.SenderPID = types.NewUInt32(5050)
		v.Message = types.NewString("s6")
		v.CreationTime = types.NewDateTime(7000098)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestAcceptInvitation_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.AcceptInvitation = func(packet nex.PacketInterface, _ game_session_types.GameSessionInvitationReceived) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.acceptInvitation(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionInvitationReceived())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestAcceptInvitation_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.AcceptInvitation = func(packet nex.PacketInterface, got0 game_session_types.GameSessionInvitationReceived) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionInvitationReceived {
			v := game_session_types.NewGameSessionInvitationReceived()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(3032)
				v.SessionID = types.NewUInt32(4041)
				return v
			}()
			v.SenderPID = types.NewUInt32(5050)
			v.Message = types.NewString("s6")
			v.CreationTime = types.NewDateTime(7000098)
			return v
		}())) {
			t.Errorf("gameSessionInvitation differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.acceptInvitation(nil, fakePacket{}, 9, func() game_session_types.GameSessionInvitationReceived {
		v := game_session_types.NewGameSessionInvitationReceived()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.SenderPID = types.NewUInt32(5050)
		v.Message = types.NewString("s6")
		v.CreationTime = types.NewDateTime(7000098)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodAcceptInvitation, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.AcceptInvitation = func(packet nex.PacketInterface, _ game_session_types.GameSessionInvitationReceived) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.acceptInvitation(nil, fakePacket{}, 9, func() game_session_types.GameSessionInvitationReceived {
		v := game_session_types.NewGameSessionInvitationReceived()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.SenderPID = types.NewUInt32(5050)
		v.Message = types.NewString("s6")
		v.CreationTime = types.NewDateTime(7000098)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestDeclineInvitation_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.declineInvitation(nil, fakePacket{}, 7, func() game_session_types.GameSessionInvitationReceived {
		v := game_session_types.NewGameSessionInvitationReceived()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.SenderPID = types.NewUInt32(5050)
		v.Message = types.NewString("s6")
		v.CreationTime = types.NewDateTime(7000098)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestDeclineInvitation_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.DeclineInvitation = func(packet nex.PacketInterface, _ game_session_types.GameSessionInvitationReceived) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.declineInvitation(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionInvitationReceived())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestDeclineInvitation_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.DeclineInvitation = func(packet nex.PacketInterface, got0 game_session_types.GameSessionInvitationReceived) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionInvitationReceived {
			v := game_session_types.NewGameSessionInvitationReceived()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(3032)
				v.SessionID = types.NewUInt32(4041)
				return v
			}()
			v.SenderPID = types.NewUInt32(5050)
			v.Message = types.NewString("s6")
			v.CreationTime = types.NewDateTime(7000098)
			return v
		}())) {
			t.Errorf("gameSessionInvitation differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.declineInvitation(nil, fakePacket{}, 9, func() game_session_types.GameSessionInvitationReceived {
		v := game_session_types.NewGameSessionInvitationReceived()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.SenderPID = types.NewUInt32(5050)
		v.Message = types.NewString("s6")
		v.CreationTime = types.NewDateTime(7000098)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodDeclineInvitation, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.DeclineInvitation = func(packet nex.PacketInterface, _ game_session_types.GameSessionInvitationReceived) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.declineInvitation(nil, fakePacket{}, 9, func() game_session_types.GameSessionInvitationReceived {
		v := game_session_types.NewGameSessionInvitationReceived()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.SenderPID = types.NewUInt32(5050)
		v.Message = types.NewString("s6")
		v.CreationTime = types.NewDateTime(7000098)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestCancelInvitation_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.cancelInvitation(nil, fakePacket{}, 7, func() game_session_types.GameSessionInvitationSent {
		v := game_session_types.NewGameSessionInvitationSent()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.RecipientPID = types.NewUInt32(5050)
		v.Message = types.NewString("s6")
		v.CreationTime = types.NewDateTime(7000098)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestCancelInvitation_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.CancelInvitation = func(packet nex.PacketInterface, _ game_session_types.GameSessionInvitationSent) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.cancelInvitation(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionInvitationSent())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestCancelInvitation_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.CancelInvitation = func(packet nex.PacketInterface, got0 game_session_types.GameSessionInvitationSent) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionInvitationSent {
			v := game_session_types.NewGameSessionInvitationSent()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(3032)
				v.SessionID = types.NewUInt32(4041)
				return v
			}()
			v.RecipientPID = types.NewUInt32(5050)
			v.Message = types.NewString("s6")
			v.CreationTime = types.NewDateTime(7000098)
			return v
		}())) {
			t.Errorf("gameSessionInvitation differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.cancelInvitation(nil, fakePacket{}, 9, func() game_session_types.GameSessionInvitationSent {
		v := game_session_types.NewGameSessionInvitationSent()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.RecipientPID = types.NewUInt32(5050)
		v.Message = types.NewString("s6")
		v.CreationTime = types.NewDateTime(7000098)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodCancelInvitation, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.CancelInvitation = func(packet nex.PacketInterface, _ game_session_types.GameSessionInvitationSent) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.cancelInvitation(nil, fakePacket{}, 9, func() game_session_types.GameSessionInvitationSent {
		v := game_session_types.NewGameSessionInvitationSent()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.RecipientPID = types.NewUInt32(5050)
		v.Message = types.NewString("s6")
		v.CreationTime = types.NewDateTime(7000098)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSendTextMessage_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.sendTextMessage(nil, fakePacket{}, 7, func() game_session_types.GameSessionMessage {
		v := game_session_types.NewGameSessionMessage()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.Message = types.NewString("s5")
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestSendTextMessage_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SendTextMessage = func(packet nex.PacketInterface, _ game_session_types.GameSessionMessage) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.sendTextMessage(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionMessage())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSendTextMessage_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SendTextMessage = func(packet nex.PacketInterface, got0 game_session_types.GameSessionMessage) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionMessage {
			v := game_session_types.NewGameSessionMessage()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(3032)
				v.SessionID = types.NewUInt32(4041)
				return v
			}()
			v.Message = types.NewString("s5")
			return v
		}())) {
			t.Errorf("gameSessionMessage differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.sendTextMessage(nil, fakePacket{}, 9, func() game_session_types.GameSessionMessage {
		v := game_session_types.NewGameSessionMessage()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.Message = types.NewString("s5")
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSendTextMessage, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SendTextMessage = func(packet nex.PacketInterface, _ game_session_types.GameSessionMessage) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.sendTextMessage(nil, fakePacket{}, 9, func() game_session_types.GameSessionMessage {
		v := game_session_types.NewGameSessionMessage()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}()
		v.Message = types.NewString("s5")
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestRegisterURLs_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.registerURLs(nil, fakePacket{}, 7, types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.2"), types.NewStationURL("prudp:/address=10.0.0.3")})
	checkNotImplemented(t, message, nexErr)
}

func TestRegisterURLs_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.RegisterURLs = func(packet nex.PacketInterface, _ types.List[types.StationURL]) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.registerURLs(errors.New("bad"), fakePacket{}, 7, types.NewList[types.StationURL]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestRegisterURLs_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.RegisterURLs = func(packet nex.PacketInterface, got0 types.List[types.StationURL]) *nex.Error {
		called = true

		if !(got0.Equals(types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.2"), types.NewStationURL("prudp:/address=10.0.0.3")})) {
			t.Errorf("stationURLs differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.registerURLs(nil, fakePacket{}, 9, types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.2"), types.NewStationURL("prudp:/address=10.0.0.3")})
	checkMessage(t, message, nexErr, ubi.MethodRegisterURLs, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.RegisterURLs = func(packet nex.PacketInterface, _ types.List[types.StationURL]) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.registerURLs(nil, fakePacket{}, 9, types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.2"), types.NewStationURL("prudp:/address=10.0.0.3")})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestJoinSession_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.joinSession(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestJoinSession_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.JoinSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.joinSession(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestJoinSession_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.JoinSession = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.joinSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodJoinSession, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.JoinSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.joinSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestAbandonSession_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.abandonSession(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestAbandonSession_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.AbandonSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.abandonSession(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestAbandonSession_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.AbandonSession = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.abandonSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodAbandonSession, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.AbandonSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.abandonSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSearchSessionsWithParticipants_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.searchSessionsWithParticipants(nil, fakePacket{}, 7, types.NewUInt32(1014), types.List[types.UInt32]{types.NewUInt32(3032), types.NewUInt32(4041)})
	checkMessage(t, message, nexErr, ubi.MethodSearchSessionsWithParticipants, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestSearchSessionsWithParticipants_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SearchSessionsWithParticipants = func(packet nex.PacketInterface, _ types.UInt32, _ types.List[types.UInt32]) (types.List[game_session_types.GameSessionSearchWithParticipantsResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[game_session_types.GameSessionSearchWithParticipantsResult](), nil
	}

	message, nexErr := commonProtocol.searchSessionsWithParticipants(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewList[types.UInt32]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSearchSessionsWithParticipants_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SearchSessionsWithParticipants = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.List[types.UInt32]) (types.List[game_session_types.GameSessionSearchWithParticipantsResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("gameSessionTypeID differs")
		}

		if !(got1.Equals(types.List[types.UInt32]{types.NewUInt32(3032), types.NewUInt32(4041)})) {
			t.Errorf("participantIDs differs")
		}

		return types.List[game_session_types.GameSessionSearchWithParticipantsResult]{func() game_session_types.GameSessionSearchWithParticipantsResult {
			v := game_session_types.NewGameSessionSearchWithParticipantsResult()
			v.ParticipantIDs = types.List[types.UInt32]{types.NewUInt32(8077), types.NewUInt32(9086)}
			return v
		}(), func() game_session_types.GameSessionSearchWithParticipantsResult {
			v := game_session_types.NewGameSessionSearchWithParticipantsResult()
			v.ParticipantIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.searchSessionsWithParticipants(nil, fakePacket{}, 9, types.NewUInt32(1014), types.List[types.UInt32]{types.NewUInt32(3032), types.NewUInt32(4041)})
	checkMessage(t, message, nexErr, ubi.MethodSearchSessionsWithParticipants, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x8d, 0x1f, 0x00, 0x00, 0x7e, 0x23, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x51, 0x2f, 0x00, 0x00, 0x42, 0x33, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SearchSessionsWithParticipants = func(packet nex.PacketInterface, _ types.UInt32, _ types.List[types.UInt32]) (types.List[game_session_types.GameSessionSearchWithParticipantsResult], *nex.Error) {
		return types.NewList[game_session_types.GameSessionSearchWithParticipantsResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.searchSessionsWithParticipants(nil, fakePacket{}, 9, types.NewUInt32(1014), types.List[types.UInt32]{types.NewUInt32(3032), types.NewUInt32(4041)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetSessions_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getSessions(nil, fakePacket{}, 7, types.List[game_session_types.GameSessionKey]{func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(3032)
		v.SessionID = types.NewUInt32(4041)
		return v
	}(), func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(6059)
		v.SessionID = types.NewUInt32(7068)
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodGetSessions, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetSessions_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetSessions = func(packet nex.PacketInterface, _ types.List[game_session_types.GameSessionKey]) (types.List[game_session_types.GameSessionSearchResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[game_session_types.GameSessionSearchResult](), nil
	}

	message, nexErr := commonProtocol.getSessions(errors.New("bad"), fakePacket{}, 7, types.NewList[game_session_types.GameSessionKey]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetSessions_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetSessions = func(packet nex.PacketInterface, got0 types.List[game_session_types.GameSessionKey]) (types.List[game_session_types.GameSessionSearchResult], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[game_session_types.GameSessionKey]{func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(3032)
			v.SessionID = types.NewUInt32(4041)
			return v
		}(), func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(6059)
			v.SessionID = types.NewUInt32(7068)
			return v
		}()})) {
			t.Errorf("gameSessionKeys differs")
		}

		return types.List[game_session_types.GameSessionSearchResult]{func() game_session_types.GameSessionSearchResult {
			v := game_session_types.NewGameSessionSearchResult()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(11104)
				v.SessionID = types.NewUInt32(12113)
				return v
			}()
			v.HostPID = types.NewUInt32(13122)
			v.HostURLs = types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.15"), types.NewStationURL("prudp:/address=10.0.0.16")}
			v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(19176)
				v.Value = types.NewUInt32(20185)
				return v
			}(), func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(22203)
				v.Value = types.NewUInt32(23212)
				return v
			}()}
			return v
		}(), func() game_session_types.GameSessionSearchResult {
			v := game_session_types.NewGameSessionSearchResult()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(26239)
				v.SessionID = types.NewUInt32(27248)
				return v
			}()
			v.HostPID = types.NewUInt32(28257)
			v.HostURLs = types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.30"), types.NewStationURL("prudp:/address=10.0.0.31")}
			v.Attributes = types.List[game_session_types.Property]{func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(34311)
				v.Value = types.NewUInt32(35320)
				return v
			}(), func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(37338)
				v.Value = types.NewUInt32(38347)
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getSessions(nil, fakePacket{}, 9, types.List[game_session_types.GameSessionKey]{func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(3032)
		v.SessionID = types.NewUInt32(4041)
		return v
	}(), func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(6059)
		v.SessionID = types.NewUInt32(7068)
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodGetSessions, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x60, 0x2b, 0x00, 0x00, 0x51, 0x2f, 0x00, 0x00, 0x42, 0x33, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x35, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x36, 0x00, 0x02, 0x00, 0x00, 0x00, 0xe8, 0x4a, 0x00, 0x00, 0xd9, 0x4e, 0x00, 0x00, 0xbb, 0x56, 0x00, 0x00, 0xac, 0x5a, 0x00, 0x00, 0x7f, 0x66, 0x00, 0x00, 0x70, 0x6a, 0x00, 0x00, 0x61, 0x6e, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x33, 0x30, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x33, 0x31, 0x00, 0x02, 0x00, 0x00, 0x00, 0x07, 0x86, 0x00, 0x00, 0xf8, 0x89, 0x00, 0x00, 0xda, 0x91, 0x00, 0x00, 0xcb, 0x95, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetSessions = func(packet nex.PacketInterface, _ types.List[game_session_types.GameSessionKey]) (types.List[game_session_types.GameSessionSearchResult], *nex.Error) {
		return types.NewList[game_session_types.GameSessionSearchResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getSessions(nil, fakePacket{}, 9, types.List[game_session_types.GameSessionKey]{func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(3032)
		v.SessionID = types.NewUInt32(4041)
		return v
	}(), func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(6059)
		v.SessionID = types.NewUInt32(7068)
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetParticipantsURLs_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getParticipantsURLs(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})
	checkMessage(t, message, nexErr, ubi.MethodGetParticipantsURLs, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetParticipantsURLs_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetParticipantsURLs = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey, _ types.List[types.UInt32]) (types.List[game_session_types.GameSessionParticipant], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[game_session_types.GameSessionParticipant](), nil
	}

	message, nexErr := commonProtocol.getParticipantsURLs(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey(), types.NewList[types.UInt32]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetParticipantsURLs_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetParticipantsURLs = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey, got1 types.List[types.UInt32]) (types.List[game_session_types.GameSessionParticipant], *nex.Error) {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		if !(got1.Equals(types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})) {
			t.Errorf("participantIDs differs")
		}

		return types.List[game_session_types.GameSessionParticipant]{func() game_session_types.GameSessionParticipant {
			v := game_session_types.NewGameSessionParticipant()
			v.PID = types.NewUInt32(9086)
			v.Name = types.NewString("s10")
			v.StationURLs = types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.12"), types.NewStationURL("prudp:/address=10.0.0.13")}
			return v
		}(), func() game_session_types.GameSessionParticipant {
			v := game_session_types.NewGameSessionParticipant()
			v.PID = types.NewUInt32(15140)
			v.Name = types.NewString("s16")
			v.StationURLs = types.List[types.StationURL]{types.NewStationURL("prudp:/address=10.0.0.18"), types.NewStationURL("prudp:/address=10.0.0.19")}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getParticipantsURLs(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})
	checkMessage(t, message, nexErr, ubi.MethodGetParticipantsURLs, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x7e, 0x23, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x30, 0x00, 0x02, 0x00, 0x00, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x32, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x33, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x36, 0x00, 0x02, 0x00, 0x00, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x38, 0x00, 0x19, 0x00, 0x70, 0x72, 0x75, 0x64, 0x70, 0x3a, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x3d, 0x31, 0x30, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x39, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetParticipantsURLs = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey, _ types.List[types.UInt32]) (types.List[game_session_types.GameSessionParticipant], *nex.Error) {
		return types.NewList[game_session_types.GameSessionParticipant](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getParticipantsURLs(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestMigrateSessionHost_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.migrateSessionHost(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestMigrateSessionHost_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.MigrateSessionHost = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.migrateSessionHost(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestMigrateSessionHost_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.MigrateSessionHost = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey) *nex.Error {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.migrateSessionHost(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodMigrateSessionHost, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.MigrateSessionHost = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.migrateSessionHost(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSplitSession_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.splitSession(nil, fakePacket{}, 7, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestSplitSession_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SplitSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) (game_session_types.GameSessionKey, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return game_session_types.NewGameSessionKey(), nil
	}

	message, nexErr := commonProtocol.splitSession(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSplitSession_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SplitSession = func(packet nex.PacketInterface, got0 game_session_types.GameSessionKey) (game_session_types.GameSessionKey, *nex.Error) {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(2023)
			v.SessionID = types.NewUInt32(3032)
			return v
		}())) {
			t.Errorf("gameSessionKey differs")
		}

		return func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(5050)
			v.SessionID = types.NewUInt32(6059)
			return v
		}(), nil
	}

	message, nexErr := commonProtocol.splitSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSplitSession, 9, []byte{0xba, 0x13, 0x00, 0x00, 0xab, 0x17, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SplitSession = func(packet nex.PacketInterface, _ game_session_types.GameSessionKey) (game_session_types.GameSessionKey, *nex.Error) {
		return game_session_types.NewGameSessionKey(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.splitSession(nil, fakePacket{}, 9, func() game_session_types.GameSessionKey {
		v := game_session_types.NewGameSessionKey()
		v.TypeID = types.NewUInt32(2023)
		v.SessionID = types.NewUInt32(3032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSearchSocialSessions_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.searchSocialSessions(nil, fakePacket{}, 7, func() game_session_types.GameSessionSocialQuery {
		v := game_session_types.NewGameSessionSocialQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.Parameters = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(6059)
			v.Value = types.NewUInt32(7068)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(9086)
			v.Value = types.NewUInt32(10095)
			return v
		}()}
		v.ParticipantIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchSocialSessions, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestSearchSocialSessions_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SearchSocialSessions = func(packet nex.PacketInterface, _ game_session_types.GameSessionSocialQuery) (types.List[game_session_types.GameSessionSearchWithParticipantsResult], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[game_session_types.GameSessionSearchWithParticipantsResult](), nil
	}

	message, nexErr := commonProtocol.searchSocialSessions(errors.New("bad"), fakePacket{}, 7, game_session_types.NewGameSessionSocialQuery())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSearchSocialSessions_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SearchSocialSessions = func(packet nex.PacketInterface, got0 game_session_types.GameSessionSocialQuery) (types.List[game_session_types.GameSessionSearchWithParticipantsResult], *nex.Error) {
		called = true

		if !(got0.Equals(func() game_session_types.GameSessionSocialQuery {
			v := game_session_types.NewGameSessionSocialQuery()
			v.TypeID = types.NewUInt32(2023)
			v.QueryID = types.NewUInt32(3032)
			v.Parameters = types.List[game_session_types.Property]{func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(6059)
				v.Value = types.NewUInt32(7068)
				return v
			}(), func() game_session_types.Property {
				v := game_session_types.NewProperty()
				v.ID = types.NewUInt32(9086)
				v.Value = types.NewUInt32(10095)
				return v
			}()}
			v.ParticipantIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
			return v
		}())) {
			t.Errorf("gameSessionSocialQuery differs")
		}

		return types.List[game_session_types.GameSessionSearchWithParticipantsResult]{func() game_session_types.GameSessionSearchWithParticipantsResult {
			v := game_session_types.NewGameSessionSearchWithParticipantsResult()
			v.ParticipantIDs = types.List[types.UInt32]{types.NewUInt32(17158), types.NewUInt32(18167)}
			return v
		}(), func() game_session_types.GameSessionSearchWithParticipantsResult {
			v := game_session_types.NewGameSessionSearchWithParticipantsResult()
			v.ParticipantIDs = types.List[types.UInt32]{types.NewUInt32(21194), types.NewUInt32(22203)}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.searchSocialSessions(nil, fakePacket{}, 9, func() game_session_types.GameSessionSocialQuery {
		v := game_session_types.NewGameSessionSocialQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.Parameters = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(6059)
			v.Value = types.NewUInt32(7068)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(9086)
			v.Value = types.NewUInt32(10095)
			return v
		}()}
		v.ParticipantIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchSocialSessions, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x06, 0x43, 0x00, 0x00, 0xf7, 0x46, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xca, 0x52, 0x00, 0x00, 0xbb, 0x56, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SearchSocialSessions = func(packet nex.PacketInterface, _ game_session_types.GameSessionSocialQuery) (types.List[game_session_types.GameSessionSearchWithParticipantsResult], *nex.Error) {
		return types.NewList[game_session_types.GameSessionSearchWithParticipantsResult](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.searchSocialSessions(nil, fakePacket{}, 9, func() game_session_types.GameSessionSocialQuery {
		v := game_session_types.NewGameSessionSocialQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.Parameters = types.List[game_session_types.Property]{func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(6059)
			v.Value = types.NewUInt32(7068)
			return v
		}(), func() game_session_types.Property {
			v := game_session_types.NewProperty()
			v.ID = types.NewUInt32(9086)
			v.Value = types.NewUInt32(10095)
			return v
		}()}
		v.ParticipantIDs = types.List[types.UInt32]{types.NewUInt32(12113), types.NewUInt32(13122)}
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReportUnsuccessfulJoinSessions_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.reportUnsuccessfulJoinSessions(nil, fakePacket{}, 7, types.List[game_session_types.GameSessionUnsuccessfulJoinSession]{func() game_session_types.GameSessionUnsuccessfulJoinSession {
		v := game_session_types.NewGameSessionUnsuccessfulJoinSession()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(4041)
			v.SessionID = types.NewUInt32(5050)
			return v
		}()
		v.ErrorCategory = types.NewUInt32(6059)
		v.ErrorCode = types.NewInt32(-218)
		return v
	}(), func() game_session_types.GameSessionUnsuccessfulJoinSession {
		v := game_session_types.NewGameSessionUnsuccessfulJoinSession()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(10095)
			v.SessionID = types.NewUInt32(11104)
			return v
		}()
		v.ErrorCategory = types.NewUInt32(12113)
		v.ErrorCode = types.NewInt32(-404)
		return v
	}()})
	checkNotImplemented(t, message, nexErr)
}

func TestReportUnsuccessfulJoinSessions_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReportUnsuccessfulJoinSessions = func(packet nex.PacketInterface, _ types.List[game_session_types.GameSessionUnsuccessfulJoinSession]) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.reportUnsuccessfulJoinSessions(errors.New("bad"), fakePacket{}, 7, types.NewList[game_session_types.GameSessionUnsuccessfulJoinSession]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReportUnsuccessfulJoinSessions_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReportUnsuccessfulJoinSessions = func(packet nex.PacketInterface, got0 types.List[game_session_types.GameSessionUnsuccessfulJoinSession]) *nex.Error {
		called = true

		if !(got0.Equals(types.List[game_session_types.GameSessionUnsuccessfulJoinSession]{func() game_session_types.GameSessionUnsuccessfulJoinSession {
			v := game_session_types.NewGameSessionUnsuccessfulJoinSession()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(4041)
				v.SessionID = types.NewUInt32(5050)
				return v
			}()
			v.ErrorCategory = types.NewUInt32(6059)
			v.ErrorCode = types.NewInt32(-218)
			return v
		}(), func() game_session_types.GameSessionUnsuccessfulJoinSession {
			v := game_session_types.NewGameSessionUnsuccessfulJoinSession()
			v.SessionKey = func() game_session_types.GameSessionKey {
				v := game_session_types.NewGameSessionKey()
				v.TypeID = types.NewUInt32(10095)
				v.SessionID = types.NewUInt32(11104)
				return v
			}()
			v.ErrorCategory = types.NewUInt32(12113)
			v.ErrorCode = types.NewInt32(-404)
			return v
		}()})) {
			t.Errorf("unsuccessfulJoinSessions differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.reportUnsuccessfulJoinSessions(nil, fakePacket{}, 9, types.List[game_session_types.GameSessionUnsuccessfulJoinSession]{func() game_session_types.GameSessionUnsuccessfulJoinSession {
		v := game_session_types.NewGameSessionUnsuccessfulJoinSession()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(4041)
			v.SessionID = types.NewUInt32(5050)
			return v
		}()
		v.ErrorCategory = types.NewUInt32(6059)
		v.ErrorCode = types.NewInt32(-218)
		return v
	}(), func() game_session_types.GameSessionUnsuccessfulJoinSession {
		v := game_session_types.NewGameSessionUnsuccessfulJoinSession()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(10095)
			v.SessionID = types.NewUInt32(11104)
			return v
		}()
		v.ErrorCategory = types.NewUInt32(12113)
		v.ErrorCode = types.NewInt32(-404)
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodReportUnsuccessfulJoinSessions, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReportUnsuccessfulJoinSessions = func(packet nex.PacketInterface, _ types.List[game_session_types.GameSessionUnsuccessfulJoinSession]) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.reportUnsuccessfulJoinSessions(nil, fakePacket{}, 9, types.List[game_session_types.GameSessionUnsuccessfulJoinSession]{func() game_session_types.GameSessionUnsuccessfulJoinSession {
		v := game_session_types.NewGameSessionUnsuccessfulJoinSession()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(4041)
			v.SessionID = types.NewUInt32(5050)
			return v
		}()
		v.ErrorCategory = types.NewUInt32(6059)
		v.ErrorCode = types.NewInt32(-218)
		return v
	}(), func() game_session_types.GameSessionUnsuccessfulJoinSession {
		v := game_session_types.NewGameSessionUnsuccessfulJoinSession()
		v.SessionKey = func() game_session_types.GameSessionKey {
			v := game_session_types.NewGameSessionKey()
			v.TypeID = types.NewUInt32(10095)
			v.SessionID = types.NewUInt32(11104)
			return v
		}()
		v.ErrorCategory = types.NewUInt32(12113)
		v.ErrorCode = types.NewInt32(-404)
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}
