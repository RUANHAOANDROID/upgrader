package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
func Auto(ctx context.Context, cancel context.CancelFunc) {
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
	if resp.Code == 1 {
		fileName := resp.Data.FileName
		version := resp.Data.VersionCode
		downloadUrl := resp.Data.DownloadUrl
		pkg.Log.Infof("发现新版本 version=%v ,file=%v", version, fileName)
		conf.Version = version
		conf.Save("config.yml")
		// 创建备份目录和临时目录
		if err := checkDIr(conf.TempDir); err != nil {
			pkg.Log.Error("创建临时目录错误：" + err.Error())
			return
		}

		// 下载更新包
		pkg.Log.Println("下载更新包...")
		filePath := filepath.Join(conf.TempDir, fileName)
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
		if err := copyFileTo(conf.BackupDir, fileName, tarFile.Name()); err != nil {
			pkg.Log.Error("拷贝更新包错误：" + err.Error())
			return
		}
		// 删除运行目录下所有文件
		if _, err := os.Stat(conf.RunnerDir); os.IsNotExist(err) {
			fmt.Printf("Directory %s does not exist\n", conf.RunnerDir)
		}
		err = os.RemoveAll(conf.RunnerDir)
		if err != nil {
			fmt.Printf("Failed to remove directory: %v\n", err)
		}
		// 解压到运行目录
		if err := extractTar(conf.RunnerDir, fileName, tarFile.Name()); err != nil {
			pkg.Log.Error("解压更新包错误：" + err.Error())
			return
		}
		//更名文件
		old := filepath.Join(conf.RunnerDir, strings.TrimSuffix(fileName, ".tar"))
		newDir := filepath.Join(conf.RunnerDir, "app")
		err = os.Rename(old, newDir)
		if err != nil {
			fmt.Println("重命名失败:", err)
			return
		}
		cancel()
		runScript(ctx)
	} else {
		pkg.Log.Println("未发现更新")
	}
}

func copyFileTo(targetDir string, targetFile string, sourceFile string) error {
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	backFile, err := os.Create(targetDir + "/" + targetFile)
	if err != nil {
		pkg.Log.Error("创建文件错误：" + err.Error())
		return err
	}
	defer backFile.Close()
	//拷贝到目录
	if _, err := io.Copy(backFile, source); err != nil {
		pkg.Log.Printf("无法拷贝 %s: %v\n", source.Name(), err)
		return err
	}
	return nil
}
func extractTar(targetDir string, fileName, tarFile string) error {
	source, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer source.Close()
	err = pkg.UnTar(targetDir, tarFile, false)
	if err != nil {
		pkg.Log.Println("解压缩失败" + err.Error())
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
