package userstorage

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-storage"
	user_storage_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-storage/types"
)

// respond builds a success RMC message for the given method
func respond(endpoint nex.EndpointInterface, methodID uint32, callID uint32, body []byte) *nex.RMCMessage {
	message := nex.NewRMCSuccess(endpoint, body)
	message.ProtocolID = ubi.ProtocolID
	message.MethodID = methodID
	message.CallID = callID

	return message
}

// encodeSearchContentsResponse encodes the SearchContents response body
func encodeSearchContentsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SearchContentsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) searchContents(err error, packet nex.PacketInterface, callID uint32, query user_storage_types.UserStorageQuery) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::SearchContents: invalid request")
	}

	if commonProtocol.SearchContents == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeSearchContentsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchContentsResponse{})

		return respond(endpoint, ubi.MethodSearchContents, callID, body), nil
	}

	r0, nexErr := commonProtocol.SearchContents(packet, query)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSearchContentsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchContentsResponse{SearchResults: r0})

	return respond(endpoint, ubi.MethodSearchContents, callID, body), nil
}

// encodeSearchContentsWithTotalResponse encodes the SearchContentsWithTotal response body
func encodeSearchContentsWithTotalResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SearchContentsWithTotalResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) searchContentsWithTotal(err error, packet nex.PacketInterface, callID uint32, query user_storage_types.UserStorageQuery) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::SearchContentsWithTotal: invalid request")
	}

	if commonProtocol.SearchContentsWithTotal == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeSearchContentsWithTotalResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchContentsWithTotalResponse{})

		return respond(endpoint, ubi.MethodSearchContentsWithTotal, callID, body), nil
	}

	r0, r1, nexErr := commonProtocol.SearchContentsWithTotal(packet, query)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSearchContentsWithTotalResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchContentsWithTotalResponse{SearchResults: r0, TotalResults: r1})

	return respond(endpoint, ubi.MethodSearchContentsWithTotal, callID, body), nil
}

func (commonProtocol *CommonProtocol) deleteContent(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::DeleteContent: invalid request")
	}

	if commonProtocol.DeleteContent == nil {
		common_globals.Logger.Warning("User Storage::DeleteContent not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::DeleteContent not implemented")
	}

	nexErr := commonProtocol.DeleteContent(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodDeleteContent, callID, nil), nil
}

// encodeSaveMetaDataResponse encodes the SaveMetaData response body
func encodeSaveMetaDataResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SaveMetaDataResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) saveMetaData(err error, packet nex.PacketInterface, callID uint32, properties types.List[user_storage_types.ContentProperty], contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::SaveMetaData: invalid request")
	}

	if commonProtocol.SaveMetaData == nil {
		common_globals.Logger.Warning("User Storage::SaveMetaData not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::SaveMetaData not implemented")
	}

	r0, nexErr := commonProtocol.SaveMetaData(packet, properties, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSaveMetaDataResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SaveMetaDataResponse{ContentKey: r0})

	return respond(endpoint, ubi.MethodSaveMetaData, callID, body), nil
}

// encodeSaveContentDBResponse encodes the SaveContentDB response body
func encodeSaveContentDBResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SaveContentDBResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) saveContentDB(err error, packet nex.PacketInterface, callID uint32, properties types.List[user_storage_types.ContentProperty], data types.Buffer, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::SaveContentDB: invalid request")
	}

	if commonProtocol.SaveContentDB == nil {
		common_globals.Logger.Warning("User Storage::SaveContentDB not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::SaveContentDB not implemented")
	}

	r0, nexErr := commonProtocol.SaveContentDB(packet, properties, data, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSaveContentDBResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SaveContentDBResponse{ContentKey: r0})

	return respond(endpoint, ubi.MethodSaveContentDB, callID, body), nil
}

// encodeSaveContentAndGetUploadInfoResponse encodes the SaveContentAndGetUploadInfo response body
func encodeSaveContentAndGetUploadInfoResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SaveContentAndGetUploadInfoResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) saveContentAndGetUploadInfo(err error, packet nex.PacketInterface, callID uint32, properties types.List[user_storage_types.ContentProperty], size types.UInt32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::SaveContentAndGetUploadInfo: invalid request")
	}

	if commonProtocol.SaveContentAndGetUploadInfo == nil {
		common_globals.Logger.Warning("User Storage::SaveContentAndGetUploadInfo not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::SaveContentAndGetUploadInfo not implemented")
	}

	r0, r1, r2, nexErr := commonProtocol.SaveContentAndGetUploadInfo(packet, properties, size, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSaveContentAndGetUploadInfoResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SaveContentAndGetUploadInfoResponse{UploadInfo: r0, PendingID: r1, Headers: r2})

	return respond(endpoint, ubi.MethodSaveContentAndGetUploadInfo, callID, body), nil
}

// encodeUploadEndResponse encodes the UploadEnd response body
func encodeUploadEndResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.UploadEndResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) uploadEnd(err error, packet nex.PacketInterface, callID uint32, pendingID types.UInt64, result types.Bool, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::UploadEnd: invalid request")
	}

	if commonProtocol.UploadEnd == nil {
		common_globals.Logger.Warning("User Storage::UploadEnd not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::UploadEnd not implemented")
	}

	r0, nexErr := commonProtocol.UploadEnd(packet, pendingID, result, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeUploadEndResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.UploadEndResponse{ContentKey: r0})

	return respond(endpoint, ubi.MethodUploadEnd, callID, body), nil
}

// encodeGetContentDBResponse encodes the GetContentDB response body
func encodeGetContentDBResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetContentDBResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getContentDB(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::GetContentDB: invalid request")
	}

	if commonProtocol.GetContentDB == nil {
		common_globals.Logger.Warning("User Storage::GetContentDB not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::GetContentDB not implemented")
	}

	r0, nexErr := commonProtocol.GetContentDB(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetContentDBResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetContentDBResponse{Data: r0})

	return respond(endpoint, ubi.MethodGetContentDB, callID, body), nil
}

// encodeGetContentURLResponse encodes the GetContentURL response body
func encodeGetContentURLResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetContentURLResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getContentURL(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::GetContentURL: invalid request")
	}

	if commonProtocol.GetContentURL == nil {
		common_globals.Logger.Warning("User Storage::GetContentURL not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::GetContentURL not implemented")
	}

	r0, nexErr := commonProtocol.GetContentURL(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetContentURLResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetContentURLResponse{DownloadInfo: r0})

	return respond(endpoint, ubi.MethodGetContentURL, callID, body), nil
}

// encodeGetSlotCountResponse encodes the GetSlotCount response body
func encodeGetSlotCountResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetSlotCountResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getSlotCount(err error, packet nex.PacketInterface, callID uint32, typeID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::GetSlotCount: invalid request")
	}

	if commonProtocol.GetSlotCount == nil {
		common_globals.Logger.Warning("User Storage::GetSlotCount not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::GetSlotCount not implemented")
	}

	r0, nexErr := commonProtocol.GetSlotCount(packet, typeID)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetSlotCountResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetSlotCountResponse{SlotCount: r0})

	return respond(endpoint, ubi.MethodGetSlotCount, callID, body), nil
}

// encodeGetMetaDataResponse encodes the GetMetaData response body
func encodeGetMetaDataResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetMetaDataResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getMetaData(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::GetMetaData: invalid request")
	}

	if commonProtocol.GetMetaData == nil {
		common_globals.Logger.Warning("User Storage::GetMetaData not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::GetMetaData not implemented")
	}

	r0, nexErr := commonProtocol.GetMetaData(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetMetaDataResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetMetaDataResponse{Content: r0})

	return respond(endpoint, ubi.MethodGetMetaData, callID, body), nil
}

func (commonProtocol *CommonProtocol) like(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::Like: invalid request")
	}

	if commonProtocol.Like == nil {
		common_globals.Logger.Warning("User Storage::Like not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::Like not implemented")
	}

	nexErr := commonProtocol.Like(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodLike, callID, nil), nil
}

func (commonProtocol *CommonProtocol) unlike(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::Unlike: invalid request")
	}

	if commonProtocol.Unlike == nil {
		common_globals.Logger.Warning("User Storage::Unlike not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::Unlike not implemented")
	}

	nexErr := commonProtocol.Unlike(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodUnlike, callID, nil), nil
}

// encodeIsLikedResponse encodes the IsLiked response body
func encodeIsLikedResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.IsLikedResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) isLiked(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::IsLiked: invalid request")
	}

	if commonProtocol.IsLiked == nil {
		common_globals.Logger.Warning("User Storage::IsLiked not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::IsLiked not implemented")
	}

	r0, nexErr := commonProtocol.IsLiked(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeIsLikedResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.IsLikedResponse{Liked: r0})

	return respond(endpoint, ubi.MethodIsLiked, callID, body), nil
}

// encodeGetFavouritesResponse encodes the GetFavourites response body
func encodeGetFavouritesResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetFavouritesResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getFavourites(err error, packet nex.PacketInterface, callID uint32, contentTypes types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::GetFavourites: invalid request")
	}

	if commonProtocol.GetFavourites == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetFavouritesResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetFavouritesResponse{})

		return respond(endpoint, ubi.MethodGetFavourites, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetFavourites(packet, contentTypes)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetFavouritesResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetFavouritesResponse{Favourites: r0})

	return respond(endpoint, ubi.MethodGetFavourites, callID, body), nil
}

func (commonProtocol *CommonProtocol) makeFavourite(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::MakeFavourite: invalid request")
	}

	if commonProtocol.MakeFavourite == nil {
		common_globals.Logger.Warning("User Storage::MakeFavourite not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::MakeFavourite not implemented")
	}

	nexErr := commonProtocol.MakeFavourite(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodMakeFavourite, callID, nil), nil
}

func (commonProtocol *CommonProtocol) removeFromFavourites(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::RemoveFromFavourites: invalid request")
	}

	if commonProtocol.RemoveFromFavourites == nil {
		common_globals.Logger.Warning("User Storage::RemoveFromFavourites not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::RemoveFromFavourites not implemented")
	}

	nexErr := commonProtocol.RemoveFromFavourites(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodRemoveFromFavourites, callID, nil), nil
}

func (commonProtocol *CommonProtocol) reportInappropriate(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey, reason types.String) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::ReportInappropriate: invalid request")
	}

	if commonProtocol.ReportInappropriate == nil {
		common_globals.Logger.Warning("User Storage::ReportInappropriate not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::ReportInappropriate not implemented")
	}

	nexErr := commonProtocol.ReportInappropriate(packet, contentKey, reason)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodReportInappropriate, callID, nil), nil
}

func (commonProtocol *CommonProtocol) incrementPlayCount(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::IncrementPlayCount: invalid request")
	}

	if commonProtocol.IncrementPlayCount == nil {
		common_globals.Logger.Warning("User Storage::IncrementPlayCount not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::IncrementPlayCount not implemented")
	}

	nexErr := commonProtocol.IncrementPlayCount(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodIncrementPlayCount, callID, nil), nil
}

func (commonProtocol *CommonProtocol) updateCustomStat(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey, statID types.UInt16, incValue types.Int64) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::UpdateCustomStat: invalid request")
	}

	if commonProtocol.UpdateCustomStat == nil {
		common_globals.Logger.Warning("User Storage::UpdateCustomStat not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::UpdateCustomStat not implemented")
	}

	nexErr := commonProtocol.UpdateCustomStat(packet, contentKey, statID, incValue)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodUpdateCustomStat, callID, nil), nil
}

// encodeGetOwnContentsResponse encodes the GetOwnContents response body
func encodeGetOwnContentsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetOwnContentsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getOwnContents(err error, packet nex.PacketInterface, callID uint32, typeID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::GetOwnContents: invalid request")
	}

	if commonProtocol.GetOwnContents == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetOwnContentsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetOwnContentsResponse{})

		return respond(endpoint, ubi.MethodGetOwnContents, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetOwnContents(packet, typeID)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetOwnContentsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetOwnContentsResponse{Results: r0})

	return respond(endpoint, ubi.MethodGetOwnContents, callID, body), nil
}

// encodeGetMostPopularTagsResponse encodes the GetMostPopularTags response body
func encodeGetMostPopularTagsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetMostPopularTagsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getMostPopularTags(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::GetMostPopularTags: invalid request")
	}

	if commonProtocol.GetMostPopularTags == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetMostPopularTagsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetMostPopularTagsResponse{})

		return respond(endpoint, ubi.MethodGetMostPopularTags, callID, body), nil
	}

	r0, r1, nexErr := commonProtocol.GetMostPopularTags(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetMostPopularTagsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetMostPopularTagsResponse{Tags: r0, TotalNumberOfTaggings: r1})

	return respond(endpoint, ubi.MethodGetMostPopularTags, callID, body), nil
}

// encodeGetTagsResponse encodes the GetTags response body
func encodeGetTagsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetTagsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getTags(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::GetTags: invalid request")
	}

	if commonProtocol.GetTags == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetTagsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetTagsResponse{})

		return respond(endpoint, ubi.MethodGetTags, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetTags(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetTagsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetTagsResponse{TagIds: r0})

	return respond(endpoint, ubi.MethodGetTags, callID, body), nil
}

func (commonProtocol *CommonProtocol) tagContent(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey, newTagIds types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::TagContent: invalid request")
	}

	if commonProtocol.TagContent == nil {
		common_globals.Logger.Warning("User Storage::TagContent not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage::TagContent not implemented")
	}

	nexErr := commonProtocol.TagContent(packet, contentKey, newTagIds)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodTagContent, callID, nil), nil
}

// encodeSearchContentsByPlayersResponse encodes the SearchContentsByPlayers response body
func encodeSearchContentsByPlayersResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SearchContentsByPlayersResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) searchContentsByPlayers(err error, packet nex.PacketInterface, callID uint32, pids types.List[types.UInt32], query user_storage_types.UserStorageQuery) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::SearchContentsByPlayers: invalid request")
	}

	if commonProtocol.SearchContentsByPlayers == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeSearchContentsByPlayersResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchContentsByPlayersResponse{})

		return respond(endpoint, ubi.MethodSearchContentsByPlayers, callID, body), nil
	}

	r0, nexErr := commonProtocol.SearchContentsByPlayers(packet, pids, query)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSearchContentsByPlayersResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchContentsByPlayersResponse{SearchResults: r0})

	return respond(endpoint, ubi.MethodSearchContentsByPlayers, callID, body), nil
}

// encodeSearchContentsByPlayersWithTotalResponse encodes the SearchContentsByPlayersWithTotal response body
func encodeSearchContentsByPlayersWithTotalResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SearchContentsByPlayersWithTotalResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) searchContentsByPlayersWithTotal(err error, packet nex.PacketInterface, callID uint32, pids types.List[types.UInt32], query user_storage_types.UserStorageQuery) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage::SearchContentsByPlayersWithTotal: invalid request")
	}

	if commonProtocol.SearchContentsByPlayersWithTotal == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeSearchContentsByPlayersWithTotalResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchContentsByPlayersWithTotalResponse{})

		return respond(endpoint, ubi.MethodSearchContentsByPlayersWithTotal, callID, body), nil
	}

	r0, r1, nexErr := commonProtocol.SearchContentsByPlayersWithTotal(packet, pids, query)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSearchContentsByPlayersWithTotalResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchContentsByPlayersWithTotalResponse{SearchResults: r0, TotalResults: r1})

	return respond(endpoint, ubi.MethodSearchContentsByPlayersWithTotal, callID, body), nil
}
