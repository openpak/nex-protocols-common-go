package titlestorage

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/title-storage"
)

// respond builds a success RMC message for the given method
func respond(endpoint nex.EndpointInterface, methodID uint32, callID uint32, body []byte) *nex.RMCMessage {
	message := nex.NewRMCSuccess(endpoint, body)
	message.ProtocolID = ubi.ProtocolID
	message.MethodID = methodID
	message.CallID = callID

	return message
}

// encodeBrowseContentsResponse encodes the BrowseContents response body
func encodeBrowseContentsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.BrowseContentsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) browseContents(err error, packet nex.PacketInterface, callID uint32, flags types.UInt32, flagsMaskType types.UInt8, extendedFlags types.UInt32, extendedFlagsMaskType types.UInt8) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage::BrowseContents: invalid request")
	}

	if commonProtocol.BrowseContents == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeBrowseContentsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.BrowseContentsResponse{})

		return respond(endpoint, ubi.MethodBrowseContents, callID, body), nil
	}

	r0, nexErr := commonProtocol.BrowseContents(packet, flags, flagsMaskType, extendedFlags, extendedFlagsMaskType)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeBrowseContentsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.BrowseContentsResponse{Contents: r0})

	return respond(endpoint, ubi.MethodBrowseContents, callID, body), nil
}

// encodeBrowseContentsByTypeResponse encodes the BrowseContentsByType response body
func encodeBrowseContentsByTypeResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.BrowseContentsByTypeResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) browseContentsByType(err error, packet nex.PacketInterface, callID uint32, contentType types.String, flags types.UInt32, flagsMaskType types.UInt8, extendedFlags types.UInt32, extendedFlagsMaskType types.UInt8) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage::BrowseContentsByType: invalid request")
	}

	if commonProtocol.BrowseContentsByType == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeBrowseContentsByTypeResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.BrowseContentsByTypeResponse{})

		return respond(endpoint, ubi.MethodBrowseContentsByType, callID, body), nil
	}

	r0, nexErr := commonProtocol.BrowseContentsByType(packet, contentType, flags, flagsMaskType, extendedFlags, extendedFlagsMaskType)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeBrowseContentsByTypeResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.BrowseContentsByTypeResponse{Contents: r0})

	return respond(endpoint, ubi.MethodBrowseContentsByType, callID, body), nil
}

// encodeGetContentDetailsResponse encodes the GetContentDetails response body
func encodeGetContentDetailsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetContentDetailsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getContentDetails(err error, packet nex.PacketInterface, callID uint32, contentID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage::GetContentDetails: invalid request")
	}

	if commonProtocol.GetContentDetails == nil {
		common_globals.Logger.Warning("Title Storage::GetContentDetails not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Title Storage::GetContentDetails not implemented")
	}

	r0, nexErr := commonProtocol.GetContentDetails(packet, contentID)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetContentDetailsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetContentDetailsResponse{Details: r0})

	return respond(endpoint, ubi.MethodGetContentDetails, callID, body), nil
}

// encodeDownloadContentResponse encodes the DownloadContent response body
func encodeDownloadContentResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.DownloadContentResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) downloadContent(err error, packet nex.PacketInterface, callID uint32, contentID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage::DownloadContent: invalid request")
	}

	if commonProtocol.DownloadContent == nil {
		common_globals.Logger.Warning("Title Storage::DownloadContent not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Title Storage::DownloadContent not implemented")
	}

	r0, r1, r2, nexErr := commonProtocol.DownloadContent(packet, contentID)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeDownloadContentResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.DownloadContentResponse{Proto: r0, Host: r1, Path: r2})

	return respond(endpoint, ubi.MethodDownloadContent, callID, body), nil
}
