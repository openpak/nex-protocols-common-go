package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
)

// buildLookupUbiAccountIDsByUsernamesGlobalResponse writes the response of Ubi Account Management::LookupUbiAccountIDsByUsernamesGlobal to the given stream
func buildLookupUbiAccountIDsByUsernamesGlobalResponse(stream *nex.ByteStreamOut, ubiAccountIDs types.Map[types.String, types.String]) {
	ubiAccountIDs.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) lookupUbiAccountIDsByUsernamesGlobal(err error, packet nex.PacketInterface, callID uint32, usernames types.List[types.String]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.LookupUbiAccountIDsByUsernamesGlobal == nil {
		common_globals.Logger.Warning("UbiAccountManagement::LookupUbiAccountIDsByUsernamesGlobal not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respUbiAccountIDs, nexError := commonProtocol.LookupUbiAccountIDsByUsernamesGlobal(packet, usernames)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildLookupUbiAccountIDsByUsernamesGlobalResponse(rmcResponseStream, respUbiAccountIDs)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodLookupUbiAccountIDsByUsernamesGlobal
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
