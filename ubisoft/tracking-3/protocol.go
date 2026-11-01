// Package tracking3 implements the common Tracking 3 protocol handlers of the Quazal Rendez-Vous (Ubisoft) family.
//
// Layouts follow the NintendoClients wiki page "Tracking Protocol 3" (MIT). Every method calls a user-provided hook;
// a method whose hook is nil answers Core::NotImplemented unless a safe default is documented on its hook field.
package tracking3

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/tracking-3"
	tracking_3_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/tracking-3/types"
)

// CommonProtocol handles the Tracking 3 protocol with user-provided hooks
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi.Interface

	// SendTag is optional: when nil the request is accepted and discarded.
	SendTag func(packet nex.PacketInterface, trackingID types.UInt32, tag types.String, attributes types.String, deltaTime types.UInt32) *nex.Error

	// SendTagAndUpdateUserInfo is optional: when nil the request is accepted and discarded.
	SendTagAndUpdateUserInfo func(packet nex.PacketInterface, trackingID types.UInt32, tag types.String, attributes types.String, deltaTime types.UInt32, userID types.String) *nex.Error

	// SendUserInfo is required: when nil the request is answered with Core::NotImplemented.
	SendUserInfo func(packet nex.PacketInterface, userInfo tracking_3_types.TrackingInformation, deltaTime types.UInt32) (tracking_3_types.TrackingInformation, types.UInt32, *nex.Error)

	// GetConfiguration is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetConfiguration func(packet nex.PacketInterface) (types.List[types.String], *nex.Error)

	// SendTags is optional: when nil the request is accepted and discarded.
	SendTags func(packet nex.PacketInterface, tagData types.List[tracking_3_types.TrackingTag]) *nex.Error
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerSendTag(commonProtocol.sendTag)
	protocol.SetHandlerSendTagAndUpdateUserInfo(commonProtocol.sendTagAndUpdateUserInfo)
	protocol.SetHandlerSendUserInfo(commonProtocol.sendUserInfo)
	protocol.SetHandlerGetConfiguration(commonProtocol.getConfiguration)
	protocol.SetHandlerSendTags(commonProtocol.sendTags)

	return commonProtocol
}
