package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
	ubi_account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management/types"
)

// buildValidatePasswordResponse writes the response of Ubi Account Management::ValidatePassword to the given stream
func buildValidatePasswordResponse(stream *nex.ByteStreamOut, failedReasons types.List[ubi_account_management_types.ValidationFailureReason]) {
	failedReasons.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) validatePassword(err error, packet nex.PacketInterface, callID uint32, password types.String, username types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.ValidatePassword == nil {
		common_globals.Logger.Warning("UbiAccountManagement::ValidatePassword not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respFailedReasons, nexError := commonProtocol.ValidatePassword(packet, password, username)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildValidatePasswordResponse(rmcResponseStream, respFailedReasons)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodValidatePassword
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
