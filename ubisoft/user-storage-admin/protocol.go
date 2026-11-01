// Package userstorageadmin implements the common User Storage Admin protocol handlers of the Quazal Rendez-Vous (Ubisoft) family.
//
// Layouts follow the NintendoClients wiki page "User Storage Admin Protocol" (MIT). Every method calls a user-provided hook;
// a method whose hook is nil answers Core::NotImplemented unless a safe default is documented on its hook field.
//
// Layout provisional:
//   - The wiki gives no protocol ID for this protocol; ProtocolID is 0 until it is known.
//   - UserContent, UserContentKey and ContentProperty are documented on the (unavailable) User Content page; the User Storage page layouts are reused.
package userstorageadmin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-storage-admin"
	user_storage_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-storage/types"
)

// CommonProtocol handles the User Storage Admin protocol with user-provided hooks
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi.Interface

	// GetContentsToModerate is required: when nil the request is answered with Core::NotImplemented.
	GetContentsToModerate func(packet nex.PacketInterface, typeID types.UInt32, offset types.UInt32, size types.UInt32) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error)

	// FlagContentAsVerified is required: when nil the request is answered with Core::NotImplemented.
	FlagContentAsVerified func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) *nex.Error

	// BanContent is required: when nil the request is answered with Core::NotImplemented.
	BanContent func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey) *nex.Error

	// BanUser is required: when nil the request is answered with Core::NotImplemented.
	BanUser func(packet nex.PacketInterface, pid types.UInt32, reason types.String, banContents types.Bool, expireDate types.DateTime) *nex.Error

	// BanUserFromContentType is required: when nil the request is answered with Core::NotImplemented.
	BanUserFromContentType func(packet nex.PacketInterface, typeID types.UInt32, pid types.UInt32, reason types.String, banContents types.Bool, expireDate types.DateTime) *nex.Error

	// UnbanUser is required: when nil the request is answered with Core::NotImplemented.
	UnbanUser func(packet nex.PacketInterface, pid types.UInt32) *nex.Error

	// UnbanUserFromContentType is required: when nil the request is answered with Core::NotImplemented.
	UnbanUserFromContentType func(packet nex.PacketInterface, typeID types.UInt32, pid types.UInt32) *nex.Error

	// GetContentsToModerateWithThreshold is required: when nil the request is answered with Core::NotImplemented.
	GetContentsToModerateWithThreshold func(packet nex.PacketInterface, typeID types.UInt32, threshold types.UInt32, offset types.UInt32, size types.UInt32) (types.List[user_storage_types.UserContent], types.UInt32, *nex.Error)

	// UpdateMetaData is required: when nil the request is answered with Core::NotImplemented.
	UpdateMetaData func(packet nex.PacketInterface, contentKey user_storage_types.UserContentKey, properties types.List[user_storage_types.ContentProperty]) *nex.Error
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerGetContentsToModerate(commonProtocol.getContentsToModerate)
	protocol.SetHandlerFlagContentAsVerified(commonProtocol.flagContentAsVerified)
	protocol.SetHandlerBanContent(commonProtocol.banContent)
	protocol.SetHandlerBanUser(commonProtocol.banUser)
	protocol.SetHandlerBanUserFromContentType(commonProtocol.banUserFromContentType)
	protocol.SetHandlerUnbanUser(commonProtocol.unbanUser)
	protocol.SetHandlerUnbanUserFromContentType(commonProtocol.unbanUserFromContentType)
	protocol.SetHandlerGetContentsToModerateWithThreshold(commonProtocol.getContentsToModerateWithThreshold)
	protocol.SetHandlerUpdateMetaData(commonProtocol.updateMetaData)

	return commonProtocol
}
