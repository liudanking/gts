package main

import (
	"flag"
	"net/http"

	"github.com/liudanking/gotranslate"

	"github.com/gin-gonic/gin"
	log "github.com/liudanking/goutil/logutil"
)

var _gt *gotranslate.GTranslate

func main() {
	var laddr string
	var debug bool
	flag.StringVar(&laddr, "l", "localhost:9080", "listen address")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.Parse()

	var err error
	_gt, err = gotranslate.New(gotranslate.TRANSLATE_CN_ADDR, nil)
	if err != nil {
		log.Error("new gotranslate error:%v", err)
		return
	}

	// if !debug {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	engine := gin.New()
	engine.Use(gin.Logger())

	engine.GET("/simple_translate", SimpleTranslateHanlder)

	log.Notice("start serving at [%s]", laddr)
	if err = engine.Run(laddr); err != nil {
		log.Error("gin engine run error:%v", err)
	}

}

type SimpleTranslateReq struct {
	SL string `form:"sl" json:"sl"`
	TL string `form:"tl" binding:"required" json:"tl"`
	Q  string `form:"q" binding:"required" json:"q"`
}

func SimpleTranslateHanlder(c *gin.Context) {
	req := &SimpleTranslateReq{}
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code": 1000,
			"err_msg":  err.Error(),
		})
		return
	}

	if req.SL == "" {
		req.SL = "auto"
	}

	ret, err := _gt.SimpleTranslate(req.SL, req.TL, req.Q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code": 1001,
			"err_msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ret": ret})
	return
}
