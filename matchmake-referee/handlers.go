package matchmakereferee

import (
	"errors"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmcutil"
	matchmake_referee "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-referee"
	matchmake_referee_types "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-referee/types"
)

func notImplemented(method string) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Warningf("MatchmakeReferee::%s missing hook!", method)
	return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "MatchmakeReferee::"+method+" not implemented")
}

func invalidArgument(err error) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Error(err.Error())
	return nil, nex.NewError(nex.ResultCodes.MatchmakeReferee.InvalidArgument, err.Error())
}

func hookFailed(method string, err error) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Errorf("MatchmakeReferee::%s hook failed: %s", method, err.Error())
	return nil, rmcutil.HookError(err, nex.ResultCodes.MatchmakeReferee.Unknown)
}

func reply(packet nex.PacketInterface, methodID uint32, callID uint32, write func(*nex.ByteStreamOut)) *nex.RMCMessage {
	return rmcutil.Success(packet.Sender().Endpoint(), matchmake_referee.ProtocolID, methodID, callID, write)
}

func (c *CommonProtocol) startRound(err error, packet nex.PacketInterface, callID uint32, param matchmake_referee_types.MatchmakeRefereeStartRoundParam) (*nex.RMCMessage, *nex.Error) {
	if c.StartRound == nil {
		return notImplemented("StartRound")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if len(param.PIDs) == 0 {
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "MatchmakeReferee::StartRound with no pids")
	}
	roundID, err := c.StartRound(packet.Sender().PID(), param)
	if err != nil {
		return hookFailed("StartRound", err)
	}

	return reply(packet, matchmake_referee.MethodStartRound, callID, func(s *nex.ByteStreamOut) { roundID.WriteTo(s) }), nil
}

func (c *CommonProtocol) getStartRoundParam(err error, packet nex.PacketInterface, callID uint32, roundID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if c.GetStartRoundParam == nil {
		return notImplemented("GetStartRoundParam")
	}
	if err != nil {
		return invalidArgument(err)
	}
	param, err := c.GetStartRoundParam(packet.Sender().PID(), roundID)
	if err != nil {
		return hookFailed("GetStartRoundParam", err)
	}

	return reply(packet, matchmake_referee.MethodGetStartRoundParam, callID, func(s *nex.ByteStreamOut) { param.WriteTo(s) }), nil
}

func (c *CommonProtocol) endRound(err error, packet nex.PacketInterface, callID uint32, endRoundParam matchmake_referee_types.MatchmakeRefereeEndRoundParam) (*nex.RMCMessage, *nex.Error) {
	if c.EndRound == nil {
		return notImplemented("EndRound")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := c.EndRound(packet.Sender().PID(), endRoundParam); err != nil {
		return hookFailed("EndRound", err)
	}

	return reply(packet, matchmake_referee.MethodEndRound, callID, nil), nil
}

func (c *CommonProtocol) endRoundWithoutReport(err error, packet nex.PacketInterface, callID uint32, roundID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if c.EndRoundWithoutReport == nil {
		return notImplemented("EndRoundWithoutReport")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := c.EndRoundWithoutReport(packet.Sender().PID(), roundID); err != nil {
		return hookFailed("EndRoundWithoutReport", err)
	}

	return reply(packet, matchmake_referee.MethodEndRoundWithoutReport, callID, nil), nil
}

func (c *CommonProtocol) getRoundParticipants(err error, packet nex.PacketInterface, callID uint32, roundID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if c.GetRoundParticipants == nil {
		return notImplemented("GetRoundParticipants")
	}
	if err != nil {
		return invalidArgument(err)
	}
	pids, err := c.GetRoundParticipants(packet.Sender().PID(), roundID)
	if err != nil {
		return hookFailed("GetRoundParticipants", err)
	}

	return reply(packet, matchmake_referee.MethodGetRoundParticipants, callID, func(s *nex.ByteStreamOut) { pids.WriteTo(s) }), nil
}

func (c *CommonProtocol) getNotSummarizedRound(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if c.GetNotSummarizedRound == nil {
		return notImplemented("GetNotSummarizedRound")
	}
	if err != nil {
		return invalidArgument(err)
	}
	rounds, err := c.GetNotSummarizedRound(packet.Sender().PID())
	if err != nil {
		return hookFailed("GetNotSummarizedRound", err)
	}

	return reply(packet, matchmake_referee.MethodGetNotSummarizedRound, callID, func(s *nex.ByteStreamOut) { rounds.WriteTo(s) }), nil
}

func (c *CommonProtocol) getRound(err error, packet nex.PacketInterface, callID uint32, roundID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if c.GetRound == nil {
		return notImplemented("GetRound")
	}
	if err != nil {
		return invalidArgument(err)
	}
	round, err := c.GetRound(packet.Sender().PID(), roundID)
	if err != nil {
		return hookFailed("GetRound", err)
	}

	return reply(packet, matchmake_referee.MethodGetRound, callID, func(s *nex.ByteStreamOut) { round.WriteTo(s) }), nil
}

func (c *CommonProtocol) getStatsPrimary(err error, packet nex.PacketInterface, callID uint32, target matchmake_referee_types.MatchmakeRefereeStatsTarget) (*nex.RMCMessage, *nex.Error) {
	if c.GetStatsPrimary == nil {
		return notImplemented("GetStatsPrimary")
	}
	if err != nil {
		return invalidArgument(err)
	}
	stats, err := c.GetStatsPrimary(packet.Sender().PID(), target)
	if err != nil {
		return hookFailed("GetStatsPrimary", err)
	}

	return reply(packet, matchmake_referee.MethodGetStatsPrimary, callID, func(s *nex.ByteStreamOut) { stats.WriteTo(s) }), nil
}

func (c *CommonProtocol) getStatsPrimaries(err error, packet nex.PacketInterface, callID uint32, targets types.List[matchmake_referee_types.MatchmakeRefereeStatsTarget]) (*nex.RMCMessage, *nex.Error) {
	if c.GetStatsPrimaries == nil && c.GetStatsPrimary == nil {
		return notImplemented("GetStatsPrimaries")
	}
	if err != nil {
		return invalidArgument(err)
	}

	pid := packet.Sender().PID()

	var stats types.List[matchmake_referee_types.MatchmakeRefereeStats]
	var results types.List[types.QResult]

	if c.GetStatsPrimaries != nil {
		stats, results, err = c.GetStatsPrimaries(pid, targets)
		if err != nil {
			return hookFailed("GetStatsPrimaries", err)
		}
	} else {
		stats = types.NewList[matchmake_referee_types.MatchmakeRefereeStats]()
		results = types.NewList[types.QResult]()

		for _, target := range targets {
			entry, err := c.GetStatsPrimary(pid, target)

			var nexErr *nex.Error
			switch {
			case err == nil:
				results = append(results, types.NewQResultSuccess(nex.ResultCodes.Core.Unknown))
			case errors.As(err, &nexErr):
				entry = matchmake_referee_types.NewMatchmakeRefereeStats()
				results = append(results, types.NewQResultError(nexErr.ResultCode))
			default:
				return hookFailed("GetStatsPrimaries", err)
			}

			stats = append(stats, entry)
		}
	}

	return reply(packet, matchmake_referee.MethodGetStatsPrimaries, callID, func(s *nex.ByteStreamOut) {
		stats.WriteTo(s)
		results.WriteTo(s)
	}), nil
}

func (c *CommonProtocol) getStatsAll(err error, packet nex.PacketInterface, callID uint32, target matchmake_referee_types.MatchmakeRefereeStatsTarget) (*nex.RMCMessage, *nex.Error) {
	if c.GetStatsAll == nil {
		return notImplemented("GetStatsAll")
	}
	if err != nil {
		return invalidArgument(err)
	}
	stats, err := c.GetStatsAll(packet.Sender().PID(), target)
	if err != nil {
		return hookFailed("GetStatsAll", err)
	}

	return reply(packet, matchmake_referee.MethodGetStatsAll, callID, func(s *nex.ByteStreamOut) { stats.WriteTo(s) }), nil
}

func (c *CommonProtocol) createStats(err error, packet nex.PacketInterface, callID uint32, param matchmake_referee_types.MatchmakeRefereeStatsInitParam) (*nex.RMCMessage, *nex.Error) {
	if c.CreateStats == nil {
		return notImplemented("CreateStats")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := c.CreateStats(packet.Sender().PID(), param); err != nil {
		return hookFailed("CreateStats", err)
	}

	return reply(packet, matchmake_referee.MethodCreateStats, callID, nil), nil
}

func (c *CommonProtocol) getOrCreateStats(err error, packet nex.PacketInterface, callID uint32, param matchmake_referee_types.MatchmakeRefereeStatsInitParam) (*nex.RMCMessage, *nex.Error) {
	if c.GetOrCreateStats == nil {
		return notImplemented("GetOrCreateStats")
	}
	if err != nil {
		return invalidArgument(err)
	}
	stats, err := c.GetOrCreateStats(packet.Sender().PID(), param)
	if err != nil {
		return hookFailed("GetOrCreateStats", err)
	}

	return reply(packet, matchmake_referee.MethodGetOrCreateStats, callID, func(s *nex.ByteStreamOut) { stats.WriteTo(s) }), nil
}

func (c *CommonProtocol) resetStats(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if c.ResetStats == nil {
		return notImplemented("ResetStats")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := c.ResetStats(packet.Sender().PID()); err != nil {
		return hookFailed("ResetStats", err)
	}

	return reply(packet, matchmake_referee.MethodResetStats, callID, nil), nil
}
