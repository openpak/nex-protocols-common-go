package useraccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	user_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-account-management"
)

// buildLookupPrincipalIDsResponse writes the response of User Account Management::LookupPrincipalIDs to the given stream
func buildLookupPrincipalIDsResponse(stream *nex.ByteStreamOut, pids types.Map[types.String, types.UInt32]) {
	pids.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) lookupPrincipalIDs(err error, packet nex.PacketInterface, callID uint32, firstPartyIds types.List[types.String], platformId types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.LookupPrincipalIDs == nil {
		common_globals.Logger.Warning("UserAccountManagement::LookupPrincipalIDs not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respPids, nexError := commonProtocol.LookupPrincipalIDs(packet, firstPartyIds, platformId)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildLookupPrincipalIDsResponse(rmcResponseStream, respPids)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = user_account_management.ProtocolID
	rmcResponse.MethodID = user_account_management.MethodLookupPrincipalIDs
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
