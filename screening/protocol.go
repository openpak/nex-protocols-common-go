// Package screening provides the hook driven common handlers for the Screening (0x7C). The wiki documents no request or response layout, so the hooks return the raw response body. With no hook set the methods answer Core::NotImplemented, unless AcceptReports is true, in which case the report is dropped with a log line and an empty success is answered protocol
package screening

import (
	"github.com/PretendoNetwork/nex-go/v2"
	screening "github.com/PretendoNetwork/nex-protocols-go/v2/screening"
	screening_types "github.com/PretendoNetwork/nex-protocols-go/v2/screening/types"
)

// CommonProtocol implements the Screening protocol on top of hooks. A method whose hook is nil answers Core::NotImplemented
type CommonProtocol struct {
	endpoint nex.EndpointInterface
	protocol screening.Interface

	// AcceptReports makes methods without a hook log the report and answer an empty success instead of Core::NotImplemented
	AcceptReports bool

	// ReportDataStoreContent is called for Screening::ReportDataStoreContent
	ReportDataStoreContent func(packet nex.PacketInterface, pContentParam screening_types.ScreeningDataStoreContentParam, pViolationParam screening_types.ScreeningUGCViolationParam) ([]byte, *nex.Error)
	// ReportUser is called for Screening::ReportUser
	ReportUser func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error)
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol screening.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerReportDataStoreContent(commonProtocol.reportDataStoreContent)
	protocol.SetHandlerReportUser(commonProtocol.reportUser)

	return commonProtocol
}
