package privileges

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	privileges "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/privileges"
)

// buildGetPrivilegeRemainDurationResponse writes the response of Privileges::GetPrivilegeRemainDuration to the given stream
func buildGetPrivilegeRemainDurationResponse(stream *nex.ByteStreamOut, seconds types.Int32) {
	seconds.WriteTo(stream)
}

func (commonProtocol *CommonProtocol) getPrivilegeRemainDuration(err error, packet nex.PacketInterface, callID uint32, privilegeID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	if commonProtocol.GetPrivilegeRemainDuration == nil {
		common_globals.Logger.Warning("Privileges::GetPrivilegeRemainDuration not implemented")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	respSeconds, nexError := commonProtocol.GetPrivilegeRemainDuration(packet, privilegeID)
	if nexError != nil {
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	endpoint := packet.Sender().Endpoint()
	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	buildGetPrivilegeRemainDurationResponse(rmcResponseStream, respSeconds)

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = privileges.ProtocolID
	rmcResponse.MethodID = privileges.MethodGetPrivilegeRemainDuration
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
