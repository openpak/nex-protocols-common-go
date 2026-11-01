// Package news implements the common News protocol handlers of the Quazal Rendez-Vous (Ubisoft) family.
//
// Layouts follow the NintendoClients wiki page "News Protocol" (MIT). Every method calls a user-provided hook;
// a method whose hook is nil answers Core::NotImplemented unless a safe default is documented on its hook field.
package news

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/news"
	news_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/news/types"
)

// CommonProtocol handles the News protocol with user-provided hooks
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi.Interface

	// GetChannels is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetChannels func(packet nex.PacketInterface, resultRange types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error)

	// GetChannelsByTypes is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetChannelsByTypes func(packet nex.PacketInterface, newsChannelTypes types.List[types.String], resultRange types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error)

	// GetSubscribableChannels is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetSubscribableChannels func(packet nex.PacketInterface, resultRange types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error)

	// GetChannelsByIDs is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetChannelsByIDs func(packet nex.PacketInterface, newsChannelIDs types.List[types.UInt32]) (types.List[news_types.NewsChannel], *nex.Error)

	// GetSubscribedChannels is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetSubscribedChannels func(packet nex.PacketInterface, resultRange types.ResultRange) (types.List[news_types.NewsChannel], *nex.Error)

	// SubscribeChannel is required: when nil the request is answered with Core::NotImplemented.
	SubscribeChannel func(packet nex.PacketInterface, newsChannelID types.UInt32) *nex.Error

	// UnsubscribeChannel is required: when nil the request is answered with Core::NotImplemented.
	UnsubscribeChannel func(packet nex.PacketInterface, newsChannelID types.UInt32) *nex.Error

	// GetNewsHeaders is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetNewsHeaders func(packet nex.PacketInterface, recipient news_types.NewsRecipient, rangeParam types.ResultRange) (types.List[news_types.NewsHeader], *nex.Error)

	// GetNewsMessages is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetNewsMessages func(packet nex.PacketInterface, newsMessageIDs types.List[types.UInt32]) (types.List[news_types.NewsMessage], *nex.Error)

	// GetNumberOfNews is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetNumberOfNews func(packet nex.PacketInterface, recipient news_types.NewsRecipient) (types.UInt32, *nex.Error)

	// GetChannelByType is required: when nil the request is answered with Core::NotImplemented.
	GetChannelByType func(packet nex.PacketInterface, newsChannelType types.String) (news_types.NewsChannel, *nex.Error)

	// GetNewsHeadersByType is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetNewsHeadersByType func(packet nex.PacketInterface, newsChannelType types.String, rangeParam types.ResultRange) (types.List[news_types.NewsHeader], *nex.Error)

	// GetNewsMessagesByType is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetNewsMessagesByType func(packet nex.PacketInterface, newsChannelType types.String, rangeParam types.ResultRange) (types.List[news_types.NewsMessage], *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerGetChannels(commonProtocol.getChannels)
	protocol.SetHandlerGetChannelsByTypes(commonProtocol.getChannelsByTypes)
	protocol.SetHandlerGetSubscribableChannels(commonProtocol.getSubscribableChannels)
	protocol.SetHandlerGetChannelsByIDs(commonProtocol.getChannelsByIDs)
	protocol.SetHandlerGetSubscribedChannels(commonProtocol.getSubscribedChannels)
	protocol.SetHandlerSubscribeChannel(commonProtocol.subscribeChannel)
	protocol.SetHandlerUnsubscribeChannel(commonProtocol.unsubscribeChannel)
	protocol.SetHandlerGetNewsHeaders(commonProtocol.getNewsHeaders)
	protocol.SetHandlerGetNewsMessages(commonProtocol.getNewsMessages)
	protocol.SetHandlerGetNumberOfNews(commonProtocol.getNumberOfNews)
	protocol.SetHandlerGetChannelByType(commonProtocol.getChannelByType)
	protocol.SetHandlerGetNewsHeadersByType(commonProtocol.getNewsHeadersByType)
	protocol.SetHandlerGetNewsMessagesByType(commonProtocol.getNewsMessagesByType)

	return commonProtocol
}
