package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
)

// buildLookupUsernamesByUbiAccountIDsGlobalResponse writes the response of Ubi Account Management::LookupUsernamesByUbiAccountIDsGlobal to the given stream
func buildLookupUsernamesByUbiAccountIDsGlobalResponse(stream *nex.ByteStreamOut, usernames types.Map[types.String, types.String]) {
	usernames.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) lookupUsernamesByUbiAccountIDsGlobal(err error, packet nex.PacketInterface, callID uint32, ubiAccountIds types.List[types.String]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.LookupUsernamesByUbiAccountIDsGlobal == nil {
		common_globals.Logger.Warning("UbiAccountManagement::LookupUsernamesByUbiAccountIDsGlobal not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respUsernames, nexError := commonProtocol.LookupUsernamesByUbiAccountIDsGlobal(packet, ubiAccountIds)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildLookupUsernamesByUbiAccountIDsGlobalResponse(rmcResponseStream, respUsernames)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodLookupUsernamesByUbiAccountIDsGlobal
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
