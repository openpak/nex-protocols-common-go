package userstorageadmin

import (
	"bytes"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-storage-admin"
	user_storage_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-storage/types"
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

func TestGetContentsToModerate_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getContentsToModerate(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewUInt32(2023), types.NewUInt32(3032))
	checkNotImplemented(t, message, nexErr)
}

func TestGetContentsToModerate_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetContentsToModerate = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.UInt32) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[user_storage_types.UserContent](), types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.getContentsToModerate(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewUInt32(0), types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetContentsToModerate_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetContentsToModerate = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.UInt32, got2 types.UInt32) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("typeID differs")
		}

		if !(got1.Equals(types.NewUInt32(2023))) {
			t.Errorf("offset differs")
		}

		if !(got2.Equals(types.NewUInt32(3032))) {
			t.Errorf("size differs")
		}

		return types.List[user_storage_types.UserContent]{func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(7068)
				v.ContentID = types.NewUInt64(800000067)
				return v
			}()
			v.Pid = types.NewUInt32(9086)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(12113)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v13")}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(15140)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-21)}
				return v
			}()}
			return v
		}(), func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(19176)
				v.ContentID = types.NewUInt64(2000000151)
				return v
			}()
			v.Pid = types.NewUInt32(21194)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(24221)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v25")}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(27248)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-33)}
				return v
			}()}
			return v
		}()}, types.NewUInt32(29266), nil
	}

	message, nexErr := commonProtocol.getContentsToModerate(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.NewUInt32(3032))
	checkMessage(t, message, nexErr, ubi.MethodGetContentsToModerate, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x43, 0x08, 0xaf, 0x2f, 0x00, 0x00, 0x00, 0x00, 0x7e, 0x23, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x51, 0x2f, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x31, 0x33, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x01, 0xeb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xe8, 0x4a, 0x00, 0x00, 0x97, 0x94, 0x35, 0x77, 0x00, 0x00, 0x00, 0x00, 0xca, 0x52, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x9d, 0x5e, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x32, 0x35, 0x00, 0x70, 0x6a, 0x00, 0x00, 0x01, 0xdf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x52, 0x72, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetContentsToModerate = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.UInt32) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		return types.NewList[user_storage_types.UserContent](), types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getContentsToModerate(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.NewUInt32(3032))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestFlagContentAsVerified_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.flagContentAsVerified(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestFlagContentAsVerified_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.FlagContentAsVerified = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.flagContentAsVerified(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestFlagContentAsVerified_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.FlagContentAsVerified = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) *nex.Error {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.flagContentAsVerified(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodFlagContentAsVerified, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.FlagContentAsVerified = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.flagContentAsVerified(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestBanContent_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.banContent(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestBanContent_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.BanContent = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.banContent(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestBanContent_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.BanContent = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) *nex.Error {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.banContent(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodBanContent, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.BanContent = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.banContent(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestBanUser_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.banUser(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewString("s2"), types.NewBool(true), types.NewDateTime(4000089))
	checkNotImplemented(t, message, nexErr)
}

func TestBanUser_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.BanUser = func(packet nex.PacketInterface, _ types.UInt32, _ types.String, _ types.Bool, _ types.DateTime) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.banUser(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewString(""), types.NewBool(false), types.NewDateTime(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestBanUser_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.BanUser = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.String, got2 types.Bool, got3 types.DateTime) *nex.Error {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("pid differs")
		}

		if !(got1.Equals(types.NewString("s2"))) {
			t.Errorf("reason differs")
		}

		if !(got2.Equals(types.NewBool(true))) {
			t.Errorf("banContents differs")
		}

		if !(got3.Equals(types.NewDateTime(4000089))) {
			t.Errorf("expireDate differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.banUser(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewString("s2"), types.NewBool(true), types.NewDateTime(4000089))
	checkMessage(t, message, nexErr, ubi.MethodBanUser, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.BanUser = func(packet nex.PacketInterface, _ types.UInt32, _ types.String, _ types.Bool, _ types.DateTime) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.banUser(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewString("s2"), types.NewBool(true), types.NewDateTime(4000089))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestBanUserFromContentType_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.banUserFromContentType(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewUInt32(2023), types.NewString("s3"), types.NewBool(false), types.NewDateTime(5000092))
	checkNotImplemented(t, message, nexErr)
}

func TestBanUserFromContentType_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.BanUserFromContentType = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.String, _ types.Bool, _ types.DateTime) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.banUserFromContentType(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewUInt32(0), types.NewString(""), types.NewBool(false), types.NewDateTime(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestBanUserFromContentType_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.BanUserFromContentType = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.UInt32, got2 types.String, got3 types.Bool, got4 types.DateTime) *nex.Error {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("typeID differs")
		}

		if !(got1.Equals(types.NewUInt32(2023))) {
			t.Errorf("pid differs")
		}

		if !(got2.Equals(types.NewString("s3"))) {
			t.Errorf("reason differs")
		}

		if !(got3.Equals(types.NewBool(false))) {
			t.Errorf("banContents differs")
		}

		if !(got4.Equals(types.NewDateTime(5000092))) {
			t.Errorf("expireDate differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.banUserFromContentType(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.NewString("s3"), types.NewBool(false), types.NewDateTime(5000092))
	checkMessage(t, message, nexErr, ubi.MethodBanUserFromContentType, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.BanUserFromContentType = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.String, _ types.Bool, _ types.DateTime) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.banUserFromContentType(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.NewString("s3"), types.NewBool(false), types.NewDateTime(5000092))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestUnbanUser_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.unbanUser(nil, fakePacket{}, 7, types.NewUInt32(1014))
	checkNotImplemented(t, message, nexErr)
}

func TestUnbanUser_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.UnbanUser = func(packet nex.PacketInterface, _ types.UInt32) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.unbanUser(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestUnbanUser_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.UnbanUser = func(packet nex.PacketInterface, got0 types.UInt32) *nex.Error {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("pid differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.unbanUser(nil, fakePacket{}, 9, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodUnbanUser, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.UnbanUser = func(packet nex.PacketInterface, _ types.UInt32) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.unbanUser(nil, fakePacket{}, 9, types.NewUInt32(1014))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestUnbanUserFromContentType_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.unbanUserFromContentType(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewUInt32(2023))
	checkNotImplemented(t, message, nexErr)
}

func TestUnbanUserFromContentType_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.UnbanUserFromContentType = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.unbanUserFromContentType(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestUnbanUserFromContentType_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.UnbanUserFromContentType = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.UInt32) *nex.Error {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("typeID differs")
		}

		if !(got1.Equals(types.NewUInt32(2023))) {
			t.Errorf("pid differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.unbanUserFromContentType(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023))
	checkMessage(t, message, nexErr, ubi.MethodUnbanUserFromContentType, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.UnbanUserFromContentType = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.unbanUserFromContentType(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetContentsToModerateWithThreshold_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getContentsToModerateWithThreshold(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewUInt32(2023), types.NewUInt32(3032), types.NewUInt32(4041))
	checkNotImplemented(t, message, nexErr)
}

func TestGetContentsToModerateWithThreshold_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetContentsToModerateWithThreshold = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.UInt32, _ types.UInt32) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[user_storage_types.UserContent](), types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.getContentsToModerateWithThreshold(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewUInt32(0), types.NewUInt32(0), types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetContentsToModerateWithThreshold_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetContentsToModerateWithThreshold = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.UInt32, got2 types.UInt32, got3 types.UInt32) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("typeID differs")
		}

		if !(got1.Equals(types.NewUInt32(2023))) {
			t.Errorf("threshold differs")
		}

		if !(got2.Equals(types.NewUInt32(3032))) {
			t.Errorf("offset differs")
		}

		if !(got3.Equals(types.NewUInt32(4041))) {
			t.Errorf("size differs")
		}

		return types.List[user_storage_types.UserContent]{func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(8077)
				v.ContentID = types.NewUInt64(900000074)
				return v
			}()
			v.Pid = types.NewUInt32(10095)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(13122)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-19)}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(16149)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v17")}
				return v
			}()}
			return v
		}(), func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(20185)
				v.ContentID = types.NewUInt64(2100000158)
				return v
			}()
			v.Pid = types.NewUInt32(22203)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(25230)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-31)}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(28257)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v29")}
				return v
			}()}
			return v
		}()}, types.NewUInt32(30275), nil
	}

	message, nexErr := commonProtocol.getContentsToModerateWithThreshold(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.NewUInt32(3032), types.NewUInt32(4041))
	checkMessage(t, message, nexErr, ubi.MethodGetContentsToModerateWithThreshold, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x8d, 0x1f, 0x00, 0x00, 0x4a, 0xe9, 0xa4, 0x35, 0x00, 0x00, 0x00, 0x00, 0x6f, 0x27, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x42, 0x33, 0x00, 0x00, 0x01, 0xed, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x15, 0x3f, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x31, 0x37, 0x00, 0xd9, 0x4e, 0x00, 0x00, 0x9e, 0x75, 0x2b, 0x7d, 0x00, 0x00, 0x00, 0x00, 0xbb, 0x56, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x8e, 0x62, 0x00, 0x00, 0x01, 0xe1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x61, 0x6e, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x32, 0x39, 0x00, 0x43, 0x76, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetContentsToModerateWithThreshold = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32, _ types.UInt32, _ types.UInt32) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		return types.NewList[user_storage_types.UserContent](), types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getContentsToModerateWithThreshold(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023), types.NewUInt32(3032), types.NewUInt32(4041))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestUpdateMetaData_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.updateMetaData(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(9086)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-15)}
		return v
	}()})
	checkNotImplemented(t, message, nexErr)
}

func TestUpdateMetaData_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.UpdateMetaData = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey, _ types.List[user_storage_types.ContentProperty]) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.updateMetaData(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey(), types.NewList[user_storage_types.ContentProperty]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestUpdateMetaData_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.UpdateMetaData = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey, got1 types.List[user_storage_types.ContentProperty]) *nex.Error {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		if !(got1.Equals(types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(6059)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(9086)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-15)}
			return v
		}()})) {
			t.Errorf("properties differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.updateMetaData(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(9086)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-15)}
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodUpdateMetaData, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.UpdateMetaData = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey, _ types.List[user_storage_types.ContentProperty]) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.updateMetaData(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(9086)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-15)}
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}
