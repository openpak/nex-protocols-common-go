package offlinegamenotifications

import (
	"bytes"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/notifications/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/offline-game-notifications"
	offline_game_notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/offline-game-notifications/types"
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

func TestPollNotifications_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.pollNotifications(nil, fakePacket{}, 7)
	checkMessage(t, message, nexErr, ubi.MethodPollNotifications, 7, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestPollNotifications_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.PollNotifications = func(packet nex.PacketInterface) (types.List[notifications_types.NotificationEvent], types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[notifications_types.NotificationEvent](), types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.pollNotifications(errors.New("bad"), fakePacket{}, 7)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestPollNotifications_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.PollNotifications = func(packet nex.PacketInterface) (types.List[notifications_types.NotificationEvent], types.UInt32, *nex.Error) {
		called = true

		return types.List[notifications_types.NotificationEvent]{func() notifications_types.NotificationEvent {
			v := notifications_types.NewNotificationEvent()
			v.PIDSource = types.NewPID(22)
			v.Type = 26
			v.Param1 = types.NewUInt64(34)
			v.Param2 = types.NewUInt64(38)
			v.StrParam = types.NewString("n2")
			return v
		}(), func() notifications_types.NotificationEvent {
			v := notifications_types.NewNotificationEvent()
			v.PIDSource = types.NewPID(33)
			v.Type = 39
			v.Param1 = types.NewUInt64(51)
			v.Param2 = types.NewUInt64(57)
			v.StrParam = types.NewString("n3")
			return v
		}()}, types.NewUInt32(4041), nil
	}

	message, nexErr := commonProtocol.pollNotifications(nil, fakePacket{}, 9)
	checkMessage(t, message, nexErr, ubi.MethodPollNotifications, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x1a, 0x00, 0x00, 0x00, 0x22, 0x00, 0x00, 0x00, 0x26, 0x00, 0x00, 0x00, 0x03, 0x00, 0x6e, 0x32, 0x00, 0x21, 0x00, 0x00, 0x00, 0x27, 0x00, 0x00, 0x00, 0x33, 0x00, 0x00, 0x00, 0x39, 0x00, 0x00, 0x00, 0x03, 0x00, 0x6e, 0x33, 0x00, 0xc9, 0x0f, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.PollNotifications = func(packet nex.PacketInterface) (types.List[notifications_types.NotificationEvent], types.UInt32, *nex.Error) {
		return types.NewList[notifications_types.NotificationEvent](), types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.pollNotifications(nil, fakePacket{}, 9)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestPollSpecificOfflineNotifications_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.pollSpecificOfflineNotifications(nil, fakePacket{}, 7, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	checkMessage(t, message, nexErr, ubi.MethodPollSpecificOfflineNotifications, 7, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestPollSpecificOfflineNotifications_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.PollSpecificOfflineNotifications = func(packet nex.PacketInterface, _ types.List[types.UInt32]) (types.List[offline_game_notifications_types.TimedNotification], types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[offline_game_notifications_types.TimedNotification](), types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.pollSpecificOfflineNotifications(errors.New("bad"), fakePacket{}, 7, types.NewList[types.UInt32]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestPollSpecificOfflineNotifications_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.PollSpecificOfflineNotifications = func(packet nex.PacketInterface, got0 types.List[types.UInt32]) (types.List[offline_game_notifications_types.TimedNotification], types.UInt32, *nex.Error) {
		called = true

		if !(got0.Equals(types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})) {
			t.Errorf("majortype differs")
		}

		return types.List[offline_game_notifications_types.TimedNotification]{func() offline_game_notifications_types.TimedNotification {
			v := offline_game_notifications_types.NewTimedNotification()
			v.Timestamp = types.NewDateTime(6000095)
			v.Notification = func() notifications_types.NotificationEvent {
				v := notifications_types.NewNotificationEvent()
				v.PIDSource = types.NewPID(77)
				v.Type = 91
				v.Param1 = types.NewUInt64(119)
				v.Param2 = types.NewUInt64(133)
				v.StrParam = types.NewString("n7")
				return v
			}()
			return v
		}(), func() offline_game_notifications_types.TimedNotification {
			v := offline_game_notifications_types.NewTimedNotification()
			v.Timestamp = types.NewDateTime(9000104)
			v.Notification = func() notifications_types.NotificationEvent {
				v := notifications_types.NewNotificationEvent()
				v.PIDSource = types.NewPID(110)
				v.Type = 130
				v.Param1 = types.NewUInt64(170)
				v.Param2 = types.NewUInt64(190)
				v.StrParam = types.NewString("n10")
				return v
			}()
			return v
		}()}, types.NewUInt32(11104), nil
	}

	message, nexErr := commonProtocol.pollSpecificOfflineNotifications(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	checkMessage(t, message, nexErr, ubi.MethodPollSpecificOfflineNotifications, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xdf, 0x8d, 0x5b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4d, 0x00, 0x00, 0x00, 0x5b, 0x00, 0x00, 0x00, 0x77, 0x00, 0x00, 0x00, 0x85, 0x00, 0x00, 0x00, 0x03, 0x00, 0x6e, 0x37, 0x00, 0xa8, 0x54, 0x89, 0x00, 0x00, 0x00, 0x00, 0x00, 0x6e, 0x00, 0x00, 0x00, 0x82, 0x00, 0x00, 0x00, 0xaa, 0x00, 0x00, 0x00, 0xbe, 0x00, 0x00, 0x00, 0x04, 0x00, 0x6e, 0x31, 0x30, 0x00, 0x60, 0x2b, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.PollSpecificOfflineNotifications = func(packet nex.PacketInterface, _ types.List[types.UInt32]) (types.List[offline_game_notifications_types.TimedNotification], types.UInt32, *nex.Error) {
		return types.NewList[offline_game_notifications_types.TimedNotification](), types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.pollSpecificOfflineNotifications(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestPollAnyOfflineNotifications_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.pollAnyOfflineNotifications(nil, fakePacket{}, 7)
	checkMessage(t, message, nexErr, ubi.MethodPollAnyOfflineNotifications, 7, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestPollAnyOfflineNotifications_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.PollAnyOfflineNotifications = func(packet nex.PacketInterface) (types.List[offline_game_notifications_types.TimedNotification], types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[offline_game_notifications_types.TimedNotification](), types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.pollAnyOfflineNotifications(errors.New("bad"), fakePacket{}, 7)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestPollAnyOfflineNotifications_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.PollAnyOfflineNotifications = func(packet nex.PacketInterface) (types.List[offline_game_notifications_types.TimedNotification], types.UInt32, *nex.Error) {
		called = true

		return types.List[offline_game_notifications_types.TimedNotification]{func() offline_game_notifications_types.TimedNotification {
			v := offline_game_notifications_types.NewTimedNotification()
			v.Timestamp = types.NewDateTime(3000086)
			v.Notification = func() notifications_types.NotificationEvent {
				v := notifications_types.NewNotificationEvent()
				v.PIDSource = types.NewPID(44)
				v.Type = 52
				v.Param1 = types.NewUInt64(68)
				v.Param2 = types.NewUInt64(76)
				v.StrParam = types.NewString("n4")
				return v
			}()
			return v
		}(), func() offline_game_notifications_types.TimedNotification {
			v := offline_game_notifications_types.NewTimedNotification()
			v.Timestamp = types.NewDateTime(6000095)
			v.Notification = func() notifications_types.NotificationEvent {
				v := notifications_types.NewNotificationEvent()
				v.PIDSource = types.NewPID(77)
				v.Type = 91
				v.Param1 = types.NewUInt64(119)
				v.Param2 = types.NewUInt64(133)
				v.StrParam = types.NewString("n7")
				return v
			}()
			return v
		}()}, types.NewUInt32(8077), nil
	}

	message, nexErr := commonProtocol.pollAnyOfflineNotifications(nil, fakePacket{}, 9)
	checkMessage(t, message, nexErr, ubi.MethodPollAnyOfflineNotifications, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x16, 0xc7, 0x2d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00, 0x34, 0x00, 0x00, 0x00, 0x44, 0x00, 0x00, 0x00, 0x4c, 0x00, 0x00, 0x00, 0x03, 0x00, 0x6e, 0x34, 0x00, 0xdf, 0x8d, 0x5b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4d, 0x00, 0x00, 0x00, 0x5b, 0x00, 0x00, 0x00, 0x77, 0x00, 0x00, 0x00, 0x85, 0x00, 0x00, 0x00, 0x03, 0x00, 0x6e, 0x37, 0x00, 0x8d, 0x1f, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.PollAnyOfflineNotifications = func(packet nex.PacketInterface) (types.List[offline_game_notifications_types.TimedNotification], types.UInt32, *nex.Error) {
		return types.NewList[offline_game_notifications_types.TimedNotification](), types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.pollAnyOfflineNotifications(nil, fakePacket{}, 9)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}
