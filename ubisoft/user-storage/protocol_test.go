package userstorage

import (
	"bytes"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-storage"
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

func TestSearchContents_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.searchContents(nil, fakePacket{}, 7, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(7068)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-13)}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchContents, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestSearchContents_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SearchContents = func(packet nex.PacketInterface, _ user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[user_storage_types.UserContent](), nil
	}

	message, nexErr := commonProtocol.searchContents(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserStorageQuery())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSearchContents_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SearchContents = func(packet nex.PacketInterface, got0 user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], *nex.Error) {
		called = true

		if !(got0.Equals(func() user_storage_types.UserStorageQuery {
			v := user_storage_types.NewUserStorageQuery()
			v.TypeID = types.NewUInt32(2023)
			v.QueryID = types.NewUInt32(3032)
			v.ResultRange = types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)}
			v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(7068)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-13)}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(10095)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
				return v
			}()}
			return v
		}())) {
			t.Errorf("query differs")
		}

		return types.List[user_storage_types.UserContent]{func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(15140)
				v.ContentID = types.NewUInt64(1600000123)
				return v
			}()
			v.Pid = types.NewUInt32(17158)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(20185)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v21")}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(23212)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-29)}
				return v
			}()}
			return v
		}(), func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(27248)
				v.ContentID = types.NewUInt64(2800000207)
				return v
			}()
			v.Pid = types.NewUInt32(29266)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(32293)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v33")}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(35320)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-41)}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.searchContents(nil, fakePacket{}, 9, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(7068)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-13)}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchContents, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x7b, 0x10, 0x5e, 0x5f, 0x00, 0x00, 0x00, 0x00, 0x06, 0x43, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xd9, 0x4e, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x32, 0x31, 0x00, 0xac, 0x5a, 0x00, 0x00, 0x01, 0xe3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x70, 0x6a, 0x00, 0x00, 0xcf, 0x9c, 0xe4, 0xa6, 0x00, 0x00, 0x00, 0x00, 0x52, 0x72, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x25, 0x7e, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x33, 0x00, 0xf8, 0x89, 0x00, 0x00, 0x01, 0xd7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SearchContents = func(packet nex.PacketInterface, _ user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], *nex.Error) {
		return types.NewList[user_storage_types.UserContent](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.searchContents(nil, fakePacket{}, 9, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(7068)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-13)}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}()}
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSearchContentsWithTotal_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.searchContentsWithTotal(nil, fakePacket{}, 7, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(7068)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-13)}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchContentsWithTotal, 7, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestSearchContentsWithTotal_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SearchContentsWithTotal = func(packet nex.PacketInterface, _ user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[user_storage_types.UserContent](), types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.searchContentsWithTotal(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserStorageQuery())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSearchContentsWithTotal_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SearchContentsWithTotal = func(packet nex.PacketInterface, got0 user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		called = true

		if !(got0.Equals(func() user_storage_types.UserStorageQuery {
			v := user_storage_types.NewUserStorageQuery()
			v.TypeID = types.NewUInt32(2023)
			v.QueryID = types.NewUInt32(3032)
			v.ResultRange = types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)}
			v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(7068)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-13)}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(10095)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
				return v
			}()}
			return v
		}())) {
			t.Errorf("query differs")
		}

		return types.List[user_storage_types.UserContent]{func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(15140)
				v.ContentID = types.NewUInt64(1600000123)
				return v
			}()
			v.Pid = types.NewUInt32(17158)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(20185)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v21")}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(23212)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-29)}
				return v
			}()}
			return v
		}(), func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(27248)
				v.ContentID = types.NewUInt64(2800000207)
				return v
			}()
			v.Pid = types.NewUInt32(29266)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(32293)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v33")}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(35320)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-41)}
				return v
			}()}
			return v
		}()}, types.NewUInt32(37338), nil
	}

	message, nexErr := commonProtocol.searchContentsWithTotal(nil, fakePacket{}, 9, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(7068)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-13)}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchContentsWithTotal, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x7b, 0x10, 0x5e, 0x5f, 0x00, 0x00, 0x00, 0x00, 0x06, 0x43, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xd9, 0x4e, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x32, 0x31, 0x00, 0xac, 0x5a, 0x00, 0x00, 0x01, 0xe3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x70, 0x6a, 0x00, 0x00, 0xcf, 0x9c, 0xe4, 0xa6, 0x00, 0x00, 0x00, 0x00, 0x52, 0x72, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x25, 0x7e, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x33, 0x00, 0xf8, 0x89, 0x00, 0x00, 0x01, 0xd7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xda, 0x91, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SearchContentsWithTotal = func(packet nex.PacketInterface, _ user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		return types.NewList[user_storage_types.UserContent](), types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.searchContentsWithTotal(nil, fakePacket{}, 9, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(2023)
		v.QueryID = types.NewUInt32(3032)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(12), Length: types.NewUInt32(22)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(7068)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-13)}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}()}
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestDeleteContent_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.deleteContent(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestDeleteContent_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.DeleteContent = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.deleteContent(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestDeleteContent_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.DeleteContent = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) *nex.Error {
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

	message, nexErr := commonProtocol.deleteContent(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodDeleteContent, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.DeleteContent = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.deleteContent(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSaveMetaData_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.saveMetaData(nil, fakePacket{}, 7, types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(3032)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}()}, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(9086)
		v.ContentID = types.NewUInt64(1000000081)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestSaveMetaData_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SaveMetaData = func(packet nex.PacketInterface, _ types.List[user_storage_types.ContentProperty], _ user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return user_storage_types.NewUserContentKey(), nil
	}

	message, nexErr := commonProtocol.saveMetaData(errors.New("bad"), fakePacket{}, 7, types.NewList[user_storage_types.ContentProperty](), user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSaveMetaData_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SaveMetaData = func(packet nex.PacketInterface, got0 types.List[user_storage_types.ContentProperty], got1 user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error) {
		called = true

		if !(got0.Equals(types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(3032)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(6059)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
			return v
		}()})) {
			t.Errorf("properties differs")
		}

		if !(got1.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(9086)
			v.ContentID = types.NewUInt64(1000000081)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(12113)
			v.ContentID = types.NewUInt64(1300000102)
			return v
		}(), nil
	}

	message, nexErr := commonProtocol.saveMetaData(nil, fakePacket{}, 9, types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(3032)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}()}, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(9086)
		v.ContentID = types.NewUInt64(1000000081)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSaveMetaData, 9, []byte{0x51, 0x2f, 0x00, 0x00, 0x66, 0x6d, 0x7c, 0x4d, 0x00, 0x00, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SaveMetaData = func(packet nex.PacketInterface, _ types.List[user_storage_types.ContentProperty], _ user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error) {
		return user_storage_types.NewUserContentKey(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.saveMetaData(nil, fakePacket{}, 9, types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(3032)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}()}, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(9086)
		v.ContentID = types.NewUInt64(1000000081)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSaveContentDB_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.saveContentDB(nil, fakePacket{}, 7, types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(3032)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}()}, types.NewBuffer([]byte{8, 9, 10}), func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(10095)
		v.ContentID = types.NewUInt64(1100000088)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestSaveContentDB_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SaveContentDB = func(packet nex.PacketInterface, _ types.List[user_storage_types.ContentProperty], _ types.Buffer, _ user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return user_storage_types.NewUserContentKey(), nil
	}

	message, nexErr := commonProtocol.saveContentDB(errors.New("bad"), fakePacket{}, 7, types.NewList[user_storage_types.ContentProperty](), types.NewBuffer(nil), user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSaveContentDB_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SaveContentDB = func(packet nex.PacketInterface, got0 types.List[user_storage_types.ContentProperty], got1 types.Buffer, got2 user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error) {
		called = true

		if !(got0.Equals(types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(3032)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(6059)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
			return v
		}()})) {
			t.Errorf("properties differs")
		}

		if !(got1.Equals(types.NewBuffer([]byte{8, 9, 10}))) {
			t.Errorf("data differs")
		}

		if !(got2.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(10095)
			v.ContentID = types.NewUInt64(1100000088)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(13122)
			v.ContentID = types.NewUInt64(1400000109)
			return v
		}(), nil
	}

	message, nexErr := commonProtocol.saveContentDB(nil, fakePacket{}, 9, types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(3032)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}()}, types.NewBuffer([]byte{8, 9, 10}), func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(10095)
		v.ContentID = types.NewUInt64(1100000088)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSaveContentDB, 9, []byte{0x42, 0x33, 0x00, 0x00, 0x6d, 0x4e, 0x72, 0x53, 0x00, 0x00, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SaveContentDB = func(packet nex.PacketInterface, _ types.List[user_storage_types.ContentProperty], _ types.Buffer, _ user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error) {
		return user_storage_types.NewUserContentKey(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.saveContentDB(nil, fakePacket{}, 9, types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(3032)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}()}, types.NewBuffer([]byte{8, 9, 10}), func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(10095)
		v.ContentID = types.NewUInt64(1100000088)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSaveContentAndGetUploadInfo_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.saveContentAndGetUploadInfo(nil, fakePacket{}, 7, types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(3032)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}()}, types.NewUInt32(8077), func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(10095)
		v.ContentID = types.NewUInt64(1100000088)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestSaveContentAndGetUploadInfo_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SaveContentAndGetUploadInfo = func(packet nex.PacketInterface, _ types.List[user_storage_types.ContentProperty], _ types.UInt32, _ user_storage_types.UserContentKey) (user_storage_types.UserContentURL, types.UInt64, types.List[types.String], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return user_storage_types.NewUserContentURL(), types.NewUInt64(0), types.NewList[types.String](), nil
	}

	message, nexErr := commonProtocol.saveContentAndGetUploadInfo(errors.New("bad"), fakePacket{}, 7, types.NewList[user_storage_types.ContentProperty](), types.NewUInt32(0), user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSaveContentAndGetUploadInfo_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SaveContentAndGetUploadInfo = func(packet nex.PacketInterface, got0 types.List[user_storage_types.ContentProperty], got1 types.UInt32, got2 user_storage_types.UserContentKey) (user_storage_types.UserContentURL, types.UInt64, types.List[types.String], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(3032)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(6059)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
			return v
		}()})) {
			t.Errorf("properties differs")
		}

		if !(got1.Equals(types.NewUInt32(8077))) {
			t.Errorf("size differs")
		}

		if !(got2.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(10095)
			v.ContentID = types.NewUInt64(1100000088)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return func() user_storage_types.UserContentURL {
			v := user_storage_types.NewUserContentURL()
			v.Protocol = types.NewString("s13")
			v.Host = types.NewString("s14")
			v.Path = types.NewString("s15")
			return v
		}(), types.NewUInt64(1600000123), types.List[types.String]{types.NewString("s18"), types.NewString("s19")}, nil
	}

	message, nexErr := commonProtocol.saveContentAndGetUploadInfo(nil, fakePacket{}, 9, types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(3032)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}()}, types.NewUInt32(8077), func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(10095)
		v.ContentID = types.NewUInt64(1100000088)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSaveContentAndGetUploadInfo, 9, []byte{0x04, 0x00, 0x73, 0x31, 0x33, 0x00, 0x04, 0x00, 0x73, 0x31, 0x34, 0x00, 0x04, 0x00, 0x73, 0x31, 0x35, 0x00, 0x7b, 0x10, 0x5e, 0x5f, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x38, 0x00, 0x04, 0x00, 0x73, 0x31, 0x39, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SaveContentAndGetUploadInfo = func(packet nex.PacketInterface, _ types.List[user_storage_types.ContentProperty], _ types.UInt32, _ user_storage_types.UserContentKey) (user_storage_types.UserContentURL, types.UInt64, types.List[types.String], *nex.Error) {
		return user_storage_types.NewUserContentURL(), types.NewUInt64(0), types.NewList[types.String](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.saveContentAndGetUploadInfo(nil, fakePacket{}, 9, types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(3032)
		v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-9)}
		return v
	}(), func() user_storage_types.ContentProperty {
		v := user_storage_types.NewContentProperty()
		v.ID = types.NewUInt32(6059)
		v.Value = types.Variant{TypeID: 4, Type: types.NewString("v7")}
		return v
	}()}, types.NewUInt32(8077), func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(10095)
		v.ContentID = types.NewUInt64(1100000088)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestUploadEnd_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.uploadEnd(nil, fakePacket{}, 7, types.NewUInt64(100000018), types.NewBool(false), func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(4041)
		v.ContentID = types.NewUInt64(500000046)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestUploadEnd_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.UploadEnd = func(packet nex.PacketInterface, _ types.UInt64, _ types.Bool, _ user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return user_storage_types.NewUserContentKey(), nil
	}

	message, nexErr := commonProtocol.uploadEnd(errors.New("bad"), fakePacket{}, 7, types.NewUInt64(0), types.NewBool(false), user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestUploadEnd_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.UploadEnd = func(packet nex.PacketInterface, got0 types.UInt64, got1 types.Bool, got2 user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt64(100000018))) {
			t.Errorf("pendingID differs")
		}

		if !(got1.Equals(types.NewBool(false))) {
			t.Errorf("result differs")
		}

		if !(got2.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(4041)
			v.ContentID = types.NewUInt64(500000046)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(7068)
			v.ContentID = types.NewUInt64(800000067)
			return v
		}(), nil
	}

	message, nexErr := commonProtocol.uploadEnd(nil, fakePacket{}, 9, types.NewUInt64(100000018), types.NewBool(false), func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(4041)
		v.ContentID = types.NewUInt64(500000046)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodUploadEnd, 9, []byte{0x9c, 0x1b, 0x00, 0x00, 0x43, 0x08, 0xaf, 0x2f, 0x00, 0x00, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.UploadEnd = func(packet nex.PacketInterface, _ types.UInt64, _ types.Bool, _ user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error) {
		return user_storage_types.NewUserContentKey(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.uploadEnd(nil, fakePacket{}, 9, types.NewUInt64(100000018), types.NewBool(false), func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(4041)
		v.ContentID = types.NewUInt64(500000046)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetContentDB_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getContentDB(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestGetContentDB_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetContentDB = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (types.Buffer, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewBuffer(nil), nil
	}

	message, nexErr := commonProtocol.getContentDB(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetContentDB_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetContentDB = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) (types.Buffer, *nex.Error) {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return types.NewBuffer([]byte{4, 5, 6}), nil
	}

	message, nexErr := commonProtocol.getContentDB(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetContentDB, 9, []byte{0x03, 0x00, 0x00, 0x00, 0x04, 0x05, 0x06})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetContentDB = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (types.Buffer, *nex.Error) {
		return types.NewBuffer(nil), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getContentDB(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetContentURL_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getContentURL(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestGetContentURL_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetContentURL = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (user_storage_types.UserContentURL, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return user_storage_types.NewUserContentURL(), nil
	}

	message, nexErr := commonProtocol.getContentURL(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetContentURL_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetContentURL = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) (user_storage_types.UserContentURL, *nex.Error) {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return func() user_storage_types.UserContentURL {
			v := user_storage_types.NewUserContentURL()
			v.Protocol = types.NewString("s5")
			v.Host = types.NewString("s6")
			v.Path = types.NewString("s7")
			return v
		}(), nil
	}

	message, nexErr := commonProtocol.getContentURL(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetContentURL, 9, []byte{0x03, 0x00, 0x73, 0x35, 0x00, 0x03, 0x00, 0x73, 0x36, 0x00, 0x03, 0x00, 0x73, 0x37, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetContentURL = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (user_storage_types.UserContentURL, *nex.Error) {
		return user_storage_types.NewUserContentURL(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getContentURL(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetSlotCount_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getSlotCount(nil, fakePacket{}, 7, types.NewUInt32(1014))
	checkNotImplemented(t, message, nexErr)
}

func TestGetSlotCount_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetSlotCount = func(packet nex.PacketInterface, _ types.UInt32) (user_storage_types.UserSlotCount, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return user_storage_types.NewUserSlotCount(), nil
	}

	message, nexErr := commonProtocol.getSlotCount(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetSlotCount_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetSlotCount = func(packet nex.PacketInterface, got0 types.UInt32) (user_storage_types.UserSlotCount, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("typeID differs")
		}

		return func() user_storage_types.UserSlotCount {
			v := user_storage_types.NewUserSlotCount()
			v.Pid = types.NewUInt32(3032)
			v.TypeID = types.NewUInt32(4041)
			v.UsedSlots = types.NewUInt32(5050)
			v.TotalSlots = types.NewUInt32(6059)
			return v
		}(), nil
	}

	message, nexErr := commonProtocol.getSlotCount(nil, fakePacket{}, 9, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodGetSlotCount, 9, []byte{0xd8, 0x0b, 0x00, 0x00, 0xc9, 0x0f, 0x00, 0x00, 0xba, 0x13, 0x00, 0x00, 0xab, 0x17, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetSlotCount = func(packet nex.PacketInterface, _ types.UInt32) (user_storage_types.UserSlotCount, *nex.Error) {
		return user_storage_types.NewUserSlotCount(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getSlotCount(nil, fakePacket{}, 9, types.NewUInt32(1014))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetMetaData_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getMetaData(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestGetMetaData_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetMetaData = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (user_storage_types.UserContent, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return user_storage_types.NewUserContent(), nil
	}

	message, nexErr := commonProtocol.getMetaData(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetMetaData_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetMetaData = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) (user_storage_types.UserContent, *nex.Error) {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(6059)
				v.ContentID = types.NewUInt64(700000060)
				return v
			}()
			v.Pid = types.NewUInt32(8077)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(11104)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-17)}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(14131)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v15")}
				return v
			}()}
			return v
		}(), nil
	}

	message, nexErr := commonProtocol.getMetaData(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetMetaData, 9, []byte{0xab, 0x17, 0x00, 0x00, 0x3c, 0x27, 0xb9, 0x29, 0x00, 0x00, 0x00, 0x00, 0x8d, 0x1f, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x60, 0x2b, 0x00, 0x00, 0x01, 0xef, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x33, 0x37, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x31, 0x35, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetMetaData = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (user_storage_types.UserContent, *nex.Error) {
		return user_storage_types.NewUserContent(), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getMetaData(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestLike_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.like(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestLike_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.Like = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.like(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestLike_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.Like = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) *nex.Error {
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

	message, nexErr := commonProtocol.like(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodLike, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.Like = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.like(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestUnlike_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.unlike(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestUnlike_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.Unlike = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.unlike(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestUnlike_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.Unlike = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) *nex.Error {
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

	message, nexErr := commonProtocol.unlike(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodUnlike, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.Unlike = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.unlike(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestIsLiked_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.isLiked(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestIsLiked_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.IsLiked = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (types.Bool, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewBool(false), nil
	}

	message, nexErr := commonProtocol.isLiked(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestIsLiked_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.IsLiked = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) (types.Bool, *nex.Error) {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return types.NewBool(false), nil
	}

	message, nexErr := commonProtocol.isLiked(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodIsLiked, 9, []byte{0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.IsLiked = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (types.Bool, *nex.Error) {
		return types.NewBool(false), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.isLiked(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetFavourites_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getFavourites(nil, fakePacket{}, 7, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	checkMessage(t, message, nexErr, ubi.MethodGetFavourites, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetFavourites_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetFavourites = func(packet nex.PacketInterface, _ types.List[types.UInt32]) (types.List[user_storage_types.UserContent], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[user_storage_types.UserContent](), nil
	}

	message, nexErr := commonProtocol.getFavourites(errors.New("bad"), fakePacket{}, 7, types.NewList[types.UInt32]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetFavourites_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetFavourites = func(packet nex.PacketInterface, got0 types.List[types.UInt32]) (types.List[user_storage_types.UserContent], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})) {
			t.Errorf("contentTypes differs")
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
		}()}, nil
	}

	message, nexErr := commonProtocol.getFavourites(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	checkMessage(t, message, nexErr, ubi.MethodGetFavourites, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x43, 0x08, 0xaf, 0x2f, 0x00, 0x00, 0x00, 0x00, 0x7e, 0x23, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x51, 0x2f, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x31, 0x33, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x01, 0xeb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xe8, 0x4a, 0x00, 0x00, 0x97, 0x94, 0x35, 0x77, 0x00, 0x00, 0x00, 0x00, 0xca, 0x52, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x9d, 0x5e, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x32, 0x35, 0x00, 0x70, 0x6a, 0x00, 0x00, 0x01, 0xdf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetFavourites = func(packet nex.PacketInterface, _ types.List[types.UInt32]) (types.List[user_storage_types.UserContent], *nex.Error) {
		return types.NewList[user_storage_types.UserContent](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getFavourites(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestMakeFavourite_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.makeFavourite(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestMakeFavourite_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.MakeFavourite = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.makeFavourite(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestMakeFavourite_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.MakeFavourite = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) *nex.Error {
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

	message, nexErr := commonProtocol.makeFavourite(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodMakeFavourite, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.MakeFavourite = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.makeFavourite(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestRemoveFromFavourites_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.removeFromFavourites(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestRemoveFromFavourites_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.RemoveFromFavourites = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.removeFromFavourites(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestRemoveFromFavourites_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.RemoveFromFavourites = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) *nex.Error {
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

	message, nexErr := commonProtocol.removeFromFavourites(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodRemoveFromFavourites, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.RemoveFromFavourites = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.removeFromFavourites(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestReportInappropriate_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.reportInappropriate(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.NewString("s4"))
	checkNotImplemented(t, message, nexErr)
}

func TestReportInappropriate_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ReportInappropriate = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey, _ types.String) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.reportInappropriate(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey(), types.NewString(""))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestReportInappropriate_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ReportInappropriate = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey, got1 types.String) *nex.Error {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		if !(got1.Equals(types.NewString("s4"))) {
			t.Errorf("reason differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.reportInappropriate(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.NewString("s4"))
	checkMessage(t, message, nexErr, ubi.MethodReportInappropriate, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ReportInappropriate = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey, _ types.String) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.reportInappropriate(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.NewString("s4"))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestIncrementPlayCount_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.incrementPlayCount(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkNotImplemented(t, message, nexErr)
}

func TestIncrementPlayCount_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.IncrementPlayCount = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.incrementPlayCount(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestIncrementPlayCount_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.IncrementPlayCount = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) *nex.Error {
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

	message, nexErr := commonProtocol.incrementPlayCount(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodIncrementPlayCount, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.IncrementPlayCount = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.incrementPlayCount(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestUpdateCustomStat_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.updateCustomStat(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.NewUInt16(407), types.NewInt64(-20000000036))
	checkNotImplemented(t, message, nexErr)
}

func TestUpdateCustomStat_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.UpdateCustomStat = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey, _ types.UInt16, _ types.Int64) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.updateCustomStat(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey(), types.NewUInt16(0), types.NewInt64(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestUpdateCustomStat_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.UpdateCustomStat = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey, got1 types.UInt16, got2 types.Int64) *nex.Error {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		if !(got1.Equals(types.NewUInt16(407))) {
			t.Errorf("statID differs")
		}

		if !(got2.Equals(types.NewInt64(-20000000036))) {
			t.Errorf("incValue differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.updateCustomStat(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.NewUInt16(407), types.NewInt64(-20000000036))
	checkMessage(t, message, nexErr, ubi.MethodUpdateCustomStat, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.UpdateCustomStat = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey, _ types.UInt16, _ types.Int64) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.updateCustomStat(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.NewUInt16(407), types.NewInt64(-20000000036))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetOwnContents_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getOwnContents(nil, fakePacket{}, 7, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodGetOwnContents, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetOwnContents_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetOwnContents = func(packet nex.PacketInterface, _ types.UInt32) (types.List[user_storage_types.UserContent], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[user_storage_types.UserContent](), nil
	}

	message, nexErr := commonProtocol.getOwnContents(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetOwnContents_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetOwnContents = func(packet nex.PacketInterface, got0 types.UInt32) (types.List[user_storage_types.UserContent], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("typeID differs")
		}

		return types.List[user_storage_types.UserContent]{func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(5050)
				v.ContentID = types.NewUInt64(600000053)
				return v
			}()
			v.Pid = types.NewUInt32(7068)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(10095)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(13122)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-19)}
				return v
			}()}
			return v
		}(), func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(17158)
				v.ContentID = types.NewUInt64(1800000137)
				return v
			}()
			v.Pid = types.NewUInt32(19176)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(22203)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v23")}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(25230)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-31)}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getOwnContents(nil, fakePacket{}, 9, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodGetOwnContents, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xba, 0x13, 0x00, 0x00, 0x35, 0x46, 0xc3, 0x23, 0x00, 0x00, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x6f, 0x27, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x31, 0x31, 0x00, 0x42, 0x33, 0x00, 0x00, 0x01, 0xed, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x06, 0x43, 0x00, 0x00, 0x89, 0xd2, 0x49, 0x6b, 0x00, 0x00, 0x00, 0x00, 0xe8, 0x4a, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xbb, 0x56, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x32, 0x33, 0x00, 0x8e, 0x62, 0x00, 0x00, 0x01, 0xe1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetOwnContents = func(packet nex.PacketInterface, _ types.UInt32) (types.List[user_storage_types.UserContent], *nex.Error) {
		return types.NewList[user_storage_types.UserContent](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getOwnContents(nil, fakePacket{}, 9, types.NewUInt32(1014))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetMostPopularTags_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getMostPopularTags(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetMostPopularTags, 7, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestGetMostPopularTags_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetMostPopularTags = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (types.List[user_storage_types.WeightedTag], types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[user_storage_types.WeightedTag](), types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.getMostPopularTags(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetMostPopularTags_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetMostPopularTags = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) (types.List[user_storage_types.WeightedTag], types.UInt32, *nex.Error) {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return types.List[user_storage_types.WeightedTag]{func() user_storage_types.WeightedTag {
			v := user_storage_types.NewWeightedTag()
			v.Id = types.NewUInt32(6059)
			v.NumberOfOccurences = types.NewUInt32(7068)
			return v
		}(), func() user_storage_types.WeightedTag {
			v := user_storage_types.NewWeightedTag()
			v.Id = types.NewUInt32(9086)
			v.NumberOfOccurences = types.NewUInt32(10095)
			return v
		}()}, types.NewUInt32(11104), nil
	}

	message, nexErr := commonProtocol.getMostPopularTags(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetMostPopularTags, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xab, 0x17, 0x00, 0x00, 0x9c, 0x1b, 0x00, 0x00, 0x7e, 0x23, 0x00, 0x00, 0x6f, 0x27, 0x00, 0x00, 0x60, 0x2b, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetMostPopularTags = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (types.List[user_storage_types.WeightedTag], types.UInt32, *nex.Error) {
		return types.NewList[user_storage_types.WeightedTag](), types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getMostPopularTags(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetTags_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getTags(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetTags, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetTags_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetTags = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (types.List[types.UInt32], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[types.UInt32](), nil
	}

	message, nexErr := commonProtocol.getTags(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetTags_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetTags = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey) (types.List[types.UInt32], *nex.Error) {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		return types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)}, nil
	}

	message, nexErr := commonProtocol.getTags(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodGetTags, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xba, 0x13, 0x00, 0x00, 0xab, 0x17, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetTags = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey) (types.List[types.UInt32], *nex.Error) {
		return types.NewList[types.UInt32](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getTags(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestTagContent_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.tagContent(nil, fakePacket{}, 7, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})
	checkNotImplemented(t, message, nexErr)
}

func TestTagContent_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.TagContent = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey, _ types.List[types.UInt32]) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.tagContent(errors.New("bad"), fakePacket{}, 7, user_storage_types.NewUserContentKey(), types.NewList[types.UInt32]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestTagContent_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.TagContent = func(packet nex.PacketInterface, got0 user_storage_types.UserContentKey, got1 types.List[types.UInt32]) *nex.Error {
		called = true

		if !(got0.Equals(func() user_storage_types.UserContentKey {
			v := user_storage_types.NewUserContentKey()
			v.TypeID = types.NewUInt32(2023)
			v.ContentID = types.NewUInt64(300000032)
			return v
		}())) {
			t.Errorf("contentKey differs")
		}

		if !(got1.Equals(types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})) {
			t.Errorf("newTagIds differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.tagContent(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})
	checkMessage(t, message, nexErr, ubi.MethodTagContent, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.TagContent = func(packet nex.PacketInterface, _ user_storage_types.UserContentKey, _ types.List[types.UInt32]) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.tagContent(nil, fakePacket{}, 9, func() user_storage_types.UserContentKey {
		v := user_storage_types.NewUserContentKey()
		v.TypeID = types.NewUInt32(2023)
		v.ContentID = types.NewUInt64(300000032)
		return v
	}(), types.List[types.UInt32]{types.NewUInt32(5050), types.NewUInt32(6059)})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSearchContentsByPlayers_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.searchContentsByPlayers(nil, fakePacket{}, 7, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(5050)
		v.QueryID = types.NewUInt32(6059)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(21), Length: types.NewUInt32(31)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(13122)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-19)}
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchContentsByPlayers, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestSearchContentsByPlayers_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SearchContentsByPlayers = func(packet nex.PacketInterface, _ types.List[types.UInt32], _ user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[user_storage_types.UserContent](), nil
	}

	message, nexErr := commonProtocol.searchContentsByPlayers(errors.New("bad"), fakePacket{}, 7, types.NewList[types.UInt32](), user_storage_types.NewUserStorageQuery())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSearchContentsByPlayers_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SearchContentsByPlayers = func(packet nex.PacketInterface, got0 types.List[types.UInt32], got1 user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], *nex.Error) {
		called = true

		if !(got0.Equals(types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})) {
			t.Errorf("pids differs")
		}

		if !(got1.Equals(func() user_storage_types.UserStorageQuery {
			v := user_storage_types.NewUserStorageQuery()
			v.TypeID = types.NewUInt32(5050)
			v.QueryID = types.NewUInt32(6059)
			v.ResultRange = types.ResultRange{Offset: types.NewUInt32(21), Length: types.NewUInt32(31)}
			v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(10095)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(13122)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-19)}
				return v
			}()}
			return v
		}())) {
			t.Errorf("query differs")
		}

		return types.List[user_storage_types.UserContent]{func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(18167)
				v.ContentID = types.NewUInt64(1900000144)
				return v
			}()
			v.Pid = types.NewUInt32(20185)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(23212)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-29)}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(26239)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v27")}
				return v
			}()}
			return v
		}(), func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(30275)
				v.ContentID = types.NewUInt64(3100000228)
				return v
			}()
			v.Pid = types.NewUInt32(32293)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(35320)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-41)}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(38347)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v39")}
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.searchContentsByPlayers(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(5050)
		v.QueryID = types.NewUInt32(6059)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(21), Length: types.NewUInt32(31)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(13122)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-19)}
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchContentsByPlayers, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xf7, 0x46, 0x00, 0x00, 0x90, 0xb3, 0x3f, 0x71, 0x00, 0x00, 0x00, 0x00, 0xd9, 0x4e, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x5a, 0x00, 0x00, 0x01, 0xe3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f, 0x66, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x32, 0x37, 0x00, 0x43, 0x76, 0x00, 0x00, 0xe4, 0x3f, 0xc6, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x25, 0x7e, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xf8, 0x89, 0x00, 0x00, 0x01, 0xd7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xcb, 0x95, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x39, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SearchContentsByPlayers = func(packet nex.PacketInterface, _ types.List[types.UInt32], _ user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], *nex.Error) {
		return types.NewList[user_storage_types.UserContent](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.searchContentsByPlayers(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(5050)
		v.QueryID = types.NewUInt32(6059)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(21), Length: types.NewUInt32(31)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(13122)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-19)}
			return v
		}()}
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSearchContentsByPlayersWithTotal_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.searchContentsByPlayersWithTotal(nil, fakePacket{}, 7, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(5050)
		v.QueryID = types.NewUInt32(6059)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(21), Length: types.NewUInt32(31)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(13122)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-19)}
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchContentsByPlayersWithTotal, 7, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestSearchContentsByPlayersWithTotal_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SearchContentsByPlayersWithTotal = func(packet nex.PacketInterface, _ types.List[types.UInt32], _ user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[user_storage_types.UserContent](), types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.searchContentsByPlayersWithTotal(errors.New("bad"), fakePacket{}, 7, types.NewList[types.UInt32](), user_storage_types.NewUserStorageQuery())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSearchContentsByPlayersWithTotal_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SearchContentsByPlayersWithTotal = func(packet nex.PacketInterface, got0 types.List[types.UInt32], got1 user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		called = true

		if !(got0.Equals(types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)})) {
			t.Errorf("pids differs")
		}

		if !(got1.Equals(func() user_storage_types.UserStorageQuery {
			v := user_storage_types.NewUserStorageQuery()
			v.TypeID = types.NewUInt32(5050)
			v.QueryID = types.NewUInt32(6059)
			v.ResultRange = types.ResultRange{Offset: types.NewUInt32(21), Length: types.NewUInt32(31)}
			v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(10095)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(13122)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-19)}
				return v
			}()}
			return v
		}())) {
			t.Errorf("query differs")
		}

		return types.List[user_storage_types.UserContent]{func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(18167)
				v.ContentID = types.NewUInt64(1900000144)
				return v
			}()
			v.Pid = types.NewUInt32(20185)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(23212)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-29)}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(26239)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v27")}
				return v
			}()}
			return v
		}(), func() user_storage_types.UserContent {
			v := user_storage_types.NewUserContent()
			v.Key = func() user_storage_types.UserContentKey {
				v := user_storage_types.NewUserContentKey()
				v.TypeID = types.NewUInt32(30275)
				v.ContentID = types.NewUInt64(3100000228)
				return v
			}()
			v.Pid = types.NewUInt32(32293)
			v.Properties = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(35320)
				v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-41)}
				return v
			}(), func() user_storage_types.ContentProperty {
				v := user_storage_types.NewContentProperty()
				v.ID = types.NewUInt32(38347)
				v.Value = types.Variant{TypeID: 4, Type: types.NewString("v39")}
				return v
			}()}
			return v
		}()}, types.NewUInt32(40365), nil
	}

	message, nexErr := commonProtocol.searchContentsByPlayersWithTotal(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(5050)
		v.QueryID = types.NewUInt32(6059)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(21), Length: types.NewUInt32(31)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(13122)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-19)}
			return v
		}()}
		return v
	}())
	checkMessage(t, message, nexErr, ubi.MethodSearchContentsByPlayersWithTotal, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xf7, 0x46, 0x00, 0x00, 0x90, 0xb3, 0x3f, 0x71, 0x00, 0x00, 0x00, 0x00, 0xd9, 0x4e, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x5a, 0x00, 0x00, 0x01, 0xe3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f, 0x66, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x32, 0x37, 0x00, 0x43, 0x76, 0x00, 0x00, 0xe4, 0x3f, 0xc6, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x25, 0x7e, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xf8, 0x89, 0x00, 0x00, 0x01, 0xd7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xcb, 0x95, 0x00, 0x00, 0x04, 0x04, 0x00, 0x76, 0x33, 0x39, 0x00, 0xad, 0x9d, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SearchContentsByPlayersWithTotal = func(packet nex.PacketInterface, _ types.List[types.UInt32], _ user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error) {
		return types.NewList[user_storage_types.UserContent](), types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.searchContentsByPlayersWithTotal(nil, fakePacket{}, 9, types.List[types.UInt32]{types.NewUInt32(2023), types.NewUInt32(3032)}, func() user_storage_types.UserStorageQuery {
		v := user_storage_types.NewUserStorageQuery()
		v.TypeID = types.NewUInt32(5050)
		v.QueryID = types.NewUInt32(6059)
		v.ResultRange = types.ResultRange{Offset: types.NewUInt32(21), Length: types.NewUInt32(31)}
		v.Parameters = types.List[user_storage_types.ContentProperty]{func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(10095)
			v.Value = types.Variant{TypeID: 4, Type: types.NewString("v11")}
			return v
		}(), func() user_storage_types.ContentProperty {
			v := user_storage_types.NewContentProperty()
			v.ID = types.NewUInt32(13122)
			v.Value = types.Variant{TypeID: 1, Type: types.NewInt64(-19)}
			return v
		}()}
		return v
	}())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}
