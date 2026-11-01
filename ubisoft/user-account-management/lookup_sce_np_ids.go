package useraccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	user_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-account-management"
)

// buildLookupSceNpIdsResponse writes the response of User Account Management::LookupSceNpIds to the given stream
func buildLookupSceNpIdsResponse(stream *nex.ByteStreamOut, npids types.Map[types.UInt32, types.QBuffer]) {
	npids.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) lookupSceNpIds(err error, packet nex.PacketInterface, callID uint32, pids types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.LookupSceNpIds == nil {
		common_globals.Logger.Warning("UserAccountManagement::LookupSceNpIds not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respNpids, nexError := commonProtocol.LookupSceNpIds(packet, pids)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildLookupSceNpIdsResponse(rmcResponseStream, respNpids)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = user_account_management.ProtocolID
	rmcResponse.MethodID = user_account_management.MethodLookupSceNpIds
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
