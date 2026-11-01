package offlinegamenotifications

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/offline-game-notifications"
)

// respond builds a success RMC message for the given method
func respond(endpoint nex.EndpointInterface, methodID uint32, callID uint32, body []byte) *nex.RMCMessage {
	message := nex.NewRMCSuccess(endpoint, body)
	message.ProtocolID = ubi.ProtocolID
	message.MethodID = methodID
	message.CallID = callID

	return message
}

// encodePollNotificationsResponse encodes the PollNotifications response body
func encodePollNotificationsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.PollNotificationsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) pollNotifications(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Offline Game Notifications::PollNotifications: invalid request")
	}

	if commonProtocol.PollNotifications == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodePollNotificationsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.PollNotificationsResponse{})

		return respond(endpoint, ubi.MethodPollNotifications, callID, body), nil
	}

	r0, r1, nexErr := commonProtocol.PollNotifications(packet)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodePollNotificationsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.PollNotificationsResponse{ListNotifications: r0, NbRemainingNotifs: r1})

	return respond(endpoint, ubi.MethodPollNotifications, callID, body), nil
}

// encodePollSpecificOfflineNotificationsResponse encodes the PollSpecificOfflineNotifications response body
func encodePollSpecificOfflineNotificationsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.PollSpecificOfflineNotificationsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) pollSpecificOfflineNotifications(err error, packet nex.PacketInterface, callID uint32, majortype types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Offline Game Notifications::PollSpecificOfflineNotifications: invalid request")
	}

	if commonProtocol.PollSpecificOfflineNotifications == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodePollSpecificOfflineNotificationsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.PollSpecificOfflineNotificationsResponse{})

		return respond(endpoint, ubi.MethodPollSpecificOfflineNotifications, callID, body), nil
	}

	r0, r1, nexErr := commonProtocol.PollSpecificOfflineNotifications(packet, majortype)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodePollSpecificOfflineNotificationsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.PollSpecificOfflineNotificationsResponse{ListTimedNotification: r0, Ret: r1})

	return respond(endpoint, ubi.MethodPollSpecificOfflineNotifications, callID, body), nil
}

// encodePollAnyOfflineNotificationsResponse encodes the PollAnyOfflineNotifications response body
func encodePollAnyOfflineNotificationsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.PollAnyOfflineNotificationsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) pollAnyOfflineNotifications(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Offline Game Notifications::PollAnyOfflineNotifications: invalid request")
	}

	if commonProtocol.PollAnyOfflineNotifications == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodePollAnyOfflineNotificationsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.PollAnyOfflineNotificationsResponse{})

		return respond(endpoint, ubi.MethodPollAnyOfflineNotifications, callID, body), nil
	}

	r0, r1, nexErr := commonProtocol.PollAnyOfflineNotifications(packet)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodePollAnyOfflineNotificationsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.PollAnyOfflineNotificationsResponse{ListTimedNotification: r0, NbRemainingNotifs: r1})

	return respond(endpoint, ubi.MethodPollAnyOfflineNotifications, callID, body), nil
}
