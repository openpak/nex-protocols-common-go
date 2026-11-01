// Package playerstats implements the common Player Stats protocol handlers of the Quazal Rendez-Vous (Ubisoft) family.
//
// Layouts follow the NintendoClients wiki page "Player Stats Protocol" (MIT). Every method calls a user-provided hook;
// a method whose hook is nil answers Core::NotImplemented unless a safe default is documented on its hook field.
//
// Layout provisional:
//   - ReadLeaderboardsByPlayers2 takes only the query list on the wiki; the player list is presumably carried in LeaderboardQuery2.m_estimatedPIDs.
package playerstats

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/player-stats"
	player_stats_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/player-stats/types"
)

// CommonProtocol handles the Player Stats protocol with user-provided hooks
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi.Interface

	// WriteStats is required: when nil the request is answered with Core::NotImplemented.
	WriteStats func(packet nex.PacketInterface, playerStatUpdates types.List[player_stats_types.PlayerStatUpdate]) *nex.Error

	// ReadStatsByPlayers is optional: when nil an empty result (empty lists, zero counters) is returned.
	ReadStatsByPlayers func(packet nex.PacketInterface, playerPIDs types.List[types.UInt32], queries types.List[player_stats_types.StatboardQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error)

	// ReadLeaderboardsNearPlayer is optional: when nil an empty result (empty lists, zero counters) is returned.
	ReadLeaderboardsNearPlayer func(packet nex.PacketInterface, playerPID types.UInt32, count types.UInt32, queries types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error)

	// ReadLeaderboardsByRank is optional: when nil an empty result (empty lists, zero counters) is returned.
	ReadLeaderboardsByRank func(packet nex.PacketInterface, startingRank types.UInt32, count types.UInt32, queries types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error)

	// ReadLeaderboardsByPlayers is optional: when nil an empty result (empty lists, zero counters) is returned.
	ReadLeaderboardsByPlayers func(packet nex.PacketInterface, playerPIDs types.List[types.UInt32], queries types.List[player_stats_types.LeaderboardQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error)

	// ReadStatboardHistory is optional: when nil an empty result (empty lists, zero counters) is returned.
	ReadStatboardHistory func(packet nex.PacketInterface, queries types.List[player_stats_types.StatboardHistoryQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error)

	// ReadLeaderboardHistory is optional: when nil an empty result (empty lists, zero counters) is returned.
	ReadLeaderboardHistory func(packet nex.PacketInterface, queries types.List[player_stats_types.LeaderboardHistoryQuery]) (types.List[player_stats_types.LeaderboardResult], *nex.Error)

	// ReadStatboardHistoryAggregated is optional: when nil an empty result (empty lists, zero counters) is returned.
	ReadStatboardHistoryAggregated func(packet nex.PacketInterface, queries types.List[player_stats_types.StatboardHistoryAggregatedQuery]) (types.List[player_stats_types.StatboardResult], *nex.Error)

	// GetStatboardNextPurgeDate is required: when nil the request is answered with Core::NotImplemented.
	GetStatboardNextPurgeDate func(packet nex.PacketInterface, boardID types.UInt32, resetFrequency types.UInt32) (types.DateTime, *nex.Error)

	// ReadLeaderboardsNearPlayer2 is optional: when nil an empty result (empty lists, zero counters) is returned.
	ReadLeaderboardsNearPlayer2 func(packet nex.PacketInterface, playerPID types.UInt32, count types.UInt32, queries types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error)

	// ReadLeaderboardsByRank2 is optional: when nil an empty result (empty lists, zero counters) is returned.
	ReadLeaderboardsByRank2 func(packet nex.PacketInterface, startingRank types.UInt32, count types.UInt32, queries types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error)

	// ReadLeaderboardsByPlayers2 is optional: when nil an empty result (empty lists, zero counters) is returned.
	ReadLeaderboardsByPlayers2 func(packet nex.PacketInterface, queries types.List[player_stats_types.LeaderboardQuery2]) (types.List[player_stats_types.LeaderboardResult], *nex.Error)

	// ReadPopulationStats is optional: when nil an empty result (empty lists, zero counters) is returned.
	ReadPopulationStats func(packet nex.PacketInterface, queries types.List[player_stats_types.PopulationStatQuery]) (types.List[player_stats_types.PopulationStatResult], *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerWriteStats(commonProtocol.writeStats)
	protocol.SetHandlerReadStatsByPlayers(commonProtocol.readStatsByPlayers)
	protocol.SetHandlerReadLeaderboardsNearPlayer(commonProtocol.readLeaderboardsNearPlayer)
	protocol.SetHandlerReadLeaderboardsByRank(commonProtocol.readLeaderboardsByRank)
	protocol.SetHandlerReadLeaderboardsByPlayers(commonProtocol.readLeaderboardsByPlayers)
	protocol.SetHandlerReadStatboardHistory(commonProtocol.readStatboardHistory)
	protocol.SetHandlerReadLeaderboardHistory(commonProtocol.readLeaderboardHistory)
	protocol.SetHandlerReadStatboardHistoryAggregated(commonProtocol.readStatboardHistoryAggregated)
	protocol.SetHandlerGetStatboardNextPurgeDate(commonProtocol.getStatboardNextPurgeDate)
	protocol.SetHandlerReadLeaderboardsNearPlayer2(commonProtocol.readLeaderboardsNearPlayer2)
	protocol.SetHandlerReadLeaderboardsByRank2(commonProtocol.readLeaderboardsByRank2)
	protocol.SetHandlerReadLeaderboardsByPlayers2(commonProtocol.readLeaderboardsByPlayers2)
	protocol.SetHandlerReadPopulationStats(commonProtocol.readPopulationStats)

	return commonProtocol
}
