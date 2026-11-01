// Package offlinegamenotifications implements the common Offline Game Notifications protocol handlers of the Quazal Rendez-Vous (Ubisoft) family.
//
// Layouts follow the NintendoClients wiki page "Offline Game Notification Protocol" (MIT). Every method calls a user-provided hook;
// a method whose hook is nil answers Core::NotImplemented unless a safe default is documented on its hook field.
//
// Layout provisional:
//   - NotificationEvent is the nex-protocols-go notifications type; its field widths depend on the endpoint library version.
package offlinegamenotifications

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/notifications/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/offline-game-notifications"
	offline_game_notifications_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/offline-game-notifications/types"
)

// CommonProtocol handles the Offline Game Notifications protocol with user-provided hooks
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi.Interface

	// PollNotifications is optional: when nil an empty result (empty lists, zero counters) is returned.
	PollNotifications func(packet nex.PacketInterface) (types.List[notifications_types.NotificationEvent], types.UInt32, *nex.Error)

	// PollSpecificOfflineNotifications is optional: when nil an empty result (empty lists, zero counters) is returned.
	PollSpecificOfflineNotifications func(packet nex.PacketInterface, majortype types.List[types.UInt32]) (types.List[offline_game_notifications_types.TimedNotification], types.UInt32, *nex.Error)

	// PollAnyOfflineNotifications is optional: when nil an empty result (empty lists, zero counters) is returned.
	PollAnyOfflineNotifications func(packet nex.PacketInterface) (types.List[offline_game_notifications_types.TimedNotification], types.UInt32, *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerPollNotifications(commonProtocol.pollNotifications)
	protocol.SetHandlerPollSpecificOfflineNotifications(commonProtocol.pollSpecificOfflineNotifications)
	protocol.SetHandlerPollAnyOfflineNotifications(commonProtocol.pollAnyOfflineNotifications)

	return commonProtocol
}
