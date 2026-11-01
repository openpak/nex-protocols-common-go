package screening

import (
	"bytes"
	"testing"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/internal/rmctest"
	screening "github.com/PretendoNetwork/nex-protocols-go/v2/screening"
	screening_types "github.com/PretendoNetwork/nex-protocols-go/v2/screening/types"
)

func TestAcceptReportsAnswersEmptySuccess(t *testing.T) {
	commonProtocol := &CommonProtocol{AcceptReports: true}

	response, nexError := commonProtocol.reportDataStoreContent(nil, rmctest.NewPacket(1), 8, screening_types.ScreeningDataStoreContentParam{}, screening_types.ScreeningUGCViolationParam{})
	if nexError != nil || !response.IsSuccess || len(response.Parameters) != 0 {
		t.Fatalf("response=%+v err=%v", response, nexError)
	}

	if response.ProtocolID != screening.ProtocolID || response.MethodID != screening.MethodReportDataStoreContent || response.CallID != 8 {
		t.Fatalf("bad header: %+v", response)
	}

	response, nexError = commonProtocol.reportUser(nil, rmctest.NewPacket(1), 9, []byte{1})
	if nexError != nil || !response.IsSuccess || len(response.Parameters) != 0 || response.MethodID != screening.MethodReportUser {
		t.Fatalf("response=%+v err=%v", response, nexError)
	}
}

func TestProtocolID(t *testing.T) {
	if screening.ProtocolID != 0x7C {
		t.Fatalf("Screening protocol ID = %#x", screening.ProtocolID)
	}
}

func TestHookBodyIsReturned(t *testing.T) {
	commonProtocol := &CommonProtocol{
		ReportUser: func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error) {
			return []byte{9, 9}, nil
		},
	}

	response, nexError := commonProtocol.reportUser(nil, rmctest.NewPacket(1), 1, nil)
	if nexError != nil || !bytes.Equal(response.Parameters, []byte{9, 9}) {
		t.Fatalf("response=%+v err=%v", response, nexError)
	}
}

func TestHookWinsOverAcceptReports(t *testing.T) {
	commonProtocol := &CommonProtocol{
		AcceptReports: true,
		ReportUser: func(packet nex.PacketInterface, packetPayload []byte) ([]byte, *nex.Error) {
			return nil, nex.NewError(nex.ResultCodes.Core.AccessDenied, "no")
		},
	}

	if _, nexError := commonProtocol.reportUser(nil, rmctest.NewPacket(1), 1, nil); !rmctest.HasCode(nexError, nex.ResultCodes.Core.AccessDenied) {
		t.Fatalf("expected the hook error, got %v", nexError)
	}
}
