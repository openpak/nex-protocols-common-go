// Package rmctest provides a fake packet for exercising RMC handlers without a network
package rmctest

import (
	"net"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

// Connection is a ConnectionInterface backed by a real endpoint with default stream settings
type Connection struct {
	endpoint nex.EndpointInterface
	pid      types.PID
}

// Endpoint returns the endpoint
func (c *Connection) Endpoint() nex.EndpointInterface { return c.endpoint }

// Address returns no address
func (c *Connection) Address() net.Addr { return nil }

// PID returns the PID
func (c *Connection) PID() types.PID { return c.pid }

// SetPID sets the PID
func (c *Connection) SetPID(pid types.PID) { c.pid = pid }

// Packet is a PacketInterface whose sender is a Connection
type Packet struct {
	sender  *Connection
	message *nex.RMCMessage
	payload []byte
}

// Sender returns the connection
func (p *Packet) Sender() nex.ConnectionInterface { return p.sender }

// Payload returns the payload
func (p *Packet) Payload() []byte { return p.payload }

// SetPayload sets the payload
func (p *Packet) SetPayload(payload []byte) { p.payload = payload }

// RMCMessage returns the request message
func (p *Packet) RMCMessage() *nex.RMCMessage { return p.message }

// SetRMCMessage sets the request message
func (p *Packet) SetRMCMessage(message *nex.RMCMessage) { p.message = message }

// NewPacket returns a packet from a connection with the given PID, on an endpoint using 4 byte PIDs and 2 byte string lengths
func NewPacket(pid types.PID) *Packet {
	endpoint := nex.NewPRUDPEndPoint(1)
	endpoint.Server = nex.NewPRUDPServer()

	return &Packet{sender: &Connection{endpoint: endpoint, pid: pid}}
}

// HasCode reports whether err carries the given result code, ignoring the error bit
func HasCode(err *nex.Error, code uint32) bool {
	return err != nil && err.ResultCode == nex.NewError(code, "").ResultCode
}
