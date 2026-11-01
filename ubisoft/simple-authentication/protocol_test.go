package simpleauthentication

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"testing"
)

func init() {
	types.RegisterObjectHolderType(types.NewData())
}

func TestBuildAuthenticateResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0100010013100000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073323100"}, {true, "01000100131000000037000000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073323100"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildAuthenticateResponse(out, types.NewQResult(0x00010001), types.NewUInt32(4115), func() types.RVConnectionData {
			v := types.NewRVConnectionData()
			v.StationURL = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			v.SpecialProtocols = types.List[types.UInt8]{7}
			v.StationURLSpecialProtocols = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			return v
		}(), types.NewString("s21"))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestAuthenticateHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.authenticate(nil, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.authenticate(errTest, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLoginWithTokenResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0100010017100000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073323500"}, {true, "01000100171000000037000000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073323500"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLoginWithTokenResponse(out, types.NewQResult(0x00010001), types.NewUInt32(4119), func() types.RVConnectionData {
			v := types.NewRVConnectionData()
			v.StationURL = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			v.SpecialProtocols = types.List[types.UInt8]{7}
			v.StationURLSpecialProtocols = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			return v
		}(), types.NewString("s25"))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLoginWithTokenHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.loginWithToken(nil, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.loginWithToken(errTest, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLoginWithTokenExResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010001001b100000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073323900"}, {true, "010001001b1000000037000000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073323900"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLoginWithTokenExResponse(out, types.NewQResult(0x00010001), types.NewUInt32(4123), func() types.RVConnectionData {
			v := types.NewRVConnectionData()
			v.StationURL = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			v.SpecialProtocols = types.List[types.UInt8]{7}
			v.StationURLSpecialProtocols = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			return v
		}(), types.NewString("s29"))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLoginWithTokenExHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.loginWithTokenEx(nil, nil, 1, types.NewString(""), types.NewDataHolder())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.loginWithTokenEx(errTest, nil, 1, types.NewString(""), types.NewDataHolder())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLoginResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010001001f100000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073333300"}, {true, "010001001f1000000037000000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073333300"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLoginResponse(out, types.NewQResult(0x00010001), types.NewUInt32(4127), func() types.RVConnectionData {
			v := types.NewRVConnectionData()
			v.StationURL = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			v.SpecialProtocols = types.List[types.UInt8]{7}
			v.StationURLSpecialProtocols = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			return v
		}(), types.NewString("s33"))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLoginHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.login(nil, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.login(errTest, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLoginWithSubAccountResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0100010023100000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073333700"}, {true, "01000100231000000037000000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073333700"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLoginWithSubAccountResponse(out, types.NewQResult(0x00010001), types.NewUInt32(4131), func() types.RVConnectionData {
			v := types.NewRVConnectionData()
			v.StationURL = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			v.SpecialProtocols = types.List[types.UInt8]{7}
			v.StationURLSpecialProtocols = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			return v
		}(), types.NewString("s37"))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLoginWithSubAccountHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.loginWithSubAccount(nil, nil, 1, types.NewDataHolder())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.loginWithSubAccount(errTest, nil, 1, types.NewDataHolder())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildRegisterResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0100010027100000170070727564703a2f616464726573733d312e322e332e3400"}, {true, "0100010027100000170070727564703a2f616464726573733d312e322e332e3400"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildRegisterResponse(out, types.NewQResult(0x00010001), types.NewUInt32(4135), types.NewStationURL(types.NewString("prudp:/address=1.2.3.4")))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestRegisterHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.register(nil, nil, 1, types.NewList[types.StationURL]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.register(errTest, nil, 1, types.NewList[types.StationURL]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildRegisterExResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010001002a100000170070727564703a2f616464726573733d312e322e332e3400"}, {true, "010001002a100000170070727564703a2f616464726573733d312e322e332e3400"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildRegisterExResponse(out, types.NewQResult(0x00010001), types.NewUInt32(4138), types.NewStationURL(types.NewString("prudp:/address=1.2.3.4")))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestRegisterExHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.registerEx(nil, nil, 1, types.NewList[types.StationURL](), types.NewDataHolder())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.registerEx(errTest, nil, 1, types.NewList[types.StationURL](), types.NewDataHolder())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLoginWithTokenCafeResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010001002d100000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073343700"}, {true, "010001002d1000000037000000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073343700"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLoginWithTokenCafeResponse(out, types.NewQResult(0x00010001), types.NewUInt32(4141), func() types.RVConnectionData {
			v := types.NewRVConnectionData()
			v.StationURL = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			v.SpecialProtocols = types.List[types.UInt8]{7}
			v.StationURLSpecialProtocols = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			return v
		}(), types.NewString("s47"))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLoginWithTokenCafeHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.loginWithTokenCafe(nil, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.loginWithTokenCafe(errTest, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLoginWithTokenCafeExResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0100010031100000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073353100"}, {true, "01000100311000000037000000170070727564703a2f616464726573733d312e322e332e34000100000007170070727564703a2f616464726573733d312e322e332e3400040073353100"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLoginWithTokenCafeExResponse(out, types.NewQResult(0x00010001), types.NewUInt32(4145), func() types.RVConnectionData {
			v := types.NewRVConnectionData()
			v.StationURL = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			v.SpecialProtocols = types.List[types.UInt8]{7}
			v.StationURLSpecialProtocols = types.NewStationURL(types.NewString("prudp:/address=1.2.3.4"))
			return v
		}(), types.NewString("s51"))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLoginWithTokenCafeExHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.loginWithTokenCafeEx(nil, nil, 1, types.NewString(""), types.NewDataHolder())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.loginWithTokenCafeEx(errTest, nil, 1, types.NewString(""), types.NewDataHolder())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}
