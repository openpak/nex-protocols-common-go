package titlestorageadmin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/title-storage-admin"
)

// respond builds a success RMC message for the given method
func respond(endpoint nex.EndpointInterface, methodID uint32, callID uint32, body []byte) *nex.RMCMessage {
	message := nex.NewRMCSuccess(endpoint, body)
	message.ProtocolID = ubi.ProtocolID
	message.MethodID = methodID
	message.CallID = callID

	return message
}

// encodeGetContentTypesDescriptionResponse encodes the GetContentTypesDescription response body
func encodeGetContentTypesDescriptionResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetContentTypesDescriptionResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getContentTypesDescription(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage Admin::GetContentTypesDescription: invalid request")
	}

	if commonProtocol.GetContentTypesDescription == nil {
		common_globals.Logger.Warning("Title Storage Admin::GetContentTypesDescription not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Title Storage Admin::GetContentTypesDescription not implemented")
	}

	r0, nexErr := commonProtocol.GetContentTypesDescription(packet)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetContentTypesDescriptionResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetContentTypesDescriptionResponse{ContentDescriptions: r0})

	return respond(endpoint, ubi.MethodGetContentTypesDescription, callID, body), nil
}

// encodeGetContentTypeDescriptionResponse encodes the GetContentTypeDescription response body
func encodeGetContentTypeDescriptionResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetContentTypeDescriptionResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getContentTypeDescription(err error, packet nex.PacketInterface, callID uint32, contentType types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage Admin::GetContentTypeDescription: invalid request")
	}

	if commonProtocol.GetContentTypeDescription == nil {
		common_globals.Logger.Warning("Title Storage Admin::GetContentTypeDescription not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Title Storage Admin::GetContentTypeDescription not implemented")
	}

	r0, nexErr := commonProtocol.GetContentTypeDescription(packet, contentType)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetContentTypeDescriptionResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetContentTypeDescriptionResponse{ExtensionFields: r0})

	return respond(endpoint, ubi.MethodGetContentTypeDescription, callID, body), nil
}

// encodeBrowseContentsResponse encodes the BrowseContents response body
func encodeBrowseContentsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.BrowseContentsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) browseContents(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage Admin::BrowseContents: invalid request")
	}

	if commonProtocol.BrowseContents == nil {
		common_globals.Logger.Warning("Title Storage Admin::BrowseContents not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Title Storage Admin::BrowseContents not implemented")
	}

	r0, nexErr := commonProtocol.BrowseContents(packet)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeBrowseContentsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.BrowseContentsResponse{Contents: r0})

	return respond(endpoint, ubi.MethodBrowseContents, callID, body), nil
}

// encodeGetContentDetailsResponse encodes the GetContentDetails response body
func encodeGetContentDetailsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetContentDetailsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getContentDetails(err error, packet nex.PacketInterface, callID uint32, contentId types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage Admin::GetContentDetails: invalid request")
	}

	if commonProtocol.GetContentDetails == nil {
		common_globals.Logger.Warning("Title Storage Admin::GetContentDetails not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Title Storage Admin::GetContentDetails not implemented")
	}

	r0, nexErr := commonProtocol.GetContentDetails(packet, contentId)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetContentDetailsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetContentDetailsResponse{Details: r0})

	return respond(endpoint, ubi.MethodGetContentDetails, callID, body), nil
}

// encodeCreateContentAndGetUploadInformationResponse encodes the CreateContentAndGetUploadInformation response body
func encodeCreateContentAndGetUploadInformationResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.CreateContentAndGetUploadInformationResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) createContentAndGetUploadInformation(err error, packet nex.PacketInterface, callID uint32, name types.String, description types.String, key types.String, dateAvailable types.DateTime, dateExpired types.DateTime, flags types.UInt32, extendedFlags types.UInt32, contentType types.String, extension types.DataHolder, mime types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage Admin::CreateContentAndGetUploadInformation: invalid request")
	}

	if commonProtocol.CreateContentAndGetUploadInformation == nil {
		common_globals.Logger.Warning("Title Storage Admin::CreateContentAndGetUploadInformation not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Title Storage Admin::CreateContentAndGetUploadInformation not implemented")
	}

	r0, r1, r2, r3, nexErr := commonProtocol.CreateContentAndGetUploadInformation(packet, name, description, key, dateAvailable, dateExpired, flags, extendedFlags, contentType, extension, mime)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeCreateContentAndGetUploadInformationResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.CreateContentAndGetUploadInformationResponse{ContentId: r0, UploadURL: r1, UploadMimeType: r2, UploadBody: r3})

	return respond(endpoint, ubi.MethodCreateContentAndGetUploadInformation, callID, body), nil
}

// encodeRemoveContentResponse encodes the RemoveContent response body
func encodeRemoveContentResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.RemoveContentResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) removeContent(err error, packet nex.PacketInterface, callID uint32, contentId types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage Admin::RemoveContent: invalid request")
	}

	if commonProtocol.RemoveContent == nil {
		common_globals.Logger.Warning("Title Storage Admin::RemoveContent not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Title Storage Admin::RemoveContent not implemented")
	}

	r0, nexErr := commonProtocol.RemoveContent(packet, contentId)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeRemoveContentResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.RemoveContentResponse{Retval: r0})

	return respond(endpoint, ubi.MethodRemoveContent, callID, body), nil
}

// encodeUpdateContentMetaDataResponse encodes the UpdateContentMetaData response body
func encodeUpdateContentMetaDataResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.UpdateContentMetaDataResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) updateContentMetaData(err error, packet nex.PacketInterface, callID uint32, contentId types.UInt32, name types.String, description types.String, dateAvailable types.DateTime, dateExpired types.DateTime, flags types.UInt32, extendedFlags types.UInt32, extension types.DataHolder) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage Admin::UpdateContentMetaData: invalid request")
	}

	if commonProtocol.UpdateContentMetaData == nil {
		common_globals.Logger.Warning("Title Storage Admin::UpdateContentMetaData not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Title Storage Admin::UpdateContentMetaData not implemented")
	}

	r0, nexErr := commonProtocol.UpdateContentMetaData(packet, contentId, name, description, dateAvailable, dateExpired, flags, extendedFlags, extension)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeUpdateContentMetaDataResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.UpdateContentMetaDataResponse{Retval: r0})

	return respond(endpoint, ubi.MethodUpdateContentMetaData, callID, body), nil
}

// encodeContentInsertionCompletedResponse encodes the ContentInsertionCompleted response body
func encodeContentInsertionCompletedResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.ContentInsertionCompletedResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) contentInsertionCompleted(err error, packet nex.PacketInterface, callID uint32, contentId types.UInt32, contentSize types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Title Storage Admin::ContentInsertionCompleted: invalid request")
	}

	if commonProtocol.ContentInsertionCompleted == nil {
		common_globals.Logger.Warning("Title Storage Admin::ContentInsertionCompleted not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Title Storage Admin::ContentInsertionCompleted not implemented")
	}

	r0, nexErr := commonProtocol.ContentInsertionCompleted(packet, contentId, contentSize)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeContentInsertionCompletedResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.ContentInsertionCompletedResponse{Retval: r0})

	return respond(endpoint, ubi.MethodContentInsertionCompleted, callID, body), nil
}
