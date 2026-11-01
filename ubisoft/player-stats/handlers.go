package playerstats

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/player-stats"
	player_stats_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/player-stats/types"
)

// respond builds a success RMC message for the given method
func respond(endpoint nex.EndpointInterface, methodID uint32, callID uint32, body []byte) *nex.RMCMessage {
	message := nex.NewRMCSuccess(endpoint, body)
	message.ProtocolID = ubi.ProtocolID
	message.MethodID = methodID
	message.CallID = callID

	return message
}

func (commonProtocol *CommonProtocol) writeStats(err error, packet nex.PacketInterface, callID uint32, playerStatUpdates types.List[player_stats_types.PlayerStatUpdate]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::WriteStats: invalid request")
	}

	if commonProtocol.WriteStats == nil {
		common_globals.Logger.Warning("Player Stats::WriteStats not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Player Stats::WriteStats not implemented")
	}

	nexErr := commonProtocol.WriteStats(packet, playerStatUpdates)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodWriteStats, callID, nil), nil
}

// encodeReadStatsByPlayersResponse encodes the ReadStatsByPlayers response body
func encodeReadStatsByPlayersResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ReadStatsByPlayersResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) readStatsByPlayers(err error, packet nex.PacketInterface, callID uint32, playerPIDs types.List[types.UInt32], queries types.List[player_stats_types.StatboardQuery]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::ReadStatsByPlayers: invalid request")
	}

	if commonProtocol.ReadStatsByPlayers == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeReadStatsByPlayersResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadStatsByPlayersResponse{})

		return respond(endpoint, ubi.MethodReadStatsByPlayers, callID, body), nil
	}

	r0, nexErr := commonProtocol.ReadStatsByPlayers(packet, playerPIDs, queries)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeReadStatsByPlayersResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadStatsByPlayersResponse{Results: r0})

	return respond(endpoint, ubi.MethodReadStatsByPlayers, callID, body), nil
}

// encodeReadLeaderboardsNearPlayerResponse encodes the ReadLeaderboardsNearPlayer response body
func encodeReadLeaderboardsNearPlayerResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ReadLeaderboardsNearPlayerResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) readLeaderboardsNearPlayer(err error, packet nex.PacketInterface, callID uint32, playerPID types.UInt32, count types.UInt32, queries types.List[player_stats_types.LeaderboardQuery]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::ReadLeaderboardsNearPlayer: invalid request")
	}

	if commonProtocol.ReadLeaderboardsNearPlayer == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeReadLeaderboardsNearPlayerResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsNearPlayerResponse{})

		return respond(endpoint, ubi.MethodReadLeaderboardsNearPlayer, callID, body), nil
	}

	r0, nexErr := commonProtocol.ReadLeaderboardsNearPlayer(packet, playerPID, count, queries)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeReadLeaderboardsNearPlayerResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsNearPlayerResponse{Results: r0})

	return respond(endpoint, ubi.MethodReadLeaderboardsNearPlayer, callID, body), nil
}

// encodeReadLeaderboardsByRankResponse encodes the ReadLeaderboardsByRank response body
func encodeReadLeaderboardsByRankResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ReadLeaderboardsByRankResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) readLeaderboardsByRank(err error, packet nex.PacketInterface, callID uint32, startingRank types.UInt32, count types.UInt32, queries types.List[player_stats_types.LeaderboardQuery]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::ReadLeaderboardsByRank: invalid request")
	}

	if commonProtocol.ReadLeaderboardsByRank == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeReadLeaderboardsByRankResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsByRankResponse{})

		return respond(endpoint, ubi.MethodReadLeaderboardsByRank, callID, body), nil
	}

	r0, nexErr := commonProtocol.ReadLeaderboardsByRank(packet, startingRank, count, queries)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeReadLeaderboardsByRankResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsByRankResponse{Results: r0})

	return respond(endpoint, ubi.MethodReadLeaderboardsByRank, callID, body), nil
}

// encodeReadLeaderboardsByPlayersResponse encodes the ReadLeaderboardsByPlayers response body
func encodeReadLeaderboardsByPlayersResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ReadLeaderboardsByPlayersResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) readLeaderboardsByPlayers(err error, packet nex.PacketInterface, callID uint32, playerPIDs types.List[types.UInt32], queries types.List[player_stats_types.LeaderboardQuery]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::ReadLeaderboardsByPlayers: invalid request")
	}

	if commonProtocol.ReadLeaderboardsByPlayers == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeReadLeaderboardsByPlayersResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsByPlayersResponse{})

		return respond(endpoint, ubi.MethodReadLeaderboardsByPlayers, callID, body), nil
	}

	r0, nexErr := commonProtocol.ReadLeaderboardsByPlayers(packet, playerPIDs, queries)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeReadLeaderboardsByPlayersResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsByPlayersResponse{Results: r0})

	return respond(endpoint, ubi.MethodReadLeaderboardsByPlayers, callID, body), nil
}

// encodeReadStatboardHistoryResponse encodes the ReadStatboardHistory response body
func encodeReadStatboardHistoryResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ReadStatboardHistoryResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) readStatboardHistory(err error, packet nex.PacketInterface, callID uint32, queries types.List[player_stats_types.StatboardHistoryQuery]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::ReadStatboardHistory: invalid request")
	}

	if commonProtocol.ReadStatboardHistory == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeReadStatboardHistoryResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadStatboardHistoryResponse{})

		return respond(endpoint, ubi.MethodReadStatboardHistory, callID, body), nil
	}

	r0, nexErr := commonProtocol.ReadStatboardHistory(packet, queries)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeReadStatboardHistoryResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadStatboardHistoryResponse{Results: r0})

	return respond(endpoint, ubi.MethodReadStatboardHistory, callID, body), nil
}

// encodeReadLeaderboardHistoryResponse encodes the ReadLeaderboardHistory response body
func encodeReadLeaderboardHistoryResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ReadLeaderboardHistoryResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) readLeaderboardHistory(err error, packet nex.PacketInterface, callID uint32, queries types.List[player_stats_types.LeaderboardHistoryQuery]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::ReadLeaderboardHistory: invalid request")
	}

	if commonProtocol.ReadLeaderboardHistory == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeReadLeaderboardHistoryResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardHistoryResponse{})

		return respond(endpoint, ubi.MethodReadLeaderboardHistory, callID, body), nil
	}

	r0, nexErr := commonProtocol.ReadLeaderboardHistory(packet, queries)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeReadLeaderboardHistoryResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardHistoryResponse{Results: r0})

	return respond(endpoint, ubi.MethodReadLeaderboardHistory, callID, body), nil
}

// encodeReadStatboardHistoryAggregatedResponse encodes the ReadStatboardHistoryAggregated response body
func encodeReadStatboardHistoryAggregatedResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ReadStatboardHistoryAggregatedResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) readStatboardHistoryAggregated(err error, packet nex.PacketInterface, callID uint32, queries types.List[player_stats_types.StatboardHistoryAggregatedQuery]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::ReadStatboardHistoryAggregated: invalid request")
	}

	if commonProtocol.ReadStatboardHistoryAggregated == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeReadStatboardHistoryAggregatedResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadStatboardHistoryAggregatedResponse{})

		return respond(endpoint, ubi.MethodReadStatboardHistoryAggregated, callID, body), nil
	}

	r0, nexErr := commonProtocol.ReadStatboardHistoryAggregated(packet, queries)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeReadStatboardHistoryAggregatedResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadStatboardHistoryAggregatedResponse{Results: r0})

	return respond(endpoint, ubi.MethodReadStatboardHistoryAggregated, callID, body), nil
}

// encodeGetStatboardNextPurgeDateResponse encodes the GetStatboardNextPurgeDate response body
func encodeGetStatboardNextPurgeDateResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetStatboardNextPurgeDateResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getStatboardNextPurgeDate(err error, packet nex.PacketInterface, callID uint32, boardID types.UInt32, resetFrequency types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::GetStatboardNextPurgeDate: invalid request")
	}

	if commonProtocol.GetStatboardNextPurgeDate == nil {
		common_globals.Logger.Warning("Player Stats::GetStatboardNextPurgeDate not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Player Stats::GetStatboardNextPurgeDate not implemented")
	}

	r0, nexErr := commonProtocol.GetStatboardNextPurgeDate(packet, boardID, resetFrequency)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetStatboardNextPurgeDateResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetStatboardNextPurgeDateResponse{PurgeDate: r0})

	return respond(endpoint, ubi.MethodGetStatboardNextPurgeDate, callID, body), nil
}

// encodeReadLeaderboardsNearPlayer2Response encodes the ReadLeaderboardsNearPlayer2 response body
func encodeReadLeaderboardsNearPlayer2Response(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ReadLeaderboardsNearPlayer2Response) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) readLeaderboardsNearPlayer2(err error, packet nex.PacketInterface, callID uint32, playerPID types.UInt32, count types.UInt32, queries types.List[player_stats_types.LeaderboardQuery2]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::ReadLeaderboardsNearPlayer2: invalid request")
	}

	if commonProtocol.ReadLeaderboardsNearPlayer2 == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeReadLeaderboardsNearPlayer2Response(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsNearPlayer2Response{})

		return respond(endpoint, ubi.MethodReadLeaderboardsNearPlayer2, callID, body), nil
	}

	r0, nexErr := commonProtocol.ReadLeaderboardsNearPlayer2(packet, playerPID, count, queries)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeReadLeaderboardsNearPlayer2Response(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsNearPlayer2Response{Results: r0})

	return respond(endpoint, ubi.MethodReadLeaderboardsNearPlayer2, callID, body), nil
}

// encodeReadLeaderboardsByRank2Response encodes the ReadLeaderboardsByRank2 response body
func encodeReadLeaderboardsByRank2Response(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ReadLeaderboardsByRank2Response) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) readLeaderboardsByRank2(err error, packet nex.PacketInterface, callID uint32, startingRank types.UInt32, count types.UInt32, queries types.List[player_stats_types.LeaderboardQuery2]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::ReadLeaderboardsByRank2: invalid request")
	}

	if commonProtocol.ReadLeaderboardsByRank2 == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeReadLeaderboardsByRank2Response(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsByRank2Response{})

		return respond(endpoint, ubi.MethodReadLeaderboardsByRank2, callID, body), nil
	}

	r0, nexErr := commonProtocol.ReadLeaderboardsByRank2(packet, startingRank, count, queries)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeReadLeaderboardsByRank2Response(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsByRank2Response{Results: r0})

	return respond(endpoint, ubi.MethodReadLeaderboardsByRank2, callID, body), nil
}

// encodeReadLeaderboardsByPlayers2Response encodes the ReadLeaderboardsByPlayers2 response body
func encodeReadLeaderboardsByPlayers2Response(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ReadLeaderboardsByPlayers2Response) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) readLeaderboardsByPlayers2(err error, packet nex.PacketInterface, callID uint32, queries types.List[player_stats_types.LeaderboardQuery2]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::ReadLeaderboardsByPlayers2: invalid request")
	}

	if commonProtocol.ReadLeaderboardsByPlayers2 == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeReadLeaderboardsByPlayers2Response(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsByPlayers2Response{})

		return respond(endpoint, ubi.MethodReadLeaderboardsByPlayers2, callID, body), nil
	}

	r0, nexErr := commonProtocol.ReadLeaderboardsByPlayers2(packet, queries)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeReadLeaderboardsByPlayers2Response(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadLeaderboardsByPlayers2Response{Results: r0})

	return respond(endpoint, ubi.MethodReadLeaderboardsByPlayers2, callID, body), nil
}

// encodeReadPopulationStatsResponse encodes the ReadPopulationStats response body
func encodeReadPopulationStatsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ReadPopulationStatsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) readPopulationStats(err error, packet nex.PacketInterface, callID uint32, queries types.List[player_stats_types.PopulationStatQuery]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Player Stats::ReadPopulationStats: invalid request")
	}

	if commonProtocol.ReadPopulationStats == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeReadPopulationStatsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadPopulationStatsResponse{})

		return respond(endpoint, ubi.MethodReadPopulationStats, callID, body), nil
	}

	r0, nexErr := commonProtocol.ReadPopulationStats(packet, queries)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeReadPopulationStatsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ReadPopulationStatsResponse{Results: r0})

	return respond(endpoint, ubi.MethodReadPopulationStats, callID, body), nil
}
