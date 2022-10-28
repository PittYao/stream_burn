package burninfocmd

import (
	"context"
	"errors"
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/PittYao/stream_burn/components/mysql"
	"github.com/PittYao/stream_burn/helper"
	"github.com/PittYao/stream_burn/internal/consts"
	"github.com/PittYao/stream_burn/internal/httpclient"
	"github.com/PittYao/stream_burn/internal/model/burninfo"
	"github.com/PittYao/stream_burn/internal/model/burnsetting"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

type BurnInfoCmd struct {
	gorm.Model
	FfmpegCmd     string
	FfmpegCmdArgs []string `gorm:"-"`
	DoneStatus    int64    // 完成状态 -1=失败 1=成功
	BurnInfoID    uint
	CloseTime     *time.Time
}

// Update 更新
func (b *BurnInfoCmd) Update() error {
	save := mysql.Instance.Save(&b)
	if save.Error != nil {
		log.L.Error("BurnInfoCmd 更新失败", zap.Error(save.Error))
		return errors.New("BurnInfoCmd 更新失败")
	}
	return save.Error

}

// FfmpegDownloadMixVideo 执行下载视频
func (b *BurnInfoCmd) FfmpegDownloadMixVideo() error {
	// 执行命令
	ctx, _ := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, config.C.Ffmpeg.LibPath, b.FfmpegCmdArgs...)

	cmd.StdinPipe()

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	if err := cmd.Start(); err != nil {
		log.L.Error("下载三合一画面ffmpeg命令失败", zap.String("ffmpeg_cmd", b.FfmpegCmd))
		return err
	}

	// 等待命令执行完成
	err := cmd.Wait()
	closeTime := time.Now()
	b.CloseTime = &closeTime

	if err == nil {
		log.L.Info("下载合成画面ffmpeg命令过程中正常结束",
			zap.String("ffmpeg_cmd", b.FfmpegCmd),
			zap.Uint("id", b.ID),
		)
		b.DoneStatus = consts.Success
	} else {
		log.L.Error("下载合成画面ffmpeg命令过程中异常结束",
			zap.String("ffmpeg_cmd", b.FfmpegCmd),
			zap.Uint("id", b.ID),
		)
		b.DoneStatus = consts.RunIngError
	}

	b.Update()

	b.CmdDoneCallBack()

	return nil
}

// FfmpegDownloadMix4Video 执行下载视频
func (b *BurnInfoCmd) FfmpegDownloadMix4Video() error {
	// 执行命令
	ctx, _ := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, config.C.Ffmpeg.LibPath, b.FfmpegCmdArgs...)

	cmd.StdinPipe()

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	if err := cmd.Start(); err != nil {
		log.L.Error("下载4合一画面ffmpeg命令失败", zap.String("ffmpeg_cmd", b.FfmpegCmd), zap.Any("状态", consts.RunIngError))
		return err
	}

	// 等待命令执行完成
	err := cmd.Wait()
	closeTime := time.Now()
	b.CloseTime = &closeTime

	if err == nil {
		log.L.Info("下载4合一画面ffmpeg命令过程中正常结束", zap.String("ffmpeg_cmd", b.FfmpegCmd), zap.Any("状态", consts.Success))
		b.DoneStatus = consts.Success
	} else {
		log.L.Error("下载4合一画面ffmpeg命令过程中异常结束", zap.String("ffmpeg_cmd", b.FfmpegCmd))
		b.DoneStatus = consts.RunIngError
	}

	b.Update()

	b.CmdDoneCallBack()

	return nil
}

// FfmpegDownloadRoomOneVideo 执行下载房间单画面视频
func (b *BurnInfoCmd) FfmpegDownloadRoomOneVideo() error {
	// 执行命令
	ctx, _ := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, config.C.Ffmpeg.LibPath, b.FfmpegCmdArgs...)

	cmd.StdinPipe()

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	if err := cmd.Start(); err != nil {
		log.L.Error("下载房间单画面ffmpeg命令失败", zap.String("ffmpeg_cmd", b.FfmpegCmd), zap.Any("状态", consts.RunIngError))
		return err
	}

	// 等待命令执行完成
	err := cmd.Wait()
	closeTime := time.Now()
	b.CloseTime = &closeTime

	if err == nil {
		log.L.Info("下载房间单画面ffmpeg命令过程中正常结束", zap.String("ffmpeg_cmd", b.FfmpegCmd), zap.Any("状态", consts.Success))
		b.DoneStatus = consts.Success
	} else {
		log.L.Error("下载房间单画面ffmpeg命令过程中异常结束", zap.String("ffmpeg_cmd", b.FfmpegCmd))
		b.DoneStatus = consts.RunIngError
	}

	b.Update()

	b.CmdDoneCallBack()

	return nil
}

// FfmpegDownloadPublicSingleOneVideo 执行下载公区单画面视频
func (b *BurnInfoCmd) FfmpegDownloadPublicSingleOneVideo() error {
	// 执行命令
	ctx, _ := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, config.C.Ffmpeg.LibPath, b.FfmpegCmdArgs...)

	cmd.StdinPipe()

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	if err := cmd.Start(); err != nil {
		log.L.Error("下载公区单画面ffmpeg命令失败", zap.String("ffmpeg_cmd", b.FfmpegCmd), zap.Any("状态", consts.RunIngError))
		return err
	}

	// 等待命令执行完成
	err := cmd.Wait()
	closeTime := time.Now()
	b.CloseTime = &closeTime

	if err == nil {
		log.L.Info("下载公区单画面ffmpeg命令过程中正常结束", zap.String("ffmpeg_cmd", b.FfmpegCmd), zap.Any("状态", consts.Success))
		b.DoneStatus = consts.Success
	} else {
		log.L.Error("下载公区单画面ffmpeg命令过程中异常结束", zap.String("ffmpeg_cmd", b.FfmpegCmd))
		b.DoneStatus = consts.RunIngError
	}

	b.Update()

	b.CmdDoneCallBack()

	return nil
}

// CmdDoneCallBack 子任务完成后续操作
func (b *BurnInfoCmd) CmdDoneCallBack() {

	// 查询burnInfo
	burnInfo := burninfo.GetById(b.BurnInfoID)
	if burnInfo.ID == 0 {
		log.L.Error("BurnInfo 不存在", zap.Uint("id", b.BurnInfoID))
		return
	}

	// 未完成任务数量减一
	burnInfo.ReduceUndoneNum(1)

	// 查看该下载根任务是否已经完成
	burnInfo = burninfo.GetById(b.BurnInfoID)
	if burnInfo.UndoneNum <= 0 {
		// 所有子任务已完成 ,剪切视频文件到oda存储路径
		helper.CopyDir(burnInfo.SaveFileTmpPath, burnInfo.OdaSavePath, true)

		// 回调业务端接口
		err := httpclient.CallBackHttpClient(burnInfo.CallbackUrl, burnInfo.Uuid)
		if err != nil {
			burnInfo.CallbackStatus = consts.RunIngError
		} else {
			burnInfo.CallbackStatus = consts.Success
		}
		burnInfo.Update()

		// 刻录任务完成数+1
		burnsetting.AddDoneTaskNum(burnInfo.BurnSettingID)

	}

}
