package privileges

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	privileges_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges/types"
	"testing"
)

func init() {
	types.RegisterObjectHolderType(types.NewData())
}

func TestBuildGetPrivilegesResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "01000000d0110000d211000005007334363700"}, {true, "01000000d0110000000b000000d211000005007334363700"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetPrivilegesResponse(out, types.Map[types.UInt32, privileges_types.Privilege]{types.NewUInt32(4560): func() privileges_types.Privilege {
			v := privileges_types.NewPrivilege()
			v.ID = types.NewUInt32(4562)
			v.Description = types.NewString("s467")
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetPrivilegesHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getPrivileges(nil, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getPrivileges(errTest, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildActivateKeyResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0500733436390002000000d811000005007334373300db11000005007334373600"}, {true, "002b0000000500733436390002000000000b000000d811000005007334373300000b000000db11000005007334373600"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildActivateKeyResponse(out, func() privileges_types.PrivilegeGroup {
			v := privileges_types.NewPrivilegeGroup()
			v.Description = types.NewString("s469")
			v.Privileges = types.List[privileges_types.Privilege]{func() privileges_types.Privilege {
				v := privileges_types.NewPrivilege()
				v.ID = types.NewUInt32(4568)
				v.Description = types.NewString("s473")
				return v
			}(), func() privileges_types.Privilege {
				v := privileges_types.NewPrivilege()
				v.ID = types.NewUInt32(4571)
				v.Description = types.NewString("s476")
				return v
			}()}
			return v
		}())
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestActivateKeyHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.activateKey(nil, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.activateKey(errTest, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildActivateKeyWithExpectedPrivilegesResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0500733437380002000000e111000005007334383200e411000005007334383500"}, {true, "002b0000000500733437380002000000000b000000e111000005007334383200000b000000e411000005007334383500"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildActivateKeyWithExpectedPrivilegesResponse(out, func() privileges_types.PrivilegeGroup {
			v := privileges_types.NewPrivilegeGroup()
			v.Description = types.NewString("s478")
			v.Privileges = types.List[privileges_types.Privilege]{func() privileges_types.Privilege {
				v := privileges_types.NewPrivilege()
				v.ID = types.NewUInt32(4577)
				v.Description = types.NewString("s482")
				return v
			}(), func() privileges_types.Privilege {
				v := privileges_types.NewPrivilege()
				v.ID = types.NewUInt32(4580)
				v.Description = types.NewString("s485")
				return v
			}()}
			return v
		}())
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestActivateKeyWithExpectedPrivilegesHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.activateKeyWithExpectedPrivileges(nil, nil, 1, types.NewString(""), types.NewString(""), types.NewList[types.UInt32]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.activateKeyWithExpectedPrivileges(errTest, nil, 1, types.NewString(""), types.NewString(""), types.NewList[types.UInt32]())
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetPrivilegeRemainDurationResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "1afeffff"}, {true, "1afeffff"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetPrivilegeRemainDurationResponse(out, types.NewInt32(-486))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetPrivilegeRemainDurationHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getPrivilegeRemainDuration(nil, nil, 1, types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getPrivilegeRemainDuration(errTest, nil, 1, types.NewUInt32(0))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetExpiredPrivilegesResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "02000000e91100000500733439300015feffffed1100000500733439340011feffff"}, {true, "02000000000f000000e91100000500733439300015feffff000f000000ed1100000500733439340011feffff"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetExpiredPrivilegesResponse(out, types.List[privileges_types.PrivilegeEx]{func() privileges_types.PrivilegeEx {
			v := privileges_types.NewPrivilegeEx()
			v.ID = types.NewUInt32(4585)
			v.Description = types.NewString("s490")
			v.Duration = types.NewInt32(-491)
			return v
		}(), func() privileges_types.PrivilegeEx {
			v := privileges_types.NewPrivilegeEx()
			v.ID = types.NewUInt32(4589)
			v.Description = types.NewString("s494")
			v.Duration = types.NewInt32(-495)
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetExpiredPrivilegesHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getExpiredPrivileges(nil, nil, 1)
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getExpiredPrivileges(errTest, nil, 1)
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetPrivilegesExResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "02000000f2110000050073343939000cfefffff61100000500733530330008feffff"}, {true, "02000000000f000000f2110000050073343939000cfeffff000f000000f61100000500733530330008feffff"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetPrivilegesExResponse(out, types.List[privileges_types.PrivilegeEx]{func() privileges_types.PrivilegeEx {
			v := privileges_types.NewPrivilegeEx()
			v.ID = types.NewUInt32(4594)
			v.Description = types.NewString("s499")
			v.Duration = types.NewInt32(-500)
			return v
		}(), func() privileges_types.PrivilegeEx {
			v := privileges_types.NewPrivilegeEx()
			v.ID = types.NewUInt32(4598)
			v.Description = types.NewString("s503")
			v.Duration = types.NewInt32(-504)
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetPrivilegesExHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getPrivilegesEx(nil, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getPrivilegesEx(errTest, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}
