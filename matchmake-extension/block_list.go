package matchmake_extension

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmc"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
)

func (commonProtocol *CommonProtocol) getMyBlockList(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	lstPrincipalID := types.NewList[types.PID]()

	if commonProtocol.GetMyBlockList != nil {
		var nexError *nex.Error

		lstPrincipalID, nexError = commonProtocol.GetMyBlockList(packet)
		if nexError != nil {
			return nil, nexError
		}
	}

	return rmc.Respond(packet, matchmake_extension.ProtocolID, matchmake_extension.MethodGetMyBlockList, callID, lstPrincipalID), nil
}

func (commonProtocol *CommonProtocol) clearMyBlockList(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.ClearMyBlockList != nil {
		if nexError := commonProtocol.ClearMyBlockList(packet); nexError != nil {
			return nil, nexError
		}
	}

	return rmc.Respond(packet, matchmake_extension.ProtocolID, matchmake_extension.MethodClearMyBlockList, callID), nil
}
