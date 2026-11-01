package uplaywin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	uplay_win_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win/types"
	"testing"
)

func init() {
	types.RegisterObjectHolderType(types.NewData())
}

func TestBuildGetActionsResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0200000005007336343800050073363439000500733635300075fdffff050073363532000200000005007336353500010500733635370005007336353900010500733636310005007336363300050073363634000500733636350066fdffff0500733636370002000000050073363730000005007336373200050073363734000005007336373600"}, {true, "02000000004c00000005007336343800050073363439000500733635300075fdffff0500733635320002000000000f000000050073363535000105007336353700000f000000050073363539000105007336363100004c00000005007336363300050073363634000500733636350066fdffff0500733636370002000000000f000000050073363730000005007336373200000f000000050073363734000005007336373600"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetActionsResponse(out, types.List[uplay_win_types.UplayAction]{func() uplay_win_types.UplayAction {
			v := uplay_win_types.NewUplayAction()
			v.Code = types.NewString("s648")
			v.Name = types.NewString("s649")
			v.Description = types.NewString("s650")
			v.Value = types.NewInt32(-651)
			v.GameCode = types.NewString("s652")
			v.Platforms = types.List[uplay_win_types.UplayActionPlatform]{func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s655")
				v.Completed = types.NewBool(true)
				v.SpecificKey = types.NewString("s657")
				return v
			}(), func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s659")
				v.Completed = types.NewBool(true)
				v.SpecificKey = types.NewString("s661")
				return v
			}()}
			return v
		}(), func() uplay_win_types.UplayAction {
			v := uplay_win_types.NewUplayAction()
			v.Code = types.NewString("s663")
			v.Name = types.NewString("s664")
			v.Description = types.NewString("s665")
			v.Value = types.NewInt32(-666)
			v.GameCode = types.NewString("s667")
			v.Platforms = types.List[uplay_win_types.UplayActionPlatform]{func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s670")
				v.Completed = types.NewBool(false)
				v.SpecificKey = types.NewString("s672")
				return v
			}(), func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s674")
				v.Completed = types.NewBool(false)
				v.SpecificKey = types.NewString("s676")
				return v
			}()}
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetActionsHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getActions(nil, nil, 1, types.NewInt32(0), types.NewInt32(0), types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getActions(errTest, nil, 1, types.NewInt32(0), types.NewInt32(0), types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetActionsCompletedResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0200000005007336373900050073363830000500733638310056fdffff050073363833000200000005007336383600000500733638380005007336393000000500733639320005007336393400050073363935000500733639360047fdffff0500733639380002000000050073373031000105007337303300050073373035000105007337303700"}, {true, "02000000004c00000005007336373900050073363830000500733638310056fdffff0500733638330002000000000f000000050073363836000005007336383800000f000000050073363930000005007336393200004c00000005007336393400050073363935000500733639360047fdffff0500733639380002000000000f000000050073373031000105007337303300000f000000050073373035000105007337303700"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetActionsCompletedResponse(out, types.List[uplay_win_types.UplayAction]{func() uplay_win_types.UplayAction {
			v := uplay_win_types.NewUplayAction()
			v.Code = types.NewString("s679")
			v.Name = types.NewString("s680")
			v.Description = types.NewString("s681")
			v.Value = types.NewInt32(-682)
			v.GameCode = types.NewString("s683")
			v.Platforms = types.List[uplay_win_types.UplayActionPlatform]{func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s686")
				v.Completed = types.NewBool(false)
				v.SpecificKey = types.NewString("s688")
				return v
			}(), func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s690")
				v.Completed = types.NewBool(false)
				v.SpecificKey = types.NewString("s692")
				return v
			}()}
			return v
		}(), func() uplay_win_types.UplayAction {
			v := uplay_win_types.NewUplayAction()
			v.Code = types.NewString("s694")
			v.Name = types.NewString("s695")
			v.Description = types.NewString("s696")
			v.Value = types.NewInt32(-697)
			v.GameCode = types.NewString("s698")
			v.Platforms = types.List[uplay_win_types.UplayActionPlatform]{func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s701")
				v.Completed = types.NewBool(true)
				v.SpecificKey = types.NewString("s703")
				return v
			}(), func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s705")
				v.Completed = types.NewBool(true)
				v.SpecificKey = types.NewString("s707")
				return v
			}()}
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetActionsCompletedHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getActionsCompleted(nil, nil, 1, types.NewInt32(0), types.NewInt32(0), types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getActionsCompleted(errTest, nil, 1, types.NewInt32(0), types.NewInt32(0), types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetActionsCountResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "3cfdffff"}, {true, "3cfdffff"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetActionsCountResponse(out, types.NewInt32(-708))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetActionsCountHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getActionsCount(nil, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getActionsCount(errTest, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetActionsCompletedCountResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "3bfdffff"}, {true, "3bfdffff"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetActionsCompletedCountResponse(out, types.NewInt32(-709))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetActionsCompletedCountHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getActionsCompletedCount(nil, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getActionsCompletedCount(errTest, nil, 1, types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetRewardsResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0200000005007337313200050073373133000500733731340035fdffff0500733731360005007337313700020000000500733732300000050073373233000105007337323600050073373237000500733732380027fdffff05007337333000050073373331000200000005007337333400000500733733370001"}, {true, "02000000004500000005007337313200050073373133000500733731340035fdffff0500733731360005007337313700020000000008000000050073373230000000080000000500733732330001004500000005007337323600050073373237000500733732380027fdffff0500733733300005007337333100020000000008000000050073373334000000080000000500733733370001"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetRewardsResponse(out, types.List[uplay_win_types.UplayReward]{func() uplay_win_types.UplayReward {
			v := uplay_win_types.NewUplayReward()
			v.Code = types.NewString("s712")
			v.Name = types.NewString("s713")
			v.Description = types.NewString("s714")
			v.Value = types.NewInt32(-715)
			v.RewardTypeName = types.NewString("s716")
			v.GameCode = types.NewString("s717")
			v.Platforms = types.List[uplay_win_types.UplayRewardPlatform]{func() uplay_win_types.UplayRewardPlatform {
				v := uplay_win_types.NewUplayRewardPlatform()
				v.PlatformCode = types.NewString("s720")
				v.Purchased = types.NewBool(false)
				return v
			}(), func() uplay_win_types.UplayRewardPlatform {
				v := uplay_win_types.NewUplayRewardPlatform()
				v.PlatformCode = types.NewString("s723")
				v.Purchased = types.NewBool(true)
				return v
			}()}
			return v
		}(), func() uplay_win_types.UplayReward {
			v := uplay_win_types.NewUplayReward()
			v.Code = types.NewString("s726")
			v.Name = types.NewString("s727")
			v.Description = types.NewString("s728")
			v.Value = types.NewInt32(-729)
			v.RewardTypeName = types.NewString("s730")
			v.GameCode = types.NewString("s731")
			v.Platforms = types.List[uplay_win_types.UplayRewardPlatform]{func() uplay_win_types.UplayRewardPlatform {
				v := uplay_win_types.NewUplayRewardPlatform()
				v.PlatformCode = types.NewString("s734")
				v.Purchased = types.NewBool(false)
				return v
			}(), func() uplay_win_types.UplayRewardPlatform {
				v := uplay_win_types.NewUplayRewardPlatform()
				v.PlatformCode = types.NewString("s737")
				v.Purchased = types.NewBool(true)
				return v
			}()}
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetRewardsHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getRewards(nil, nil, 1, types.NewInt32(0), types.NewInt32(0), types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getRewards(errTest, nil, 1, types.NewInt32(0), types.NewInt32(0), types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetRewardsPurchasedResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0200000005007337343100050073373432000500733734330018fdffff050073373435000500733734360002000000050073373439000105007337353200000500733735350005007337353600050073373537000afdffff05007337353900050073373630000200000005007337363300010500733736360000"}, {true, "02000000004500000005007337343100050073373432000500733734330018fdffff050073373435000500733734360002000000000800000005007337343900010008000000050073373532000000450000000500733735350005007337353600050073373537000afdffff0500733735390005007337363000020000000008000000050073373633000100080000000500733736360000"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetRewardsPurchasedResponse(out, types.List[uplay_win_types.UplayReward]{func() uplay_win_types.UplayReward {
			v := uplay_win_types.NewUplayReward()
			v.Code = types.NewString("s741")
			v.Name = types.NewString("s742")
			v.Description = types.NewString("s743")
			v.Value = types.NewInt32(-744)
			v.RewardTypeName = types.NewString("s745")
			v.GameCode = types.NewString("s746")
			v.Platforms = types.List[uplay_win_types.UplayRewardPlatform]{func() uplay_win_types.UplayRewardPlatform {
				v := uplay_win_types.NewUplayRewardPlatform()
				v.PlatformCode = types.NewString("s749")
				v.Purchased = types.NewBool(true)
				return v
			}(), func() uplay_win_types.UplayRewardPlatform {
				v := uplay_win_types.NewUplayRewardPlatform()
				v.PlatformCode = types.NewString("s752")
				v.Purchased = types.NewBool(false)
				return v
			}()}
			return v
		}(), func() uplay_win_types.UplayReward {
			v := uplay_win_types.NewUplayReward()
			v.Code = types.NewString("s755")
			v.Name = types.NewString("s756")
			v.Description = types.NewString("s757")
			v.Value = types.NewInt32(-758)
			v.RewardTypeName = types.NewString("s759")
			v.GameCode = types.NewString("s760")
			v.Platforms = types.List[uplay_win_types.UplayRewardPlatform]{func() uplay_win_types.UplayRewardPlatform {
				v := uplay_win_types.NewUplayRewardPlatform()
				v.PlatformCode = types.NewString("s763")
				v.Purchased = types.NewBool(true)
				return v
			}(), func() uplay_win_types.UplayRewardPlatform {
				v := uplay_win_types.NewUplayRewardPlatform()
				v.PlatformCode = types.NewString("s766")
				v.Purchased = types.NewBool(false)
				return v
			}()}
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetRewardsPurchasedHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getRewardsPurchased(nil, nil, 1, types.NewInt32(0), types.NewInt32(0), types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getRewardsPurchased(errTest, nil, 1, types.NewInt32(0), types.NewInt32(0), types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildUplayWelcomeResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "02000000050073373730000500733737310005007337373200fbfcffff0500733737340002000000050073373737000105007337373900050073373831000105007337383300050073373835000500733738360005007337383700ecfcffff0500733738390002000000050073373932000005007337393400050073373936000005007337393800"}, {true, "02000000004c000000050073373730000500733737310005007337373200fbfcffff0500733737340002000000000f000000050073373737000105007337373900000f000000050073373831000105007337383300004c000000050073373835000500733738360005007337383700ecfcffff0500733738390002000000000f000000050073373932000005007337393400000f000000050073373936000005007337393800"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildUplayWelcomeResponse(out, types.List[uplay_win_types.UplayAction]{func() uplay_win_types.UplayAction {
			v := uplay_win_types.NewUplayAction()
			v.Code = types.NewString("s770")
			v.Name = types.NewString("s771")
			v.Description = types.NewString("s772")
			v.Value = types.NewInt32(-773)
			v.GameCode = types.NewString("s774")
			v.Platforms = types.List[uplay_win_types.UplayActionPlatform]{func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s777")
				v.Completed = types.NewBool(true)
				v.SpecificKey = types.NewString("s779")
				return v
			}(), func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s781")
				v.Completed = types.NewBool(true)
				v.SpecificKey = types.NewString("s783")
				return v
			}()}
			return v
		}(), func() uplay_win_types.UplayAction {
			v := uplay_win_types.NewUplayAction()
			v.Code = types.NewString("s785")
			v.Name = types.NewString("s786")
			v.Description = types.NewString("s787")
			v.Value = types.NewInt32(-788)
			v.GameCode = types.NewString("s789")
			v.Platforms = types.List[uplay_win_types.UplayActionPlatform]{func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s792")
				v.Completed = types.NewBool(false)
				v.SpecificKey = types.NewString("s794")
				return v
			}(), func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s796")
				v.Completed = types.NewBool(false)
				v.SpecificKey = types.NewString("s798")
				return v
			}()}
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestUplayWelcomeHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.uplayWelcome(nil, nil, 1, types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.uplayWelcome(errTest, nil, 1, types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildSetActionCompletedResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "050073383030000500733830310005007338303200ddfcffff0500733830340002000000050073383037000105007338303900050073383131000105007338313300"}, {true, "004c000000050073383030000500733830310005007338303200ddfcffff0500733830340002000000000f000000050073383037000105007338303900000f000000050073383131000105007338313300"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildSetActionCompletedResponse(out, func() uplay_win_types.UplayAction {
			v := uplay_win_types.NewUplayAction()
			v.Code = types.NewString("s800")
			v.Name = types.NewString("s801")
			v.Description = types.NewString("s802")
			v.Value = types.NewInt32(-803)
			v.GameCode = types.NewString("s804")
			v.Platforms = types.List[uplay_win_types.UplayActionPlatform]{func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s807")
				v.Completed = types.NewBool(true)
				v.SpecificKey = types.NewString("s809")
				return v
			}(), func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s811")
				v.Completed = types.NewBool(true)
				v.SpecificKey = types.NewString("s813")
				return v
			}()}
			return v
		}())
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestSetActionCompletedHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.setActionCompleted(nil, nil, 1, types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.setActionCompleted(errTest, nil, 1, types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildSetActionsCompletedResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "02000000050073383136000500733831370005007338313800cdfcffff0500733832300002000000050073383233000105007338323500050073383237000105007338323900050073383331000500733833320005007338333300befcffff0500733833350002000000050073383338000005007338343000050073383432000005007338343400"}, {true, "02000000004c000000050073383136000500733831370005007338313800cdfcffff0500733832300002000000000f000000050073383233000105007338323500000f000000050073383237000105007338323900004c000000050073383331000500733833320005007338333300befcffff0500733833350002000000000f000000050073383338000005007338343000000f000000050073383432000005007338343400"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildSetActionsCompletedResponse(out, types.List[uplay_win_types.UplayAction]{func() uplay_win_types.UplayAction {
			v := uplay_win_types.NewUplayAction()
			v.Code = types.NewString("s816")
			v.Name = types.NewString("s817")
			v.Description = types.NewString("s818")
			v.Value = types.NewInt32(-819)
			v.GameCode = types.NewString("s820")
			v.Platforms = types.List[uplay_win_types.UplayActionPlatform]{func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s823")
				v.Completed = types.NewBool(true)
				v.SpecificKey = types.NewString("s825")
				return v
			}(), func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s827")
				v.Completed = types.NewBool(true)
				v.SpecificKey = types.NewString("s829")
				return v
			}()}
			return v
		}(), func() uplay_win_types.UplayAction {
			v := uplay_win_types.NewUplayAction()
			v.Code = types.NewString("s831")
			v.Name = types.NewString("s832")
			v.Description = types.NewString("s833")
			v.Value = types.NewInt32(-834)
			v.GameCode = types.NewString("s835")
			v.Platforms = types.List[uplay_win_types.UplayActionPlatform]{func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s838")
				v.Completed = types.NewBool(false)
				v.SpecificKey = types.NewString("s840")
				return v
			}(), func() uplay_win_types.UplayActionPlatform {
				v := uplay_win_types.NewUplayActionPlatform()
				v.PlatformCode = types.NewString("s842")
				v.Completed = types.NewBool(false)
				v.SpecificKey = types.NewString("s844")
				return v
			}()}
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestSetActionsCompletedHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.setActionsCompleted(nil, nil, 1, types.NewList[types.String](), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.setActionsCompleted(errTest, nil, 1, types.NewList[types.String](), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetUserTokenResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "05007338343500"}, {true, "05007338343500"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetUserTokenResponse(out, types.NewString("s845"))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetUserTokenHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getUserToken(nil, nil, 1)
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getUserToken(errTest, nil, 1)
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetVirtualCurrencyUserBalanceResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "b2fcffff"}, {true, "b2fcffff"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetVirtualCurrencyUserBalanceResponse(out, types.NewInt32(-846))
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetVirtualCurrencyUserBalanceHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getVirtualCurrencyUserBalance(nil, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getVirtualCurrencyUserBalance(errTest, nil, 1, types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}

func TestBuildGetSectionsByKeyResponse(t *testing.T) {
	cases := []struct {
		header bool
		data   string
	}{{false, "0200000005007338343900050073383530000500733835310005007338353200020000000500733835350005007338353600bd030500733835380005007338363000050073383631000500733836320005007338363300a0fcffff0500733836350005007338363600050073383637000500733836390005007338373000cb03050073383732000500733837340005007338373500050073383736000500733837370092fcffff050073383739000500733838300005007338383100050073383832000500733838330005007338383500050073383836000500733838370005007338383800020000000500733839310005007338393200e10305007338393400050073383936000500733839370005007338393800050073383939007cfcffff0500733930310005007339303200050073393033000500733930350005007339303600ef0305007339303800050073393130000500733931310005007339313200050073393133006efcffff0500733931350005007339313600050073393137000500733931380005007339313900"}, {true, "0200000000da000000050073383439000500733835300005007338353100050073383532000200000000510000000500733835350005007338353600bd0305007338353800003500000005007338363000050073383631000500733836320005007338363300a0fcffff05007338363500050073383636000500733836370000510000000500733836390005007338373000cb030500733837320000350000000500733837340005007338373500050073383736000500733837370092fcffff050073383739000500733838300005007338383100050073383832000500733838330000da000000050073383835000500733838360005007338383700050073383838000200000000510000000500733839310005007338393200e103050073383934000035000000050073383936000500733839370005007338393800050073383939007cfcffff05007339303100050073393032000500733930330000510000000500733930350005007339303600ef03050073393038000035000000050073393130000500733931310005007339313200050073393133006efcffff0500733931350005007339313600050073393137000500733931380005007339313900"}}

	for _, c := range cases {
		out := newOut(c.header)
		buildGetSectionsByKeyResponse(out, types.List[uplay_win_types.UplaySection]{func() uplay_win_types.UplaySection {
			v := uplay_win_types.NewUplaySection()
			v.Key = types.NewString("s849")
			v.Name = types.NewString("s850")
			v.TypeName = types.NewString("s851")
			v.MenuTypeName = types.NewString("s852")
			v.ContentList = types.List[uplay_win_types.UplaySectionContent]{func() uplay_win_types.UplaySectionContent {
				v := uplay_win_types.NewUplaySectionContent()
				v.Key = types.NewString("s855")
				v.Name = types.NewString("s856")
				v.Order = types.NewInt16(957)
				v.TypeName = types.NewString("s858")
				v.LocalizedInfo = func() uplay_win_types.UplaySectionContentLocalized {
					v := uplay_win_types.NewUplaySectionContentLocalized()
					v.Key = types.NewString("s860")
					v.Culture = types.NewString("s861")
					v.Text = types.NewString("s862")
					v.Url = types.NewString("s863")
					v.Duration = types.NewInt32(-864)
					v.Size = types.NewString("s865")
					v.Width = types.NewString("s866")
					v.Height = types.NewString("s867")
					return v
				}()
				return v
			}(), func() uplay_win_types.UplaySectionContent {
				v := uplay_win_types.NewUplaySectionContent()
				v.Key = types.NewString("s869")
				v.Name = types.NewString("s870")
				v.Order = types.NewInt16(971)
				v.TypeName = types.NewString("s872")
				v.LocalizedInfo = func() uplay_win_types.UplaySectionContentLocalized {
					v := uplay_win_types.NewUplaySectionContentLocalized()
					v.Key = types.NewString("s874")
					v.Culture = types.NewString("s875")
					v.Text = types.NewString("s876")
					v.Url = types.NewString("s877")
					v.Duration = types.NewInt32(-878)
					v.Size = types.NewString("s879")
					v.Width = types.NewString("s880")
					v.Height = types.NewString("s881")
					return v
				}()
				return v
			}()}
			v.GameCode = types.NewString("s882")
			v.PlatformCode = types.NewString("s883")
			return v
		}(), func() uplay_win_types.UplaySection {
			v := uplay_win_types.NewUplaySection()
			v.Key = types.NewString("s885")
			v.Name = types.NewString("s886")
			v.TypeName = types.NewString("s887")
			v.MenuTypeName = types.NewString("s888")
			v.ContentList = types.List[uplay_win_types.UplaySectionContent]{func() uplay_win_types.UplaySectionContent {
				v := uplay_win_types.NewUplaySectionContent()
				v.Key = types.NewString("s891")
				v.Name = types.NewString("s892")
				v.Order = types.NewInt16(993)
				v.TypeName = types.NewString("s894")
				v.LocalizedInfo = func() uplay_win_types.UplaySectionContentLocalized {
					v := uplay_win_types.NewUplaySectionContentLocalized()
					v.Key = types.NewString("s896")
					v.Culture = types.NewString("s897")
					v.Text = types.NewString("s898")
					v.Url = types.NewString("s899")
					v.Duration = types.NewInt32(-900)
					v.Size = types.NewString("s901")
					v.Width = types.NewString("s902")
					v.Height = types.NewString("s903")
					return v
				}()
				return v
			}(), func() uplay_win_types.UplaySectionContent {
				v := uplay_win_types.NewUplaySectionContent()
				v.Key = types.NewString("s905")
				v.Name = types.NewString("s906")
				v.Order = types.NewInt16(1007)
				v.TypeName = types.NewString("s908")
				v.LocalizedInfo = func() uplay_win_types.UplaySectionContentLocalized {
					v := uplay_win_types.NewUplaySectionContentLocalized()
					v.Key = types.NewString("s910")
					v.Culture = types.NewString("s911")
					v.Text = types.NewString("s912")
					v.Url = types.NewString("s913")
					v.Duration = types.NewInt32(-914)
					v.Size = types.NewString("s915")
					v.Width = types.NewString("s916")
					v.Height = types.NewString("s917")
					return v
				}()
				return v
			}()}
			v.GameCode = types.NewString("s918")
			v.PlatformCode = types.NewString("s919")
			return v
		}()})
		checkBytes(t, out.Bytes(), c.data)
	}
}

func TestGetSectionsByKeyHooks(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, err := commonProtocol.getSectionsByKey(nil, nil, 1, types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.NotImplemented, "x").ResultCode {
		t.Fatalf("missing hook: got %v, want Core::NotImplemented", err)
	}

	_, err = commonProtocol.getSectionsByKey(errTest, nil, 1, types.NewString(""), types.NewString(""), types.NewString(""), types.NewString(""))
	if err == nil || err.ResultCode != nex.NewError(nex.ResultCodes.Core.InvalidArgument, "x").ResultCode {
		t.Fatalf("parse error: got %v, want Core::InvalidArgument", err)
	}
}
