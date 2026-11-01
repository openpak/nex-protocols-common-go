package uplaywin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	uplay_win "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win"
)

// buildGetActionsCountResponse writes the response of Uplay Win::GetActionsCount to the given stream
func buildGetActionsCountResponse(stream *nex.ByteStreamOut, actionsCount types.Int32) {
	actionsCount.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getActionsCount(err error, packet nex.PacketInterface, callID uint32, platformCode types.String, gameCode types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetActionsCount == nil {
		common_globals.Logger.Warning("UplayWin::GetActionsCount not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respActionsCount, nexError := commonProtocol.GetActionsCount(packet, platformCode, gameCode)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetActionsCountResponse(rmcResponseStream, respActionsCount)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = uplay_win.ProtocolID
	rmcResponse.MethodID = uplay_win.MethodGetActionsCount
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
