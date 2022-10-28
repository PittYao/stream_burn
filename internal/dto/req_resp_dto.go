package dto

type BurnSettingDTO struct {
	TaskNum        int    `json:"taskNum" example:"1" binding:"required"`
	EncryptionType int    `json:"encryptionType"`
	Password       string `json:"password"`
	OdaSavePath    string `json:"odaSavePath" example:"D:/videodata" binding:"required"`
	VideoName      string `json:"videoName" example:"video"`
}

type BurnMix3VideoDTO struct {
	RtspUrlMiddle string `json:"rtspUrlMiddle" example:"rtsp://admin:cebon61332433@192.168.99.215:554/cam/realmonitor?channel=1&subtype=0" binding:"required"`
	RtspUrlLeft   string `json:"rtspUrlLeft" example:"rtsp://admin:cebon61332433@192.168.99.215:554/cam/realmonitor?channel=1&subtype=1" binding:"required"`
	RtspUrlRight  string `json:"rtspUrlRight" example:"rtsp://admin:cebon61332433@192.168.99.215:554/cam/realmonitor?channel=1&subtype=1" binding:"required"`
	Temperature   string `json:"temperature" example:""`
	StartTime     string `json:"startTime" example:"2022-05-11 15:20:00" binding:"required" `
	EndTime       string `json:"endTime" example:"2022-05-11 15:25:00" binding:"required"`
	CallbackUrl   string `json:"callBackUrl" example:"http://localhost:8010/api/v1/callback"`
	OdaSavePath   string `json:"odaSavePath" example:"D:/downloadVideo"`
	TaskId        uint   `json:"taskId" example:"1"`
	FileSavePath  string `json:"fileSavePath"`
	VideoName     string `json:"videoName"`
}
