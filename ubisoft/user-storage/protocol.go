// Package userstorage implements the common User Storage protocol handlers of the Quazal Rendez-Vous (Ubisoft) family.
//
// Layouts follow the NintendoClients wiki page "User Storage Protocol" (MIT). Every method calls a user-provided hook;
// a method whose hook is nil answers Core::NotImplemented unless a safe default is documented on its hook field.
package userstorage

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-storage"
	user_storage_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-storage/types"
)

// CommonProtocol handles the User Storage protocol with user-provided hooks
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi.Interface

	// SearchContents is optional: when nil an empty result (empty lists, zero counters) is returned.
	SearchContents func(packet nex.PacketInterface, query user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], *nex.Error)

	// SearchContentsWithTotal is optional: when nil an empty result (empty lists, zero counters) is returned.
	SearchContentsWithTotal func(packet nex.PacketInterface, query user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error)

	// DeleteContent is required: when nil the request is answered with Core::NotImplemented.
	DeleteContent func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) *nex.Error

	// SaveMetaData is required: when nil the request is answered with Core::NotImplemented.
	SaveMetaData func(packet nex.PacketInterface, properties types.List[user_storage_types.ContentProperty], contentKey user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error)

	// SaveContentDB is required: when nil the request is answered with Core::NotImplemented.
	SaveContentDB func(packet nex.PacketInterface, properties types.List[user_storage_types.ContentProperty], data types.Buffer, contentKey user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error)

	// SaveContentAndGetUploadInfo is required: when nil the request is answered with Core::NotImplemented.
	SaveContentAndGetUploadInfo func(packet nex.PacketInterface, properties types.List[user_storage_types.ContentProperty], size types.UInt32, contentKey user_storage_types.UserContentKey) (user_storage_types.UserContentURL, types.UInt64, types.List[types.String], *nex.Error)

	// UploadEnd is required: when nil the request is answered with Core::NotImplemented.
	UploadEnd func(packet nex.PacketInterface, pendingID types.UInt64, result types.Bool, contentKey user_storage_types.UserContentKey) (user_storage_types.UserContentKey, *nex.Error)

	// GetContentDB is required: when nil the request is answered with Core::NotImplemented.
	GetContentDB func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) (types.Buffer, *nex.Error)

	// GetContentURL is required: when nil the request is answered with Core::NotImplemented.
	GetContentURL func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) (user_storage_types.UserContentURL, *nex.Error)

	// GetSlotCount is required: when nil the request is answered with Core::NotImplemented.
	GetSlotCount func(packet nex.PacketInterface, typeID types.UInt32) (user_storage_types.UserSlotCount, *nex.Error)

	// GetMetaData is required: when nil the request is answered with Core::NotImplemented.
	GetMetaData func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) (user_storage_types.UserContent, *nex.Error)

	// Like is required: when nil the request is answered with Core::NotImplemented.
	Like func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) *nex.Error

	// Unlike is required: when nil the request is answered with Core::NotImplemented.
	Unlike func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) *nex.Error

	// IsLiked is required: when nil the request is answered with Core::NotImplemented.
	IsLiked func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) (types.Bool, *nex.Error)

	// GetFavourites is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetFavourites func(packet nex.PacketInterface, contentTypes types.List[types.UInt32]) (types.List[user_storage_types.UserContent], *nex.Error)

	// MakeFavourite is required: when nil the request is answered with Core::NotImplemented.
	MakeFavourite func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) *nex.Error

	// RemoveFromFavourites is required: when nil the request is answered with Core::NotImplemented.
	RemoveFromFavourites func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) *nex.Error

	// ReportInappropriate is required: when nil the request is answered with Core::NotImplemented.
	ReportInappropriate func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey, reason types.String) *nex.Error

	// IncrementPlayCount is required: when nil the request is answered with Core::NotImplemented.
	IncrementPlayCount func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) *nex.Error

	// UpdateCustomStat is required: when nil the request is answered with Core::NotImplemented.
	UpdateCustomStat func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey, statID types.UInt16, incValue types.Int64) *nex.Error

	// GetOwnContents is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetOwnContents func(packet nex.PacketInterface, typeID types.UInt32) (types.List[user_storage_types.UserContent], *nex.Error)

	// GetMostPopularTags is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetMostPopularTags func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) (types.List[user_storage_types.WeightedTag], types.UInt32, *nex.Error)

	// GetTags is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetTags func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) (types.List[types.UInt32], *nex.Error)

	// TagContent is required: when nil the request is answered with Core::NotImplemented.
	TagContent func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey, newTagIds types.List[types.UInt32]) *nex.Error

	// SearchContentsByPlayers is optional: when nil an empty result (empty lists, zero counters) is returned.
	SearchContentsByPlayers func(packet nex.PacketInterface, pids types.List[types.UInt32], query user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], *nex.Error)

	// SearchContentsByPlayersWithTotal is optional: when nil an empty result (empty lists, zero counters) is returned.
	SearchContentsByPlayersWithTotal func(packet nex.PacketInterface, pids types.List[types.UInt32], query user_storage_types.UserStorageQuery) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerSearchContents(commonProtocol.searchContents)
	protocol.SetHandlerSearchContentsWithTotal(commonProtocol.searchContentsWithTotal)
	protocol.SetHandlerDeleteContent(commonProtocol.deleteContent)
	protocol.SetHandlerSaveMetaData(commonProtocol.saveMetaData)
	protocol.SetHandlerSaveContentDB(commonProtocol.saveContentDB)
	protocol.SetHandlerSaveContentAndGetUploadInfo(commonProtocol.saveContentAndGetUploadInfo)
	protocol.SetHandlerUploadEnd(commonProtocol.uploadEnd)
	protocol.SetHandlerGetContentDB(commonProtocol.getContentDB)
	protocol.SetHandlerGetContentURL(commonProtocol.getContentURL)
	protocol.SetHandlerGetSlotCount(commonProtocol.getSlotCount)
	protocol.SetHandlerGetMetaData(commonProtocol.getMetaData)
	protocol.SetHandlerLike(commonProtocol.like)
	protocol.SetHandlerUnlike(commonProtocol.unlike)
	protocol.SetHandlerIsLiked(commonProtocol.isLiked)
	protocol.SetHandlerGetFavourites(commonProtocol.getFavourites)
	protocol.SetHandlerMakeFavourite(commonProtocol.makeFavourite)
	protocol.SetHandlerRemoveFromFavourites(commonProtocol.removeFromFavourites)
	protocol.SetHandlerReportInappropriate(commonProtocol.reportInappropriate)
	protocol.SetHandlerIncrementPlayCount(commonProtocol.incrementPlayCount)
	protocol.SetHandlerUpdateCustomStat(commonProtocol.updateCustomStat)
	protocol.SetHandlerGetOwnContents(commonProtocol.getOwnContents)
	protocol.SetHandlerGetMostPopularTags(commonProtocol.getMostPopularTags)
	protocol.SetHandlerGetTags(commonProtocol.getTags)
	protocol.SetHandlerTagContent(commonProtocol.tagContent)
	protocol.SetHandlerSearchContentsByPlayers(commonProtocol.searchContentsByPlayers)
	protocol.SetHandlerSearchContentsByPlayersWithTotal(commonProtocol.searchContentsByPlayersWithTotal)

	return commonProtocol
}
