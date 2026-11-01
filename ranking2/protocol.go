// Package ranking2 implements a hook-driven common handler set for the
// Ranking2 protocol (0x7A). The game server supplies storage and ranking logic
// through the function fields on CommonProtocol; this package only decodes the
// requests, calls the hooks and encodes the responses.
//
// A method whose hook is nil answers Core::NotImplemented.
package ranking2

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking2 "github.com/PretendoNetwork/nex-protocols-go/v2/ranking2"
	"github.com/PretendoNetwork/nex-protocols-go/v2/ranking2/constants"
	ranking2_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking2/types"
)

// CommonProtocol holds the hooks for the Ranking2 protocol. Hooks that return
// an error answer Ranking2::Unknown, unless the error is a *nex.Error, which
// is sent as is.
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ranking2.Interface

	// PutScore stores the scores uploaded by pid.
	PutScore func(pid types.PID, scoreDataList types.List[ranking2_types.Ranking2ScoreData], nexUniqueID types.UInt64) error
	// GetCommonData returns the common data of a user. pid is the caller.
	GetCommonData func(pid types.PID, optionFlags constants.Ranking2GetOptionFlags, principalID types.PID, nexUniqueID types.UInt64) (ranking2_types.Ranking2CommonData, error)
	// PutCommonData stores the common data of the caller.
	PutCommonData func(pid types.PID, commonData ranking2_types.Ranking2CommonData, nexUniqueID types.UInt64) error
	// DeleteCommonData deletes the common data of the caller.
	DeleteCommonData func(pid types.PID, nexUniqueID types.UInt64) error
	// GetRanking returns a page of a ranking.
	GetRanking func(pid types.PID, getParam ranking2_types.Ranking2GetParam) (ranking2_types.Ranking2Info, error)
	// GetRankingByPrincipalID returns the ranking entries of the listed users.
	GetRankingByPrincipalID func(pid types.PID, getParam ranking2_types.Ranking2GetByListParam, principalIDList types.List[types.PID]) (ranking2_types.Ranking2Info, error)
	// GetCategorySetting returns the settings of a category.
	GetCategorySetting func(category types.UInt32) (ranking2_types.Ranking2CategorySetting, error)
	// GetRankingChart returns a score distribution chart.
	GetRankingChart func(pid types.PID, info ranking2_types.Ranking2ChartInfoInput) (ranking2_types.Ranking2ChartInfo, error)
	// GetRankingCharts returns several charts. When nil but GetRankingChart is
	// set, GetRankingChart is called once per requested chart.
	GetRankingCharts func(pid types.PID, infoArray types.List[ranking2_types.Ranking2ChartInfoInput]) (types.List[ranking2_types.Ranking2ChartInfo], error)
	// GetEstimateScoreRank estimates the rank a score would get.
	GetEstimateScoreRank func(pid types.PID, input ranking2_types.Ranking2EstimateScoreRankInput) (ranking2_types.Ranking2EstimateScoreRankOutput, error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ranking2.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerPutScore(commonProtocol.putScore)
	protocol.SetHandlerGetCommonData(commonProtocol.getCommonData)
	protocol.SetHandlerPutCommonData(commonProtocol.putCommonData)
	protocol.SetHandlerDeleteCommonData(commonProtocol.deleteCommonData)
	protocol.SetHandlerGetRanking(commonProtocol.getRanking)
	protocol.SetHandlerGetRankingByPrincipalID(commonProtocol.getRankingByPrincipalID)
	protocol.SetHandlerGetCategorySetting(commonProtocol.getCategorySetting)
	protocol.SetHandlerGetRankingChart(commonProtocol.getRankingChart)
	protocol.SetHandlerGetRankingCharts(commonProtocol.getRankingCharts)
	protocol.SetHandlerGetEstimateScoreRank(commonProtocol.getEstimateScoreRank)

	return commonProtocol
}
