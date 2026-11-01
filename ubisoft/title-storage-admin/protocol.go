// Package titlestorageadmin implements the common Title Storage Admin protocol handlers of the Quazal Rendez-Vous (Ubisoft) family.
//
// Layouts follow the NintendoClients wiki page "Title Storage Admin Protocol" (MIT). Every method calls a user-provided hook;
// a method whose hook is nil answers Core::NotImplemented unless a safe default is documented on its hook field.
//
// Layout provisional:
//   - The wiki gives no protocol ID for this protocol; ProtocolID is 0 until it is known.
package titlestorageadmin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/title-storage-admin"
	title_storage_admin_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/title-storage-admin/types"
)

// CommonProtocol handles the Title Storage Admin protocol with user-provided hooks
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi.Interface

	// GetContentTypesDescription is required: when nil the request is answered with Core::NotImplemented.
	GetContentTypesDescription func(packet nex.PacketInterface) (types.List[title_storage_admin_types.ContentTypeDescription], *nex.Error)

	// GetContentTypeDescription is required: when nil the request is answered with Core::NotImplemented.
	GetContentTypeDescription func(packet nex.PacketInterface, contentType types.String) (types.List[title_storage_admin_types.ContentTypeExtensionField], *nex.Error)

	// BrowseContents is required: when nil the request is answered with Core::NotImplemented.
	BrowseContents func(packet nex.PacketInterface) (types.List[title_storage_admin_types.AdminTitleContent], *nex.Error)

	// GetContentDetails is required: when nil the request is answered with Core::NotImplemented.
	GetContentDetails func(packet nex.PacketInterface, contentId types.UInt32) (types.DataHolder, *nex.Error)

	// CreateContentAndGetUploadInformation is required: when nil the request is answered with Core::NotImplemented.
	CreateContentAndGetUploadInformation func(packet nex.PacketInterface, name types.String, description types.String, key types.String, dateAvailable types.DateTime, dateExpired types.DateTime, flags types.UInt32, extendedFlags types.UInt32, contentType types.String, extension types.DataHolder, mime types.String) (types.UInt32, types.String, types.String, types.Map[types.String, types.String], *nex.Error)

	// RemoveContent is required: when nil the request is answered with Core::NotImplemented.
	RemoveContent func(packet nex.PacketInterface, contentId types.UInt32) (types.Bool, *nex.Error)

	// UpdateContentMetaData is required: when nil the request is answered with Core::NotImplemented.
	UpdateContentMetaData func(packet nex.PacketInterface, contentId types.UInt32, name types.String, description types.String, dateAvailable types.DateTime, dateExpired types.DateTime, flags types.UInt32, extendedFlags types.UInt32, extension types.DataHolder) (types.Bool, *nex.Error)

	// ContentInsertionCompleted is required: when nil the request is answered with Core::NotImplemented.
	ContentInsertionCompleted func(packet nex.PacketInterface, contentId types.UInt32, contentSize types.UInt32) (types.Bool, *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerGetContentTypesDescription(commonProtocol.getContentTypesDescription)
	protocol.SetHandlerGetContentTypeDescription(commonProtocol.getContentTypeDescription)
	protocol.SetHandlerBrowseContents(commonProtocol.browseContents)
	protocol.SetHandlerGetContentDetails(commonProtocol.getContentDetails)
	protocol.SetHandlerCreateContentAndGetUploadInformation(commonProtocol.createContentAndGetUploadInformation)
	protocol.SetHandlerRemoveContent(commonProtocol.removeContent)
	protocol.SetHandlerUpdateContentMetaData(commonProtocol.updateContentMetaData)
	protocol.SetHandlerContentInsertionCompleted(commonProtocol.contentInsertionCompleted)

	return commonProtocol
}
