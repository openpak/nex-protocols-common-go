package matchmakereferee

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/fakenet"
	matchmake_referee "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-referee"
	mr_types "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-referee/types"
)

const callID = 0x66

func TestRegistersEveryMethod(t *testing.T) {
	p := matchmake_referee.NewProtocol()
	NewCommonProtocol(p)

	if p.StartRound == nil || p.GetStartRoundParam == nil || p.EndRound == nil || p.EndRoundWithoutReport == nil ||
		p.GetRoundParticipants == nil || p.GetNotSummarizedRound == nil || p.GetRound == nil || p.GetStatsPrimary == nil ||
		p.GetStatsPrimaries == nil || p.GetStatsAll == nil || p.CreateStats == nil || p.GetOrCreateStats == nil || p.ResetStats == nil {
		t.Fatal("a MatchmakeReferee method was left without a handler")
	}
}

func TestNotImplementedWithoutHooks(t *testing.T) {
	c := &CommonProtocol{}
	p := fakenet.NewPacket(3)
	code := nex.ResultCodes.Core.NotImplemented
	start := mr_types.NewMatchmakeRefereeStartRoundParam()
	target := mr_types.NewMatchmakeRefereeStatsTarget()
	initParam := mr_types.NewMatchmakeRefereeStatsInitParam()

	msg, err := c.startRound(nil, p, callID, start)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getStartRoundParam(nil, p, callID, 1)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.endRound(nil, p, callID, mr_types.NewMatchmakeRefereeEndRoundParam())
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.endRoundWithoutReport(nil, p, callID, 1)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getRoundParticipants(nil, p, callID, 1)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getNotSummarizedRound(nil, p, callID)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getRound(nil, p, callID, 1)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getStatsPrimary(nil, p, callID, target)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getStatsPrimaries(nil, p, callID, nil)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getStatsAll(nil, p, callID, target)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.createStats(nil, p, callID, initParam)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getOrCreateStats(nil, p, callID, initParam)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.resetStats(nil, p, callID)
	fakenet.RequireError(t, msg, err, code)
}

func TestStartRound(t *testing.T) {
	c := &CommonProtocol{StartRound: func(types.PID, mr_types.MatchmakeRefereeStartRoundParam) (types.UInt64, error) {
		return 0x0102030405060708, nil
	}}
	p := fakenet.NewPacket(3)

	empty := mr_types.NewMatchmakeRefereeStartRoundParam()
	msg, err := c.startRound(nil, p, callID, empty)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.InvalidArgument)

	param := mr_types.NewMatchmakeRefereeStartRoundParam()
	param.PIDs = types.List[types.PID]{3}
	msg, err = c.startRound(nil, p, callID, param)
	fakenet.RequireBody(t, msg, err, matchmake_referee.ProtocolID, matchmake_referee.MethodStartRound, callID,
		[]byte{8, 7, 6, 5, 4, 3, 2, 1})
}

func TestGetStartRoundParamLayout(t *testing.T) {
	c := &CommonProtocol{GetStartRoundParam: func(types.PID, types.UInt64) (mr_types.MatchmakeRefereeStartRoundParam, error) {
		r := mr_types.NewMatchmakeRefereeStartRoundParam()
		r.PersonalDataCategory = 1
		r.GID = 2
		r.PIDs = types.List[types.PID]{7, 8}
		return r, nil
	}}

	msg, err := c.getStartRoundParam(nil, fakenet.NewPacket(3), callID, 9)
	fakenet.RequireBody(t, msg, err, matchmake_referee.ProtocolID, matchmake_referee.MethodGetStartRoundParam, callID, []byte{
		1, 0, 0, 0, 2, 0, 0, 0,
		2, 0, 0, 0, 7, 0, 0, 0, 8, 0, 0, 0,
	})
}

func TestEmptyResponses(t *testing.T) {
	c := &CommonProtocol{
		EndRound:              func(types.PID, mr_types.MatchmakeRefereeEndRoundParam) error { return nil },
		EndRoundWithoutReport: func(types.PID, types.UInt64) error { return nil },
		CreateStats:           func(types.PID, mr_types.MatchmakeRefereeStatsInitParam) error { return nil },
		ResetStats:            func(types.PID) error { return nil },
	}
	p := fakenet.NewPacket(3)

	msg, err := c.endRound(nil, p, callID, mr_types.NewMatchmakeRefereeEndRoundParam())
	fakenet.RequireBody(t, msg, err, matchmake_referee.ProtocolID, matchmake_referee.MethodEndRound, callID, nil)
	msg, err = c.endRoundWithoutReport(nil, p, callID, 1)
	fakenet.RequireBody(t, msg, err, matchmake_referee.ProtocolID, matchmake_referee.MethodEndRoundWithoutReport, callID, nil)
	msg, err = c.createStats(nil, p, callID, mr_types.NewMatchmakeRefereeStatsInitParam())
	fakenet.RequireBody(t, msg, err, matchmake_referee.ProtocolID, matchmake_referee.MethodCreateStats, callID, nil)
	msg, err = c.resetStats(nil, p, callID)
	fakenet.RequireBody(t, msg, err, matchmake_referee.ProtocolID, matchmake_referee.MethodResetStats, callID, nil)
}

func TestParticipantsAndRoundLayouts(t *testing.T) {
	c := &CommonProtocol{
		GetRoundParticipants: func(types.PID, types.UInt64) (types.List[types.UInt64], error) {
			return types.List[types.UInt64]{5}, nil
		},
		GetRound: func(types.PID, types.UInt64) (mr_types.MatchmakeRefereeRound, error) {
			r := mr_types.NewMatchmakeRefereeRound()
			r.RoundID = 1
			r.GID = 2
			r.State = 3
			r.PersonalDataCategory = 4
			res := mr_types.NewMatchmakeRefereePersonalRoundResult()
			res.PID = 9
			res.PersonalRoundResultFlag = 10
			res.RoundWinLoss = 11
			res.RatingValueChange = -1
			res.Buffer = types.QBuffer{0xAB}
			r.NormalizedPersonalRoundResults = types.List[mr_types.MatchmakeRefereePersonalRoundResult]{res}
			return r, nil
		},
	}
	p := fakenet.NewPacket(3)

	msg, err := c.getRoundParticipants(nil, p, callID, 1)
	fakenet.RequireBody(t, msg, err, matchmake_referee.ProtocolID, matchmake_referee.MethodGetRoundParticipants, callID,
		[]byte{1, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0})

	msg, err = c.getRound(nil, p, callID, 1)
	fakenet.RequireBody(t, msg, err, matchmake_referee.ProtocolID, matchmake_referee.MethodGetRound, callID, []byte{
		1, 0, 0, 0, 0, 0, 0, 0, // roundId
		2, 0, 0, 0, 3, 0, 0, 0, 4, 0, 0, 0,
		1, 0, 0, 0, // result list length
		9, 0, 0, 0, 10, 0, 0, 0, 11, 0, 0, 0, 0xFF, 0xFF, 0xFF, 0xFF,
		1, 0, 0xAB,
	})
}

func TestStatsLayout(t *testing.T) {
	stats := mr_types.NewMatchmakeRefereeStats()
	stats.UniqueID = 1
	stats.Category = 2
	stats.PID = 3
	stats.RatingValue = 99

	c := &CommonProtocol{GetStatsPrimary: func(types.PID, mr_types.MatchmakeRefereeStatsTarget) (mr_types.MatchmakeRefereeStats, error) {
		return stats, nil
	}}

	want := make([]byte, 0, 64)
	want = append(want, 1, 0, 0, 0, 0, 0, 0, 0) // uniqueId
	want = append(want, 2, 0, 0, 0)             // category
	want = append(want, 3, 0, 0, 0)             // pid
	want = append(want, make([]byte, 12*4)...)  // twelve zero counters
	want = append(want, 99, 0, 0, 0)            // ratingValue

	msg, err := c.getStatsPrimary(nil, fakenet.NewPacket(3), callID, mr_types.NewMatchmakeRefereeStatsTarget())
	fakenet.RequireBody(t, msg, err, matchmake_referee.ProtocolID, matchmake_referee.MethodGetStatsPrimary, callID, want)
}

func TestGetStatsPrimariesFallbackAndResults(t *testing.T) {
	c := &CommonProtocol{GetStatsPrimary: func(pid types.PID, target mr_types.MatchmakeRefereeStatsTarget) (mr_types.MatchmakeRefereeStats, error) {
		if target.Category == 2 {
			return mr_types.NewMatchmakeRefereeStats(), nex.NewError(nex.ResultCodes.MatchmakeReferee.StatsNotFound, "none")
		}
		s := mr_types.NewMatchmakeRefereeStats()
		s.Category = target.Category
		return s, nil
	}}

	a, b := mr_types.NewMatchmakeRefereeStatsTarget(), mr_types.NewMatchmakeRefereeStatsTarget()
	a.Category, b.Category = 1, 2

	msg, err := c.getStatsPrimaries(nil, fakenet.NewPacket(3), callID, types.List[mr_types.MatchmakeRefereeStatsTarget]{a, b})
	if err != nil {
		t.Fatal(err)
	}

	// Each stats entry is 8 (uniqueId) + 4 + 4 + 12*4 + 4 = 68 bytes.
	want := []byte{2, 0, 0, 0}
	first := make([]byte, 68)
	first[8] = 1 // category
	want = append(want, first...)
	want = append(want, make([]byte, 68)...)
	want = append(want, 2, 0, 0, 0)          // results length
	want = append(want, 1, 0, 1, 0)          // success
	want = append(want, 0x06, 0, 0x6F, 0x80) // MatchmakeReferee::StatsNotFound

	fakenet.RequireBody(t, msg, nil, matchmake_referee.ProtocolID, matchmake_referee.MethodGetStatsPrimaries, callID, want)
}

func TestHookErrors(t *testing.T) {
	c := &CommonProtocol{ResetStats: func(types.PID) error { return errors.New("x") }}
	p := fakenet.NewPacket(3)

	msg, err := c.resetStats(nil, p, callID)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.MatchmakeReferee.Unknown)
	msg, err = c.resetStats(errors.New("short"), p, callID)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.MatchmakeReferee.InvalidArgument)

	c.ResetStats = func(types.PID) error {
		return nex.NewError(nex.ResultCodes.MatchmakeReferee.RoundNotFound, "r")
	}
	msg, err = c.resetStats(nil, p, callID)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.MatchmakeReferee.RoundNotFound)
}
