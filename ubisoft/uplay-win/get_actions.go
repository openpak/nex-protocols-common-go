package uplaywin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	uplay_win "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win"
	uplay_win_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win/types"
)

// buildGetActionsResponse writes the response of Uplay Win::GetActions to the given stream
func buildGetActionsResponse(stream *nex.ByteStreamOut, actionList types.List[uplay_win_types.UplayAction]) {
	actionList.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getActions(err error, packet nex.PacketInterface, callID uint32, startRowIndex types.Int32, maximumRows types.Int32, sortExpression types.String, cultureName types.String, platformCode types.String, gameCode types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetActions == nil {
		common_globals.Logger.Warning("UplayWin::GetActions not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respActionList, nexError := commonProtocol.GetActions(packet, startRowIndex, maximumRows, sortExpression, cultureName, platformCode, gameCode)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetActionsResponse(rmcResponseStream, respActionList)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = uplay_win.ProtocolID
	rmcResponse.MethodID = uplay_win.MethodGetActions
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
