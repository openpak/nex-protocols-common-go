package matchmake_extension

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmc"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension/database"
	match_making_constants "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/constants"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
)

const (
	// * MAX_PRINCIPALID_SIZE_TO_FIND_MATCHMAKE_SESSION, the same bound GetPlayingSession applies
	maxPrincipalIDsToFindMatchmakeSession = 300

	// * MAX_MATCHMAKE_SESSION_BY_PARTICIPANT, the same bound GetPlayingSession applies
	maxMatchmakeSessionsByParticipant = 1000
)

// buildParticipantResults turns the playing sessions of the requested principals into the response list.
// Only MatchmakeSessions are reported, and the ApplicationBuffer and MatchmakeParam are only populated when
// the client asked for them through the result options
func buildParticipantResults(playingSessions types.List[match_making_types.PlayingSession], options match_making_constants.FindMatchmakeSessionResultOption) types.List[match_making_types.FindMatchmakeSessionByParticipantResult] {
	results := types.NewList[match_making_types.FindMatchmakeSessionByParticipantResult]()

	for _, playingSession := range playingSessions {
		if len(results) >= maxMatchmakeSessionsByParticipant {
			break
		}

		matchmakeSession, ok := playingSession.Gathering.Object.(match_making_types.MatchmakeSession)
		if !ok {
			continue
		}

		if !options.HasFlag(match_making_constants.FindMatchmakeSessionResultOptionApplicationBuffer) {
			matchmakeSession.ApplicationBuffer = types.NewBuffer(nil)
		}

		if !options.HasFlag(match_making_constants.FindMatchmakeSessionResultOptionMatchmakeParam) {
			matchmakeSession.MatchmakeParam = match_making_types.NewMatchmakeParam()
		}

		result := match_making_types.NewFindMatchmakeSessionByParticipantResult()
		result.PrincipalID = playingSession.PrincipalID
		result.Session = matchmakeSession

		results = append(results, result)
	}

	return results
}

func (commonProtocol *CommonProtocol) findMatchmakeSessionByParticipant(err error, packet nex.PacketInterface, callID uint32, param match_making_types.FindMatchmakeSessionByParticipantParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if len(param.PrincipalIDList) > maxPrincipalIDsToFindMatchmakeSession {
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "PrincipalIDList is bigger than MAX_PRINCIPALID_SIZE_TO_FIND_MATCHMAKE_SESSION")
	}

	commonProtocol.manager.Mutex.RLock()

	playingSessions, nexError := database.GetPlayingSession(commonProtocol.manager, param.PrincipalIDList)
	if nexError != nil {
		commonProtocol.manager.Mutex.RUnlock()
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	commonProtocol.manager.Mutex.RUnlock()

	// * TODO - param.BlockListParam is ignored, the common library has no block list to consult yet
	results := buildParticipantResults(playingSessions, param.ResultOptions)

	if commonProtocol.OnAfterFindMatchmakeSessionByParticipant != nil {
		go commonProtocol.OnAfterFindMatchmakeSessionByParticipant(packet, param)
	}

	return rmc.Respond(packet, matchmake_extension.ProtocolID, matchmake_extension.MethodFindMatchmakeSessionByParticipant, callID, results), nil
}
