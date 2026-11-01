package matchmaking

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmc"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making/database"
	match_making "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
)

// maxFindByParticipants is the most PIDs FindByParticipants looks up in one call
const maxFindByParticipants = 300

// loadGatheringHolders resolves gathering IDs to the GatheringHolders sent back by the Find* methods.
// Gatherings that disappeared in the meantime (or that the manager cannot describe) are skipped
func (commonProtocol *CommonProtocol) loadGatheringHolders(sourcePID types.PID, ids []uint32) types.List[match_making_types.GatheringHolder] {
	holders := types.NewList[match_making_types.GatheringHolder]()

	for _, id := range ids {
		gathering, _, nexError := commonProtocol.manager.GetDetailedGatheringByID(commonProtocol.manager, uint64(sourcePID), id)
		if nexError != nil {
			common_globals.Logger.Warningf("Skipping gathering %d: %s", id, nexError.Error())
			continue
		}

		gatheringInterface, ok := gathering.Copy().(match_making_types.GatheringInterface)
		if !ok {
			continue
		}

		holder := match_making_types.NewGatheringHolder()
		holder.Object = gatheringInterface

		holders = append(holders, holder)
	}

	return holders
}

// findGatherings runs a query under the manager lock and answers with the matching gatherings
func (commonProtocol *CommonProtocol) findGatherings(packet nex.PacketInterface, callID uint32, methodID uint32, query func() ([]uint32, *nex.Error)) (*nex.RMCMessage, *nex.Error) {
	commonProtocol.manager.Mutex.RLock()

	ids, nexError := query()
	if nexError != nil {
		commonProtocol.manager.Mutex.RUnlock()
		common_globals.Logger.Error(nexError.Error())
		return nil, nexError
	}

	holders := commonProtocol.loadGatheringHolders(packet.Sender().PID(), ids)

	commonProtocol.manager.Mutex.RUnlock()

	return rmc.Respond(packet, match_making.ProtocolID, methodID, callID, holders), nil
}

func (commonProtocol *CommonProtocol) findByOwner(err error, packet nex.PacketInterface, callID uint32, id types.PID, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	return commonProtocol.findGatherings(packet, callID, match_making.MethodFindByOwner, func() ([]uint32, *nex.Error) {
		return database.FindGatheringIDsByOwner(commonProtocol.manager, id, resultRange)
	})
}

func (commonProtocol *CommonProtocol) findByType(err error, packet nex.PacketInterface, callID uint32, strType types.String, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	return commonProtocol.findGatherings(packet, callID, match_making.MethodFindByType, func() ([]uint32, *nex.Error) {
		return database.FindGatheringIDsByType(commonProtocol.manager, string(strType), resultRange)
	})
}

func (commonProtocol *CommonProtocol) findByDescription(err error, packet nex.PacketInterface, callID uint32, strDescription types.String, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	return commonProtocol.findGatherings(packet, callID, match_making.MethodFindByDescription, func() ([]uint32, *nex.Error) {
		return database.FindGatheringIDsByDescription(commonProtocol.manager, string(strDescription), resultRange)
	})
}

func (commonProtocol *CommonProtocol) findByDescriptionRegex(err error, packet nex.PacketInterface, callID uint32, strDescriptionRegex types.String, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if nexError := database.ValidateDescriptionRegex(string(strDescriptionRegex)); nexError != nil {
		return nil, nexError
	}

	return commonProtocol.findGatherings(packet, callID, match_making.MethodFindByDescriptionRegex, func() ([]uint32, *nex.Error) {
		return database.FindGatheringIDsByDescriptionRegex(commonProtocol.manager, string(strDescriptionRegex), resultRange)
	})
}

func (commonProtocol *CommonProtocol) findByID(err error, packet nex.PacketInterface, callID uint32, lstID types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if len(lstID) > database.MaxFindResults {
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "lstID is bigger than the maximum number of results")
	}

	return commonProtocol.findGatherings(packet, callID, match_making.MethodFindByID, func() ([]uint32, *nex.Error) {
		ids := make([]uint32, 0, len(lstID))
		for _, id := range lstID {
			ids = append(ids, uint32(id))
		}

		return ids, nil
	})
}

func (commonProtocol *CommonProtocol) findByParticipants(err error, packet nex.PacketInterface, callID uint32, pid types.List[types.PID]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if len(pid) > maxFindByParticipants {
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "pid is bigger than the maximum number of participants")
	}

	return commonProtocol.findGatherings(packet, callID, match_making.MethodFindByParticipants, func() ([]uint32, *nex.Error) {
		return database.FindGatheringIDsByParticipants(commonProtocol.manager, pid)
	})
}

// findInvitations answers an empty list. The common library does not implement Invite, so no invitations are
// ever stored that could be found
func (commonProtocol *CommonProtocol) findInvitations(err error, packet nex.PacketInterface, callID uint32, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	return rmc.Respond(packet, match_making.ProtocolID, match_making.MethodFindInvitations, callID, types.NewList[match_making_types.GatheringHolder]()), nil
}

// findBySQLQuery never executes the query. The client sends a raw SQL string meant for the original Rendez-Vous
// server's own database schema, which has nothing in common with ours. Running it would be an SQL injection by
// design, and translating it is not possible in general, so the query is only logged and an empty list is answered
func (commonProtocol *CommonProtocol) findBySQLQuery(err error, packet nex.PacketInterface, callID uint32, strQuery types.String, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	query := string(strQuery)
	if len(query) > 200 {
		query = query[:200] + "..."
	}

	common_globals.Logger.Warningf("FindBySQLQuery is not executed, answering an empty list. Query: %q", query)

	return rmc.Respond(packet, match_making.ProtocolID, match_making.MethodFindBySQLQuery, callID, types.NewList[match_making_types.GatheringHolder]()), nil
}
