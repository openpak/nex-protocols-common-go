// Package matchmakereferee implements a hook-driven common handler set for the
// MatchmakeReferee protocol (0x78, the standard variant; the libeagle titles use
// a customised version that this package does not cover). The game server
// supplies round and stats storage through the function fields on
// CommonProtocol.
//
// A method whose hook is nil answers Core::NotImplemented. Hooks that return an
// error answer MatchmakeReferee::Unknown unless the error is a *nex.Error, which
// is sent as is (use it for RoundNotFound, StatsNotFound, NotParticipatedGathering
// and so on). In every hook pid is the caller.
package matchmakereferee

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	matchmake_referee "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-referee"
	matchmake_referee_types "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-referee/types"
)

// CommonProtocol holds the hooks for the MatchmakeReferee protocol.
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol matchmake_referee.Interface

	// StartRound opens a round and returns its id. The handler already rejects
	// an empty pid list with Core::InvalidArgument; checking that the pids
	// belong to the gathering (MatchmakeReferee::NotParticipatedGathering) is up to the hook.
	StartRound func(pid types.PID, param matchmake_referee_types.MatchmakeRefereeStartRoundParam) (types.UInt64, error)
	// GetStartRoundParam returns the parameters a round was started with.
	GetStartRoundParam func(pid types.PID, roundID types.UInt64) (matchmake_referee_types.MatchmakeRefereeStartRoundParam, error)
	// EndRound closes a round with the reported results.
	EndRound func(pid types.PID, param matchmake_referee_types.MatchmakeRefereeEndRoundParam) error
	// EndRoundWithoutReport closes a round without results.
	EndRoundWithoutReport func(pid types.PID, roundID types.UInt64) error
	// GetRoundParticipants returns the participant ids of a round.
	GetRoundParticipants func(pid types.PID, roundID types.UInt64) (types.List[types.UInt64], error)
	// GetNotSummarizedRound returns the rounds that were not summarized yet.
	GetNotSummarizedRound func(pid types.PID) (types.List[matchmake_referee_types.MatchmakeRefereeRound], error)
	// GetRound returns a round.
	GetRound func(pid types.PID, roundID types.UInt64) (matchmake_referee_types.MatchmakeRefereeRound, error)
	// GetStatsPrimary returns the stats of one user in one category.
	GetStatsPrimary func(pid types.PID, target matchmake_referee_types.MatchmakeRefereeStatsTarget) (matchmake_referee_types.MatchmakeRefereeStats, error)
	// GetStatsPrimaries returns the stats of several targets, with one result
	// code per target. When nil but GetStatsPrimary is set, GetStatsPrimary is
	// called per target; a *nex.Error from it becomes that target's result code
	// and an empty stats entry.
	GetStatsPrimaries func(pid types.PID, targets types.List[matchmake_referee_types.MatchmakeRefereeStatsTarget]) (types.List[matchmake_referee_types.MatchmakeRefereeStats], types.List[types.QResult], error)
	// GetStatsAll returns the stats of a user in every category.
	GetStatsAll func(pid types.PID, target matchmake_referee_types.MatchmakeRefereeStatsTarget) (types.List[matchmake_referee_types.MatchmakeRefereeStats], error)
	// CreateStats creates the caller's stats in a category.
	CreateStats func(pid types.PID, param matchmake_referee_types.MatchmakeRefereeStatsInitParam) error
	// GetOrCreateStats returns the caller's stats in a category, creating them first if missing.
	GetOrCreateStats func(pid types.PID, param matchmake_referee_types.MatchmakeRefereeStatsInitParam) (matchmake_referee_types.MatchmakeRefereeStats, error)
	// ResetStats resets the caller's stats.
	ResetStats func(pid types.PID) error
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol matchmake_referee.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerStartRound(commonProtocol.startRound)
	protocol.SetHandlerGetStartRoundParam(commonProtocol.getStartRoundParam)
	protocol.SetHandlerEndRound(commonProtocol.endRound)
	protocol.SetHandlerEndRoundWithoutReport(commonProtocol.endRoundWithoutReport)
	protocol.SetHandlerGetRoundParticipants(commonProtocol.getRoundParticipants)
	protocol.SetHandlerGetNotSummarizedRound(commonProtocol.getNotSummarizedRound)
	protocol.SetHandlerGetRound(commonProtocol.getRound)
	protocol.SetHandlerGetStatsPrimary(commonProtocol.getStatsPrimary)
	protocol.SetHandlerGetStatsPrimaries(commonProtocol.getStatsPrimaries)
	protocol.SetHandlerGetStatsAll(commonProtocol.getStatsAll)
	protocol.SetHandlerCreateStats(commonProtocol.createStats)
	protocol.SetHandlerGetOrCreateStats(commonProtocol.getOrCreateStats)
	protocol.SetHandlerResetStats(commonProtocol.resetStats)

	return commonProtocol
}
