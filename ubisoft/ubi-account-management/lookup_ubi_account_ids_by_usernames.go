package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
)

// buildLookupUbiAccountIDsByUsernamesResponse writes the response of Ubi Account Management::LookupUbiAccountIDsByUsernames to the given stream
func buildLookupUbiAccountIDsByUsernamesResponse(stream *nex.ByteStreamOut, ubiAccountIDs types.Map[types.String, types.String]) {
	ubiAccountIDs.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) lookupUbiAccountIDsByUsernames(err error, packet nex.PacketInterface, callID uint32, usernames types.List[types.String]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.LookupUbiAccountIDsByUsernames == nil {
		common_globals.Logger.Warning("UbiAccountManagement::LookupUbiAccountIDsByUsernames not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respUbiAccountIDs, nexError := commonProtocol.LookupUbiAccountIDsByUsernames(packet, usernames)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildLookupUbiAccountIDsByUsernamesResponse(rmcResponseStream, respUbiAccountIDs)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodLookupUbiAccountIDsByUsernames
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
