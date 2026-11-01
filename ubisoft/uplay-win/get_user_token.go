package uplaywin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	uplay_win "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/uplay-win"
)

// buildGetUserTokenResponse writes the response of Uplay Win::GetUserToken to the given stream
func buildGetUserTokenResponse(stream *nex.ByteStreamOut, token types.String) {
	token.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getUserToken(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetUserToken == nil {
		common_globals.Logger.Warning("UplayWin::GetUserToken not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respToken, nexError := commonProtocol.GetUserToken(packet)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetUserTokenResponse(rmcResponseStream, respToken)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = uplay_win.ProtocolID
	rmcResponse.MethodID = uplay_win.MethodGetUserToken
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
