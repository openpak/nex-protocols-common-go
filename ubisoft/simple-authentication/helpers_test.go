package simpleauthentication

import (
	"encoding/hex"
	"errors"
	"testing"

	nex "github.com/PretendoNetwork/nex-go/v2"
)

func newSettings(header bool) *nex.ByteStreamSettings {
	settings := nex.NewByteStreamSettings()
	settings.UseStructureHeader = header

	return settings
}

func newOut(header bool) *nex.ByteStreamOut {
	return nex.NewByteStreamOut(nex.NewLibraryVersions(), newSettings(header))
}

func newIn(t *testing.T, data string, header bool) *nex.ByteStreamIn {
	t.Helper()

	raw, err := hex.DecodeString(data)
	if err != nil {
		t.Fatal(err)
	}

	return nex.NewByteStreamIn(raw, nex.NewLibraryVersions(), newSettings(header))
}

func checkBytes(t *testing.T, got []byte, want string) {
	t.Helper()

	if hex.EncodeToString(got) != want {
		t.Fatalf("bytes mismatch\n got: %x\nwant: %s", got, want)
	}
}

var errTest = errors.New("test")
