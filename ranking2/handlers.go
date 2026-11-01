package ranking2

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmcutil"
	ranking2 "github.com/PretendoNetwork/nex-protocols-go/v2/ranking2"
	"github.com/PretendoNetwork/nex-protocols-go/v2/ranking2/constants"
	ranking2_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking2/types"
)

func notImplemented(method string) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Warningf("Ranking2::%s missing hook!", method)
	return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Ranking2::"+method+" not implemented")
}

func invalidArgument(err error) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Error(err.Error())
	return nil, nex.NewError(nex.ResultCodes.Ranking2.InvalidArgument, err.Error())
}

func hookFailed(method string, err error) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Errorf("Ranking2::%s hook failed: %s", method, err.Error())
	return nil, rmcutil.HookError(err, nex.ResultCodes.Ranking2.Unknown)
}

func reply(packet nex.PacketInterface, methodID uint32, callID uint32, write func(*nex.ByteStreamOut)) *nex.RMCMessage {
	return rmcutil.Success(packet.Sender().Endpoint(), ranking2.ProtocolID, methodID, callID, write)
}

func (commonProtocol *CommonProtocol) putScore(err error, packet nex.PacketInterface, callID uint32, scoreDataList types.List[ranking2_types.Ranking2ScoreData], nexUniqueID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.PutScore == nil {
		return notImplemented("PutScore")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := commonProtocol.PutScore(packet.Sender().PID(), scoreDataList, nexUniqueID); err != nil {
		return hookFailed("PutScore", err)
	}

	return reply(packet, ranking2.MethodPutScore, callID, nil), nil
}

func (commonProtocol *CommonProtocol) getCommonData(err error, packet nex.PacketInterface, callID uint32, optionFlags constants.Ranking2GetOptionFlags, principalID types.PID, nexUniqueID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetCommonData == nil {
		return notImplemented("GetCommonData")
	}
	if err != nil {
		return invalidArgument(err)
	}
	commonData, err := commonProtocol.GetCommonData(packet.Sender().PID(), optionFlags, principalID, nexUniqueID)
	if err != nil {
		return hookFailed("GetCommonData", err)
	}

	return reply(packet, ranking2.MethodGetCommonData, callID, func(s *nex.ByteStreamOut) { commonData.WriteTo(s) }), nil
}

func (commonProtocol *CommonProtocol) putCommonData(err error, packet nex.PacketInterface, callID uint32, commonData ranking2_types.Ranking2CommonData, nexUniqueID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.PutCommonData == nil {
		return notImplemented("PutCommonData")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := commonProtocol.PutCommonData(packet.Sender().PID(), commonData, nexUniqueID); err != nil {
		return hookFailed("PutCommonData", err)
	}

	return reply(packet, ranking2.MethodPutCommonData, callID, nil), nil
}

func (commonProtocol *CommonProtocol) deleteCommonData(err error, packet nex.PacketInterface, callID uint32, nexUniqueID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.DeleteCommonData == nil {
		return notImplemented("DeleteCommonData")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := commonProtocol.DeleteCommonData(packet.Sender().PID(), nexUniqueID); err != nil {
		return hookFailed("DeleteCommonData", err)
	}

	return reply(packet, ranking2.MethodDeleteCommonData, callID, nil), nil
}

func (commonProtocol *CommonProtocol) getRanking(err error, packet nex.PacketInterface, callID uint32, getParam ranking2_types.Ranking2GetParam) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetRanking == nil {
		return notImplemented("GetRanking")
	}
	if err != nil {
		return invalidArgument(err)
	}
	info, err := commonProtocol.GetRanking(packet.Sender().PID(), getParam)
	if err != nil {
		return hookFailed("GetRanking", err)
	}

	return reply(packet, ranking2.MethodGetRanking, callID, func(s *nex.ByteStreamOut) { info.WriteTo(s) }), nil
}

func (commonProtocol *CommonProtocol) getRankingByPrincipalID(err error, packet nex.PacketInterface, callID uint32, getParam ranking2_types.Ranking2GetByListParam, principalIDList types.List[types.PID]) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetRankingByPrincipalID == nil {
		return notImplemented("GetRankingByPrincipalID")
	}
	if err != nil {
		return invalidArgument(err)
	}
	info, err := commonProtocol.GetRankingByPrincipalID(packet.Sender().PID(), getParam, principalIDList)
	if err != nil {
		return hookFailed("GetRankingByPrincipalID", err)
	}

	return reply(packet, ranking2.MethodGetRankingByPrincipalID, callID, func(s *nex.ByteStreamOut) { info.WriteTo(s) }), nil
}

func (commonProtocol *CommonProtocol) getCategorySetting(err error, packet nex.PacketInterface, callID uint32, category types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetCategorySetting == nil {
		return notImplemented("GetCategorySetting")
	}
	if err != nil {
		return invalidArgument(err)
	}
	setting, err := commonProtocol.GetCategorySetting(category)
	if err != nil {
		return hookFailed("GetCategorySetting", err)
	}

	return reply(packet, ranking2.MethodGetCategorySetting, callID, func(s *nex.ByteStreamOut) { setting.WriteTo(s) }), nil
}

func (commonProtocol *CommonProtocol) getRankingChart(err error, packet nex.PacketInterface, callID uint32, info ranking2_types.Ranking2ChartInfoInput) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetRankingChart == nil {
		return notImplemented("GetRankingChart")
	}
	if err != nil {
		return invalidArgument(err)
	}
	chart, err := commonProtocol.GetRankingChart(packet.Sender().PID(), info)
	if err != nil {
		return hookFailed("GetRankingChart", err)
	}

	return reply(packet, ranking2.MethodGetRankingChart, callID, func(s *nex.ByteStreamOut) { chart.WriteTo(s) }), nil
}

func (commonProtocol *CommonProtocol) getRankingCharts(err error, packet nex.PacketInterface, callID uint32, infoArray types.List[ranking2_types.Ranking2ChartInfoInput]) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetRankingCharts == nil && commonProtocol.GetRankingChart == nil {
		return notImplemented("GetRankingCharts")
	}
	if err != nil {
		return invalidArgument(err)
	}

	pid := packet.Sender().PID()

	var charts types.List[ranking2_types.Ranking2ChartInfo]
	if commonProtocol.GetRankingCharts != nil {
		charts, err = commonProtocol.GetRankingCharts(pid, infoArray)
		if err != nil {
			return hookFailed("GetRankingCharts", err)
		}
	} else {
		charts = types.NewList[ranking2_types.Ranking2ChartInfo]()
		for _, info := range infoArray {
			chart, err := commonProtocol.GetRankingChart(pid, info)
			if err != nil {
				return hookFailed("GetRankingCharts", err)
			}
			charts = append(charts, chart)
		}
	}

	return reply(packet, ranking2.MethodGetRankingCharts, callID, func(s *nex.ByteStreamOut) { charts.WriteTo(s) }), nil
}

func (commonProtocol *CommonProtocol) getEstimateScoreRank(err error, packet nex.PacketInterface, callID uint32, input ranking2_types.Ranking2EstimateScoreRankInput) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetEstimateScoreRank == nil {
		return notImplemented("GetEstimateScoreRank")
	}
	if err != nil {
		return invalidArgument(err)
	}
	output, err := commonProtocol.GetEstimateScoreRank(packet.Sender().PID(), input)
	if err != nil {
		return hookFailed("GetEstimateScoreRank", err)
	}

	return reply(packet, ranking2.MethodGetEstimateScoreRank, callID, func(s *nex.ByteStreamOut) { output.WriteTo(s) }), nil
}
