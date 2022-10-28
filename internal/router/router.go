// 在这个文件中注册 URL handler

package router

import (
	"github.com/PittYao/stream_burn/components/gin/response"
	"github.com/PittYao/stream_burn/internal/api"
	"github.com/PittYao/stream_burn/internal/httpclient"
	"github.com/gin-gonic/gin"
)

// Routes 注册 API URL 路由
func Routes(app *gin.Engine) {
	group := app.Group("/api")
	{
		v1 := group.Group("/v1")
		{
			v1.POST("burnTask", api.BurnTask)
			v1.POST("burnMixVideo", api.BurnMix3Video)
			//v1.POST("burnMixVideo4to1", api.BurnMixVideo4To1)
			v1.POST("burnSingleVideo", api.BurnSingleVideo)
			//v1.POST("burnOtherSingleVideo", api.BurnOtherSingleVideo)
			//v1.POST("burnFile", api.BurnFile)
			//// 3合一 放开此注释 [兼容老版本]
			//v1.POST("burnParams", api.BurnParams)
			//// 4合一是该接口
			////v1.POST("burnParams", api.Burn41Params)
			//v1.POST("burnSingleParams", api.BurnSingleParams)
			//v1.POST("burnOtherSingleParams", api.BurnOtherSingleParams)

			v1.POST("callback", func(c *gin.Context) {
				var callbackDTO httpclient.CallbackDTO
				err := c.ShouldBindJSON(&callbackDTO)
				if err != nil {
					response.Err(c, err.Error())
					return
				}

				response.OK(c, nil)
			})

		}

	}
}
