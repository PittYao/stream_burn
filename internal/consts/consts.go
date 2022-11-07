package consts

const (
	LogFormatConsole = "console"
	LogFormatJson    = "json"

	EOFError    = "EOF"
	EOFErrorMsg = "没有传入body参数"
)

const (
	// 本机服务类型

	Single          string = "single"
	Mix3            string = "mix3"
	Mix3Temperature string = "mix3Temperature"
	Mix4            string = "mix4"
	Mix4Temperature string = "mix4Temperature"
	PublicSingle    string = "publicSingle"
)

const (
	// 服务端口

	Mix3Port   = ":8007"
	Mix4Port   = ":8006"
	SinglePort = ":8005"
	PublicPort = ":8004"
	RtspPort   = "554"
)

const (
	// 接口地址

	CopyM3u8 = "/copyM3u8"
)

const (
	// nginx检测

	Http          = "http://"
	Localhost     = "127.0.0.1"
	RtmpPort      = "1935"
	M3u8UrlPort   = ":8880"
	NginxDisk     = "root \\w:" // 替换nginx盘符表达式
	NginxConfName = "/conf/nginx.conf"
)

const (
	// 流命令类型

	Transform string = "transform"
	Save      string = "save"
	Reboot    string = "reboot"
)

const (
	// 任务运行状态 -1=异常结束 1=正常结束

	Success     int64 = 1  //正常结束
	RunIng      int64 = 1  // 正在运行
	RunIngError int64 = -1 //  异常结束
)

const (
	TsFile        = "ts"
	FirstTsFile   = "video000.ts"
	TsFilePrefix  = "video"
	EmptyTsFile   = ""
	M3u8File      = "playlist.m3u8"
	SplitFileName = "_%03d.mp4"
)
