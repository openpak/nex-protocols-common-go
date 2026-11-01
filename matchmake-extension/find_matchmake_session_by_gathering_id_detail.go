package matchmake_extension

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmc"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension/database"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
)

// scrubSessionDetail removes what only joining a session should hand out. This method can be called by anyone who
// knows a gathering ID, so the session key and the user password must not be part of the answer
func scrubSessionDetail(session match_making_types.MatchmakeSession) match_making_types.MatchmakeSession {
	session.SessionKey = types.NewBuffer(nil)
	session.UserPassword = ""

	return session
}

func (commonProtocol *CommonProtocol) findMatchmakeSessionByGatheringIDDetail(err error, packet nex.PacketInterface, callID uint32, gid types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	endpoint, ok := packet.Sender().Endpoint().(*nex.PRUDPEndPoint)
	if !ok {
		return nil, nex.NewError(nex.ResultCodes.Core.Exception, "endpoint is not a PRUDP endpoint")
	}

	commonProtocol.manager.Mutex.RLock()

	session, _, nexError := database.GetMatchmakeSessionByID(commonProtocol.manager, endpoint, uint32(gid))
	if nexError != nil {
		commonProtocol.manager.Mutex.RUnlock()
		return nil, nexError
	}

	commonProtocol.manager.Mutex.RUnlock()

	if commonProtocol.OnAfterFindMatchmakeSessionByGatheringIDDetail != nil {
		go commonProtocol.OnAfterFindMatchmakeSessionByGatheringIDDetail(packet, gid)
	}

	return rmc.Respond(packet, matchmake_extension.ProtocolID, matchmake_extension.MethodFindMatchmakeSessionByGatheringIDDetail, callID, scrubSessionDetail(session)), nil
}
