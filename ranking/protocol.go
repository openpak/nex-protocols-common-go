package ranking

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking "github.com/PretendoNetwork/nex-protocols-go/v2/ranking"
	"github.com/PretendoNetwork/nex-protocols-go/v2/ranking/constants"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
)

type CommonProtocol struct {
	endpoint                                                       nex.EndpointInterface
	protocol                                                       ranking.Interface
	emptyBoards                                                    bool
	GetCommonData                                                  func(uniqueID types.UInt64) (types.Buffer, error)
	UploadCommonData                                               func(pid types.PID, uniqueID types.UInt64, commonData types.Buffer) error
	InsertRankingByPIDAndRankingScoreData                          func(pid types.PID, rankingScoreData ranking_types.RankingScoreData, uniqueID types.UInt64) error
	GetRankingsAndCountByCategoryAndRankingOrderParam              func(category types.UInt32, rankingOrderParam ranking_types.RankingOrderParam) (types.List[ranking_types.RankingRankData], uint32, error)
	GetNearbyRankingsAndCountByCategoryAndRankingOrderParam        func(pid types.PID, category types.UInt32, rankingOrderParam ranking_types.RankingOrderParam) (types.List[ranking_types.RankingRankData], uint32, error)
	GetFriendsRankingsAndCountByCategoryAndRankingOrderParam       func(pid types.PID, category types.UInt32, rankingOrderParam ranking_types.RankingOrderParam) (types.List[ranking_types.RankingRankData], uint32, error)
	GetNearbyFriendsRankingsAndCountByCategoryAndRankingOrderParam func(pid types.PID, category types.UInt32, rankingOrderParam ranking_types.RankingOrderParam) (types.List[ranking_types.RankingRankData], uint32, error)
	GetOwnRankingByCategoryAndRankingOrderParam                    func(pid types.PID, category types.UInt32, rankingOrderParam ranking_types.RankingOrderParam) (types.List[ranking_types.RankingRankData], uint32, error)
	OnAfterGetCachedTopXRanking                                    func(packet nex.PacketInterface, category types.UInt32, orderParam ranking_types.RankingOrderParam)
	OnAfterGetCachedTopXRankings                                   func(packet nex.PacketInterface, categories types.List[types.UInt32], orderParams types.List[ranking_types.RankingOrderParam])
	OnAfterGetCommonData                                           func(packet nex.PacketInterface, uniqueID types.UInt64)
	OnAfterGetRanking                                              func(packet nex.PacketInterface, rankingMode constants.RankingMode, category types.UInt32, orderParam ranking_types.RankingOrderParam, uniqueID types.UInt64, principalID types.PID)
	OnAfterUploadCommonData                                        func(packet nex.PacketInterface, commonData types.Buffer, uniqueID types.UInt64)
	OnAfterUploadScore                                             func(packet nex.PacketInterface, scoreData ranking_types.RankingScoreData, uniqueID types.UInt64)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ranking.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerGetCachedTopXRanking(commonProtocol.getCachedTopXRanking)
	protocol.SetHandlerGetCachedTopXRankings(commonProtocol.getCachedTopXRankings)
	protocol.SetHandlerGetCommonData(commonProtocol.getCommonData)
	protocol.SetHandlerGetRanking(commonProtocol.getRanking)
	protocol.SetHandlerUploadCommonData(commonProtocol.uploadCommonData)
	protocol.SetHandlerUploadScore(commonProtocol.uploadScore)

	return commonProtocol
}

// UseEmptyBoards is a stopgap for titles that register the Ranking protocol
// with no storage behind it. It installs hooks that accept score and common
// data uploads and throw them away, and that answer every read with an empty
// board: no entries, a total count of zero and empty common data. Instead of
// Ranking::NotFound, the ranking getters then send an empty result.
//
// Only hooks that are still nil are filled, so call it after setting the hooks
// the game server really implements. Nothing is stored: scores uploaded
// through the installed hooks are lost, and a title that needs real rankings
// must supply its own hooks instead.
func (commonProtocol *CommonProtocol) UseEmptyBoards() {
	commonProtocol.emptyBoards = true

	emptyList := func() types.List[ranking_types.RankingRankData] {
		return types.NewList[ranking_types.RankingRankData]()
	}

	if commonProtocol.GetCommonData == nil {
		commonProtocol.GetCommonData = func(types.UInt64) (types.Buffer, error) {
			return types.Buffer{}, nil
		}
	}

	if commonProtocol.UploadCommonData == nil {
		commonProtocol.UploadCommonData = func(types.PID, types.UInt64, types.Buffer) error { return nil }
	}

	if commonProtocol.InsertRankingByPIDAndRankingScoreData == nil {
		commonProtocol.InsertRankingByPIDAndRankingScoreData = func(types.PID, ranking_types.RankingScoreData, types.UInt64) error {
			return nil
		}
	}

	if commonProtocol.GetRankingsAndCountByCategoryAndRankingOrderParam == nil {
		commonProtocol.GetRankingsAndCountByCategoryAndRankingOrderParam = func(types.UInt32, ranking_types.RankingOrderParam) (types.List[ranking_types.RankingRankData], uint32, error) {
			return emptyList(), 0, nil
		}
	}

	byPID := func(types.PID, types.UInt32, ranking_types.RankingOrderParam) (types.List[ranking_types.RankingRankData], uint32, error) {
		return emptyList(), 0, nil
	}

	if commonProtocol.GetNearbyRankingsAndCountByCategoryAndRankingOrderParam == nil {
		commonProtocol.GetNearbyRankingsAndCountByCategoryAndRankingOrderParam = byPID
	}

	if commonProtocol.GetFriendsRankingsAndCountByCategoryAndRankingOrderParam == nil {
		commonProtocol.GetFriendsRankingsAndCountByCategoryAndRankingOrderParam = byPID
	}

	if commonProtocol.GetNearbyFriendsRankingsAndCountByCategoryAndRankingOrderParam == nil {
		commonProtocol.GetNearbyFriendsRankingsAndCountByCategoryAndRankingOrderParam = byPID
	}

	if commonProtocol.GetOwnRankingByCategoryAndRankingOrderParam == nil {
		commonProtocol.GetOwnRankingByCategoryAndRankingOrderParam = byPID
	}
}
