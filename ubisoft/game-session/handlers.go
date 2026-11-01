package gamesession

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/game-session"
	game_session_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/game-session/types"
)

// respond builds a success RMC message for the given method
func respond(endpoint nex.EndpointInterface, methodID uint32, callID uint32, body []byte) *nex.RMCMessage {
	message := nex.NewRMCSuccess(endpoint, body)
	message.ProtocolID = ubi.ProtocolID
	message.MethodID = methodID
	message.CallID = callID

	return message
}

// encodeCreateSessionResponse encodes the CreateSession response body
func encodeCreateSessionResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.CreateSessionResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) createSession(err error, packet nex.PacketInterface, callID uint32, gameSession game_session_types.GameSession) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::CreateSession: invalid request")
	}

	if commonProtocol.CreateSession == nil {
		common_globals.Logger.Warning("Game Session::CreateSession not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::CreateSession not implemented")
	}

	r0, nexErr := commonProtocol.CreateSession(packet, gameSession)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeCreateSessionResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.CreateSessionResponse{GameSessionKey: r0})

	return respond(endpoint, ubi.MethodCreateSession, callID, body), nil
}

func (commonProtocol *CommonProtocol) updateSession(err error, packet nex.PacketInterface, callID uint32, gameSessionUpdate game_session_types.GameSessionUpdate) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::UpdateSession: invalid request")
	}

	if commonProtocol.UpdateSession == nil {
		common_globals.Logger.Warning("Game Session::UpdateSession not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::UpdateSession not implemented")
	}

	nexErr := commonProtocol.UpdateSession(packet, gameSessionUpdate)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodUpdateSession, callID, nil), nil
}

func (commonProtocol *CommonProtocol) deleteSession(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::DeleteSession: invalid request")
	}

	if commonProtocol.DeleteSession == nil {
		common_globals.Logger.Warning("Game Session::DeleteSession not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::DeleteSession not implemented")
	}

	nexErr := commonProtocol.DeleteSession(packet, gameSessionKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodDeleteSession, callID, nil), nil
}

// encodeMigrateSessionResponse encodes the MigrateSession response body
func encodeMigrateSessionResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.MigrateSessionResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) migrateSession(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::MigrateSession: invalid request")
	}

	if commonProtocol.MigrateSession == nil {
		common_globals.Logger.Warning("Game Session::MigrateSession not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::MigrateSession not implemented")
	}

	r0, nexErr := commonProtocol.MigrateSession(packet, gameSessionKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeMigrateSessionResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.MigrateSessionResponse{GameSessionKeyMigrated: r0})

	return respond(endpoint, ubi.MethodMigrateSession, callID, body), nil
}

func (commonProtocol *CommonProtocol) leaveSession(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::LeaveSession: invalid request")
	}

	if commonProtocol.LeaveSession == nil {
		common_globals.Logger.Warning("Game Session::LeaveSession not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::LeaveSession not implemented")
	}

	nexErr := commonProtocol.LeaveSession(packet, gameSessionKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodLeaveSession, callID, nil), nil
}

// encodeGetSessionResponse encodes the GetSession response body
func encodeGetSessionResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetSessionResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getSession(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::GetSession: invalid request")
	}

	if commonProtocol.GetSession == nil {
		common_globals.Logger.Warning("Game Session::GetSession not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::GetSession not implemented")
	}

	r0, nexErr := commonProtocol.GetSession(packet, gameSessionKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetSessionResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetSessionResponse{SearchResult: r0})

	return respond(endpoint, ubi.MethodGetSession, callID, body), nil
}

// encodeSearchSessionsResponse encodes the SearchSessions response body
func encodeSearchSessionsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SearchSessionsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) searchSessions(err error, packet nex.PacketInterface, callID uint32, gameSessionQuery game_session_types.GameSessionQuery) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::SearchSessions: invalid request")
	}

	if commonProtocol.SearchSessions == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeSearchSessionsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchSessionsResponse{})

		return respond(endpoint, ubi.MethodSearchSessions, callID, body), nil
	}

	r0, nexErr := commonProtocol.SearchSessions(packet, gameSessionQuery)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSearchSessionsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchSessionsResponse{SearchResults: r0})

	return respond(endpoint, ubi.MethodSearchSessions, callID, body), nil
}

func (commonProtocol *CommonProtocol) addParticipants(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey, publicParticipantIDs types.List[types.UInt32], privateParticipantIDs types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::AddParticipants: invalid request")
	}

	if commonProtocol.AddParticipants == nil {
		common_globals.Logger.Warning("Game Session::AddParticipants not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::AddParticipants not implemented")
	}

	nexErr := commonProtocol.AddParticipants(packet, gameSessionKey, publicParticipantIDs, privateParticipantIDs)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodAddParticipants, callID, nil), nil
}

func (commonProtocol *CommonProtocol) removeParticipants(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey, participantIDs types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::RemoveParticipants: invalid request")
	}

	if commonProtocol.RemoveParticipants == nil {
		common_globals.Logger.Warning("Game Session::RemoveParticipants not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::RemoveParticipants not implemented")
	}

	nexErr := commonProtocol.RemoveParticipants(packet, gameSessionKey, participantIDs)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodRemoveParticipants, callID, nil), nil
}

// encodeGetParticipantCountResponse encodes the GetParticipantCount response body
func encodeGetParticipantCountResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetParticipantCountResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getParticipantCount(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::GetParticipantCount: invalid request")
	}

	if commonProtocol.GetParticipantCount == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetParticipantCountResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetParticipantCountResponse{})

		return respond(endpoint, ubi.MethodGetParticipantCount, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetParticipantCount(packet, gameSessionKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetParticipantCountResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetParticipantCountResponse{Count: r0})

	return respond(endpoint, ubi.MethodGetParticipantCount, callID, body), nil
}

// encodeGetParticipantsResponse encodes the GetParticipants response body
func encodeGetParticipantsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetParticipantsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getParticipants(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::GetParticipants: invalid request")
	}

	if commonProtocol.GetParticipants == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetParticipantsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetParticipantsResponse{})

		return respond(endpoint, ubi.MethodGetParticipants, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetParticipants(packet, gameSessionKey, resultRange)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetParticipantsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetParticipantsResponse{Participants: r0})

	return respond(endpoint, ubi.MethodGetParticipants, callID, body), nil
}

func (commonProtocol *CommonProtocol) sendInvitation(err error, packet nex.PacketInterface, callID uint32, invitation game_session_types.GameSessionInvitation) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::SendInvitation: invalid request")
	}

	if commonProtocol.SendInvitation == nil {
		common_globals.Logger.Warning("Game Session::SendInvitation not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::SendInvitation not implemented")
	}

	nexErr := commonProtocol.SendInvitation(packet, invitation)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodSendInvitation, callID, nil), nil
}

// encodeGetInvitationReceivedCountResponse encodes the GetInvitationReceivedCount response body
func encodeGetInvitationReceivedCountResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetInvitationReceivedCountResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getInvitationReceivedCount(err error, packet nex.PacketInterface, callID uint32, gameSessionTypeID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::GetInvitationReceivedCount: invalid request")
	}

	if commonProtocol.GetInvitationReceivedCount == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetInvitationReceivedCountResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetInvitationReceivedCountResponse{})

		return respond(endpoint, ubi.MethodGetInvitationReceivedCount, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetInvitationReceivedCount(packet, gameSessionTypeID)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetInvitationReceivedCountResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetInvitationReceivedCountResponse{Count: r0})

	return respond(endpoint, ubi.MethodGetInvitationReceivedCount, callID, body), nil
}

// encodeGetInvitationsReceivedResponse encodes the GetInvitationsReceived response body
func encodeGetInvitationsReceivedResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetInvitationsReceivedResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getInvitationsReceived(err error, packet nex.PacketInterface, callID uint32, gameSessionTypeID types.UInt32, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::GetInvitationsReceived: invalid request")
	}

	if commonProtocol.GetInvitationsReceived == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetInvitationsReceivedResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetInvitationsReceivedResponse{})

		return respond(endpoint, ubi.MethodGetInvitationsReceived, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetInvitationsReceived(packet, gameSessionTypeID, resultRange)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetInvitationsReceivedResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetInvitationsReceivedResponse{Invitations: r0})

	return respond(endpoint, ubi.MethodGetInvitationsReceived, callID, body), nil
}

// encodeGetInvitationSentCountResponse encodes the GetInvitationSentCount response body
func encodeGetInvitationSentCountResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetInvitationSentCountResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getInvitationSentCount(err error, packet nex.PacketInterface, callID uint32, gameSessionTypeID types.UInt32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::GetInvitationSentCount: invalid request")
	}

	if commonProtocol.GetInvitationSentCount == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetInvitationSentCountResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetInvitationSentCountResponse{})

		return respond(endpoint, ubi.MethodGetInvitationSentCount, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetInvitationSentCount(packet, gameSessionTypeID)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetInvitationSentCountResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetInvitationSentCountResponse{Count: r0})

	return respond(endpoint, ubi.MethodGetInvitationSentCount, callID, body), nil
}

// encodeGetInvitationsSentResponse encodes the GetInvitationsSent response body
func encodeGetInvitationsSentResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetInvitationsSentResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getInvitationsSent(err error, packet nex.PacketInterface, callID uint32, gameSessionTypeID types.UInt32, resultRange types.ResultRange) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::GetInvitationsSent: invalid request")
	}

	if commonProtocol.GetInvitationsSent == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetInvitationsSentResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetInvitationsSentResponse{})

		return respond(endpoint, ubi.MethodGetInvitationsSent, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetInvitationsSent(packet, gameSessionTypeID, resultRange)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetInvitationsSentResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetInvitationsSentResponse{Invitations: r0})

	return respond(endpoint, ubi.MethodGetInvitationsSent, callID, body), nil
}

func (commonProtocol *CommonProtocol) acceptInvitation(err error, packet nex.PacketInterface, callID uint32, gameSessionInvitation game_session_types.GameSessionInvitationReceived) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::AcceptInvitation: invalid request")
	}

	if commonProtocol.AcceptInvitation == nil {
		common_globals.Logger.Warning("Game Session::AcceptInvitation not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::AcceptInvitation not implemented")
	}

	nexErr := commonProtocol.AcceptInvitation(packet, gameSessionInvitation)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodAcceptInvitation, callID, nil), nil
}

func (commonProtocol *CommonProtocol) declineInvitation(err error, packet nex.PacketInterface, callID uint32, gameSessionInvitation game_session_types.GameSessionInvitationReceived) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::DeclineInvitation: invalid request")
	}

	if commonProtocol.DeclineInvitation == nil {
		common_globals.Logger.Warning("Game Session::DeclineInvitation not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::DeclineInvitation not implemented")
	}

	nexErr := commonProtocol.DeclineInvitation(packet, gameSessionInvitation)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodDeclineInvitation, callID, nil), nil
}

func (commonProtocol *CommonProtocol) cancelInvitation(err error, packet nex.PacketInterface, callID uint32, gameSessionInvitation game_session_types.GameSessionInvitationSent) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::CancelInvitation: invalid request")
	}

	if commonProtocol.CancelInvitation == nil {
		common_globals.Logger.Warning("Game Session::CancelInvitation not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::CancelInvitation not implemented")
	}

	nexErr := commonProtocol.CancelInvitation(packet, gameSessionInvitation)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodCancelInvitation, callID, nil), nil
}

func (commonProtocol *CommonProtocol) sendTextMessage(err error, packet nex.PacketInterface, callID uint32, gameSessionMessage game_session_types.GameSessionMessage) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::SendTextMessage: invalid request")
	}

	if commonProtocol.SendTextMessage == nil {
		common_globals.Logger.Warning("Game Session::SendTextMessage not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::SendTextMessage not implemented")
	}

	nexErr := commonProtocol.SendTextMessage(packet, gameSessionMessage)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodSendTextMessage, callID, nil), nil
}

func (commonProtocol *CommonProtocol) registerURLs(err error, packet nex.PacketInterface, callID uint32, stationURLs types.List[types.StationURL]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::RegisterURLs: invalid request")
	}

	if commonProtocol.RegisterURLs == nil {
		common_globals.Logger.Warning("Game Session::RegisterURLs not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::RegisterURLs not implemented")
	}

	nexErr := commonProtocol.RegisterURLs(packet, stationURLs)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodRegisterURLs, callID, nil), nil
}

func (commonProtocol *CommonProtocol) joinSession(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::JoinSession: invalid request")
	}

	if commonProtocol.JoinSession == nil {
		common_globals.Logger.Warning("Game Session::JoinSession not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::JoinSession not implemented")
	}

	nexErr := commonProtocol.JoinSession(packet, gameSessionKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodJoinSession, callID, nil), nil
}

func (commonProtocol *CommonProtocol) abandonSession(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::AbandonSession: invalid request")
	}

	if commonProtocol.AbandonSession == nil {
		common_globals.Logger.Warning("Game Session::AbandonSession not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::AbandonSession not implemented")
	}

	nexErr := commonProtocol.AbandonSession(packet, gameSessionKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodAbandonSession, callID, nil), nil
}

// encodeSearchSessionsWithParticipantsResponse encodes the SearchSessionsWithParticipants response body
func encodeSearchSessionsWithParticipantsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SearchSessionsWithParticipantsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) searchSessionsWithParticipants(err error, packet nex.PacketInterface, callID uint32, gameSessionTypeID types.UInt32, participantIDs types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::SearchSessionsWithParticipants: invalid request")
	}

	if commonProtocol.SearchSessionsWithParticipants == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeSearchSessionsWithParticipantsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchSessionsWithParticipantsResponse{})

		return respond(endpoint, ubi.MethodSearchSessionsWithParticipants, callID, body), nil
	}

	r0, nexErr := commonProtocol.SearchSessionsWithParticipants(packet, gameSessionTypeID, participantIDs)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSearchSessionsWithParticipantsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchSessionsWithParticipantsResponse{SearchResults: r0})

	return respond(endpoint, ubi.MethodSearchSessionsWithParticipants, callID, body), nil
}

// encodeGetSessionsResponse encodes the GetSessions response body
func encodeGetSessionsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetSessionsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getSessions(err error, packet nex.PacketInterface, callID uint32, gameSessionKeys types.List[game_session_types.GameSessionKey]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::GetSessions: invalid request")
	}

	if commonProtocol.GetSessions == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetSessionsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetSessionsResponse{})

		return respond(endpoint, ubi.MethodGetSessions, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetSessions(packet, gameSessionKeys)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetSessionsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetSessionsResponse{SearchResults: r0})

	return respond(endpoint, ubi.MethodGetSessions, callID, body), nil
}

// encodeGetParticipantsURLsResponse encodes the GetParticipantsURLs response body
func encodeGetParticipantsURLsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.GetParticipantsURLsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) getParticipantsURLs(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey, participantIDs types.List[types.UInt32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::GetParticipantsURLs: invalid request")
	}

	if commonProtocol.GetParticipantsURLs == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeGetParticipantsURLsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetParticipantsURLsResponse{})

		return respond(endpoint, ubi.MethodGetParticipantsURLs, callID, body), nil
	}

	r0, nexErr := commonProtocol.GetParticipantsURLs(packet, gameSessionKey, participantIDs)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeGetParticipantsURLsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.GetParticipantsURLsResponse{Participants: r0})

	return respond(endpoint, ubi.MethodGetParticipantsURLs, callID, body), nil
}

func (commonProtocol *CommonProtocol) migrateSessionHost(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::MigrateSessionHost: invalid request")
	}

	if commonProtocol.MigrateSessionHost == nil {
		common_globals.Logger.Warning("Game Session::MigrateSessionHost not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::MigrateSessionHost not implemented")
	}

	nexErr := commonProtocol.MigrateSessionHost(packet, gameSessionKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodMigrateSessionHost, callID, nil), nil
}

// encodeSplitSessionResponse encodes the SplitSession response body
func encodeSplitSessionResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SplitSessionResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) splitSession(err error, packet nex.PacketInterface, callID uint32, gameSessionKey game_session_types.GameSessionKey) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::SplitSession: invalid request")
	}

	if commonProtocol.SplitSession == nil {
		common_globals.Logger.Warning("Game Session::SplitSession not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::SplitSession not implemented")
	}

	r0, nexErr := commonProtocol.SplitSession(packet, gameSessionKey)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSplitSessionResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SplitSessionResponse{GameSessionKeyMigrated: r0})

	return respond(endpoint, ubi.MethodSplitSession, callID, body), nil
}

// encodeSearchSocialSessionsResponse encodes the SearchSocialSessions response body
func encodeSearchSocialSessionsResponse(versions *nex.LibraryVersions, settings *nex.ByteStreamSettings, response ubi.SearchSocialSessionsResponse) []byte {
	stream := nex.NewByteStreamOut(versions, settings)
	response.WriteTo(stream)

	return stream.Bytes()
}

func (commonProtocol *CommonProtocol) searchSocialSessions(err error, packet nex.PacketInterface, callID uint32, gameSessionSocialQuery game_session_types.GameSessionSocialQuery) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::SearchSocialSessions: invalid request")
	}

	if commonProtocol.SearchSocialSessions == nil {
		endpoint := packet.Sender().Endpoint()

		body := encodeSearchSocialSessionsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchSocialSessionsResponse{})

		return respond(endpoint, ubi.MethodSearchSocialSessions, callID, body), nil
	}

	r0, nexErr := commonProtocol.SearchSocialSessions(packet, gameSessionSocialQuery)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	body := encodeSearchSocialSessionsResponse(endpoint.LibraryVersions(), endpoint.ByteStreamSettings(), ubi.SearchSocialSessionsResponse{SearchResults: r0})

	return respond(endpoint, ubi.MethodSearchSocialSessions, callID, body), nil
}

func (commonProtocol *CommonProtocol) reportUnsuccessfulJoinSessions(err error, packet nex.PacketInterface, callID uint32, unsuccessfulJoinSessions types.List[game_session_types.GameSessionUnsuccessfulJoinSession]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())

		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "Game Session::ReportUnsuccessfulJoinSessions: invalid request")
	}

	if commonProtocol.ReportUnsuccessfulJoinSessions == nil {
		common_globals.Logger.Warning("Game Session::ReportUnsuccessfulJoinSessions not implemented")

		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "Game Session::ReportUnsuccessfulJoinSessions not implemented")
	}

	nexErr := commonProtocol.ReportUnsuccessfulJoinSessions(packet, unsuccessfulJoinSessions)
	if nexErr != nil {
		common_globals.Logger.Error(nexErr.Error())

		return nil, nexErr
	}

	endpoint := packet.Sender().Endpoint()

	return respond(endpoint, ubi.MethodReportUnsuccessfulJoinSessions, callID, nil), nil
}
