package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
)

// buildForgetPasswordResponse writes the response of Ubi Account Management::ForgetPassword to the given stream
func buildForgetPasswordResponse(stream *nex.ByteStreamOut) {
}

func (commonProtocol *CommonProtocol) forgetPassword(err error, packet nex.PacketInterface, callID uint32, usernameOrEmail types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.ForgetPassword == nil {
		common_globals.Logger.Warning("UbiAccountManagement::ForgetPassword not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	nexError := commonProtocol.ForgetPassword(packet, usernameOrEmail)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildForgetPasswordResponse(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodForgetPassword
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
