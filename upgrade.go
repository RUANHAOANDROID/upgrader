package main

import (
	"context"
	"encoding/json"
	"github.com/mholt/archiver/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"upgrader/config"
	"upgrader/pkg"
)

type VersionResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		VersionCode string `json:"versionCode"`
		DownloadUrl string `json:"downloadUrl"`
		FileName    string `json:"fileName"`
	} `json:"data"`
}

// Auto 假定初始时没有包
func Auto() {
	// 读取本地版本信息
	conf, err := config.Load("./config.yml")
	if err != nil {
		panic("配置文件错误")
	}

	pkg.Log.Infoln(conf)
	// 获取版本信息
	url := conf.URL + "/gateMachine/queryVersion/" + conf.AuthCode + "/" + conf.Version
	pkg.Log.Println(url)
	response, err := http.Get(url)
	if err != nil {
		pkg.Log.Error(err)
		return
	}
	defer response.Body.Close()

	var resp VersionResponse
	body, err := io.ReadAll(response.Body)
	if err != nil {
		pkg.Log.Error(err)
		return
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		pkg.Log.Error(err)
		return
	}
	//pkg.Log.Fatalf("check update resp code=%d,msg=%s", resp.Code, resp.Msg)
	//pkg.Log.Fatalf("check update resp code=%d,msg=%s,%s", resp.Code, resp.Msg, resp.Data)
	pkg.Log.Infof("check update resp %v", resp)
	// 解析JSON响应
	if resp.Code == 1 {
		fileNam := resp.Data.FileName
		version := resp.Data.VersionCode
		downloadUrl := resp.Data.DownloadUrl
		pkg.Log.Infof("发现新版本 version=%v ,file=%v", version, fileNam)

		// 创建备份目录和临时目录
		if err := checkDIr(conf.TempDir); err != nil {
			pkg.Log.Error("创建临时目录错误：" + err.Error())
			return
		}

		// 下载更新包
		pkg.Log.Println("下载更新包...")
		filePath := filepath.Join(conf.TempDir, fileNam)
		tarFile, err := os.Create(filePath)
		if err != nil {
			pkg.Log.Error("创建文件错误：" + err.Error())
			return
		}
		defer tarFile.Close()

		if err := DownloadFile(tarFile, downloadUrl, filePath); err != nil {
			pkg.Log.Error("下载文件错误：" + err.Error())
			return
		}

		// 创建备份目录并拷贝更新包
		if err := checkDIr(conf.BackupDir); err != nil {
			pkg.Log.Error("创建备份目录错误：" + err.Error())
			return
		}
		if err := copyFileTo(conf.BackupDir, fileNam, tarFile); err != nil {
			pkg.Log.Error("拷贝更新包错误：" + err.Error())
			return
		}

		// 解压到运行目录
		if err := extractTar(tarFile); err != nil {
			pkg.Log.Error("解压更新包错误：" + err.Error())
			return
		}
	} else {
		pkg.Log.Println("未发现更新")
	}
}

func copyFileTo(targetDir string, fileName string, source *os.File) error {
	// 创建目标文件
	destination, err := os.Create(targetDir + "/" + fileName)
	if err != nil {
		pkg.Log.Printf("无法创建目标文件 %s: %v\n", fileName, err)
		return err
	}
	defer destination.Close()
	//拷贝到目录
	if _, err := io.Copy(destination, source); err != nil {
		pkg.Log.Printf("无法拷贝 %s: %v\n", source.Name(), err)
		return err
	}
	return nil
}
func extractTar(tarFile *os.File) error {
	format := archiver.Tar{}
	err := format.Extract(context.Background(), tarFile, []string{}, nil)
	if err != nil {
		return err
	}
	pkg.Log.Println("解压缩完成")
	return nil
}
func checkDIr(dir string) error {
	// 检查文件夹是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 如果文件夹不存在，则创建它
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			pkg.Log.Printf("无法创建文件夹 %s: %v\n", dir, err)
			return err
		}
		pkg.Log.Printf("已创建文件夹 %s\n", dir)
	} else {
		pkg.Log.Printf("文件夹 %s 已存在\n", dir)
	}
	return nil
}
func DownloadFile(file *os.File, url, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	return err
}
