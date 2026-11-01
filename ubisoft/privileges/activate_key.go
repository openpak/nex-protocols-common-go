package privileges

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	privileges "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges"
	privileges_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges/types"
)

// buildActivateKeyResponse writes the response of Privileges::ActivateKey to the given stream
func buildActivateKeyResponse(stream *nex.ByteStreamOut, privilege privileges_types.PrivilegeGroup) {
	privilege.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) activateKey(err error, packet nex.PacketInterface, callID uint32, uniqueKey types.String, languageCode types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.ActivateKey == nil {
		common_globals.Logger.Warning("Privileges::ActivateKey not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respPrivilege, nexError := commonProtocol.ActivateKey(packet, uniqueKey, languageCode)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildActivateKeyResponse(rmcResponseStream, respPrivilege)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = privileges.ProtocolID
	rmcResponse.MethodID = privileges.MethodActivateKey
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
