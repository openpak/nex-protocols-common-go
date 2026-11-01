package matchmake_extension

import (
	"bytes"
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
	match_making_constants "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/constants"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
)

func TestGetMyBlockListDefaultIsEmptyList(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	response, nexError := commonProtocol.getMyBlockList(nil, rmctest.NewPacket(1), 7)
	if nexError != nil {
		t.Fatalf("unexpected error: %v", nexError)
	}

	// * List<PID>: u32 count only
	if want := []byte{0, 0, 0, 0}; !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}

	if response.ProtocolID != matchmake_extension.ProtocolID || response.MethodID != matchmake_extension.MethodGetMyBlockList || response.CallID != 7 || !response.IsSuccess {
		t.Fatalf("bad response header: %+v", response)
	}
}

func TestGetMyBlockListHook(t *testing.T) {
	commonProtocol := &CommonProtocol{
		GetMyBlockList: func(packet nex.PacketInterface) (types.List[types.PID], *nex.Error) {
			return types.List[types.PID]{1, 0x01020304}, nil
		},
	}

	response, nexError := commonProtocol.getMyBlockList(nil, rmctest.NewPacket(1), 1)
	if nexError != nil {
		t.Fatalf("unexpected error: %v", nexError)
	}

	// * count=2, then two 4 byte little endian PIDs
	want := []byte{2, 0, 0, 0, 1, 0, 0, 0, 4, 3, 2, 1}
	if !bytes.Equal(response.Parameters, want) {
		t.Fatalf("parameters = %x, want %x", response.Parameters, want)
	}
}

func TestGetMyBlockListHookError(t *testing.T) {
	commonProtocol := &CommonProtocol{
		GetMyBlockList: func(packet nex.PacketInterface) (types.List[types.PID], *nex.Error) {
			return nil, nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
		},
	}

	if _, nexError := commonProtocol.getMyBlockList(nil, rmctest.NewPacket(1), 1); !rmctest.HasCode(nexError, nex.ResultCodes.Core.Unknown) {
		t.Fatalf("expected the hook error, got %v", nexError)
	}
}

func TestGetMyBlockListDecodeError(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	if _, nexError := commonProtocol.getMyBlockList(errors.New("bad"), rmctest.NewPacket(1), 1); !rmctest.HasCode(nexError, nex.ResultCodes.Core.InvalidArgument) {
		t.Fatalf("expected InvalidArgument, got %v", nexError)
	}
}

func TestClearMyBlockList(t *testing.T) {
	called := false
	commonProtocol := &CommonProtocol{}

	// * No hook: plain success without a body
	response, nexError := commonProtocol.clearMyBlockList(nil, rmctest.NewPacket(1), 9)
	if nexError != nil || len(response.Parameters) != 0 || response.MethodID != matchmake_extension.MethodClearMyBlockList || response.CallID != 9 {
		t.Fatalf("default clear failed: %+v %v", response, nexError)
	}

	commonProtocol.ClearMyBlockList = func(packet nex.PacketInterface) *nex.Error {
		called = true
		return nil
	}

	if _, nexError = commonProtocol.clearMyBlockList(nil, rmctest.NewPacket(1), 9); nexError != nil || !called {
		t.Fatalf("hook not called: %v", nexError)
	}

	commonProtocol.ClearMyBlockList = func(packet nex.PacketInterface) *nex.Error {
		return nex.NewError(nex.ResultCodes.Core.Unknown, "boom")
	}

	if _, nexError = commonProtocol.clearMyBlockList(nil, rmctest.NewPacket(1), 9); nexError == nil {
		t.Fatal("expected the hook error")
	}
}

func playingSession(pid types.PID, gathering match_making_types.GatheringInterface) match_making_types.PlayingSession {
	playing := match_making_types.NewPlayingSession()
	playing.PrincipalID = pid
	playing.Gathering.Object = gathering

	return playing
}

func TestBuildParticipantResults(t *testing.T) {
	session := match_making_types.NewMatchmakeSession()
	session.ApplicationBuffer = types.NewBuffer([]byte{1, 2, 3})
	session.MatchmakeParam.Params = types.NewMap[types.String, types.Variant]()
	session.MatchmakeParam.Params["k"] = types.NewVariant()

	sessions := types.List[match_making_types.PlayingSession]{
		playingSession(10, session),
		playingSession(11, match_making_types.NewGathering()), // * not a MatchmakeSession, skipped
	}

	none := buildParticipantResults(sessions, match_making_constants.FindMatchmakeSessionResultOptionNone)
	if len(none) != 1 || none[0].PrincipalID != 10 {
		t.Fatalf("expected only principal 10, got %v", none)
	}

	if len(none[0].Session.ApplicationBuffer) != 0 || len(none[0].Session.MatchmakeParam.Params) != 0 {
		t.Fatal("ApplicationBuffer and MatchmakeParam must be dropped without the result options")
	}

	both := buildParticipantResults(sessions, match_making_constants.FindMatchmakeSessionResultOptionApplicationBuffer|match_making_constants.FindMatchmakeSessionResultOptionMatchmakeParam)
	if !bytes.Equal(both[0].Session.ApplicationBuffer, []byte{1, 2, 3}) || len(both[0].Session.MatchmakeParam.Params) != 1 {
		t.Fatal("ApplicationBuffer and MatchmakeParam must be kept when requested")
	}

	// * The original session must not have been modified
	if len(session.ApplicationBuffer) != 3 {
		t.Fatal("source session was modified")
	}
}

func TestBuildParticipantResultsEmptyEncoding(t *testing.T) {
	results := buildParticipantResults(nil, match_making_constants.FindMatchmakeSessionResultOptionNone)

	packet := rmctest.NewPacket(1)
	stream := nex.NewByteStreamOut(packet.Sender().Endpoint().LibraryVersions(), packet.Sender().Endpoint().ByteStreamSettings())
	results.WriteTo(stream)

	if want := []byte{0, 0, 0, 0}; !bytes.Equal(stream.Bytes(), want) {
		t.Fatalf("bytes = %x, want %x", stream.Bytes(), want)
	}
}

func TestFindMatchmakeSessionByParticipantRejectsTooManyPrincipals(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	param := match_making_types.NewFindMatchmakeSessionByParticipantParam()
	param.PrincipalIDList = make(types.List[types.PID], maxPrincipalIDsToFindMatchmakeSession+1)

	if _, nexError := commonProtocol.findMatchmakeSessionByParticipant(nil, rmctest.NewPacket(1), 1, param); !rmctest.HasCode(nexError, nex.ResultCodes.Core.InvalidArgument) {
		t.Fatalf("expected InvalidArgument, got %v", nexError)
	}

	if _, nexError := commonProtocol.findMatchmakeSessionByParticipant(errors.New("bad"), rmctest.NewPacket(1), 1, match_making_types.NewFindMatchmakeSessionByParticipantParam()); !rmctest.HasCode(nexError, nex.ResultCodes.Core.InvalidArgument) {
		t.Fatalf("expected InvalidArgument on decode error, got %v", nexError)
	}
}

func TestScrubSessionDetail(t *testing.T) {
	session := match_making_types.NewMatchmakeSession()
	session.SessionKey = types.NewBuffer([]byte{9, 9})
	session.UserPassword = "secret"
	session.CodeWord = "visible"

	scrubbed := scrubSessionDetail(session)
	if len(scrubbed.SessionKey) != 0 || scrubbed.UserPassword != "" || scrubbed.CodeWord != "visible" {
		t.Fatalf("bad scrub: %+v", scrubbed)
	}
}
