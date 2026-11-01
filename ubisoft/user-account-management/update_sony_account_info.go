package useraccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	user_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-account-management"
)

// buildUpdateSonyAccountInfoResponse writes the response of User Account Management::UpdateSonyAccountInfo to the given stream
func buildUpdateSonyAccountInfoResponse(stream *nex.ByteStreamOut) {
}

func (commonProtocol *CommonProtocol) updateSonyAccountInfo(err error, packet nex.PacketInterface, callID uint32, ticketData types.QBuffer, ticketSize types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.UpdateSonyAccountInfo == nil {
		common_globals.Logger.Warning("UserAccountManagement::UpdateSonyAccountInfo not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	nexError := commonProtocol.UpdateSonyAccountInfo(packet, ticketData, ticketSize)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildUpdateSonyAccountInfoResponse(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = user_account_management.ProtocolID
	rmcResponse.MethodID = user_account_management.MethodUpdateSonyAccountInfo
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
