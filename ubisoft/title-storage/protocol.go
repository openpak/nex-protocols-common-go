// Package titlestorage implements the common Title Storage protocol handlers of the Quazal Rendez-Vous (Ubisoft) family.
//
// Layouts follow the NintendoClients wiki page "Title Storage Protocol" (MIT). Every method calls a user-provided hook;
// a method whose hook is nil answers Core::NotImplemented unless a safe default is documented on its hook field.
//
// Layout provisional:
//   - BrowseContentsByType returns List<any> on the wiki; modelled as a list of AnyDataHolder (DataHolder).
package titlestorage

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/title-storage"
	title_storage_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/title-storage/types"
)

// CommonProtocol handles the Title Storage protocol with user-provided hooks
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi.Interface

	// BrowseContents is optional: when nil an empty result (empty lists, zero counters) is returned.
	BrowseContents func(packet nex.PacketInterface, flags types.UInt32, flagsMaskType types.UInt8, extendedFlags types.UInt32, extendedFlagsMaskType types.UInt8) (types.List[title_storage_types.TitleContent], *nex.Error)

	// BrowseContentsByType is optional: when nil an empty result (empty lists, zero counters) is returned.
	BrowseContentsByType func(packet nex.PacketInterface, contentType types.String, flags types.UInt32, flagsMaskType types.UInt8, extendedFlags types.UInt32, extendedFlagsMaskType types.UInt8) (types.List[types.DataHolder], *nex.Error)

	// GetContentDetails is required: when nil the request is answered with Core::NotImplemented.
	GetContentDetails func(packet nex.PacketInterface, contentID types.UInt32) (types.DataHolder, *nex.Error)

	// DownloadContent is required: when nil the request is answered with Core::NotImplemented.
	DownloadContent func(packet nex.PacketInterface, contentID types.UInt32) (types.String, types.String, types.String, *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerBrowseContents(commonProtocol.browseContents)
	protocol.SetHandlerBrowseContentsByType(commonProtocol.browseContentsByType)
	protocol.SetHandlerGetContentDetails(commonProtocol.getContentDetails)
	protocol.SetHandlerDownloadContent(commonProtocol.downloadContent)

	return commonProtocol
}
