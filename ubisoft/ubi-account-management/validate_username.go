package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
	ubi_account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management/types"
)

// buildValidateUsernameResponse writes the response of Ubi Account Management::ValidateUsername to the given stream
func buildValidateUsernameResponse(stream *nex.ByteStreamOut, usernameValidation ubi_account_management_types.UsernameValidation) {
	usernameValidation.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) validateUsername(err error, packet nex.PacketInterface, callID uint32, username types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.ValidateUsername == nil {
		common_globals.Logger.Warning("UbiAccountManagement::ValidateUsername not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respUsernameValidation, nexError := commonProtocol.ValidateUsername(packet, username)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildValidateUsernameResponse(rmcResponseStream, respUsernameValidation)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodValidateUsername
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
