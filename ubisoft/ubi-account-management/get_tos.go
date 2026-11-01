package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
	ubi_account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management/types"
)

// buildGetTOSResponse writes the response of Ubi Account Management::GetTOS to the given stream
func buildGetTOSResponse(stream *nex.ByteStreamOut, tos ubi_account_management_types.TOS) {
	tos.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getTOS(err error, packet nex.PacketInterface, callID uint32, countryCode types.String, languageCode types.String, htmlVersion types.Bool) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetTOS == nil {
		common_globals.Logger.Warning("UbiAccountManagement::GetTOS not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respTos, nexError := commonProtocol.GetTOS(packet, countryCode, languageCode, htmlVersion)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetTOSResponse(rmcResponseStream, respTos)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodGetTOS
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
