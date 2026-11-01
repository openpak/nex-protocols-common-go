package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
	ubi_account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management/types"
)

// buildHasAcceptedLatestTOSResponse writes the response of Ubi Account Management::HasAcceptedLatestTOS to the given stream
func buildHasAcceptedLatestTOSResponse(stream *nex.ByteStreamOut, hasAccepted types.Bool, failedReasons types.List[ubi_account_management_types.ValidationFailureReason]) {
	hasAccepted.WriteTo(stream)
	failedReasons.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) hasAcceptedLatestTOS(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.HasAcceptedLatestTOS == nil {
		common_globals.Logger.Warning("UbiAccountManagement::HasAcceptedLatestTOS not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respHasAccepted, respFailedReasons, nexError := commonProtocol.HasAcceptedLatestTOS(packet)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildHasAcceptedLatestTOSResponse(rmcResponseStream, respHasAccepted, respFailedReasons)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodHasAcceptedLatestTOS
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
