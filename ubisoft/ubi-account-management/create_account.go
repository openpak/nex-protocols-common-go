package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
	ubi_account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management/types"
)

// buildCreateAccountResponse writes the response of Ubi Account Management::CreateAccount to the given stream
func buildCreateAccountResponse(stream *nex.ByteStreamOut, ubiAccount ubi_account_management_types.UbiAccount, failedReasons types.List[ubi_account_management_types.ValidationFailureReason]) {
	ubiAccount.WriteTo(stream)
	failedReasons.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) createAccount(err error, packet nex.PacketInterface, callID uint32, ubiAccount ubi_account_management_types.UbiAccount) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.CreateAccount == nil {
		common_globals.Logger.Warning("UbiAccountManagement::CreateAccount not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respUbiAccount, respFailedReasons, nexError := commonProtocol.CreateAccount(packet, ubiAccount)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildCreateAccountResponse(rmcResponseStream, respUbiAccount, respFailedReasons)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodCreateAccount
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
