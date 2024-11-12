package services

import (
	"bytes"
	"fmt"
	daoInterfaces "parsing-service/apps/decoder/dao_interfaces"
	"parsing-service/apps/decoder/models"
	"parsing-service/pkg/logger"
)

type DecoderService struct {
	CommandMappingDAO    daoInterfaces.ICommandMappingDAO
	DeserializeLogicsDAO daoInterfaces.IDeserializeLogicsDAO
	logger               logger.ILogger
	// mu                   sync.Mutex
}

func NewDecoder(
	commandMappingDAO daoInterfaces.ICommandMappingDAO,
	deserializeLogicsDAO daoInterfaces.IDeserializeLogicsDAO,
	logger logger.ILogger) *DecoderService {
	return &DecoderService{
		CommandMappingDAO:    commandMappingDAO,
		DeserializeLogicsDAO: deserializeLogicsDAO,
		logger:               logger,
	}
}

func (s *DecoderService) GetCommandMappingFromCmdID(cmdID int, requestID string) []models.CommandMapping {
	return s.CommandMappingDAO.GetCommandMappingByCmdID(cmdID, requestID)
}

func (s *DecoderService) GetCommandMappingFromCmdName(cmdName string, requestID string) []models.CommandMapping {
	return s.CommandMappingDAO.GetCommandMappingByCmdName(cmdName, requestID)
}

func (s *DecoderService) GetCommandMappingFromSwVersionAndCmdID(swVersion string, cmdID int, requestID string) []models.CommandMapping {
	return s.CommandMappingDAO.GetCommandMappingBySwVersionAndCmdID(swVersion, cmdID, requestID)
}

func (s *DecoderService) GetDeserializeLogicsFromCmdId(cmdID int, requestID string) []models.DeserializeLogics {
	return s.DeserializeLogicsDAO.GetDeserializeLogicsByCmdId(cmdID, requestID)
}

func (s *DecoderService) GetDeserializeLogicsFromSwVersionAndCmdID(swVersion string, cmdID int, requestID string) []models.DeserializeLogics {
	return s.DeserializeLogicsDAO.GetDeserializeLogicsBySwVersionAndCmdID(swVersion, cmdID, requestID)
}

func checkPacketIntegrity(payload []byte) bool {

	// check length
	splittedPayload := bytes.Split(payload, []byte("RECT"))

	for _, part := range splittedPayload {
		var offset int
		if dcuPort == 0 {
			offset = 4
			dcuPort = int(part[len(part)-4])<<24 | int(part[len(part)-3])<<16 | int(part[len(part)-2])<<8 | int(part[len(part)-1])
		} else {
			offset = 0
			// dcuPort remains unchanged
		}

		if (len(part) < 18 && offset == 4) || (len(part) < 14 && offset == 0) {
			if len(part) > 0 {

			} else {
				continue
			}

		}

		if len(part) >= int(part[11])+13+1 {
			startByte := part[0]
			if startByte != 0xAA {
				continue
			}
			part = part[1:]

			if len(part) >= 11 {
				stopBytePos := 11 + int(part[10])

				if len(part) >= stopBytePos+2 {
					crcInPacket, err := DeserializeUInt16(part[stopBytePos:stopBytePos + 2])
					
				} else {
					fmt.Println("stop byte issue")
				}
			}

		}
	}

	// check crc
}
