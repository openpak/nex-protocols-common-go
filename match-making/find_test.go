package matchmaking

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
	match_making "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
)

func emptyListBytes() []byte { return []byte{0, 0, 0, 0} }

func TestFindInvitationsIsEmpty(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	response, nexError := commonProtocol.findInvitations(nil, rmctest.NewPacket(1), 4, types.NewResultRange())
	if nexError != nil {
		t.Fatalf("unexpected error: %v", nexError)
	}

	if !bytes.Equal(response.Parameters, emptyListBytes()) {
		t.Fatalf("parameters = %x", response.Parameters)
	}

	if response.ProtocolID != match_making.ProtocolID || response.MethodID != match_making.MethodFindInvitations || response.CallID != 4 {
		t.Fatalf("bad header: %+v", response)
	}
}

func TestFindBySQLQueryIsNotExecuted(t *testing.T) {
	// * A nil manager proves the query never reaches the database
	commonProtocol := &CommonProtocol{}

	response, nexError := commonProtocol.findBySQLQuery(nil, rmctest.NewPacket(1), 5, "SELECT 1; DROP TABLE matchmaking.gatherings", types.NewResultRange())
	if nexError != nil {
		t.Fatalf("unexpected error: %v", nexError)
	}

	if !bytes.Equal(response.Parameters, emptyListBytes()) || response.MethodID != match_making.MethodFindBySQLQuery {
		t.Fatalf("bad response: %+v", response)
	}
}

func TestFindDecodeErrors(t *testing.T) {
	commonProtocol := &CommonProtocol{}
	packet := rmctest.NewPacket(1)
	bad := errors.New("bad")
	rr := types.NewResultRange()

	calls := map[string]func() *nex.Error{
		"FindByOwner":            func() *nex.Error { _, e := commonProtocol.findByOwner(bad, packet, 1, 1, rr); return e },
		"FindByType":             func() *nex.Error { _, e := commonProtocol.findByType(bad, packet, 1, "", rr); return e },
		"FindByDescription":      func() *nex.Error { _, e := commonProtocol.findByDescription(bad, packet, 1, "", rr); return e },
		"FindByDescriptionRegex": func() *nex.Error { _, e := commonProtocol.findByDescriptionRegex(bad, packet, 1, "", rr); return e },
		"FindByID":               func() *nex.Error { _, e := commonProtocol.findByID(bad, packet, 1, nil); return e },
		"FindByParticipants":     func() *nex.Error { _, e := commonProtocol.findByParticipants(bad, packet, 1, nil); return e },
		"FindInvitations":        func() *nex.Error { _, e := commonProtocol.findInvitations(bad, packet, 1, rr); return e },
		"FindBySQLQuery":         func() *nex.Error { _, e := commonProtocol.findBySQLQuery(bad, packet, 1, "", rr); return e },
	}

	for name, call := range calls {
		if !rmctest.HasCode(call(), nex.ResultCodes.Core.InvalidArgument) {
			t.Errorf("%s: expected InvalidArgument", name)
		}
	}
}

func TestFindByDescriptionRegexRejectsBadPatterns(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	for _, pattern := range []string{"(", `(a)\1`, strings.Repeat("a", 257)} {
		_, nexError := commonProtocol.findByDescriptionRegex(nil, rmctest.NewPacket(1), 1, types.String(pattern), types.NewResultRange())
		if !rmctest.HasCode(nexError, nex.ResultCodes.Core.InvalidArgument) {
			t.Errorf("pattern %q: expected InvalidArgument, got %v", pattern, nexError)
		}
	}
}

func TestFindByIDLimit(t *testing.T) {
	commonProtocol := &CommonProtocol{}

	_, nexError := commonProtocol.findByID(nil, rmctest.NewPacket(1), 1, make(types.List[types.UInt32], 101))
	if !rmctest.HasCode(nexError, nex.ResultCodes.Core.InvalidArgument) {
		t.Fatalf("expected InvalidArgument, got %v", nexError)
	}
}

func fakeManager() *common_globals.MatchmakingManager {
	return &common_globals.MatchmakingManager{
		Mutex: &sync.RWMutex{},
		GetDetailedGatheringByID: func(manager *common_globals.MatchmakingManager, sourcePID uint64, gatheringID uint32) (types.RVType, string, *nex.Error) {
			if gatheringID == 404 {
				return nil, "", nex.NewError(nex.ResultCodes.RendezVous.SessionVoid, "gone")
			}

			gathering := match_making_types.NewGathering()
			gathering.ID = types.NewUInt32(gatheringID)
			gathering.Description = "room"

			return gathering, "Gathering", nil
		},
	}
}

func TestFindByIDEncoding(t *testing.T) {
	commonProtocol := &CommonProtocol{manager: fakeManager()}
	packet := rmctest.NewPacket(1)

	// * 404 is skipped, 7 is returned
	response, nexError := commonProtocol.findByID(nil, packet, 3, types.List[types.UInt32]{404, 7})
	if nexError != nil {
		t.Fatalf("unexpected error: %v", nexError)
	}

	if response.MethodID != match_making.MethodFindByID || response.CallID != 3 {
		t.Fatalf("bad header: %+v", response)
	}

	// * List<Data<Gathering>>: u32 count=1, then the holder starts with the class name String "Gathering\0" (u16 length 10)
	prefix := append([]byte{1, 0, 0, 0, 10, 0}, []byte("Gathering\x00")...)
	if !bytes.HasPrefix(response.Parameters, prefix) {
		t.Fatalf("parameters = %x, want prefix %x", response.Parameters, prefix)
	}

	// * It must also round trip through the protocol types
	endpoint := packet.Sender().Endpoint()
	stream := nex.NewByteStreamIn(response.Parameters, endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	var holders types.List[match_making_types.GatheringHolder]
	if err := holders.ExtractFrom(stream); err != nil {
		t.Fatalf("response does not decode: %v", err)
	}

	gathering, ok := holders[0].Object.(match_making_types.Gathering)
	if len(holders) != 1 || !ok || gathering.ID != 7 || gathering.Description != "room" {
		t.Fatalf("decoded %v", holders)
	}

	if stream.Remaining() != 0 {
		t.Fatalf("%d trailing bytes", stream.Remaining())
	}
}
