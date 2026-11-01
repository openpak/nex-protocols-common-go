package useraccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	user_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-account-management"
)

// buildLookupFirstPartyIdsResponse writes the response of User Account Management::LookupFirstPartyIds to the given stream
func buildLookupFirstPartyIdsResponse(stream *nex.ByteStreamOut, firstPartyIds types.Map[types.UInt32, types.String]) {
	firstPartyIds.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) lookupFirstPartyIds(err error, packet nex.PacketInterface, callID uint32, pids types.List[types.UInt32], platformId types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.LookupFirstPartyIds == nil {
		common_globals.Logger.Warning("UserAccountManagement::LookupFirstPartyIds not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respFirstPartyIds, nexError := commonProtocol.LookupFirstPartyIds(packet, pids, platformId)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildLookupFirstPartyIdsResponse(rmcResponseStream, respFirstPartyIds)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = user_account_management.ProtocolID
	rmcResponse.MethodID = user_account_management.MethodLookupFirstPartyIds
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
