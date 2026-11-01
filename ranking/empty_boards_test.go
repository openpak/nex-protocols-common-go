package ranking

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/fakenet"
	ranking "github.com/PretendoNetwork/nex-protocols-go/v2/ranking"
	"github.com/PretendoNetwork/nex-protocols-go/v2/ranking/constants"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
)

const callID = 0x77

// An empty RankingResult: list length 0, total count 0, sinceTime 0x1F40420000.
var emptyResult = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0x00, 0x00, 0x42, 0x40, 0x1F, 0x00, 0x00, 0x00}

func TestWithoutUseEmptyBoardsNothingIsImplemented(t *testing.T) {
	c := &CommonProtocol{}
	p := fakenet.NewPacket(3)

	msg, err := c.uploadScore(nil, p, callID, ranking_types.NewRankingScoreData(), 1)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
	msg, err = c.getCommonData(nil, p, callID, 1)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
	msg, err = c.getRanking(nil, p, callID, constants.RankingModeRange, 1, ranking_types.NewRankingOrderParam(), 0, 0)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.NotImplemented)
}

func TestUseEmptyBoardsAnswersEveryMethod(t *testing.T) {
	c := &CommonProtocol{}
	c.UseEmptyBoards()
	p := fakenet.NewPacket(3)

	msg, err := c.uploadScore(nil, p, callID, ranking_types.NewRankingScoreData(), 1)
	fakenet.RequireBody(t, msg, err, ranking.ProtocolID, ranking.MethodUploadScore, callID, nil)

	msg, err = c.uploadCommonData(nil, p, callID, types.Buffer{1, 2}, 1)
	fakenet.RequireBody(t, msg, err, ranking.ProtocolID, ranking.MethodUploadCommonData, callID, nil)

	msg, err = c.getCommonData(nil, p, callID, 1)
	fakenet.RequireBody(t, msg, err, ranking.ProtocolID, ranking.MethodGetCommonData, callID, []byte{0, 0, 0, 0})

	for _, mode := range []constants.RankingMode{constants.RankingModeRange, constants.RankingModeNear, constants.RankingModeFriendRange, constants.RankingModeFriendNear, constants.RankingModeUser} {
		msg, err = c.getRanking(nil, p, callID, mode, 1, ranking_types.NewRankingOrderParam(), 0, 0)
		fakenet.RequireBody(t, msg, err, ranking.ProtocolID, ranking.MethodGetRanking, callID, emptyResult)
	}

	msg, err = c.getCachedTopXRanking(nil, p, callID, 1, ranking_types.NewRankingOrderParam())
	if err != nil || msg == nil || len(msg.Parameters) < len(emptyResult) {
		t.Fatalf("GetCachedTopXRanking: %v", err)
	}
	if string(msg.Parameters[:len(emptyResult)]) != string(emptyResult) {
		t.Fatalf("cached result does not start with an empty RankingResult: % X", msg.Parameters)
	}

	msg, err = c.getCachedTopXRankings(nil, p, callID, types.List[types.UInt32]{1}, types.List[ranking_types.RankingOrderParam]{ranking_types.NewRankingOrderParam()})
	if err != nil || msg == nil {
		t.Fatalf("GetCachedTopXRankings: %v", err)
	}
}

func TestUseEmptyBoardsKeepsExistingHooks(t *testing.T) {
	uploaded := 0
	c := &CommonProtocol{
		InsertRankingByPIDAndRankingScoreData: func(types.PID, ranking_types.RankingScoreData, types.UInt64) error {
			uploaded++
			return errors.New("real hook")
		},
	}
	c.UseEmptyBoards()

	msg, err := c.uploadScore(nil, fakenet.NewPacket(3), callID, ranking_types.NewRankingScoreData(), 1)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Ranking.Unknown)
	if uploaded != 1 {
		t.Fatalf("the existing hook was replaced (calls=%d)", uploaded)
	}
}
