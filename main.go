package main

import (
	"regexp"
	"flag"
	"os"
	"bufio"
	"io"
	"path/filepath"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"./utils"
	"bytes"
	"strings"
)

var (
	idCtt = []string{}
	output = "./output"
)

func main() {
	var date, filePath string

	flag.StringVar(&date, "date", "", "20180101")
	flag.StringVar(&filePath, "path", "", "/path/to/your/lof/file")
	flag.StringVar(&output, "output", "./output", "/path/to/output/file")
	flag.Parse()

	if filePath == "" {
		fmt.Println("Invild log path string.")
	}

	if date == "" {
		date = utils.GetDateString(time.Now())
	}

	root := utils.GetCwd()
	if !path.IsAbs(filePath) {
		filePath = filepath.Join(root, filePath)
	}
	if !path.IsAbs(output) {
		output = filepath.Join(root, output)
	}

	// 创建out文件夹
	utils.EnsureDir(output)

	readDir(filePath, date)
}

// 读取文件内容   *.gzhxy
func readFiles(fullPath string) {
	fi, err := os.Open(fullPath)
	utils.ErrHadle(err)
	defer fi.Close()

	// idReg := regexp.MustCompile("qid=(\\d+).+(\"error\":0)+")
	idReg := regexp.MustCompile("qid=(\\d+).+(%22error%22%3A0)+")

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		logsArr := bytes.Split(a, []byte(" "))
		if len(logsArr) > 6 {
			logByte := logsArr[6]
			ok, err := regexp.Match("/mpwb.gif", logByte)
			if ok && err == nil {
				ids := idReg.FindAllSubmatch(logByte, -1)
				if ids != nil && len(ids[0]) > 0{
					id := string(ids[0][1])
					idCtt = append(idCtt, id)
				}
			}
		}
	}
}


// 读取指定目录
// wb_log
func readDir(filePath string, date string) {
	var prefix = "access_webb_wise_mip.log."
	var nameRe = regexp.MustCompile(prefix + date + "(\\d{2})")

	// 清理日志存储目录
	//cleanTmp(path)
	files, err := ioutil.ReadDir(filePath)
	utils.ErrHadle(err)
	for _, file := range files {
		fileName := file.Name()
		if nameRe.MatchString(fileName) {
			readList(filepath.Join(filePath, fileName))
		}
	}
}

// 获取机器列表
// m.2018*
func readList(name string) {
	nameArr := strings.Split(name, ".")
	oFileName := nameArr[len(nameArr)-1]
	files, err := ioutil.ReadDir(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range files {
		fileName := file.Name()
		readFiles(filepath.Join(name, fileName))
	}

	fmt.Println(idCtt)
	writeData(filepath.Join(output, oFileName), &idCtt)
	idCtt = []string{}
}

func writeData(path string, data *[]string) {
	rsData := strings.Join(*data, "\n")

	e := ioutil.WriteFile(path, []byte(rsData), 0777)
	if e != nil {
		fmt.Println("Write " + path + "successfully.")
	}
}
