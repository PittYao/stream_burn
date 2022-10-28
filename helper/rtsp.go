package helper

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/internal/consts"
	"io"
	"strings"
)

// GetIpFormRtspUrl rtsp中获取ip为 192.168.99.152
func GetIpFormRtspUrl(rtspUrl string) string {
	i := strings.Index(rtspUrl, "@")
	ip := rtspUrl[(i + 1):]

	// 如果是rtsp://admin:CEBON123@192.168.99.116:554/cam/realmonitor?channel=1%26subtype=0; 去除:554/cam...
	j := strings.Index(ip, ":")
	if j != -1 {
		ip = ip[0:j]
	}
	return ip
}

// GetIpFormRtspUrls rtspUrls中获取ip为 192.168.99.152
func GetIpFormRtspUrls(rtspUrls ...string) []string {
	var ips []string
	if len(rtspUrls) != 0 {
		for _, rtspUrl := range rtspUrls {
			ip := GetIpFormRtspUrl(rtspUrl)
			ips = append(ips, ip)
		}
	}

	return ips
}

// CheckCameraIpNetWork 检测摄像头ip网络
func CheckCameraIpNetWork(rtspUrls ...string) error {
	for _, rtspUrl := range rtspUrls {
		cameraIp := GetIpFormRtspUrl(rtspUrl)
		err := CheckPortRunning(cameraIp, consts.RtspPort)
		if err != nil {
			log.L.Sugar().Errorf("摄像头网络不通 ip:%s", cameraIp)
			return errors.New(fmt.Sprintf("摄像头网络不通 ip:%s", cameraIp))
		}
	}

	return nil
}

func Md5ByString(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		panic(err)
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
}

// GetIpDirPathFormRtspUrls 从rtsp中获取ip名称的文件夹
func GetIpDirPathFormRtspUrls(split string, rtspUrls ...string) (string, error) {
	ips := GetIpFormRtspUrls(rtspUrls...)
	if len(ips) == 0 {
		return "", errors.New("根据RTSP_URL获取ip失败")
	}

	var ipDir string
	for i, ip := range ips {
		if i == 0 {
			ipDir = ip
		} else {
			ipDir = ipDir + split + ip
		}
	}
	return ipDir, nil
}
