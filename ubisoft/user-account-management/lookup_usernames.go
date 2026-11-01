package useraccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	user_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-account-management"
)

// buildLookupUsernamesResponse writes the response of User Account Management::LookupUsernames to the given stream
func buildLookupUsernamesResponse(stream *nex.ByteStreamOut, userNames types.Map[types.UInt32, types.String]) {
	userNames.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) lookupUsernames(err error, packet nex.PacketInterface, callID uint32, pids types.List[types.UInt32], platformId types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.LookupUsernames == nil {
		common_globals.Logger.Warning("UserAccountManagement::LookupUsernames not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respUserNames, nexError := commonProtocol.LookupUsernames(packet, pids, platformId)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildLookupUsernamesResponse(rmcResponseStream, respUserNames)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = user_account_management.ProtocolID
	rmcResponse.MethodID = user_account_management.MethodLookupUsernames
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
