package uplaywin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	uplay_win "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win"
)

// buildGetActionsCompletedCountResponse writes the response of Uplay Win::GetActionsCompletedCount to the given stream
func buildGetActionsCompletedCountResponse(stream *nex.ByteStreamOut, actionsCount types.Int32) {
	actionsCount.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getActionsCompletedCount(err error, packet nex.PacketInterface, callID uint32, platformCode types.String, gameCode types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetActionsCompletedCount == nil {
		common_globals.Logger.Warning("UplayWin::GetActionsCompletedCount not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respActionsCount, nexError := commonProtocol.GetActionsCompletedCount(packet, platformCode, gameCode)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetActionsCompletedCountResponse(rmcResponseStream, respActionsCount)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = uplay_win.ProtocolID
	rmcResponse.MethodID = uplay_win.MethodGetActionsCompletedCount
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
