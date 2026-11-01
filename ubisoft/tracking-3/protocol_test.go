package tracking3

import (
	"bytes"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/tracking-3"
	tracking_3_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/tracking-3/types"
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

func TestSendTag_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.sendTag(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewString("s2"), types.NewString("s3"), types.NewUInt32(4041))
	checkMessage(t, message, nexErr, ubi.MethodSendTag, 7, nil)
}

func TestSendTag_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SendTag = func(packet nex.PacketInterface, _ types.UInt32, _ types.String, _ types.String, _ types.UInt32) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.sendTag(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewString(""), types.NewString(""), types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSendTag_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SendTag = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.String, got2 types.String, got3 types.UInt32) *nex.Error {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("trackingID differs")
		}

		if !(got1.Equals(types.NewString("s2"))) {
			t.Errorf("tag differs")
		}

		if !(got2.Equals(types.NewString("s3"))) {
			t.Errorf("attributes differs")
		}

		if !(got3.Equals(types.NewUInt32(4041))) {
			t.Errorf("deltaTime differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.sendTag(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewString("s2"), types.NewString("s3"), types.NewUInt32(4041))
	checkMessage(t, message, nexErr, ubi.MethodSendTag, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SendTag = func(packet nex.PacketInterface, _ types.UInt32, _ types.String, _ types.String, _ types.UInt32) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.sendTag(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewString("s2"), types.NewString("s3"), types.NewUInt32(4041))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSendTagAndUpdateUserInfo_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.sendTagAndUpdateUserInfo(nil, fakePacket{}, 7, types.NewUInt32(1014), types.NewString("s2"), types.NewString("s3"), types.NewUInt32(4041), types.NewString("s5"))
	checkMessage(t, message, nexErr, ubi.MethodSendTagAndUpdateUserInfo, 7, nil)
}

func TestSendTagAndUpdateUserInfo_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SendTagAndUpdateUserInfo = func(packet nex.PacketInterface, _ types.UInt32, _ types.String, _ types.String, _ types.UInt32, _ types.String) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.sendTagAndUpdateUserInfo(errors.New("bad"), fakePacket{}, 7, types.NewUInt32(0), types.NewString(""), types.NewString(""), types.NewUInt32(0), types.NewString(""))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSendTagAndUpdateUserInfo_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SendTagAndUpdateUserInfo = func(packet nex.PacketInterface, got0 types.UInt32, got1 types.String, got2 types.String, got3 types.UInt32, got4 types.String) *nex.Error {
		called = true

		if !(got0.Equals(types.NewUInt32(1014))) {
			t.Errorf("trackingID differs")
		}

		if !(got1.Equals(types.NewString("s2"))) {
			t.Errorf("tag differs")
		}

		if !(got2.Equals(types.NewString("s3"))) {
			t.Errorf("attributes differs")
		}

		if !(got3.Equals(types.NewUInt32(4041))) {
			t.Errorf("deltaTime differs")
		}

		if !(got4.Equals(types.NewString("s5"))) {
			t.Errorf("userID differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.sendTagAndUpdateUserInfo(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewString("s2"), types.NewString("s3"), types.NewUInt32(4041), types.NewString("s5"))
	checkMessage(t, message, nexErr, ubi.MethodSendTagAndUpdateUserInfo, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SendTagAndUpdateUserInfo = func(packet nex.PacketInterface, _ types.UInt32, _ types.String, _ types.String, _ types.UInt32, _ types.String) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.sendTagAndUpdateUserInfo(nil, fakePacket{}, 9, types.NewUInt32(1014), types.NewString("s2"), types.NewString("s3"), types.NewUInt32(4041), types.NewString("s5"))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSendUserInfo_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.sendUserInfo(nil, fakePacket{}, 7, func() tracking_3_types.TrackingInformation {
		v := tracking_3_types.NewTrackingInformation()
		v.Ipn = types.NewUInt32(2023)
		v.UserID = types.NewString("s3")
		v.MachineID = types.NewString("s4")
		v.VisitorID = types.NewString("s5")
		v.UtsVersion = types.NewString("s6")
		return v
	}(), types.NewUInt32(7068))
	checkNotImplemented(t, message, nexErr)
}

func TestSendUserInfo_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SendUserInfo = func(packet nex.PacketInterface, _ tracking_3_types.TrackingInformation, _ types.UInt32) (tracking_3_types.TrackingInformation, types.UInt32, *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return tracking_3_types.NewTrackingInformation(), types.NewUInt32(0), nil
	}

	message, nexErr := commonProtocol.sendUserInfo(errors.New("bad"), fakePacket{}, 7, tracking_3_types.NewTrackingInformation(), types.NewUInt32(0))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSendUserInfo_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SendUserInfo = func(packet nex.PacketInterface, got0 tracking_3_types.TrackingInformation, got1 types.UInt32) (tracking_3_types.TrackingInformation, types.UInt32, *nex.Error) {
		called = true

		if !(got0.Equals(func() tracking_3_types.TrackingInformation {
			v := tracking_3_types.NewTrackingInformation()
			v.Ipn = types.NewUInt32(2023)
			v.UserID = types.NewString("s3")
			v.MachineID = types.NewString("s4")
			v.VisitorID = types.NewString("s5")
			v.UtsVersion = types.NewString("s6")
			return v
		}())) {
			t.Errorf("userInfo differs")
		}

		if !(got1.Equals(types.NewUInt32(7068))) {
			t.Errorf("deltaTime differs")
		}

		return func() tracking_3_types.TrackingInformation {
			v := tracking_3_types.NewTrackingInformation()
			v.Ipn = types.NewUInt32(9086)
			v.UserID = types.NewString("s10")
			v.MachineID = types.NewString("s11")
			v.VisitorID = types.NewString("s12")
			v.UtsVersion = types.NewString("s13")
			return v
		}(), types.NewUInt32(14131), nil
	}

	message, nexErr := commonProtocol.sendUserInfo(nil, fakePacket{}, 9, func() tracking_3_types.TrackingInformation {
		v := tracking_3_types.NewTrackingInformation()
		v.Ipn = types.NewUInt32(2023)
		v.UserID = types.NewString("s3")
		v.MachineID = types.NewString("s4")
		v.VisitorID = types.NewString("s5")
		v.UtsVersion = types.NewString("s6")
		return v
	}(), types.NewUInt32(7068))
	checkMessage(t, message, nexErr, ubi.MethodSendUserInfo, 9, []byte{0x7e, 0x23, 0x00, 0x00, 0x04, 0x00, 0x73, 0x31, 0x30, 0x00, 0x04, 0x00, 0x73, 0x31, 0x31, 0x00, 0x04, 0x00, 0x73, 0x31, 0x32, 0x00, 0x04, 0x00, 0x73, 0x31, 0x33, 0x00, 0x33, 0x37, 0x00, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SendUserInfo = func(packet nex.PacketInterface, _ tracking_3_types.TrackingInformation, _ types.UInt32) (tracking_3_types.TrackingInformation, types.UInt32, *nex.Error) {
		return tracking_3_types.NewTrackingInformation(), types.NewUInt32(0), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.sendUserInfo(nil, fakePacket{}, 9, func() tracking_3_types.TrackingInformation {
		v := tracking_3_types.NewTrackingInformation()
		v.Ipn = types.NewUInt32(2023)
		v.UserID = types.NewString("s3")
		v.MachineID = types.NewString("s4")
		v.VisitorID = types.NewString("s5")
		v.UtsVersion = types.NewString("s6")
		return v
	}(), types.NewUInt32(7068))
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestGetConfiguration_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.getConfiguration(nil, fakePacket{}, 7)
	checkMessage(t, message, nexErr, ubi.MethodGetConfiguration, 7, []byte{0x00, 0x00, 0x00, 0x00})
}

func TestGetConfiguration_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.GetConfiguration = func(packet nex.PacketInterface) (types.List[types.String], *nex.Error) {
		t.Fatal("hook must not run on a decode error")

		return types.NewList[types.String](), nil
	}

	message, nexErr := commonProtocol.getConfiguration(errors.New("bad"), fakePacket{}, 7)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestGetConfiguration_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.GetConfiguration = func(packet nex.PacketInterface) (types.List[types.String], *nex.Error) {
		called = true

		return types.List[types.String]{types.NewString("s2"), types.NewString("s3")}, nil
	}

	message, nexErr := commonProtocol.getConfiguration(nil, fakePacket{}, 9)
	checkMessage(t, message, nexErr, ubi.MethodGetConfiguration, 9, []byte{0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x73, 0x32, 0x00, 0x03, 0x00, 0x73, 0x33, 0x00})

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.GetConfiguration = func(packet nex.PacketInterface) (types.List[types.String], *nex.Error) {
		return types.NewList[types.String](), nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.getConfiguration(nil, fakePacket{}, 9)
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}

func TestSendTags_NilHook(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	message, nexErr := commonProtocol.sendTags(nil, fakePacket{}, 7, types.List[tracking_3_types.TrackingTag]{func() tracking_3_types.TrackingTag {
		v := tracking_3_types.NewTrackingTag()
		v.TrackingID = types.NewUInt32(3032)
		v.Tag = types.NewString("s4")
		v.Attributes = types.NewString("s5")
		v.DeltaTime = types.NewUInt32(6059)
		v.NewUserId = types.NewString("s7")
		return v
	}(), func() tracking_3_types.TrackingTag {
		v := tracking_3_types.NewTrackingTag()
		v.TrackingID = types.NewUInt32(9086)
		v.Tag = types.NewString("s10")
		v.Attributes = types.NewString("s11")
		v.DeltaTime = types.NewUInt32(12113)
		v.NewUserId = types.NewString("s13")
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodSendTags, 7, nil)
}

func TestSendTags_InvalidRequest(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	commonProtocol.SendTags = func(packet nex.PacketInterface, _ types.List[tracking_3_types.TrackingTag]) *nex.Error {
		t.Fatal("hook must not run on a decode error")

		return nil
	}

	message, nexErr := commonProtocol.sendTags(errors.New("bad"), fakePacket{}, 7, types.NewList[tracking_3_types.TrackingTag]())
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "").ResultCode {
		t.Fatalf("want Core::InvalidArgument, got %v %v", message, nexErr)
	}
}

func TestSendTags_Hook(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	called := false

	commonProtocol.SendTags = func(packet nex.PacketInterface, got0 types.List[tracking_3_types.TrackingTag]) *nex.Error {
		called = true

		if !(got0.Equals(types.List[tracking_3_types.TrackingTag]{func() tracking_3_types.TrackingTag {
			v := tracking_3_types.NewTrackingTag()
			v.TrackingID = types.NewUInt32(3032)
			v.Tag = types.NewString("s4")
			v.Attributes = types.NewString("s5")
			v.DeltaTime = types.NewUInt32(6059)
			v.NewUserId = types.NewString("s7")
			return v
		}(), func() tracking_3_types.TrackingTag {
			v := tracking_3_types.NewTrackingTag()
			v.TrackingID = types.NewUInt32(9086)
			v.Tag = types.NewString("s10")
			v.Attributes = types.NewString("s11")
			v.DeltaTime = types.NewUInt32(12113)
			v.NewUserId = types.NewString("s13")
			return v
		}()})) {
			t.Errorf("tagData differs")
		}

		return nil
	}

	message, nexErr := commonProtocol.sendTags(nil, fakePacket{}, 9, types.List[tracking_3_types.TrackingTag]{func() tracking_3_types.TrackingTag {
		v := tracking_3_types.NewTrackingTag()
		v.TrackingID = types.NewUInt32(3032)
		v.Tag = types.NewString("s4")
		v.Attributes = types.NewString("s5")
		v.DeltaTime = types.NewUInt32(6059)
		v.NewUserId = types.NewString("s7")
		return v
	}(), func() tracking_3_types.TrackingTag {
		v := tracking_3_types.NewTrackingTag()
		v.TrackingID = types.NewUInt32(9086)
		v.Tag = types.NewString("s10")
		v.Attributes = types.NewString("s11")
		v.DeltaTime = types.NewUInt32(12113)
		v.NewUserId = types.NewString("s13")
		return v
	}()})
	checkMessage(t, message, nexErr, ubi.MethodSendTags, 9, nil)

	if !called {
		t.Fatal("hook not called")
	}

	commonProtocol.SendTags = func(packet nex.PacketInterface, _ types.List[tracking_3_types.TrackingTag]) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	message, nexErr = commonProtocol.sendTags(nil, fakePacket{}, 9, types.List[tracking_3_types.TrackingTag]{func() tracking_3_types.TrackingTag {
		v := tracking_3_types.NewTrackingTag()
		v.TrackingID = types.NewUInt32(3032)
		v.Tag = types.NewString("s4")
		v.Attributes = types.NewString("s5")
		v.DeltaTime = types.NewUInt32(6059)
		v.NewUserId = types.NewString("s7")
		return v
	}(), func() tracking_3_types.TrackingTag {
		v := tracking_3_types.NewTrackingTag()
		v.TrackingID = types.NewUInt32(9086)
		v.Tag = types.NewString("s10")
		v.Attributes = types.NewString("s11")
		v.DeltaTime = types.NewUInt32(12113)
		v.NewUserId = types.NewString("s13")
		return v
	}()})
	if message != nil || nexErr == nil || nexErr.ResultCode != nex.NewError(nex.ResultCodes.Core.Unknown, "").ResultCode {
		t.Fatalf("hook error not propagated: %v %v", message, nexErr)
	}
}
