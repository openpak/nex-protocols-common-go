package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
)

// buildLookupUbiAccountIDsByPidsResponse writes the response of Ubi Account Management::LookupUbiAccountIDsByPids to the given stream
func buildLookupUbiAccountIDsByPidsResponse(stream *nex.ByteStreamOut, ubiaccountIDs types.Map[types.UInt32, types.String]) {
	ubiaccountIDs.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) lookupUbiAccountIDsByPids(err error, packet nex.PacketInterface, callID uint32, pids types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.LookupUbiAccountIDsByPids == nil {
		common_globals.Logger.Warning("UbiAccountManagement::LookupUbiAccountIDsByPids not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respUbiaccountIDs, nexError := commonProtocol.LookupUbiAccountIDsByPids(packet, pids)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildLookupUbiAccountIDsByPidsResponse(rmcResponseStream, respUbiaccountIDs)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodLookupUbiAccountIDsByPids
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
