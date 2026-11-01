package database

import (
	"database/sql"
	"math"
	"regexp"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
)

const (
	// MaxFindResults is the most gathering IDs a single Find* call returns
	MaxFindResults = 100

	// MaxFindRegexLength is the longest description regex accepted from a client
	MaxFindRegexLength = 256
)

// NormalizeResultRange converts a client ResultRange into a SQL limit and offset.
// A zero length means "no preference" and an offset of 0xFFFFFFFF means "from the start".
// The limit is always capped to MaxFindResults
func NormalizeResultRange(resultRange types.ResultRange) (uint32, uint32) {
	limit := uint32(resultRange.Length)
	offset := uint32(resultRange.Offset)

	if limit == 0 || limit > MaxFindResults {
		limit = MaxFindResults
	}

	if offset == math.MaxUint32 {
		offset = 0
	}

	return limit, offset
}

// ValidateDescriptionRegex checks a client supplied description regex before it is handed to the database.
// The pattern is length limited and has to be valid RE2 syntax, which rejects backreferences and other
// constructs that make a backtracking regex engine expensive
func ValidateDescriptionRegex(pattern string) *nex.Error {
	if len(pattern) > MaxFindRegexLength {
		return nex.NewError(nex.ResultCodes.Core.InvalidArgument, "description regex is too long")
	}

	if _, err := regexp.Compile(pattern); err != nil {
		return nex.NewError(nex.ResultCodes.Core.InvalidArgument, "description regex is invalid")
	}

	return nil
}

func scanGatheringIDs(rows *sql.Rows) ([]uint32, *nex.Error) {
	defer rows.Close()

	ids := make([]uint32, 0)

	for rows.Next() {
		var id uint32

		if err := rows.Scan(&id); err != nil {
			common_globals.Logger.Error(err.Error())
			continue
		}

		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, nex.NewError(nex.ResultCodes.Core.Unknown, err.Error())
	}

	return ids, nil
}

func queryGatheringIDs(manager *common_globals.MatchmakingManager, where string, arg any, resultRange types.ResultRange) ([]uint32, *nex.Error) {
	limit, offset := NormalizeResultRange(resultRange)

	rows, err := manager.Database.Query(`SELECT id FROM matchmaking.gatherings WHERE registered=true AND `+where+` ORDER BY id LIMIT $2 OFFSET $3`, arg, limit, offset)
	if err != nil {
		return nil, nex.NewError(nex.ResultCodes.Core.Unknown, err.Error())
	}

	return scanGatheringIDs(rows)
}

// FindGatheringIDsByOwner returns the IDs of the registered gatherings owned by the given PID
func FindGatheringIDsByOwner(manager *common_globals.MatchmakingManager, owner types.PID, resultRange types.ResultRange) ([]uint32, *nex.Error) {
	return queryGatheringIDs(manager, `owner_pid=$1`, owner, resultRange)
}

// FindGatheringIDsByType returns the IDs of the registered gatherings of the given class name
// (Gathering, MatchmakeSession, PersistentGathering)
func FindGatheringIDsByType(manager *common_globals.MatchmakingManager, gatheringType string, resultRange types.ResultRange) ([]uint32, *nex.Error) {
	return queryGatheringIDs(manager, `type=$1`, gatheringType, resultRange)
}

// FindGatheringIDsByDescription returns the IDs of the registered gatherings with exactly the given description
func FindGatheringIDsByDescription(manager *common_globals.MatchmakingManager, description string, resultRange types.ResultRange) ([]uint32, *nex.Error) {
	return queryGatheringIDs(manager, `description=$1`, description, resultRange)
}

// FindGatheringIDsByDescriptionRegex returns the IDs of the registered gatherings whose description matches the regex.
// The pattern must have passed ValidateDescriptionRegex
func FindGatheringIDsByDescriptionRegex(manager *common_globals.MatchmakingManager, pattern string, resultRange types.ResultRange) ([]uint32, *nex.Error) {
	return queryGatheringIDs(manager, `description ~ $1`, pattern, resultRange)
}

// FindGatheringIDsByParticipants returns the IDs of the registered gatherings any of the given PIDs participates in
func FindGatheringIDsByParticipants(manager *common_globals.MatchmakingManager, participants types.List[types.PID]) ([]uint32, *nex.Error) {
	ids := make([]uint32, 0)
	seen := make(map[uint32]bool)

	for _, pid := range participants {
		found, nexError := queryGatheringIDs(manager, `$1=ANY(participants)`, pid, types.ResultRange{})
		if nexError != nil {
			return nil, nexError
		}

		for _, id := range found {
			if !seen[id] && len(ids) < MaxFindResults {
				seen[id] = true
				ids = append(ids, id)
			}
		}
	}

	return ids, nil
}
