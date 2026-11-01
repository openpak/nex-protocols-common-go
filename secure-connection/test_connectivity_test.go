package secureconnection

import (
	"errors"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/fakenet"
	secure_connection "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
)

func TestTestConnectivityAcknowledgesByDefault(t *testing.T) {
	c := &CommonProtocol{}

	msg, err := c.testConnectivity(nil, fakenet.NewPacket(5), 0x42)
	fakenet.RequireBody(t, msg, err, secure_connection.ProtocolID, secure_connection.MethodTestConnectivity, 0x42, nil)
}

func TestTestConnectivityHookOverrides(t *testing.T) {
	var got types.PID
	c := &CommonProtocol{TestConnectivity: func(pid types.PID) error {
		got = pid
		return errors.New("refused")
	}}

	msg, err := c.testConnectivity(nil, fakenet.NewPacket(5), 0x42)
	fakenet.RequireError(t, msg, err, nex.ResultCodes.Core.Unknown)

	if got != 5 {
		t.Fatalf("hook pid %d", got)
	}
}

func TestNewCommonProtocolRegistersTestConnectivity(t *testing.T) {
	p := secure_connection.NewProtocol()
	NewCommonProtocol(p)

	if p.TestConnectivity == nil {
		t.Fatal("TestConnectivity was not registered")
	}
}
