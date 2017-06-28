package txt

import (
	"os"
	"encoding/base64"
	"strings"
	"fmt"
	"io"
)

//对txt文件内容进行base64加密生成.encode文件
func Base64EncodeTxtFile(filePath string)(err error){
	f,err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return err
	}
	writePath := filePath[:strings.LastIndex(filePath,".")]
	writePath = fmt.Sprint(writePath,".encode")
	wf,err := os.Create(writePath)
	if err != nil{
		return err
	}
	defer wf.Close()
	wc := base64.NewEncoder(base64.StdEncoding,wf)
	_,err = io.Copy(wc,f)
	return err
}
//对encode文件内容进行base64解密生成.txt文件
func Base64DecodeEncodeFile(filePath string)(err error){
	f,err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return err
	}
	r := base64.NewDecoder(base64.StdEncoding,f)
	if err != nil {
		return err
	}
	writePath := filePath[:strings.LastIndex(filePath,".")]
	writePath = fmt.Sprint(writePath,".txt")
	wf,err := os.Create(writePath)
	defer wf.Close()
	if err != nil{
		return err
	}
	_,err = io.Copy(wf,r)
	return err
}