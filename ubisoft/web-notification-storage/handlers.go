package webnotificationstorage

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/web-notification-storage"
)

// respond builds a success RMC message for the given method
func respond(endpoint nex.EndpointInterface, methodID uint32, callID uint32, body []byte) *nex.RMCMessage {
	message := nex.NewRMCSuccess(endpoint, body)
	message.ProtocolID = ubi.ProtocolID
	message.MethodID = methodID
	message.CallID = callID

	return message
}

func (commonProtocol *CommonProtocol) registerUser(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Web Notifications Storage::RegisterUser: invalid request")
	}

	if commonProtocol.RegisterUser == nil {
		common_globals.Logger.Warning("Web Notifications Storage::RegisterUser not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Web Notifications Storage::RegisterUser not implemented")
	}

	nexErr := commonProtocol.RegisterUser(packet)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodRegisterUser, callID, nil), nil
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

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Web Notifications Storage::PollNotifications: invalid request")
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

	body := encodePollNotificationsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.PollNotificationsResponse{ListNotifications: r0, NbNotifications: r1})

	return respond(endpoint, ubi.MethodPollNotifications, callID, body), nil
}

func (commonProtocol *CommonProtocol) unregisterUser(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Web Notifications Storage::UnregisterUser: invalid request")
	}

	if commonProtocol.UnregisterUser == nil {
		common_globals.Logger.Warning("Web Notifications Storage::UnregisterUser not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Web Notifications Storage::UnregisterUser not implemented")
	}

	nexErr := commonProtocol.UnregisterUser(packet)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodUnregisterUser, callID, nil), nil
}
