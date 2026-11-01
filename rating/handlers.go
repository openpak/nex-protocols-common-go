package rating

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmcutil"
	rating "github.com/PretendoNetwork/nex-protocols-go/v2/rating"
	rating_types "github.com/PretendoNetwork/nex-protocols-go/v2/rating/types"
)

func notImplemented(method string) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Warningf("Rating::%s missing hook!", method)
	return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Rating::"+method+" not implemented")
}

func invalidArgument(err error) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Error(err.Error())
	return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, err.Error())
}

func hookFailed(method string, err error) (*nex.RMCMessage, *nex.Error) {
	common_globals.Logger.Errorf("Rating::%s hook failed: %s", method, err.Error())
	return nil, rmcutil.HookError(err, nex.ResultCodes.Core.Unknown)
}

func reply(packet nex.PacketInterface, methodID uint32, callID uint32, write func(*nex.ByteStreamOut)) *nex.RMCMessage {
	return rmcutil.Success(packet.Sender().Endpoint(), rating.ProtocolID, methodID, callID, write)
}

func (commonProtocol *CommonProtocol) unk1(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.Unk1 == nil {
		return notImplemented("Unk1")
	}
	if err != nil {
		return invalidArgument(err)
	}
	number, token, err := commonProtocol.Unk1(packet.Sender().PID())
	if err != nil {
		return hookFailed("Unk1", err)
	}

	return reply(packet, rating.MethodUnk1, callID, func(s *nex.ByteStreamOut) {
		number.WriteTo(s)
		token.WriteTo(s)
	}), nil
}

func (commonProtocol *CommonProtocol) unk2(err error, packet nex.PacketInterface, callID uint32, sessionToken rating_types.RatingSessionToken) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.Unk2 == nil {
		return notImplemented("Unk2")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := commonProtocol.Unk2(packet.Sender().PID(), sessionToken); err != nil {
		return hookFailed("Unk2", err)
	}

	return reply(packet, rating.MethodUnk2, callID, nil), nil
}

func (commonProtocol *CommonProtocol) reportRatingStats(err error, packet nex.PacketInterface, callID uint32, sessionToken rating_types.RatingSessionToken, stats types.List[rating_types.RatingStats]) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.ReportRatingStats == nil {
		return notImplemented("ReportRatingStats")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := commonProtocol.ReportRatingStats(packet.Sender().PID(), sessionToken, stats); err != nil {
		return hookFailed("ReportRatingStats", err)
	}

	return reply(packet, rating.MethodReportRatingStats, callID, nil), nil
}

func (commonProtocol *CommonProtocol) getRanking(err error, packet nex.PacketInterface, callID uint32, category types.UInt32, uniqueID types.UInt64, principalID types.PID) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetRanking == nil {
		return notImplemented("GetRanking")
	}
	if err != nil {
		return invalidArgument(err)
	}
	rankData, err := commonProtocol.GetRanking(packet.Sender().PID(), category, uniqueID, principalID)
	if err != nil {
		return hookFailed("GetRanking", err)
	}

	return reply(packet, rating.MethodGetRanking, callID, func(s *nex.ByteStreamOut) { rankData.WriteTo(s) }), nil
}

func (commonProtocol *CommonProtocol) deleteScore(err error, packet nex.PacketInterface, callID uint32, category types.UInt32, uniqueID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.DeleteScore == nil {
		return notImplemented("DeleteScore")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := commonProtocol.DeleteScore(packet.Sender().PID(), category, uniqueID); err != nil {
		return hookFailed("DeleteScore", err)
	}

	return reply(packet, rating.MethodDeleteScore, callID, nil), nil
}

func (commonProtocol *CommonProtocol) uploadCommonData(err error, packet nex.PacketInterface, callID uint32, commonData types.Buffer, uniqueID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.UploadCommonData == nil {
		return notImplemented("UploadCommonData")
	}
	if err != nil {
		return invalidArgument(err)
	}
	if err := commonProtocol.UploadCommonData(packet.Sender().PID(), commonData, uniqueID); err != nil {
		return hookFailed("UploadCommonData", err)
	}

	return reply(packet, rating.MethodUploadCommonData, callID, nil), nil
}

func (commonProtocol *CommonProtocol) getCommonData(err error, packet nex.PacketInterface, callID uint32, uniqueID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetCommonData == nil {
		return notImplemented("GetCommonData")
	}
	if err != nil {
		return invalidArgument(err)
	}
	commonData, err := commonProtocol.GetCommonData(packet.Sender().PID(), uniqueID)
	if err != nil {
		return hookFailed("GetCommonData", err)
	}

	return reply(packet, rating.MethodGetCommonData, callID, func(s *nex.ByteStreamOut) { commonData.WriteTo(s) }), nil
}
