package useraccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"testing"
)

func init() {
	types.RegisterObjectHolderType(types.NewData())
}

func TestBuildLookupSceNpIdsResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "01000000a41100000300a50203"}, {true, "01000000a41100000300a50203"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupSceNpIdsResponse(out, types.Map[types.UInt32, types.QBuffer]{types.NewUInt32(4516): types.NewQBuffer([]byte{165, 2, 3})})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupSceNpIdsHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupSceNpIds(nil, nil, 1, types.NewList[types.UInt32]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupSceNpIds(errTest, nil, 1, types.NewList[types.UInt32]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLookupPrincipalIDsResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0100000005007334323300a8110000"}, {true, "0100000005007334323300a8110000"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupPrincipalIDsResponse(out, types.Map[types.String, types.UInt32]{types.NewString("s423"): types.NewUInt32(4520)})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupPrincipalIDsHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupPrincipalIDs(nil, nil, 1, types.NewList[types.String](), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupPrincipalIDs(errTest, nil, 1, types.NewList[types.String](), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLookupFirstPartyIdsResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "01000000aa11000005007334323700"}, {true, "01000000aa11000005007334323700"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupFirstPartyIdsResponse(out, types.Map[types.UInt32, types.String]{types.NewUInt32(4522): types.NewString("s427")})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupFirstPartyIdsHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupFirstPartyIds(nil, nil, 1, types.NewList[types.UInt32](), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupFirstPartyIds(errTest, nil, 1, types.NewList[types.UInt32](), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildUserHasPlayedResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000000500733432390001"}, {true, "010000000500733432390001"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildUserHasPlayedResponse(out, types.Map[types.String, types.Bool]{types.NewString("s429"): types.NewBool(true)})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestUserHasPlayedHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.userHasPlayed(nil, nil, 1, types.NewList[types.String](), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.userHasPlayed(errTest, nil, 1, types.NewList[types.String](), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildIsUserPlayingResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "010000000500733433320000"}, {true, "010000000500733433320000"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildIsUserPlayingResponse(out, types.Map[types.String, types.Bool]{types.NewString("s432"): types.NewBool(false)})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestIsUserPlayingHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.isUserPlaying(nil, nil, 1, types.NewList[types.String](), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.isUserPlaying(errTest, nil, 1, types.NewList[types.String](), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildUpdateSonyAccountInfoResponse(t *testing.T) {
	out := newOut(false)
	buildUpdateSonyAccountInfoResponse(out)
	checkBytes(t, out.Bytes(), "")
}

func TestUpdateSonyAccountInfoHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.updateSonyAccountInfo(nil, nil, 1, types.NewQBuffer(nil), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.updateSonyAccountInfo(errTest, nil, 1, types.NewQBuffer(nil), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildLookupUsernamesResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "01000000b311000005007334333600"}, {true, "01000000b311000005007334333600"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildLookupUsernamesResponse(out, types.Map[types.UInt32, types.String]{types.NewUInt32(4531): types.NewString("s436")})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestLookupUsernamesHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.lookupUsernames(nil, nil, 1, types.NewList[types.UInt32](), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.lookupUsernames(errTest, nil, 1, types.NewList[types.UInt32](), types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}
