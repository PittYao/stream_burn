package httpclient

import (
	"errors"
	"fmt"
	"github.com/PittYao/stream_burn/components/gin/response"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/guonaihong/gout"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type CallbackDTO struct {
	Code    int         `json:"code"`
	Message string      `json:"message" `
	Data    interface{} `json:"data"`
}

func CopyM3u8HttpClient(url string, id uint) (err error, response response.Response) {
	globalWithOpt := gout.NewWithOpt(gout.WithTimeout(time.Second * 15))
	err = globalWithOpt.
		// POST请求
		POST(url).
		// 打开debug模式
		Debug(true).
		SetJSON(gout.H{
			"id": id,
		}).
		// BindJSON解析返回的body内容
		// 同类函数有BindBody, BindYAML, BindXML
		BindJSON(&response).
		// 结束函数
		Do()

	// 判断错误
	if err != nil {
		log.L.Error("请求临时m3u8文件异常", zap.String("url", url), zap.Error(err))
		return
	}

	if response.Code != 200 {
		log.L.Error("请求临时m3u8文件异常", zap.String("url", url), zap.Any("response", response))
		err = errors.New(response.Msg)
		return
	}

	return
}

func CallBackHttpClient(url string, burnInfoUuId string) error {
	globalWithOpt := gout.NewWithOpt(gout.WithTimeout(time.Second * 15))
	err := globalWithOpt.
		// POST请求
		POST(url).
		// 打开debug模式
		Debug(true).
		SetJSON(CallbackDTO{
			Code:    http.StatusOK,
			Message: "",
			Data: map[string]interface{}{
				"taskId": burnInfoUuId,
				// TODO: 兼容以前的api 添加此参数
				"status": 1,
			},
		}).
		// 结束函数
		Do()

	if err != nil {
		log.L.Error("回调业务端接口失败", zap.String("url", url), zap.Error(err))
		return err
	}

	log.L.Info("回调业务端接口成功", zap.String("url", url), zap.String("taskId", burnInfoUuId))
	return nil

}

func CallBackFileHttpClient(url string, burnInfoUuId string, error error) error {

	var callbackDTO CallbackDTO
	if error != nil {
		callbackDTO = CallbackDTO{
			Code:    http.StatusInternalServerError,
			Message: "下载失败",
			Data: map[string]interface{}{
				"taskId": burnInfoUuId,
				// TODO: 兼容以前的api 添加此参数
				"status": 1,
			},
		}
	} else {
		callbackDTO = CallbackDTO{
			Code:    http.StatusOK,
			Message: "下载完成",
			Data: map[string]interface{}{
				"taskId": burnInfoUuId,
				// TODO: 兼容以前的api 添加此参数
				"status": 1,
			},
		}
	}

	globalWithOpt := gout.NewWithOpt(gout.WithTimeout(time.Second * 3))
	err := globalWithOpt.
		// POST请求
		POST(url).
		// 打开debug模式
		Debug(true).
		SetJSON(callbackDTO).
		// 结束函数
		Do()

	if err != nil {
		log.L.Error("回调业务端文件接口失败", zap.String("url", url), zap.Error(err))
		return err
	}

	log.L.Info("回调业务端文件接口成功", zap.String("url", url), zap.String("taskId", burnInfoUuId))
	return nil

}

func CheckM3u8Available(url string) error {
	globalWithOpt := gout.NewWithOpt(gout.WithTimeout(time.Second * 5))
	resp, err := globalWithOpt.
		GET(url).
		Debug(true).
		Response()
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		errMsg := fmt.Sprintf("m3u8文件不存在:%s ", url)
		log.L.Error(errMsg)
		return errors.New(errMsg)
	}

	return err
}
