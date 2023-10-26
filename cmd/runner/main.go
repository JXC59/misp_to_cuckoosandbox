package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jx259/misp"
	"github.com/yeka/zip"
)

func main() {
	// 连接MISP （所谓的病毒库）
	mispClient, err := misp.NewClient(&misp.ClientConfig{
		Endpoint:      "your url",
		SkipTLSVerify: true,
		APIKey:        "your key",
	})
	print("MispClient:")
	println(mispClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*60)
	defer cancel()

	// 查找事件通过Tag，同时可以获取id
	Response, err := mispClient.SearchEvents(ctx, &misp.SearchEventsRequest{
		Limit: 500,
		Page:  1,
		Tag:   tagElfFile,
		// 限定日期（可自行更改）
		From: "2023-10-15",
		To:   "2023-10-18",
	})
	if err != nil {
		fmt.Print(err)
	}

	// 打印返回来的病毒个数。
	fmt.Print("this is len:")
	fmt.Println(len(Response))

	// 下载病毒文件.
	for _, value := range Response {
		for _, obj := range value.Event.Objects {
			for _, atb := range obj.Attributes {
				if strings.Contains(atb.Type, "malware-sample") {
					downAssistant(mispClient, ctx, atb.ID, atb.EventID, atb.Value)
				}
			}
		}
	}

	// 下载完成，稍等一分钟。
	fmt.Println("Downloaded , Please wait 1 Minute")
	time.Sleep(1 * time.Minute)

	// 解压病毒。
	zipDir := "/root/jx_codes/mispToCuckoo/zip/"
	destDir := "/root/jx_codes/mispToCuckoo/elf/"
	password := "infected"
	files, err := filepath.Glob(filepath.Join(zipDir, "*.zip"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, zipFile := range files {
		err := unzip(zipFile, destDir, password)
		if err != nil {
			fmt.Printf("Failed to unzip %s: %v\n", zipFile, err)
		}
	}

	// 清理删除病毒压缩包
	err = clearDirectory(zipDir)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Directory cleared successfully!")
	}

	// 读取目录下的病毒文件（已解压），上传至Sandbox
	entries, err := ioutil.ReadDir(destDir)
	if err != nil {
		fmt.Println(err)
	}
	i := 1
	for _, entry := range entries {
		if i > 3 {
			i = 1
			time.Sleep(5 * time.Minute)
			fmt.Println("Next round")
		}
		filePath := filepath.Join(destDir, entry.Name())
		taskID, err := CreateTask(filePath, "your cuckoo", "your cuckoo key", entry.Name(), 1)
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println(taskID)
		i = i + 1
		time.Sleep(10 * time.Second)
	}
	err = clearDirectory(destDir)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Directory cleared successfully!")
	}
	fmt.Println("Finished Finished Finished Finished")
}

func downAssistant(c *misp.Client, ctx context.Context, atbID string, eveID string, atbValue string) {
	dl, err := c.DownloadAttribute(ctx, atbID)
	if err != nil {
		fmt.Println(err)
	}
	path := "mispToCuckoo/zip/"
	file, err := os.Create(path + atbID + ".zip")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, err = io.Copy(file, dl)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("Downloaded eveID:：" + eveID)
}

func unzip(zipFile, destDir, password string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(password)
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fPath := filepath.Join(destDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
		} else {
			// 只解压没有后缀的文件
			if strings.Contains(f.Name, ".") {
				continue
			}

			var dir string
			if lastIndex := strings.LastIndex(fPath, string(os.PathSeparator)); lastIndex > -1 {
				dir = fPath[:lastIndex]
			}

			err = os.MkdirAll(dir, f.Mode())
			if err != nil {
				return err
			}

			outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CreateTask 通过给定的文件路径和REST URL创建一个任务，并返回任务ID和可能的错误
func CreateTask(filePath string, restURL string, token string, fileName string, machineNum int) (float64, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// 创建 multipart 写入器
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建 form 文件
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return 0, fmt.Errorf("error creating form file: %w", err)
	}

	// 将文件写入 form 文件
	_, err = io.Copy(part, file)
	if err != nil {
		return 0, fmt.Errorf("error writing to form file: %w", err)
	}

	str := "Ubuntu1804x64_"
	strs := str + strconv.Itoa(machineNum)
	fmt.Println(strs)
	err = writer.WriteField("machine", strs)
	if err != nil {
		return 0, fmt.Errorf("errorwriting machine field: %w", err)
	}

	err = writer.WriteField("timeout", strconv.Itoa(480))
	if err != nil {
		return 0, fmt.Errorf("errorwriting timeout field: %w", err)
	}

	// 关闭写入器
	err = writer.Close()
	if err != nil {
		return 0, fmt.Errorf("error closing writer: %w", err)
	}

	// 创建一个新的HTTP请求
	req, err := http.NewRequest("POST", restURL, body)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}

	// 设置Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 添加token到请求头
	req.Header.Set("Authorization", "Bearer "+token)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 创建一个HTTP客户端
	client := &http.Client{Transport: tr}

	// 发送POST请求
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error response status: %s", resp.Status)
	}

	// 读取响应体的内容
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %w", err)
	}

	// 尝试解析 JSON 响应
	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return 0, fmt.Errorf("error decoding JSON: %w", err)
	}
	// 获取 task_id
	taskID, ok := result["task_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("error getting task_id")
	}
	return taskID, nil
}

func clearDirectory(dir string) error {
	// 读取目录下的所有文件和子目录
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	// 遍历并删除每一个文件或子目录
	for _, entry := range entries {
		filePath := filepath.Join(dir, entry.Name())
		err = os.RemoveAll(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	trafficCaptured          = `secheart:analysis="traffic-captured"`
	tagSandboxPrefix         = `secheart:sandbox=`
	tagFileTypeAndroid       = `file-type:type="android"`
	tagSandboxAnalyzed       = `secheart:sandbox="avd-trace-ok"`
	tagSandboxAnalyzeFailed  = `secheart:sandbox="avd-trace-failed"`
	tagSandboxAnalyzedLegacy = `secheart:sandbox="success"`
	tagElfFile               = `file-type:type="elf64"`
)

const (
	CuckoosandboxUrl = "http://f5fc13080d.duckdns.org:9500/tasks/create/file"
)
