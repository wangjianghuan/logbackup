package main

import (
	"fmt"
	"os"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"time"
	"sync"
	"io/ioutil"
	"strings"
	"runtime"
)

var endpoint, accessKeyId, accessKeySecret, bucketName string

func main() {
	fmt.Println("OSS Go SDK Version: ", oss.Version)
	// 配置信息检查
	envCheck()
	// 文件上传备份
	backup()
}

func envCheck() {
	endpoint = os.Getenv("ENDPOINT")
	accessKeyId = os.Getenv("ACCESSKEYID")
	accessKeySecret = os.Getenv("ACCESSKEYSECRET")
	bucketName = os.Getenv("BUCKETNAME")
	if endpoint == "" {
		panic("The ENDPOINT environment variable is not set.")
	}
	if accessKeyId == "" {
		panic("The ACCESSKEYID environment variable is not set.")
	}
	if accessKeySecret == "" {
		panic("The ACCESSKEYSECRET environment variable is not set.")
	}
	if bucketName == "" {
		panic("The BUCKETNAME environment variable is not set.")
	}
}

// 日志备份
func backup() {
	// 定时任务
	ticker := time.NewTicker(60 * time.Second)
	timeFormat := "2006-01-02 15:04:05"
	quit := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	// 业务变量
	var flag bool
	// 程序启动第一次计算时间用
	next := time.Now().Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), 1, 0, 0, 0, next.Location())
	fmt.Println("The next time is " + next.Format(timeFormat))
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				now := time.Now().Local()
				if !flag {
					// 计算下一次时间点
					next = now.Add(time.Hour * 24)
					//next = time.Now().Local()
					next = time.Date(next.Year(), next.Month(), next.Day(), 1, 0, 0, 0, next.Location())
					flag = true
					fmt.Println("The next time is " + next.Format(timeFormat))
				}
				// 时间比较做业务
				//dValue := next.Sub(now).Seconds()
				//fmt.Println(dValue)
				fmt.Println(now)
				if now.Format("2006-01-02 15:04") == next.Format("2006-01-02 15:04") {
					walkDir("D:\\")
					//walkDir("/home/logs")
					flag = false
				}
				fmt.Println("Wait for to execute,next time is " + next.Format(timeFormat))
			case <-quit:
				fmt.Println("work well .")
				ticker.Stop()
				return
			}
		}
		fmt.Println("child goroutine bootstrap end")
	}()
	wg.Wait()
}

// 读取文件目录
func walkDir(dirpath string) {
	files, err := ioutil.ReadDir(dirpath) //读取目录下文件
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			fileName := file.Name()
			lastTime := time.Now();
			lastDay := time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day()-1, 0, 0, 0, 0, lastTime.Location()).Format("20060102")
			fmt.Println("last:" + lastDay)
			moveFile := "alarm-" + lastDay
			fmt.Println("fileName:" + fileName)
			fmt.Println("moveFile:" + moveFile)
			if strings.Index(fileName, moveFile) >= 0 {
				// TODO 上传文件
				//upload(dirpath,fileName)
				go ReadFile(dirpath + "//" + fileName)
				fmt.Println(fileName + " upload success!")
			}
		}
	}
}

// 异常输出
func handleError(err error) {
	fmt.Println("Error:", err)
}

// 文件上传 <yourLocalFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
func upload(localFilePath string, localFileName string) {

	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	objectName := "<yourObjectName>"

	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		handleError(err)
	}

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFilePath+"//"+localFileName)
	if err != nil {
		handleError(err)
	} else {
		fmt.Println("文件上传成功")
		runtime.Goexit()
		// TODO 删除本地文件
	}
}

// 读取文件并输出
func ReadFile(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("File reading error", err)
	}
	fmt.Println(string(data))
	fmt.Println("=========================" + filePath + "读取完毕======================")
	runtime.Goexit()
}
