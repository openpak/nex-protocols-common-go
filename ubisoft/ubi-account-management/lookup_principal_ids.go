package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
)

// buildLookupPrincipalIdsResponse writes the response of Ubi Account Management::LookupPrincipalIds to the given stream
func buildLookupPrincipalIdsResponse(stream *nex.ByteStreamOut, pids types.Map[types.String, types.UInt32]) {
	pids.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) lookupPrincipalIds(err error, packet nex.PacketInterface, callID uint32, ubiAccountIds types.List[types.String]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.LookupPrincipalIds == nil {
		common_globals.Logger.Warning("UbiAccountManagement::LookupPrincipalIds not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respPids, nexError := commonProtocol.LookupPrincipalIds(packet, ubiAccountIds)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildLookupPrincipalIdsResponse(rmcResponseStream, respPids)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodLookupPrincipalIds
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
