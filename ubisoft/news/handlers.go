package news

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/news"
	news_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/news/types"
)

// respond builds a success RMC message for the given method
func respond(endpoint nex.EndpointInterface, methodID uint32, callID uint32, body []byte) *nex.RMCMessage {
	message := nex.NewRMCSuccess(endpoint, body)
	message.ProtocolID = ubi.ProtocolID
	message.MethodID = methodID
	message.CallID = callID

	return message
}

// encodeGetChannelsResponse encodes the GetChannels response body
func encodeGetChannelsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetChannelsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getChannels(err error, packet nex.PacketInterface, callID uint32, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::GetChannels: invalid request")
	}

	if commonProtocol.GetChannels == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetChannelsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetChannelsResponse{})

		return respond(endpoint, ubi.MethodGetChannels, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetChannels(packet, resultRange)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetChannelsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetChannelsResponse{Channels: r0})

	return respond(endpoint, ubi.MethodGetChannels, callID, body), nil
}

// encodeGetChannelsByTypesResponse encodes the GetChannelsByTypes response body
func encodeGetChannelsByTypesResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetChannelsByTypesResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getChannelsByTypes(err error, packet nex.PacketInterface, callID uint32, newsChannelTypes types.List[types.String], resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::GetChannelsByTypes: invalid request")
	}

	if commonProtocol.GetChannelsByTypes == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetChannelsByTypesResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetChannelsByTypesResponse{})

		return respond(endpoint, ubi.MethodGetChannelsByTypes, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetChannelsByTypes(packet, newsChannelTypes, resultRange)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetChannelsByTypesResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetChannelsByTypesResponse{Channels: r0})

	return respond(endpoint, ubi.MethodGetChannelsByTypes, callID, body), nil
}

// encodeGetSubscribableChannelsResponse encodes the GetSubscribableChannels response body
func encodeGetSubscribableChannelsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetSubscribableChannelsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getSubscribableChannels(err error, packet nex.PacketInterface, callID uint32, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::GetSubscribableChannels: invalid request")
	}

	if commonProtocol.GetSubscribableChannels == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetSubscribableChannelsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetSubscribableChannelsResponse{})

		return respond(endpoint, ubi.MethodGetSubscribableChannels, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetSubscribableChannels(packet, resultRange)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetSubscribableChannelsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetSubscribableChannelsResponse{Channels: r0})

	return respond(endpoint, ubi.MethodGetSubscribableChannels, callID, body), nil
}

// encodeGetChannelsByIDsResponse encodes the GetChannelsByIDs response body
func encodeGetChannelsByIDsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetChannelsByIDsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getChannelsByIDs(err error, packet nex.PacketInterface, callID uint32, newsChannelIDs types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::GetChannelsByIDs: invalid request")
	}

	if commonProtocol.GetChannelsByIDs == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetChannelsByIDsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetChannelsByIDsResponse{})

		return respond(endpoint, ubi.MethodGetChannelsByIDs, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetChannelsByIDs(packet, newsChannelIDs)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetChannelsByIDsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetChannelsByIDsResponse{Channels: r0})

	return respond(endpoint, ubi.MethodGetChannelsByIDs, callID, body), nil
}

// encodeGetSubscribedChannelsResponse encodes the GetSubscribedChannels response body
func encodeGetSubscribedChannelsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetSubscribedChannelsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getSubscribedChannels(err error, packet nex.PacketInterface, callID uint32, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::GetSubscribedChannels: invalid request")
	}

	if commonProtocol.GetSubscribedChannels == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetSubscribedChannelsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetSubscribedChannelsResponse{})

		return respond(endpoint, ubi.MethodGetSubscribedChannels, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetSubscribedChannels(packet, resultRange)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetSubscribedChannelsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetSubscribedChannelsResponse{Channels: r0})

	return respond(endpoint, ubi.MethodGetSubscribedChannels, callID, body), nil
}

func (commonProtocol *CommonProtocol) subscribeChannel(err error, packet nex.PacketInterface, callID uint32, newsChannelID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::SubscribeChannel: invalid request")
	}

	if commonProtocol.SubscribeChannel == nil {
		common_globals.Logger.Warning("News::SubscribeChannel not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "News::SubscribeChannel not implemented")
	}

	nexErr := commonProtocol.SubscribeChannel(packet, newsChannelID)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodSubscribeChannel, callID, nil), nil
}

func (commonProtocol *CommonProtocol) unsubscribeChannel(err error, packet nex.PacketInterface, callID uint32, newsChannelID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::UnsubscribeChannel: invalid request")
	}

	if commonProtocol.UnsubscribeChannel == nil {
		common_globals.Logger.Warning("News::UnsubscribeChannel not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "News::UnsubscribeChannel not implemented")
	}

	nexErr := commonProtocol.UnsubscribeChannel(packet, newsChannelID)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodUnsubscribeChannel, callID, nil), nil
}

// encodeGetNewsHeadersResponse encodes the GetNewsHeaders response body
func encodeGetNewsHeadersResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetNewsHeadersResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getNewsHeaders(err error, packet nex.PacketInterface, callID uint32, recipient news_types.NewsRecipient, rangeParam types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::GetNewsHeaders: invalid request")
	}

	if commonProtocol.GetNewsHeaders == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetNewsHeadersResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetNewsHeadersResponse{})

		return respond(endpoint, ubi.MethodGetNewsHeaders, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetNewsHeaders(packet, recipient, rangeParam)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetNewsHeadersResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetNewsHeadersResponse{NewsHeaders: r0})

	return respond(endpoint, ubi.MethodGetNewsHeaders, callID, body), nil
}

// encodeGetNewsMessagesResponse encodes the GetNewsMessages response body
func encodeGetNewsMessagesResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetNewsMessagesResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getNewsMessages(err error, packet nex.PacketInterface, callID uint32, newsMessageIDs types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::GetNewsMessages: invalid request")
	}

	if commonProtocol.GetNewsMessages == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetNewsMessagesResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetNewsMessagesResponse{})

		return respond(endpoint, ubi.MethodGetNewsMessages, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetNewsMessages(packet, newsMessageIDs)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetNewsMessagesResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetNewsMessagesResponse{NewsMessages: r0})

	return respond(endpoint, ubi.MethodGetNewsMessages, callID, body), nil
}

// encodeGetNumberOfNewsResponse encodes the GetNumberOfNews response body
func encodeGetNumberOfNewsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetNumberOfNewsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getNumberOfNews(err error, packet nex.PacketInterface, callID uint32, recipient news_types.NewsRecipient) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::GetNumberOfNews: invalid request")
	}

	if commonProtocol.GetNumberOfNews == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetNumberOfNewsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetNumberOfNewsResponse{})

		return respond(endpoint, ubi.MethodGetNumberOfNews, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetNumberOfNews(packet, recipient)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetNumberOfNewsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetNumberOfNewsResponse{NumberOfNews: r0})

	return respond(endpoint, ubi.MethodGetNumberOfNews, callID, body), nil
}

// encodeGetChannelByTypeResponse encodes the GetChannelByType response body
func encodeGetChannelByTypeResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetChannelByTypeResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getChannelByType(err error, packet nex.PacketInterface, callID uint32, newsChannelType types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::GetChannelByType: invalid request")
	}

	if commonProtocol.GetChannelByType == nil {
		common_globals.Logger.Warning("News::GetChannelByType not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "News::GetChannelByType not implemented")
	}

	r0, nexErr := commonProtocol.GetChannelByType(packet, newsChannelType)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetChannelByTypeResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetChannelByTypeResponse{Channel: r0})

	return respond(endpoint, ubi.MethodGetChannelByType, callID, body), nil
}

// encodeGetNewsHeadersByTypeResponse encodes the GetNewsHeadersByType response body
func encodeGetNewsHeadersByTypeResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetNewsHeadersByTypeResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getNewsHeadersByType(err error, packet nex.PacketInterface, callID uint32, newsChannelType types.String, rangeParam types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::GetNewsHeadersByType: invalid request")
	}

	if commonProtocol.GetNewsHeadersByType == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetNewsHeadersByTypeResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetNewsHeadersByTypeResponse{})

		return respond(endpoint, ubi.MethodGetNewsHeadersByType, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetNewsHeadersByType(packet, newsChannelType, rangeParam)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetNewsHeadersByTypeResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetNewsHeadersByTypeResponse{NewsHeaders: r0})

	return respond(endpoint, ubi.MethodGetNewsHeadersByType, callID, body), nil
}

// encodeGetNewsMessagesByTypeResponse encodes the GetNewsMessagesByType response body
func encodeGetNewsMessagesByTypeResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetNewsMessagesByTypeResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getNewsMessagesByType(err error, packet nex.PacketInterface, callID uint32, newsChannelType types.String, rangeParam types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "News::GetNewsMessagesByType: invalid request")
	}

	if commonProtocol.GetNewsMessagesByType == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetNewsMessagesByTypeResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetNewsMessagesByTypeResponse{})

		return respond(endpoint, ubi.MethodGetNewsMessagesByType, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetNewsMessagesByType(packet, newsChannelType, rangeParam)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetNewsMessagesByTypeResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetNewsMessagesByTypeResponse{NewsMessages: r0})

	return respond(endpoint, ubi.MethodGetNewsMessagesByType, callID, body), nil
}
