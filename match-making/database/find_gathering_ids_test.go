package database

import (
	"math"
	"strings"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2/types"
)

func TestNormalizeResultRange(t *testing.T) {
	cases := []struct {
		offset, length        uint32
		wantLimit, wantOffset uint32
	}{
		{0, 0, MaxFindResults, 0},
		{math.MaxUint32, 10, 10, 0},
		{5, 10, 10, 5},
		{0, 5000, MaxFindResults, 0},
	}

	for _, c := range cases {
		limit, offset := NormalizeResultRange(types.ResultRange{Offset: types.NewUInt32(c.offset), Length: types.NewUInt32(c.length)})
		if limit != c.wantLimit || offset != c.wantOffset {
			t.Errorf("range(%d,%d) = (%d,%d), want (%d,%d)", c.offset, c.length, limit, offset, c.wantLimit, c.wantOffset)
		}
	}
}

func TestValidateDescriptionRegex(t *testing.T) {
	for _, ok := range []string{"", "^room [0-9]+$", "a|b"} {
		if err := ValidateDescriptionRegex(ok); err != nil {
			t.Errorf("%q should be accepted: %v", ok, err)
		}
	}

	for _, bad := range []string{"(", `(a)\1`, "[", strings.Repeat("x", MaxFindRegexLength+1)} {
		if err := ValidateDescriptionRegex(bad); err == nil {
			t.Errorf("%q should be rejected", bad)
		}
	}
}
