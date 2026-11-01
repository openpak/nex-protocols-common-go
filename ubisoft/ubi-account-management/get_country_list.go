package ubiaccountmanagement

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi_account_management "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management"
	ubi_account_management_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/ubi-account-management/types"
)

// buildGetCountryListResponse writes the response of Ubi Account Management::GetCountryList to the given stream
func buildGetCountryListResponse(stream *nex.ByteStreamOut, countries types.List[ubi_account_management_types.Country]) {
	countries.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getCountryList(err error, packet nex.PacketInterface, callID uint32, languageCode types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetCountryList == nil {
		common_globals.Logger.Warning("UbiAccountManagement::GetCountryList not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respCountries, nexError := commonProtocol.GetCountryList(packet, languageCode)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetCountryListResponse(rmcResponseStream, respCountries)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = ubi_account_management.ProtocolID
	rmcResponse.MethodID = ubi_account_management.MethodGetCountryList
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
