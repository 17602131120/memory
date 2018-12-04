package setting

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type MMSettings struct {
	MMLogger *log.Logger

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	}
	Mongo struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Db   string `yaml:"db"`
	}

	Config struct {
		//从settings.yaml 解析
		Botname                 string   `yaml:"botname"`
		ConcurrentRequest       int      `yaml:"concurrentRequest"`
		ConcurrentRequestSleep  int      `yaml:"concurrentRequestSleep"`
		ConcurrentPipeline      int      `yaml:"concurrentPipeline"`
		ConcurrentPipelineSleep int      `yaml:"concurrentPipelineSleep"`
		Debug                   bool     `yaml:"debug"`
		LogPath                 string   `yaml:"logPath"`
		//Keys                    []string `yaml:"keys,flow"` //yaml配置 keys : [a,b,c,d]
	}
}

// 判断文件夹是否存在
func (this *MMSettings) PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (this *MMSettings) PathCreate(path string) bool {

	if !this.PathExists(path) {

		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return true
		} else {
			return false
		}
	} else {
		return true
	}

}

var settings *MMSettings
var settingsOnce sync.Once

func MMSettingsSington() *MMSettings {

	settingsOnce.Do(func() {
		//log.Println("MMSettingsSington")

		settings = new(MMSettings)

		//读取yaml配置文件

		yamlFile, err := ioutil.ReadFile("settings.yaml")

		if err != nil {
			//读取错误
			log.Fatalf("未找到配置文件:#%v ", err)

		} else {
			//读取成功
			err = yaml.Unmarshal(yamlFile, settings)
			if err != nil {
				//转换错误
				log.Fatalf("解析yaml文件失败: %v", err)
			} else {
				logPath := settings.Config.LogPath
				settings.PathCreate(logPath)

				writers := []io.Writer{}
				if settings.Config.Debug {
					//debug模式，打印到控制台
					writers = append(writers, os.Stdout)

				} else {
					//文件日志
					logFile, err := os.Create(logPath + "/debug_" + time.Now().Format("20060102_150405") + ".log")
					//defer logFile.Close()
					if err != nil {
						log.Fatalln("open log file error !")
					}
					writers = append(writers, logFile)
				}

				multWriter := io.MultiWriter(writers...)

				settings.MMLogger = log.New(multWriter, "", log.LstdFlags)
				settings.MMLogger.SetFlags(log.Lshortfile)
				settings.MMLogger.SetFlags(settings.MMLogger.Flags() | log.LstdFlags)

			}

		}

	})

	return settings

}
