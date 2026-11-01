package screening

import (
	"github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmc"
	screening "github.com/PretendoNetwork/nex-protocols-go/v2/screening"
	screening_types "github.com/PretendoNetwork/nex-protocols-go/v2/screening/types"
)

func (commonProtocol *CommonProtocol) reportDataStoreContent(err error, packet nex.PacketInterface, callID uint32, pContentParam screening_types.ScreeningDataStoreContentParam, pViolationParam screening_types.ScreeningUGCViolationParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.ReportDataStoreContent == nil {
		if commonProtocol.AcceptReports {
			common_globals.Logger.Warning("Screening::ReportDataStoreContent dropped: no hook set")

			return rmc.RespondRaw(packet, screening.ProtocolID, screening.MethodReportDataStoreContent, callID, nil), nil
		}

		return nil, rmc.NotImplemented("Screening::ReportDataStoreContent")
	}

	body, nexError := commonProtocol.ReportDataStoreContent(packet, pContentParam, pViolationParam)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.RespondRaw(packet, screening.ProtocolID, screening.MethodReportDataStoreContent, callID, body), nil
}

func (commonProtocol *CommonProtocol) reportUser(err error, packet nex.PacketInterface, callID uint32, packetPayload []byte) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		return nil, rmc.InvalidArgument(err)
	}

	if commonProtocol.ReportUser == nil {
		if commonProtocol.AcceptReports {
			common_globals.Logger.Warning("Screening::ReportUser dropped: no hook set")

			return rmc.RespondRaw(packet, screening.ProtocolID, screening.MethodReportUser, callID, nil), nil
		}

		return nil, rmc.NotImplemented("Screening::ReportUser")
	}

	body, nexError := commonProtocol.ReportUser(packet, packetPayload)
	if nexError != nil {
		return nil, nexError
	}

	return rmc.RespondRaw(packet, screening.ProtocolID, screening.MethodReportUser, callID, body), nil
}
