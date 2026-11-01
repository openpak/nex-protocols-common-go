package titlestorageadmin

import (
	"bytes"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/title-storage-admin"
	title_storage_admin_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/title-storage-admin/types"
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

func TestGetContentTypesDescription_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getContentTypesDescription(nil, fakePacket{}, 7)
	checkNotImplemented(t, message, nexErr)
}

func TestGetContentTypesDescription_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetContentTypesDescription = func(packet nex.PacketInterface) (types.List[title_storage_admin_types.ContentTypeDescription], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[title_storage_admin_types.ContentTypeDescription](), nil
	}

	message, nexErr := commonProtocol.getContentTypesDescription(errors.New("bad"), fakePacket{}, 7)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetContentTypesDescription_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetContentTypesDescription = func(packet nex.PacketInterface) (types.List[title_storage_admin_types.ContentTypeDescription], *nex.Error) {
		called = true

		return types.List[title_storage_admin_types.ContentTypeDescription]{func() title_storage_admin_types.ContentTypeDescription {
			v := title_storage_admin_types.NewContentTypeDescription()
			v.ContentType = types.NewString("s3")
			v.ExtensionFields = types.List[title_storage_admin_types.ContentTypeExtensionField]{func() title_storage_admin_types.ContentTypeExtensionField {
				v := title_storage_admin_types.NewContentTypeExtensionField()
				v.FieldName = types.NewString("s6")
				v.FieldType = types.NewString("s7")
				return v
			}(), func() title_storage_admin_types.ContentTypeExtensionField {
				v := title_storage_admin_types.NewContentTypeExtensionField()
				v.FieldName = types.NewString("s9")
				v.FieldType = types.NewString("s10")
				return v
			}()}
			return v
		}(), func() title_storage_admin_types.ContentTypeDescription {
			v := title_storage_admin_types.NewContentTypeDescription()
			v.ContentType = types.NewString("s12")
			v.ExtensionFields = types.List[title_storage_admin_types.ContentTypeExtensionField]{func() title_storage_admin_types.ContentTypeExtensionField {
				v := title_storage_admin_types.NewContentTypeExtensionField()
				v.FieldName = types.NewString("s15")
				v.FieldType = types.NewString("s16")
				return v
			}(), func() title_storage_admin_types.ContentTypeExtensionField {
				v := title_storage_admin_types.NewContentTypeExtensionField()
				v.FieldName = types.NewString("s18")
				v.FieldType = types.NewString("s19")
				return v
			}()}
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getContentTypesDescription(nil, fakePacket{}, 9)
	checkMessage(t, message, nexErr, ubi.MethodGetContentTypesDescription, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x73, 0x33, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x73, 0x36, 0x00, 0x03, 0x00, 0x73, 0x37, 0x00, 0x03, 0x00, 0x73, 0x39, 0x00, 0x04, 0x00, 0x73, 0x31, 0x30, 0x00, 0x04, 0x00, 0x73, 0x31, 0x32, 0x00, 0x02, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x35, 0x00, 0x04, 0x00, 0x73, 0x31, 0x36, 0x00, 0x04, 0x00, 0x73, 0x31, 0x38, 0x00, 0x04, 0x00, 0x73, 0x31, 0x39, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetContentTypesDescription = func(packet nex.PacketInterface) (types.List[title_storage_admin_types.ContentTypeDescription], *nex.Error) {
		return types.NewList[title_storage_admin_types.ContentTypeDescription](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getContentTypesDescription(nil, fakePacket{}, 9)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetContentTypeDescription_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getContentTypeDescription(nil, fakePacket{}, 7, types.NewString("s1"))
	checkNotImplemented(t, message, nexErr)
}

func TestGetContentTypeDescription_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetContentTypeDescription = func(packet nex.PacketInterface, _ types.String) (types.List[title_storage_admin_types.ContentTypeExtensionField], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[title_storage_admin_types.ContentTypeExtensionField](), nil
	}

	message, nexErr := commonProtocol.getContentTypeDescription(errors.New("bad"), fakePacket{}, 7, types.NewString(""))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetContentTypeDescription_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetContentTypeDescription = func(packet nex.PacketInterface, got0 types.String) (types.List[title_storage_admin_types.ContentTypeExtensionField], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewString("s1"))) {
			t.Errorf("contentType differs")
		}

		return types.List[title_storage_admin_types.ContentTypeExtensionField]{func() title_storage_admin_types.ContentTypeExtensionField {
			v := title_storage_admin_types.NewContentTypeExtensionField()
			v.FieldName = types.NewString("s4")
			v.FieldType = types.NewString("s5")
			return v
		}(), func() title_storage_admin_types.ContentTypeExtensionField {
			v := title_storage_admin_types.NewContentTypeExtensionField()
			v.FieldName = types.NewString("s7")
			v.FieldType = types.NewString("s8")
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.getContentTypeDescription(nil, fakePacket{}, 9, types.NewString("s1"))
	checkMessage(t, message, nexErr, ubi.MethodGetContentTypeDescription, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x73, 0x34, 0x00, 0x03, 0x00, 0x73, 0x35, 0x00, 0x03, 0x00, 0x73, 0x37, 0x00, 0x03, 0x00, 0x73, 0x38, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetContentTypeDescription = func(packet nex.PacketInterface, _ types.String) (types.List[title_storage_admin_types.ContentTypeExtensionField], *nex.Error) {
		return types.NewList[title_storage_admin_types.ContentTypeExtensionField](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getContentTypeDescription(nil, fakePacket{}, 9, types.NewString("s1"))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestBrowseContents_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.browseContents(nil, fakePacket{}, 7)
	checkNotImplemented(t, message, nexErr)
}

func TestBrowseContents_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.BrowseContents = func(packet nex.PacketInterface) (types.List[title_storage_admin_types.AdminTitleContent], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[title_storage_admin_types.AdminTitleContent](), nil
	}

	message, nexErr := commonProtocol.browseContents(errors.New("bad"), fakePacket{}, 7)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestBrowseContents_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.BrowseContents = func(packet nex.PacketInterface) (types.List[title_storage_admin_types.AdminTitleContent], *nex.Error) {
		called = true

		return types.List[title_storage_admin_types.AdminTitleContent]{func() title_storage_admin_types.AdminTitleContent {
			v := title_storage_admin_types.NewAdminTitleContent()
			v.Id = types.NewUInt32(3032)
			v.Name = types.NewString("s4")
			v.ContentDescription = func() title_storage_admin_types.ContentDescription {
				v := title_storage_admin_types.NewContentDescription()
				v.ContentDisplayName = types.NewString("s6")
				v.ContentDescription = types.NewString("s7")
				return v
			}()
			v.Flags = types.NewUInt32(8077)
			v.ExtendedFlags = types.NewUInt32(9086)
			v.DateAvailable = types.NewDateTime(10000107)
			v.DateExpired = types.NewDateTime(11000110)
			v.ContentType = types.NewString("s12")
			v.DownloadURL = types.NewString("s13")
			return v
		}(), func() title_storage_admin_types.AdminTitleContent {
			v := title_storage_admin_types.NewAdminTitleContent()
			v.Id = types.NewUInt32(15140)
			v.Name = types.NewString("s16")
			v.ContentDescription = func() title_storage_admin_types.ContentDescription {
				v := title_storage_admin_types.NewContentDescription()
				v.ContentDisplayName = types.NewString("s18")
				v.ContentDescription = types.NewString("s19")
				return v
			}()
			v.Flags = types.NewUInt32(20185)
			v.ExtendedFlags = types.NewUInt32(21194)
			v.DateAvailable = types.NewDateTime(22000143)
			v.DateExpired = types.NewDateTime(23000146)
			v.ContentType = types.NewString("s24")
			v.DownloadURL = types.NewString("s25")
			return v
		}()}, nil
	}

	message, nexErr := commonProtocol.browseContents(nil, fakePacket{}, 9)
	checkMessage(t, message, nexErr, ubi.MethodBrowseContents, 9, []byte{0x02, 0x00, 0x00, 0x00, 0xd8, 0x0b, 0x00, 0x00, 0x03, 0x00, 0x73, 0x34, 0x00, 0x03, 0x00, 0x73, 0x36, 0x00, 0x03, 0x00, 0x73, 0x37, 0x00, 0x8d, 0x1f, 0x00, 0x00, 0x7e, 0x23, 0x00, 0x00, 0xeb, 0x96, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0xd9, 0xa7, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x32, 0x00, 0x04, 0x00, 0x73, 0x31, 0x33, 0x00, 0x24, 0x3b, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x36, 0x00, 0x04, 0x00, 0x73, 0x31, 0x38, 0x00, 0x04, 0x00, 0x73, 0x31, 0x39, 0x00, 0xd9, 0x4e, 0x00, 0x00, 0xca, 0x52, 0x00, 0x00, 0x0f, 0xb2, 0x4f, 0x01, 0x00, 0x00, 0x00, 0x00, 0x52, 0xf4, 0x5e, 0x01, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x32, 0x34, 0x00, 0x04, 0x00, 0x73, 0x32, 0x35, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.BrowseContents = func(packet nex.PacketInterface) (types.List[title_storage_admin_types.AdminTitleContent], *nex.Error) {
		return types.NewList[title_storage_admin_types.AdminTitleContent](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.browseContents(nil, fakePacket{}, 9)
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
			t.Errorf("contentId differs")
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

func TestCreateContentAndGetUploadInformation_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.createContentAndGetUploadInformation(nil, fakePacket{}, 7, types.NewString("s1"), types.NewString("s2"), types.NewString("s3"), types.NewDateTime(4000089), types.NewDateTime(5000092), types.NewUInt32(6059), types.NewUInt32(7068), types.NewString("s8"), types.DataHolder{Object: types.NewData()}, types.NewString("s10"))
	checkNotImplemented(t, message, nexErr)
}

func TestCreateContentAndGetUploadInformation_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.CreateContentAndGetUploadInformation = func(packet nex.PacketInterface, _ types.String, _ types.String, _ types.String, _ types.DateTime, _ types.DateTime, _ types.UInt32, _ types.UInt32, _ types.String, _ types.DataHolder, _ types.String) (types.UInt32, types.String, types.String, types.Map[types.String, types.String], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewUInt32(0), types.NewString(""), types.NewString(""), types.NewMap[types.String, types.String](), nil
	}

	message, nexErr := commonProtocol.createContentAndGetUploadInformation(errors.New("bad"), fakePacket{}, 7, types.NewString(""), types.NewString(""), types.NewString(""), types.NewDateTime(0), types.NewDateTime(0), types.NewUInt32(0), types.NewUInt32(0), types.NewString(""), types.NewDataHolder(), types.NewString(""))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestCreateContentAndGetUploadInformation_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.CreateContentAndGetUploadInformation = func(packet nex.PacketInterface, got0 types.String, got1 types.String, got2 types.String, got3 types.DateTime, got4 types.DateTime, got5 types.UInt32, got6 types.UInt32, got7 types.String, got8 types.DataHolder, got9 types.String) (types.UInt32, types.String, types.String, types.Map[types.String, types.String], *nex.Error) {
		called = true

		if !(got0.Equals(types.NewString("s1"))) {
			t.Errorf("name differs")
		}

		if !(got1.Equals(types.NewString("s2"))) {
			t.Errorf("description differs")
		}

		if !(got2.Equals(types.NewString("s3"))) {
			t.Errorf("key differs")
		}

		if !(got3.Equals(types.NewDateTime(4000089))) {
			t.Errorf("dateAvailable differs")
		}

		if !(got4.Equals(types.NewDateTime(5000092))) {
			t.Errorf("dateExpired differs")
		}

		if !(got5.Equals(types.NewUInt32(6059))) {
			t.Errorf("flags differs")
		}

		if !(got6.Equals(types.NewUInt32(7068))) {
			t.Errorf("extendedFlags differs")
		}

		if !(got7.Equals(types.NewString("s8"))) {
			t.Errorf("contentType differs")
		}

		if !(got8.Equals(types.DataHolder{Object: types.NewData()})) {
			t.Errorf("extension differs")
		}

		if !(got9.Equals(types.NewString("s10"))) {
			t.Errorf("mime differs")
		}

		return types.NewUInt32(11104), types.NewString("s12"), types.NewString("s13"), types.Map[types.String, types.String]{types.NewString("s15"): types.NewString("s16")}, nil
	}

	message, nexErr := commonProtocol.createContentAndGetUploadInformation(nil, fakePacket{}, 9, types.NewString("s1"), types.NewString("s2"), types.NewString("s3"), types.NewDateTime(4000089), types.NewDateTime(5000092), types.NewUInt32(6059), types.NewUInt32(7068), types.NewString("s8"), types.DataHolder{Object: types.NewData()}, types.NewString("s10"))
	checkMessage(t, message, nexErr, ubi.MethodCreateContentAndGetUploadInformation, 9, []byte{0x60, 0x2b, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x32, 0x00, 0x04, 0x00, 0x73, 0x31, 0x33, 0x00, 0x01, 0x00, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x35, 0x00, 0x04, 0x00, 0x73, 0x31, 0x36, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.CreateContentAndGetUploadInformation = func(packet nex.PacketInterface, _ types.String, _ types.String, _ types.String, _ types.DateTime, _ types.DateTime, _ types.UInt32, _ types.UInt32, _ types.String, _ types.DataHolder, _ types.String) (types.UInt32, types.String, types.String, types.Map[types.String, types.String], *nex.Error) {
		return types.NewUInt32(0), types.NewString(""), types.NewString(""), types.NewMap[types.String, types.String](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.createContentAndGetUploadInformation(nil, fakePacket{}, 9, types.NewString("s1"), types.NewString("s2"), types.NewString("s3"), types.NewDateTime(4000089), types.NewDateTime(5000092), types.NewUInt32(6059), types.NewUInt32(7068), types.NewString("s8"), types.DataHolder{Object: types.NewData()}, types.NewString("s10"))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestRemoveContent_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.removeContent(nil, fakePacket{}, 7, types.NewUInt32(1014))
	checkNotImplemented(t, message, nexErr)
}

func TestRemoveContent_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.RemoveContent = func(packet nex.PacketInterface, _ types.UInt32) (types.Bool, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewBool(false), nil
	}

	message, nexErr := commonProtocol.removeContent(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestRemoveContent_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.RemoveContent = func(packet nex.PacketInterface, got0 types.UInt32) (types.Bool, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("contentId differs")
		}

		return types.NewBool(false), nil
	}

	message, nexErr := commonProtocol.removeContent(nil, fakePacket{}, 9, types.NewUInt32(1014))
	checkMessage(t, message, nexErr, ubi.MethodRemoveContent, 9, []byte{0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.RemoveContent = func(packet nex.PacketInterface, _ types.UInt32) (types.Bool, *nex.Error) {
		return types.NewBool(false), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.removeContent(nil, fakePacket{}, 9, types.NewUInt32(1014))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestUpdateContentMetaData_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.updateContentMetaData(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewString("s2"), types.NewString("s3"), types.NewDateTime(4000089), types.NewDateTime(5000092), types.NewUInt32(6059), types.NewUInt32(7068), types.DataHolder{Object: types.NewData()})
	checkNotImplemented(t, message, nexErr)
}

func TestUpdateContentMetaData_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.UpdateContentMetaData = func(packet nex.PacketInterface, _ types.UInt32, _ types.String, _ types.String, _ types.DateTime, _ types.DateTime, _ types.UInt32, _ types.UInt32, _ types.DataHolder) (types.Bool, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewBool(false), nil
	}

	message, nexErr := commonProtocol.updateContentMetaData(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewString(""), types.NewString(""), types.NewDateTime(0), types.NewDateTime(0), types.NewUInt32(0), types.NewUInt32(0), types.NewDataHolder())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestUpdateContentMetaData_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.UpdateContentMetaData = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.String, got2 types.String, got3 types.DateTime, got4 types.DateTime, got5 types.UInt32, got6 types.UInt32, got7 types.DataHolder) (types.Bool, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("contentId differs")
		}

		if !(got1.Equals(types.NewString("s2"))) {
			t.Errorf("name differs")
		}

		if !(got2.Equals(types.NewString("s3"))) {
			t.Errorf("description differs")
		}

		if !(got3.Equals(types.NewDateTime(4000089))) {
			t.Errorf("dateAvailable differs")
		}

		if !(got4.Equals(types.NewDateTime(5000092))) {
			t.Errorf("dateExpired differs")
		}

		if !(got5.Equals(types.NewUInt32(6059))) {
			t.Errorf("flags differs")
		}

		if !(got6.Equals(types.NewUInt32(7068))) {
			t.Errorf("extendedFlags differs")
		}

		if !(got7.Equals(types.DataHolder{Object: types.NewData()})) {
			t.Errorf("extension differs")
		}

		return types.NewBool(true), nil
	}

	message, nexErr := commonProtocol.updateContentMetaData(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewString("s2"), types.NewString("s3"), types.NewDateTime(4000089), types.NewDateTime(5000092), types.NewUInt32(6059), types.NewUInt32(7068), types.DataHolder{Object: types.NewData()})
	checkMessage(t, message, nexErr, ubi.MethodUpdateContentMetaData, 9, []byte{0x01})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.UpdateContentMetaData = func(packet nex.PacketInterface, _ types.UInt32, _ types.String, _ types.String, _ types.DateTime, _ types.DateTime, _ types.UInt32, _ types.UInt32, _ types.DataHolder) (types.Bool, *nex.Error) {
		return types.NewBool(false), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.updateContentMetaData(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewString("s2"), types.NewString("s3"), types.NewDateTime(4000089), types.NewDateTime(5000092), types.NewUInt32(6059), types.NewUInt32(7068), types.DataHolder{Object: types.NewData()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestContentInsertionCompleted_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.contentInsertionCompleted(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewUInt32(2023))
	checkNotImplemented(t, message, nexErr)
}

func TestContentInsertionCompleted_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.ContentInsertionCompleted = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32) (types.Bool, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewBool(false), nil
	}

	message, nexErr := commonProtocol.contentInsertionCompleted(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestContentInsertionCompleted_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.ContentInsertionCompleted = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.UInt32) (types.Bool, *nex.Error) {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("contentId differs")
		}

		if !(got1.Equals(types.NewUInt32(2023))) {
			t.Errorf("contentSize differs")
		}

		return types.NewBool(true), nil
	}

	message, nexErr := commonProtocol.contentInsertionCompleted(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023))
	checkMessage(t, message, nexErr, ubi.MethodContentInsertionCompleted, 9, []byte{0x01})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.ContentInsertionCompleted = func(packet nex.PacketInterface, _ types.UInt32, _ types.UInt32) (types.Bool, *nex.Error) {
		return types.NewBool(false), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.contentInsertionCompleted(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewUInt32(2023))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}
