package common_globals

import (
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
)

func TestCandidateClause(t *testing.T) {
	m := &MatchmakingManager{}
	if got := m.CandidateClause(nil); got != "" {
		t.Fatalf("unset = %q", got)
	}
	m.CandidateGatherings = func(*nex.PRUDPConnection) []uint32 { return nil }
	if got := m.CandidateClause(nil); got != "" {
		t.Fatalf("nil = %q", got)
	}
	m.CandidateGatherings = func(*nex.PRUDPConnection) []uint32 { return []uint32{} }
	if got := m.CandidateClause(nil); got != " AND false" {
		t.Fatalf("empty = %q", got)
	}
	m.CandidateGatherings = func(*nex.PRUDPConnection) []uint32 { return []uint32{3, 41} }
	if got := m.CandidateClause(nil); got != " AND g.id IN (3,41)" {
		t.Fatalf("ids = %q", got)
	}
}
