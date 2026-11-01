package privileges

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	privileges "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges"
	privileges_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges/types"
)

// buildGetPrivilegesExResponse writes the response of Privileges::GetPrivilegesEx to the given stream
func buildGetPrivilegesExResponse(stream *nex.ByteStreamOut, privilegesEx types.List[privileges_types.PrivilegeEx]) {
	privilegesEx.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getPrivilegesEx(err error, packet nex.PacketInterface, callID uint32, localeCode types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetPrivilegesEx == nil {
		common_globals.Logger.Warning("Privileges::GetPrivilegesEx not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respPrivilegesEx, nexError := commonProtocol.GetPrivilegesEx(packet, localeCode)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetPrivilegesExResponse(rmcResponseStream, respPrivilegesEx)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = privileges.ProtocolID
	rmcResponse.MethodID = privileges.MethodGetPrivilegesEx
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
