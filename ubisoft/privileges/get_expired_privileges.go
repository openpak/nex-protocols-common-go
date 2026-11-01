package privileges

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	privileges "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges"
	privileges_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges/types"
)

// buildGetExpiredPrivilegesResponse writes the response of Privileges::GetExpiredPrivileges to the given stream
func buildGetExpiredPrivilegesResponse(stream *nex.ByteStreamOut, expiredPrivileges types.List[privileges_types.PrivilegeEx]) {
	expiredPrivileges.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getExpiredPrivileges(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetExpiredPrivileges == nil {
		common_globals.Logger.Warning("Privileges::GetExpiredPrivileges not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respExpiredPrivileges, nexError := commonProtocol.GetExpiredPrivileges(packet)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetExpiredPrivilegesResponse(rmcResponseStream, respExpiredPrivileges)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = privileges.ProtocolID
	rmcResponse.MethodID = privileges.MethodGetExpiredPrivileges
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
