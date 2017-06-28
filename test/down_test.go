package demotest

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yikeso/goDemo/down"
)

func TestDownloadUrlFile(t *testing.T){
	Convey("测试文件下载方法",t,func(){
		err := down.DownloadUrlFile("http://mirrors.sohu.com/centos/7/isos/x86_64/CentOS-7-x86_64-Minimal-1611.iso","e:/dlp")
		So(err,ShouldBeNil)
	})
}