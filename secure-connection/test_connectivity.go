package secureconnection

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmcutil"
	secure_connection "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
)

// testConnectivity takes no parameters and returns nothing: the console only
// wants to know that the connection still carries RMC, so the reply arriving
// is the answer. It acknowledges every call unless the TestConnectivity hook
// returns an error.
func (commonProtocol *CommonProtocol) testConnectivity(_ error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.TestConnectivity != nil {
		if err := commonProtocol.TestConnectivity(packet.Sender().PID()); err != nil {
			return nil, rmcutil.HookError(err, nex.ResultCodes.Core.Unknown)
		}
	}

	return rmcutil.Success(packet.Sender().Endpoint(), secure_connection.ProtocolID, secure_connection.MethodTestConnectivity, callID, nil), nil
}
