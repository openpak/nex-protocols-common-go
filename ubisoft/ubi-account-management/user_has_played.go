package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
)

// buildUserHasPlayedResponse writes the response of Ubi Account Management::UserHasPlayed to the given stream
func buildUserHasPlayedResponse(stream *nex.ByteStreamOut, userPresence types.Map[types.String, types.Bool]) {
	userPresence.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) userHasPlayed(err error, packet nex.PacketInterface, callID uint32, ubiAccountIDs types.List[types.String]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.UserHasPlayed == nil {
		common_globals.Logger.Warning("UbiAccountManagement::UserHasPlayed not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respUserPresence, nexError := commonProtocol.UserHasPlayed(packet, ubiAccountIDs)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildUserHasPlayedResponse(rmcResponseStream, respUserPresence)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodUserHasPlayed
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
