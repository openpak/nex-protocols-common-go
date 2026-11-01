package rating

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/fakenet"
	rating "github.com/PretendoNetwork/nex-protocols-go/v2/rating"
	rating_types "github.com/PretendoNetwork/nex-protocols-go/v2/rating/types"
)

const callID = 0x0A0B0C0D

func TestRegistersEveryMethod(t *testing.T) {
	protocol := rating.NewProtocol()
	NewCommonProtocol(protocol)

	if protocol.Unk1 == nil || protocol.Unk2 == nil || protocol.ReportRatingStats == nil || protocol.GetRanking == nil ||
		protocol.DeleteScore == nil || protocol.UploadCommonData == nil || protocol.GetCommonData == nil {
		t.Fatal("a Rating method was left without a handler")
	}
}

func TestNotImplementedWithoutHooks(t *testing.T) {
	c := &CommonProtocol{}
	p := fakenet.NewPacket(7)
	tok := rating_types.NewRatingSessionToken()
	code := nex.ResultCodes.Core.NotImplemented

	msg, err := c.unk1(nil, p, callID)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.unk2(nil, p, callID, tok)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.reportRatingStats(nil, p, callID, tok, nil)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getRanking(nil, p, callID, 1, 2, 3)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.deleteScore(nil, p, callID, 1, 2)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.uploadCommonData(nil, p, callID, nil, 2)
	fakenet.RequireError(t, msg, err, code)
	msg, err = c.getCommonData(nil, p, callID, 2)
	fakenet.RequireError(t, msg, err, code)
}

func TestUnk1Layout(t *testing.T) {
	c := &CommonProtocol{Unk1: func(pid types.PID) (types.UInt64, types.String, error) {
		return 0x0102030405060708, "ab", nil
	}}

	msg, err := c.unk1(nil, fakenet.NewPacket(7), callID)
	fakenet.RequireBody(t, msg, err, rating.ProtocolID, rating.MethodUnk1, callID, []byte{
		8, 7, 6, 5, 4, 3, 2, 1,
		0x03, 0x00, 'a', 'b', 0x00,
	})
}

func TestEmptyResponses(t *testing.T) {
	ok := func() error { return nil }
	c := &CommonProtocol{
		Unk2: func(types.PID, rating_types.RatingSessionToken) error { return ok() },
		ReportRatingStats: func(types.PID, rating_types.RatingSessionToken, types.List[rating_types.RatingStats]) error {
			return ok()
		},
		DeleteScore:      func(types.PID, types.UInt32, types.UInt64) error { return ok() },
		UploadCommonData: func(types.PID, types.Buffer, types.UInt64) error { return ok() },
	}
	p := fakenet.NewPacket(7)
	tok := rating_types.NewRatingSessionToken()

	msg, err := c.unk2(nil, p, callID, tok)
	fakenet.RequireBody(t, msg, err, rating.ProtocolID, rating.MethodUnk2, callID, nil)
	msg, err = c.reportRatingStats(nil, p, callID, tok, nil)
	fakenet.RequireBody(t, msg, err, rating.ProtocolID, rating.MethodReportRatingStats, callID, nil)
	msg, err = c.deleteScore(nil, p, callID, 1, 2)
	fakenet.RequireBody(t, msg, err, rating.ProtocolID, rating.MethodDeleteScore, callID, nil)
	msg, err = c.uploadCommonData(nil, p, callID, types.Buffer{1}, 2)
	fakenet.RequireBody(t, msg, err, rating.ProtocolID, rating.MethodUploadCommonData, callID, nil)
}

func TestHookArgumentsAndCallerPID(t *testing.T) {
	var gotPID, gotPrincipal types.PID
	var gotCategory types.UInt32
	var gotUnique types.UInt64

	c := &CommonProtocol{GetRanking: func(pid types.PID, category types.UInt32, uniqueID types.UInt64, principalID types.PID) (rating_types.RatingRankData, error) {
		gotPID, gotCategory, gotUnique, gotPrincipal = pid, category, uniqueID, principalID
		return rating_types.NewRatingRankData(), nil
	}}

	_, err := c.getRanking(nil, fakenet.NewPacket(11), callID, 4, 5, 6)
	if err != nil || gotPID != 11 || gotCategory != 4 || gotUnique != 5 || gotPrincipal != 6 {
		t.Fatalf("err=%v pid=%d cat=%d unique=%d principal=%d", err, gotPID, gotCategory, gotUnique, gotPrincipal)
	}
}

func TestGetRankingLayout(t *testing.T) {
	c := &CommonProtocol{GetRanking: func(types.PID, types.UInt32, types.UInt64, types.PID) (rating_types.RatingRankData, error) {
		d := rating_types.NewRatingRankData()
		d.PrincipalID = 1
		d.UniqueID = 2
		d.Order = 3
		d.Category = 4
		d.Score = 5
		d.Unknown1 = 6
		d.Unknown2 = 7
		d.Unknown3 = 8
		d.CommonData = types.Buffer{0xEE}
		d.UpdateTime = types.NewDateTime(0x0102030405060708)
		return d, nil
	}}

	msg, err := c.getRanking(nil, fakenet.NewPacket(1), callID, 4, 2, 1)
	fakenet.RequireBody(t, msg, err, rating.ProtocolID, rating.MethodGetRanking, callID, []byte{
		1, 0, 0, 0, // pid
		2, 0, 0, 0, 0, 0, 0, 0, // uniqueId
		3, 0, 0, 0, 4, 0, 0, 0, 5, 0, 0, 0, 6, 0, 0, 0, 7, 0, 0, 0, 8, 0, 0, 0,
		1, 0, 0, 0, 0xEE, // Buffer: u32 length + data
		8, 7, 6, 5, 4, 3, 2, 1, // DateTime
	})
}

func TestGetCommonDataLayout(t *testing.T) {
	c := &CommonProtocol{GetCommonData: func(types.PID, types.UInt64) (types.Buffer, error) {
		return types.Buffer{9, 8}, nil
	}}

	msg, err := c.getCommonData(nil, fakenet.NewPacket(1), callID, 2)
	fakenet.RequireBody(t, msg, err, rating.ProtocolID, rating.MethodGetCommonData, callID, []byte{2, 0, 0, 0, 9, 8})
}

func TestErrors(t *testing.T) {
	c := &CommonProtocol{DeleteScore: func(types.PID, types.UInt32, types.UInt64) error { return errors.New("x") }}
	p := fakenet.NewPacket(1)

	msg, err := c.deleteScore(nil, p, callID, 1, 2)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.Unknown)

	msg, err = c.deleteScore(errors.New("short read"), p, callID, 1, 2)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.InvalidArgument)

	c.DeleteScore = func(types.PID, types.UInt32, types.UInt64) error {
		return nex.NewError(nex.ResultCodes.Core.AccessDenied, "denied")
	}
	msg, err = c.deleteScore(nil, p, callID, 1, 2)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.AccessDenied)
}
