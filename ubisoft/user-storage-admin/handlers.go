package userstorageadmin

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/user-storage-admin"
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

// encodeGetContentsToModerateResponse encodes the GetContentsToModerate response body
func encodeGetContentsToModerateResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetContentsToModerateResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getContentsToModerate(err error, packet nex.PacketInterface, callID uint32, typeID types.UInt32, offset types.UInt32, size types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage Admin::GetContentsToModerate: invalid request")
	}

	if commonProtocol.GetContentsToModerate == nil {
		common_globals.Logger.Warning("User Storage Admin::GetContentsToModerate not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage Admin::GetContentsToModerate not implemented")
	}

	r0, r1, nexErr := commonProtocol.GetContentsToModerate(packet, typeID, offset, size)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetContentsToModerateResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetContentsToModerateResponse{Contents: r0, TotalResults: r1})

	return respond(endpoint, ubi.MethodGetContentsToModerate, callID, body), nil
}

func (commonProtocol *CommonProtocol) flagContentAsVerified(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage Admin::FlagContentAsVerified: invalid request")
	}

	if commonProtocol.FlagContentAsVerified == nil {
		common_globals.Logger.Warning("User Storage Admin::FlagContentAsVerified not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage Admin::FlagContentAsVerified not implemented")
	}

	nexErr := commonProtocol.FlagContentAsVerified(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodFlagContentAsVerified, callID, nil), nil
}

func (commonProtocol *CommonProtocol) banContent(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage Admin::BanContent: invalid request")
	}

	if commonProtocol.BanContent == nil {
		common_globals.Logger.Warning("User Storage Admin::BanContent not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage Admin::BanContent not implemented")
	}

	nexErr := commonProtocol.BanContent(packet, contentKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodBanContent, callID, nil), nil
}

func (commonProtocol *CommonProtocol) banUser(err error, packet nex.PacketInterface, callID uint32, pid types.UInt32, reason types.String, banContents types.Bool, expireDate types.DateTime) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage Admin::BanUser: invalid request")
	}

	if commonProtocol.BanUser == nil {
		common_globals.Logger.Warning("User Storage Admin::BanUser not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage Admin::BanUser not implemented")
	}

	nexErr := commonProtocol.BanUser(packet, pid, reason, banContents, expireDate)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodBanUser, callID, nil), nil
}

func (commonProtocol *CommonProtocol) banUserFromContentType(err error, packet nex.PacketInterface, callID uint32, typeID types.UInt32, pid types.UInt32, reason types.String, banContents types.Bool, expireDate types.DateTime) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage Admin::BanUserFromContentType: invalid request")
	}

	if commonProtocol.BanUserFromContentType == nil {
		common_globals.Logger.Warning("User Storage Admin::BanUserFromContentType not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage Admin::BanUserFromContentType not implemented")
	}

	nexErr := commonProtocol.BanUserFromContentType(packet, typeID, pid, reason, banContents, expireDate)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodBanUserFromContentType, callID, nil), nil
}

func (commonProtocol *CommonProtocol) unbanUser(err error, packet nex.PacketInterface, callID uint32, pid types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage Admin::UnbanUser: invalid request")
	}

	if commonProtocol.UnbanUser == nil {
		common_globals.Logger.Warning("User Storage Admin::UnbanUser not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage Admin::UnbanUser not implemented")
	}

	nexErr := commonProtocol.UnbanUser(packet, pid)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodUnbanUser, callID, nil), nil
}

func (commonProtocol *CommonProtocol) unbanUserFromContentType(err error, packet nex.PacketInterface, callID uint32, typeID types.UInt32, pid types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage Admin::UnbanUserFromContentType: invalid request")
	}

	if commonProtocol.UnbanUserFromContentType == nil {
		common_globals.Logger.Warning("User Storage Admin::UnbanUserFromContentType not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage Admin::UnbanUserFromContentType not implemented")
	}

	nexErr := commonProtocol.UnbanUserFromContentType(packet, typeID, pid)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodUnbanUserFromContentType, callID, nil), nil
}

// encodeGetContentsToModerateWithThresholdResponse encodes the GetContentsToModerateWithThreshold response body
func encodeGetContentsToModerateWithThresholdResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetContentsToModerateWithThresholdResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getContentsToModerateWithThreshold(err error, packet nex.PacketInterface, callID uint32, typeID types.UInt32, threshold types.UInt32, offset types.UInt32, size types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage Admin::GetContentsToModerateWithThreshold: invalid request")
	}

	if commonProtocol.GetContentsToModerateWithThreshold == nil {
		common_globals.Logger.Warning("User Storage Admin::GetContentsToModerateWithThreshold not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage Admin::GetContentsToModerateWithThreshold not implemented")
	}

	r0, r1, nexErr := commonProtocol.GetContentsToModerateWithThreshold(packet, typeID, threshold, offset, size)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetContentsToModerateWithThresholdResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetContentsToModerateWithThresholdResponse{Contents: r0, TotalResults: r1})

	return respond(endpoint, ubi.MethodGetContentsToModerateWithThreshold, callID, body), nil
}

func (commonProtocol *CommonProtocol) updateMetaData(err error, packet nex.PacketInterface, callID uint32, contentKey user_storage_types.UserContentKey, properties types.List[user_storage_types.ContentProperty]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "User Storage Admin::UpdateMetaData: invalid request")
	}

	if commonProtocol.UpdateMetaData == nil {
		common_globals.Logger.Warning("User Storage Admin::UpdateMetaData not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "User Storage Admin::UpdateMetaData not implemented")
	}

	nexErr := commonProtocol.UpdateMetaData(packet, contentKey, properties)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodUpdateMetaData, callID, nil), nil
}
