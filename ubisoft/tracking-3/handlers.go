package tracking3

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/tracking-3"
	tracking_3_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/tracking-3/types"
)

// respond builds a success RMC message for the given method
func respond(endpoint nex.EndpointInterface, methodID uint32, callID uint32, body []byte) *nex.RMCMessage {
	message := nex.NewRMCSuccess(endpoint, body)
	message.ProtocolID = ubi.ProtocolID
	message.MethodID = methodID
	message.CallID = callID

	return message
}

func (commonProtocol *CommonProtocol) sendTag(err error, packet nex.PacketInterface, callID uint32, trackingID types.UInt32, tag types.String, attributes types.String, deltaTime types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Tracking 3::SendTag: invalid request")
	}

	if commonProtocol.SendTag == nil {
		endpoint := packet.Sender().Endpoint()

		return respond(endpoint, ubi.MethodSendTag, callID, nil), nil
	}

	nexErr := commonProtocol.SendTag(packet, trackingID, tag, attributes, deltaTime)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodSendTag, callID, nil), nil
}

func (commonProtocol *CommonProtocol) sendTagAndUpdateUserInfo(err error, packet nex.PacketInterface, callID uint32, trackingID types.UInt32, tag types.String, attributes types.String, deltaTime types.UInt32, userID types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Tracking 3::SendTagAndUpdateUserInfo: invalid request")
	}

	if commonProtocol.SendTagAndUpdateUserInfo == nil {
		endpoint := packet.Sender().Endpoint()

		return respond(endpoint, ubi.MethodSendTagAndUpdateUserInfo, callID, nil), nil
	}

	nexErr := commonProtocol.SendTagAndUpdateUserInfo(packet, trackingID, tag, attributes, deltaTime, userID)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodSendTagAndUpdateUserInfo, callID, nil), nil
}

// encodeSendUserInfoResponse encodes the SendUserInfo response body
func encodeSendUserInfoResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SendUserInfoResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) sendUserInfo(err error, packet nex.PacketInterface, callID uint32, userInfo tracking_3_types.TrackingInformation, deltaTime types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Tracking 3::SendUserInfo: invalid request")
	}

	if commonProtocol.SendUserInfo == nil {
		common_globals.Logger.Warning("Tracking 3::SendUserInfo not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Tracking 3::SendUserInfo not implemented")
	}

	r0, r1, nexErr := commonProtocol.SendUserInfo(packet, userInfo, deltaTime)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSendUserInfoResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SendUserInfoResponse{UserInfo: r0, TrackingID: r1})

	return respond(endpoint, ubi.MethodSendUserInfo, callID, body), nil
}

// encodeGetConfigurationResponse encodes the GetConfiguration response body
func encodeGetConfigurationResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetConfigurationResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getConfiguration(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Tracking 3::GetConfiguration: invalid request")
	}

	if commonProtocol.GetConfiguration == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetConfigurationResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetConfigurationResponse{})

		return respond(endpoint, ubi.MethodGetConfiguration, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetConfiguration(packet)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetConfigurationResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetConfigurationResponse{Tags: r0})

	return respond(endpoint, ubi.MethodGetConfiguration, callID, body), nil
}

func (commonProtocol *CommonProtocol) sendTags(err error, packet nex.PacketInterface, callID uint32, tagData types.List[tracking_3_types.TrackingTag]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Tracking 3::SendTags: invalid request")
	}

	if commonProtocol.SendTags == nil {
		endpoint := packet.Sender().Endpoint()

		return respond(endpoint, ubi.MethodSendTags, callID, nil), nil
	}

	nexErr := commonProtocol.SendTags(packet, tagData)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodSendTags, callID, nil), nil
}
