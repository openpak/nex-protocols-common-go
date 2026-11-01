package useraccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	user_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-account-management"
)

// buildIsUserPlayingResponse writes the response of User Account Management::IsUserPlaying to the given stream
func buildIsUserPlayingResponse(stream *nex.ByteStreamOut, userPresence types.Map[types.String, types.Bool]) {
	userPresence.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) isUserPlaying(err error, packet nex.PacketInterface, callID uint32, firstPartyIds types.List[types.String], platformId types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.IsUserPlaying == nil {
		common_globals.Logger.Warning("UserAccountManagement::IsUserPlaying not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respUserPresence, nexError := commonProtocol.IsUserPlaying(packet, firstPartyIds, platformId)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildIsUserPlayingResponse(rmcResponseStream, respUserPresence)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = user_account_management.ProtocolID
	rmcResponse.MethodID = user_account_management.MethodIsUserPlaying
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
