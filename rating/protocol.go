// Package rating implements a hook-driven common handler set for the Rating
// protocol. The game server supplies storage through the function fields on
// CommonProtocol; this package decodes requests, calls the hooks and encodes
// the responses.
//
// A method whose hook is nil answers Core::NotImplemented. Hooks that return an
// error answer Core::Unknown unless the error is a *nex.Error, which is sent
// as is.
//
// The response layouts of GetRanking (a single RatingRankData) and the two
// unnamed methods come from observed traffic, not from a full specification.
package rating

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	rating "github.com/PretendoNetwork/nex-protocols-go/v2/rating"
	rating_types "github.com/PretendoNetwork/nex-protocols-go/v2/rating/types"
)

// CommonProtocol holds the hooks for the Rating protocol.
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol rating.Interface

	// Unk1 is method 1 (name unknown). It takes no parameters and returns a
	// number and a string; observed clients get a session id and a hex token.
	Unk1 func(pid types.PID) (types.UInt64, types.String, error)
	// Unk2 is method 2 (name unknown). It receives the session token and returns nothing.
	Unk2 func(pid types.PID, sessionToken rating_types.RatingSessionToken) error
	// ReportRatingStats stores the stats reported by the caller.
	ReportRatingStats func(pid types.PID, sessionToken rating_types.RatingSessionToken, stats types.List[rating_types.RatingStats]) error
	// GetRanking returns the rank entry for a user in a category.
	GetRanking func(pid types.PID, category types.UInt32, uniqueID types.UInt64, principalID types.PID) (rating_types.RatingRankData, error)
	// DeleteScore deletes the caller's score in a category.
	DeleteScore func(pid types.PID, category types.UInt32, uniqueID types.UInt64) error
	// UploadCommonData stores the caller's common data.
	UploadCommonData func(pid types.PID, commonData types.Buffer, uniqueID types.UInt64) error
	// GetCommonData returns the common data of a unique id.
	GetCommonData func(pid types.PID, uniqueID types.UInt64) (types.Buffer, error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol rating.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerUnk1(commonProtocol.unk1)
	protocol.SetHandlerUnk2(commonProtocol.unk2)
	protocol.SetHandlerReportRatingStats(commonProtocol.reportRatingStats)
	protocol.SetHandlerGetRanking(commonProtocol.getRanking)
	protocol.SetHandlerDeleteScore(commonProtocol.deleteScore)
	protocol.SetHandlerUploadCommonData(commonProtocol.uploadCommonData)
	protocol.SetHandlerGetCommonData(commonProtocol.getCommonData)

	return commonProtocol
}
