// Package gamesession implements the common Game Session protocol handlers of the Quazal Rendez-Vous (Ubisoft) family.
//
// Layouts follow the NintendoClients wiki page "Game Session Protocol" (MIT). Every method calls a user-provided hook;
// a method whose hook is nil answers Core::NotImplemented unless a safe default is documented on its hook field.
package gamesession

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	ubi "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/game-session"
	game_session_types "github.com/PretendoNetwork/nex-protocols-go/v2/ubisoft/game-session/types"
)

// CommonProtocol handles the Game Session protocol with user-provided hooks
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol ubi.Interface

	// CreateSession is required: when nil the request is answered with Core::NotImplemented.
	CreateSession func(packet nex.PacketInterface, gameSession game_session_types.GameSession) (game_session_types.GameSessionKey, *nex.Error)

	// UpdateSession is required: when nil the request is answered with Core::NotImplemented.
	UpdateSession func(packet nex.PacketInterface, gameSessionUpdate game_session_types.GameSessionUpdate) *nex.Error

	// DeleteSession is required: when nil the request is answered with Core::NotImplemented.
	DeleteSession func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey) *nex.Error

	// MigrateSession is required: when nil the request is answered with Core::NotImplemented.
	MigrateSession func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey) (game_session_types.GameSessionKey, *nex.Error)

	// LeaveSession is required: when nil the request is answered with Core::NotImplemented.
	LeaveSession func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey) *nex.Error

	// GetSession is required: when nil the request is answered with Core::NotImplemented.
	GetSession func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey) (game_session_types.GameSessionSearchResult, *nex.Error)

	// SearchSessions is optional: when nil an empty result (empty lists, zero counters) is returned.
	SearchSessions func(packet nex.PacketInterface, gameSessionQuery game_session_types.GameSessionQuery) (types.List[game_session_types.GameSessionSearchResult], *nex.Error)

	// AddParticipants is required: when nil the request is answered with Core::NotImplemented.
	AddParticipants func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey, publicParticipantIDs types.List[types.UInt32], privateParticipantIDs types.List[types.UInt32]) *nex.Error

	// RemoveParticipants is required: when nil the request is answered with Core::NotImplemented.
	RemoveParticipants func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey, participantIDs types.List[types.UInt32]) *nex.Error

	// GetParticipantCount is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetParticipantCount func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey) (types.UInt32, *nex.Error)

	// GetParticipants is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetParticipants func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey, resultRange types.ResultRange) (types.List[game_session_types.GameSessionParticipant], *nex.Error)

	// SendInvitation is required: when nil the request is answered with Core::NotImplemented.
	SendInvitation func(packet nex.PacketInterface, invitation game_session_types.GameSessionInvitation) *nex.Error

	// GetInvitationReceivedCount is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetInvitationReceivedCount func(packet nex.PacketInterface, gameSessionTypeID types.UInt32) (types.UInt32, *nex.Error)

	// GetInvitationsReceived is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetInvitationsReceived func(packet nex.PacketInterface, gameSessionTypeID types.UInt32, resultRange types.ResultRange) (types.List[game_session_types.GameSessionInvitationReceived], *nex.Error)

	// GetInvitationSentCount is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetInvitationSentCount func(packet nex.PacketInterface, gameSessionTypeID types.UInt32) (types.UInt32, *nex.Error)

	// GetInvitationsSent is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetInvitationsSent func(packet nex.PacketInterface, gameSessionTypeID types.UInt32, resultRange types.ResultRange) (types.List[game_session_types.GameSessionInvitationSent], *nex.Error)

	// AcceptInvitation is required: when nil the request is answered with Core::NotImplemented.
	AcceptInvitation func(packet nex.PacketInterface, gameSessionInvitation game_session_types.GameSessionInvitationReceived) *nex.Error

	// DeclineInvitation is required: when nil the request is answered with Core::NotImplemented.
	DeclineInvitation func(packet nex.PacketInterface, gameSessionInvitation game_session_types.GameSessionInvitationReceived) *nex.Error

	// CancelInvitation is required: when nil the request is answered with Core::NotImplemented.
	CancelInvitation func(packet nex.PacketInterface, gameSessionInvitation game_session_types.GameSessionInvitationSent) *nex.Error

	// SendTextMessage is required: when nil the request is answered with Core::NotImplemented.
	SendTextMessage func(packet nex.PacketInterface, gameSessionMessage game_session_types.GameSessionMessage) *nex.Error

	// RegisterURLs is required: when nil the request is answered with Core::NotImplemented.
	RegisterURLs func(packet nex.PacketInterface, stationURLs types.List[types.StationURL]) *nex.Error

	// JoinSession is required: when nil the request is answered with Core::NotImplemented.
	JoinSession func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey) *nex.Error

	// AbandonSession is required: when nil the request is answered with Core::NotImplemented.
	AbandonSession func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey) *nex.Error

	// SearchSessionsWithParticipants is optional: when nil an empty result (empty lists, zero counters) is returned.
	SearchSessionsWithParticipants func(packet nex.PacketInterface, gameSessionTypeID types.UInt32, participantIDs types.List[types.UInt32]) (types.List[game_session_types.GameSessionSearchWithParticipantsResult], *nex.Error)

	// GetSessions is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetSessions func(packet nex.PacketInterface, gameSessionKeys types.List[game_session_types.GameSessionKey]) (types.List[game_session_types.GameSessionSearchResult], *nex.Error)

	// GetParticipantsURLs is optional: when nil an empty result (empty lists, zero counters) is returned.
	GetParticipantsURLs func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey, participantIDs types.List[types.UInt32]) (types.List[game_session_types.GameSessionParticipant], *nex.Error)

	// MigrateSessionHost is required: when nil the request is answered with Core::NotImplemented.
	MigrateSessionHost func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey) *nex.Error

	// SplitSession is required: when nil the request is answered with Core::NotImplemented.
	SplitSession func(packet nex.PacketInterface, gameSessionKey game_session_types.GameSessionKey) (game_session_types.GameSessionKey, *nex.Error)

	// SearchSocialSessions is optional: when nil an empty result (empty lists, zero counters) is returned.
	SearchSocialSessions func(packet nex.PacketInterface, gameSessionSocialQuery game_session_types.GameSessionSocialQuery) (types.List[game_session_types.GameSessionSearchWithParticipantsResult], *nex.Error)

	// ReportUnsuccessfulJoinSessions is required: when nil the request is answered with Core::NotImplemented.
	ReportUnsuccessfulJoinSessions func(packet nex.PacketInterface, unsuccessfulJoinSessions types.List[game_session_types.GameSessionUnsuccessfulJoinSession]) *nex.Error
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol ubi.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerCreateSession(commonProtocol.createSession)
	protocol.SetHandlerUpdateSession(commonProtocol.updateSession)
	protocol.SetHandlerDeleteSession(commonProtocol.deleteSession)
	protocol.SetHandlerMigrateSession(commonProtocol.migrateSession)
	protocol.SetHandlerLeaveSession(commonProtocol.leaveSession)
	protocol.SetHandlerGetSession(commonProtocol.getSession)
	protocol.SetHandlerSearchSessions(commonProtocol.searchSessions)
	protocol.SetHandlerAddParticipants(commonProtocol.addParticipants)
	protocol.SetHandlerRemoveParticipants(commonProtocol.removeParticipants)
	protocol.SetHandlerGetParticipantCount(commonProtocol.getParticipantCount)
	protocol.SetHandlerGetParticipants(commonProtocol.getParticipants)
	protocol.SetHandlerSendInvitation(commonProtocol.sendInvitation)
	protocol.SetHandlerGetInvitationReceivedCount(commonProtocol.getInvitationReceivedCount)
	protocol.SetHandlerGetInvitationsReceived(commonProtocol.getInvitationsReceived)
	protocol.SetHandlerGetInvitationSentCount(commonProtocol.getInvitationSentCount)
	protocol.SetHandlerGetInvitationsSent(commonProtocol.getInvitationsSent)
	protocol.SetHandlerAcceptInvitation(commonProtocol.acceptInvitation)
	protocol.SetHandlerDeclineInvitation(commonProtocol.declineInvitation)
	protocol.SetHandlerCancelInvitation(commonProtocol.cancelInvitation)
	protocol.SetHandlerSendTextMessage(commonProtocol.sendTextMessage)
	protocol.SetHandlerRegisterURLs(commonProtocol.registerURLs)
	protocol.SetHandlerJoinSession(commonProtocol.joinSession)
	protocol.SetHandlerAbandonSession(commonProtocol.abandonSession)
	protocol.SetHandlerSearchSessionsWithParticipants(commonProtocol.searchSessionsWithParticipants)
	protocol.SetHandlerGetSessions(commonProtocol.getSessions)
	protocol.SetHandlerGetParticipantsURLs(commonProtocol.getParticipantsURLs)
	protocol.SetHandlerMigrateSessionHost(commonProtocol.migrateSessionHost)
	protocol.SetHandlerSplitSession(commonProtocol.splitSession)
	protocol.SetHandlerSearchSocialSessions(commonProtocol.searchSocialSessions)
	protocol.SetHandlerReportUnsuccessfulJoinSessions(commonProtocol.reportUnsuccessfulJoinSessions)

	return commonProtocol
}
