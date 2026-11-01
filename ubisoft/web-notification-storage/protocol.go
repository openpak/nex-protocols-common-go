// Package webnotificationstorage implements the common Web Notifications Storage protocol handlers of the Quazal Rendez-Vous (Ubisoft) family.
//
// Layouts follow the NintendoClients wiki page "Web Notification Storage Protocol" (MIT). Every method calls a user-provided hook;
// a method whose hook is nil answers Core::NotImplemented unless a safe default is documented on its hook field.
//
// Layout provisional:
//   - The wiki lists the protocol ID as unknown ("?"); ProtocolID is 0 until it is known.
package webnotificationstorage

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/web-notification-storage"
)

// CommonProtocol handles the Web Notifications Storage protocol with user-provided hooks
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi.Interface

	// RegisterUser is required: when nil the request is answered with Core::NotImplemented.
	RegisterUser func(packet nex.PacketInterface) *nex.Error

	// PollNotifications is optional: when nil an empty result (empty lists, zero counters) is returned.
	PollNotifications func(packet nex.PacketInterface) (types.String, types.Int32, *nex.Error)

	// UnregisterUser is required: when nil the request is answered with Core::NotImplemented.
	UnregisterUser func(packet nex.PacketInterface) *nex.Error
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerRegisterUser(commonProtocol.registerUser)
	protocol.SetHandlerPollNotifications(commonProtocol.pollNotifications)
	protocol.SetHandlerUnregisterUser(commonProtocol.unregisterUser)

	return commonProtocol
}
