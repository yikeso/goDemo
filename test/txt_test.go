package demotest

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yikeso/goDemo/txt"
)

func TestBase64EncodeTxtFile(t *testing.T){
	Convey("测试base64加密txt文件",t,func(){
		err := txt.Base64EncodeTxtFile("K:/computer.txt")
		So(err,ShouldBeNil)
	})
}

func TestBase64DecodeTxtFile(t *testing.T){
	Convey("测试base64解密encode文件",t,func(){
		err := txt.Base64DecodeEncodeFile("K:/computer.encode")
		So(err,ShouldBeNil)
	})
}