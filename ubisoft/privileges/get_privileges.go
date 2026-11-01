package privileges

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	privileges "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges"
	privileges_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges/types"
)

// buildGetPrivilegesResponse writes the response of Privileges::GetPrivileges to the given stream
func buildGetPrivilegesResponse(stream *nex.ByteStreamOut, privileges types.Map[types.UInt32, privileges_types.Privilege]) {
	privileges.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getPrivileges(err error, packet nex.PacketInterface, callID uint32, localeCode types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetPrivileges == nil {
		common_globals.Logger.Warning("Privileges::GetPrivileges not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respPrivileges, nexError := commonProtocol.GetPrivileges(packet, localeCode)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetPrivilegesResponse(rmcResponseStream, respPrivileges)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = privileges.ProtocolID
	rmcResponse.MethodID = privileges.MethodGetPrivileges
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
