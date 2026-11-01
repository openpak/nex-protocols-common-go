package ranking2

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/fakenet"
	ranking2 "github.com/PretendoNetwork/nex-protocols-go/v2/ranking2"
	"github.com/PretendoNetwork/nex-protocols-go/v2/ranking2/constants"
	ranking2_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking2/types"
)

const callID = 0x01020304

func TestRegistersEveryMethod(t *testing.T) {
	protocol := ranking2.NewProtocol()
	NewCommonProtocol(protocol)

	if protocol.PutScore == nil || protocol.GetCommonData == nil || protocol.PutCommonData == nil ||
		protocol.DeleteCommonData == nil || protocol.GetRanking == nil || protocol.GetRankingByPrincipalID == nil ||
		protocol.GetCategorySetting == nil || protocol.GetRankingChart == nil || protocol.GetRankingCharts == nil ||
		protocol.GetEstimateScoreRank == nil {
		t.Fatal("a Ranking2 method was left without a handler")
	}
}

func TestNotImplementedWithoutHooks(t *testing.T) {
	c := &CommonProtocol{}
	p := fakenet.NewPacket(7)

	msg, err := c.putScore(nil, p, callID, nil, 0)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
	msg, err = c.getCommonData(nil, p, callID, 0, 0, 0)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
	msg, err = c.putCommonData(nil, p, callID, ranking2_types.NewRanking2CommonData(), 0)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
	msg, err = c.deleteCommonData(nil, p, callID, 0)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
	msg, err = c.getRanking(nil, p, callID, ranking2_types.NewRanking2GetParam())
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
	msg, err = c.getRankingByPrincipalID(nil, p, callID, ranking2_types.NewRanking2GetByListParam(), nil)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
	msg, err = c.getCategorySetting(nil, p, callID, 0)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
	msg, err = c.getRankingChart(nil, p, callID, ranking2_types.NewRanking2ChartInfoInput())
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
	msg, err = c.getRankingCharts(nil, p, callID, nil)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
	msg, err = c.getEstimateScoreRank(nil, p, callID, ranking2_types.NewRanking2EstimateScoreRankInput())
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
}

func TestPutScoreEmptyResponseAndArguments(t *testing.T) {
	var gotPID types.PID
	var gotList types.List[ranking2_types.Ranking2ScoreData]
	var gotUnique types.UInt64

	c := &CommonProtocol{PutScore: func(pid types.PID, l types.List[ranking2_types.Ranking2ScoreData], u types.UInt64) error {
		gotPID, gotList, gotUnique = pid, l, u
		return nil
	}}

	score := ranking2_types.NewRanking2ScoreData()
	score.Category = 3
	msg, err := c.putScore(nil, fakenet.NewPacket(42), callID, types.List[ranking2_types.Ranking2ScoreData]{score}, 99)
	fakenet.RequireBody(t, msg, err, ranking2.ProtocolID, ranking2.MethodPutScore, callID, nil)

	if gotPID != 42 || len(gotList) != 1 || gotList[0].Category != 3 || gotUnique != 99 {
		t.Fatalf("hook got pid=%d list=%v unique=%d", gotPID, gotList, gotUnique)
	}
}

func TestDecodeErrorIsInvalidArgument(t *testing.T) {
	c := &CommonProtocol{PutScore: func(types.PID, types.List[ranking2_types.Ranking2ScoreData], types.UInt64) error { return nil }}
	msg, err := c.putScore(errors.New("bad"), fakenet.NewPacket(1), callID, nil, 0)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Ranking2.InvalidArgument)
}

func TestHookErrors(t *testing.T) {
	c := &CommonProtocol{DeleteCommonData: func(types.PID, types.UInt64) error { return errors.New("db down") }}
	msg, err := c.deleteCommonData(nil, fakenet.NewPacket(1), callID, 1)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Ranking2.Unknown)

	c.DeleteCommonData = func(types.PID, types.UInt64) error {
		return nex.NewError(nex.ResultCodes.Ranking2.InvalidScore, "nope")
	}
	msg, err = c.deleteCommonData(nil, fakenet.NewPacket(1), callID, 1)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Ranking2.InvalidScore)
}

func TestGetCommonDataLayout(t *testing.T) {
	c := &CommonProtocol{GetCommonData: func(pid types.PID, flags constants.Ranking2GetOptionFlags, principal types.PID, unique types.UInt64) (ranking2_types.Ranking2CommonData, error) {
		d := ranking2_types.NewRanking2CommonData()
		d.UserName = "Ab"
		d.Mii = types.QBuffer{0xAA}
		d.BinaryData = types.QBuffer{1, 2, 3}
		return d, nil
	}}

	msg, err := c.getCommonData(nil, fakenet.NewPacket(5), callID, 0, 5, 1)
	// String: u16 length (3 with NUL) + "Ab\0"; QBuffer: u16 length + bytes
	fakenet.RequireBody(t, msg, err, ranking2.ProtocolID, ranking2.MethodGetCommonData, callID, []byte{
		0x03, 0x00, 'A', 'b', 0x00,
		0x01, 0x00, 0xAA,
		0x03, 0x00, 1, 2, 3,
	})
}

func TestGetRankingInfoLayout(t *testing.T) {
	c := &CommonProtocol{GetRanking: func(pid types.PID, p ranking2_types.Ranking2GetParam) (ranking2_types.Ranking2Info, error) {
		rd := ranking2_types.NewRanking2RankData()
		rd.Misc = 0x1122334455667788
		rd.NexUniqueID = 2
		rd.PrincipalID = 0x0A0B0C0D
		rd.Rank = 1
		rd.Score = 500

		info := ranking2_types.NewRanking2Info()
		info.RankDataList = types.List[ranking2_types.Ranking2RankData]{rd}
		info.LowestRank = 10
		info.NumRankedIn = 20
		info.Season = -1
		return info, nil
	}}

	msg, err := c.getRanking(nil, fakenet.NewPacket(5), callID, ranking2_types.NewRanking2GetParam())
	fakenet.RequireBody(t, msg, err, ranking2.ProtocolID, ranking2.MethodGetRanking, callID, []byte{
		0x01, 0x00, 0x00, 0x00, // list length
		0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, // misc
		0x02, 0, 0, 0, 0, 0, 0, 0, // nexUniqueId
		0x0D, 0x0C, 0x0B, 0x0A, // pid (32 bit)
		0x01, 0, 0, 0, // rank
		0xF4, 0x01, 0, 0, // score 500
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // empty commonData: "" (u16 len 1 + NUL), two empty qBuffers
		0x0A, 0, 0, 0, // lowestRank
		0x14, 0, 0, 0, // numRankedIn
		0xFF, 0xFF, 0xFF, 0xFF, // season -1
	})
}

func TestGetRankingByPrincipalIDPassesList(t *testing.T) {
	var got types.List[types.PID]
	c := &CommonProtocol{GetRankingByPrincipalID: func(pid types.PID, p ranking2_types.Ranking2GetByListParam, l types.List[types.PID]) (ranking2_types.Ranking2Info, error) {
		got = l
		return ranking2_types.NewRanking2Info(), nil
	}}

	msg, err := c.getRankingByPrincipalID(nil, fakenet.NewPacket(5), callID, ranking2_types.NewRanking2GetByListParam(), types.List[types.PID]{1, 2})
	fakenet.RequireBody(t, msg, err, ranking2.ProtocolID, ranking2.MethodGetRankingByPrincipalID, callID, []byte{
		0, 0, 0, 0, // empty rankDataList
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // lowestRank, numRankedIn, season
	})

	if len(got) != 2 || got[1] != 2 {
		t.Fatalf("list %v", got)
	}
}

func TestGetCategorySettingLayout(t *testing.T) {
	c := &CommonProtocol{GetCategorySetting: func(category types.UInt32) (ranking2_types.Ranking2CategorySetting, error) {
		s := ranking2_types.NewRanking2CategorySetting()
		s.MinScore = 1
		s.MaxScore = 2
		s.LowestRank = 3
		s.ResetMonth = 0x0102
		s.ResetDay = 4
		s.ResetHour = 5
		s.ResetMode = 6
		s.MaxSeasonsToGoBack = 7
		s.ScoreOrder = true
		return s, nil
	}}

	msg, err := c.getCategorySetting(nil, fakenet.NewPacket(5), callID, 9)
	fakenet.RequireBody(t, msg, err, ranking2.ProtocolID, ranking2.MethodGetCategorySetting, callID, []byte{
		1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0,
		0x02, 0x01, // resetMonth u16
		4, 5, 6, 7, 1,
	})
}

func TestGetRankingChartsFallsBackToSingleChartHook(t *testing.T) {
	calls := 0
	c := &CommonProtocol{GetRankingChart: func(pid types.PID, in ranking2_types.Ranking2ChartInfoInput) (ranking2_types.Ranking2ChartInfo, error) {
		calls++
		ci := ranking2_types.NewRanking2ChartInfo()
		ci.Index = in.ChartIndex
		return ci, nil
	}}

	a, b := ranking2_types.NewRanking2ChartInfoInput(), ranking2_types.NewRanking2ChartInfoInput()
	a.ChartIndex, b.ChartIndex = 7, 8

	msg, err := c.getRankingCharts(nil, fakenet.NewPacket(5), callID, types.List[ranking2_types.Ranking2ChartInfoInput]{a, b})
	if err != nil || calls != 2 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}

	stream := fakenet.NewPacket(5).Stream(msg.Parameters)
	n, _ := stream.ReadUInt32LE()
	if n != 2 {
		t.Fatalf("chart count %d", n)
	}
}

func TestGetEstimateScoreRankLayout(t *testing.T) {
	c := &CommonProtocol{GetEstimateScoreRank: func(pid types.PID, in ranking2_types.Ranking2EstimateScoreRankInput) (ranking2_types.Ranking2EstimateScoreRankOutput, error) {
		o := ranking2_types.NewRanking2EstimateScoreRankOutput()
		o.Rank = 1
		o.Length = 2
		o.Score = in.Score
		o.Category = in.Category
		o.Season = -2
		o.SamplingRate = 9
		return o, nil
	}}

	in := ranking2_types.NewRanking2EstimateScoreRankInput()
	in.Category, in.Score = 4, 300

	msg, err := c.getEstimateScoreRank(nil, fakenet.NewPacket(5), callID, in)
	fakenet.RequireBody(t, msg, err, ranking2.ProtocolID, ranking2.MethodGetEstimateScoreRank, callID, []byte{
		1, 0, 0, 0, 2, 0, 0, 0, 0x2C, 0x01, 0, 0, 4, 0, 0, 0, 0xFE, 0xFF, 0xFF, 0xFF, 9,
	})
}
