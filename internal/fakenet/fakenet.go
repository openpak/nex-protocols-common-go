// Package fakenet provides in-memory stand-ins for the NEX packet, connection
// and endpoint interfaces so protocol handlers can be exercised without a
// network or a PRUDP server. It is only meant for tests.
package fakenet

import (
	"net"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

// Endpoint implements the parts of nex.EndpointInterface the common handlers use.
type Endpoint struct {
	nex.EndpointInterface
	versions *nex.LibraryVersions
	settings *nex.ByteStreamSettings
}

// LibraryVersions returns default library versions.
func (e *Endpoint) LibraryVersions() *nex.LibraryVersions { return e.versions }

// ByteStreamSettings returns the default stream settings (u16 strings, 32 bit PIDs, no structure headers).
func (e *Endpoint) ByteStreamSettings() *nex.ByteStreamSettings { return e.settings }

// Connection implements nex.ConnectionInterface.
type Connection struct {
	endpoint *Endpoint
	pid      types.PID
}

// Endpoint returns the fake endpoint.
func (c *Connection) Endpoint() nex.EndpointInterface { return c.endpoint }

// Address returns nil.
func (c *Connection) Address() net.Addr { return nil }

// PID returns the connection PID.
func (c *Connection) PID() types.PID { return c.pid }

// SetPID sets the connection PID.
func (c *Connection) SetPID(pid types.PID) { c.pid = pid }

// Packet implements nex.PacketInterface.
type Packet struct {
	conn *Connection
	rmc  *nex.RMCMessage
}

// Sender returns the fake connection.
func (p *Packet) Sender() nex.ConnectionInterface { return p.conn }

// Payload returns nil.
func (p *Packet) Payload() []byte { return nil }

// SetPayload does nothing.
func (p *Packet) SetPayload([]byte) {}

// RMCMessage returns the request message.
func (p *Packet) RMCMessage() *nex.RMCMessage { return p.rmc }

// SetRMCMessage replaces the request message.
func (p *Packet) SetRMCMessage(message *nex.RMCMessage) { p.rmc = message }

// NewPacket returns a packet sent by pid, with the default stream settings.
func NewPacket(pid types.PID) *Packet {
	endpoint := &Endpoint{versions: nex.NewLibraryVersions(), settings: nex.NewByteStreamSettings()}
	conn := &Connection{endpoint: endpoint, pid: pid}

	return &Packet{conn: conn}
}

// Stream returns an input stream over data using the packet endpoint settings.
func (p *Packet) Stream(data []byte) *nex.ByteStreamIn {
	return nex.NewByteStreamIn(data, p.conn.endpoint.versions, p.conn.endpoint.settings)
}

// Out returns an empty output stream using the packet endpoint settings.
func (p *Packet) Out() *nex.ByteStreamOut {
	return nex.NewByteStreamOut(p.conn.endpoint.versions, p.conn.endpoint.settings)
}
