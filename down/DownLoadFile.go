package down

import (
	"container/list"
	"sync"
	"fmt"
	"time"
	"os"
	"log"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
)

var mutex sync.Mutex

/**
 * 线程池执行方法的接口
 */
type runnable interface {
	run()
}

/**
 * 下载任务观察者，监听下载任务的进度,
 * 下载完成后，合并缓存文件
 */
type listenDownloadLength struct {
	Size int64
	DownloadFilePath string
	FileName string
	PreSecondLength int64
	TempFileList list.List
	DownloadLength int64
}

/**
 * 该方法让下载任务消费者，汇报下载量
 */
func (l *listenDownloadLength)add(p int64){
	atomic.AddInt64(&l.DownloadLength, p)

}

func (l *listenDownloadLength) run()  {
	for l.DownloadLength < l.Size {
		var v float64
		mutex.Lock()
		v = float64(l.DownloadLength - l.PreSecondLength)/1024.0
		l.PreSecondLength = l.DownloadLength
		mutex.Unlock()
		if v > 1024 {
			v /= 1024
			log.Println(l.FileName + "\\/" + fmt.Sprint("%.3f", v) + "M/s")
		}else{
			log.Println(l.FileName + "\\/" + fmt.Sprint("%.3f", v) + "KB/s")
		}
		time.Sleep(time.Second)
	}
	log.Println(l.FileName + " 网络资源拉取到本地完成。开始合并本地缓存文件")
	_,err := os.Stat(l.DownloadFilePath)
	if err != nil{
		os.Create(l.DownloadFilePath)
	}
	f,err := os.OpenFile(l.DownloadFilePath,os.O_APPEND,0666)
	defer f.Close()
	if err != nil {
		log.Println(l.FileName + " 下载失败")
		log.Println(err.Error())
		return
	}
	var subf *os.File
	for e := l.TempFileList.Front(); e != nil; e = e.Next() {
		subf,err = os.OpenFile(fmt.Sprint(e.Value),os.O_RDONLY,0666)
		if err != nil{
			log.Println(l.FileName + " 缓存文件合并失败")
			break
		}
		io.Copy(f,subf)
		subf.Close()
	}
	log.Println(l.FileName + " 下载成功")
}

type downloadConsumer struct {
	FileUrl string
	TempFilePath string
	Position int64
	End int64
	Listen *listenDownloadLength
}

func (d *downloadConsumer) run(){
	createFileIfNotExists(d.TempFilePath)
	f,err := os.OpenFile(d.TempFilePath,os.O_WRONLY,0666)
	defer f.Close()
	if err != nil {
		log.Println(d.TempFilePath + " 下载失败")
		log.Println(d.FileUrl + " 下载失败")
		log.Println(err.Error())
		return
	}
	req,err := http.NewRequest("GET",d.FileUrl,nil)
	if err != nil {
		log.Println(d.TempFilePath + " 下载失败")
		log.Println(d.FileUrl + " 下载失败")
		log.Println(err.Error())
		return
	}
	client := http.DefaultClient
	v := fmt.Sprint("bytes=",d.Position,"-" ,d.End)
	req.Header.Add("Range", v)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(d.TempFilePath + " 下载失败")
		log.Println(d.FileUrl + " 下载失败")
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	io.Copy(f,resp.Body)
}

func createFileIfNotExists(path string){
	_,err := os.Stat(path)
	if err != nil{
		os.Create(path)
	}
}

func DownloadUrlFile(fileUrl string,dir string)(err error){
	makeDirIfNotExists(dir)
	downConsumer := downloadConsumer{}
	downConsumer.FileUrl = fileUrl
	fileName := fileUrl[strings.LastIndex(fileUrl,"/")+1:]
	if strings.LastIndex(fileName,"?") > 0 {
		fileName = fileName[0:strings.LastIndex(fileName,"?")-1]
	}
	downConsumer.TempFilePath = fmt.Sprint(dir,"/",fileName)
	downConsumer.Position = 0
	resp,err := http.Get(fileUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	downConsumer.End = resp.ContentLength
	runTask(&downConsumer)
	return err
}

func makeDirIfNotExists(path string){
	_,err := os.Stat(path)
	if err != nil{
		os.MkdirAll(path,0777)
	}
}

func runTask(t runnable){
	t.run()
}