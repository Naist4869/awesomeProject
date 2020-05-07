package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Naist4869/awesomeProject/tool"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func Init() {
	var err error
	_, err = LoadAppConfig()
	if err != nil {
		log.Fatal(err)
	}
}

/*LoadAppConfig
加载服务配置  3个文件会相互覆盖

*	*Config	*Config
*	error  	error
*/
func LoadAppConfig() (*AppConfig, error) {
	// 对于debug模式，其实就是测试模式，工作目录是单个的模块目录，那么需要进入具备main.go的目录
	files := make(map[string]*os.File, 10)
	defer func() {
		for _, file := range files {
			file.Close()
		}
	}()

	if tool.IsDebug() {
		ok := false
		for wd, _ := os.Getwd(); !ok; wd = filepath.Dir(wd) {
			log.Printf("判断目录[%s]是否为主目录\n", wd)
			file, err := os.Open(wd)
			if err != nil {
				return nil, errors.Wrap(err, "定位主目录")
			}
			files[wd] = file
			// 读取文件夹里的文件
			fileInfos, err := file.Readdir(-1)
			// 读取目录中间有错误发生
			if err != nil {
				return nil, errors.Wrap(err, "读取目录信息")
			}
			// 一路读到目录的末尾
			for _, info := range fileInfos {
				if strings.Contains(info.Name(), "main.go") {
					err := os.Chdir(wd)
					if err != nil {
						return nil, errors.Wrap(err, "更改目录信息")
					}
					ok = true
					break
				}
			}
		}
	}
	appConfig := &AppConfig{}
	if err := loadServerConfig(filePath, appConfig, false); err != nil {
		return nil, err
	}
	if err := loadServerConfig(customFilePath, appConfig, true); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	if err := loadServerConfig(productionFilePath, appConfig, true); err != nil {
		return nil, err
	}
	return appConfig, nil
}

/*loadServerConfig 加载并解析指定文件中的配置
参数:
*       path            string          文件路径
*       cfg             *ServerConfig   写入目标
*       allowError      bool            是否允许错误
返回值:
*       error   error
*/
func loadServerConfig(path string, cfg *AppConfig, allowError bool) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		if allowError {
			return nil
		}
		log.Fatal("ReadFile: ", err.Error())
		return err
	}
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		log.Fatal("Unmarshal: ", err.Error())
		return err
	}
	return nil
}
