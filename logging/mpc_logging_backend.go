package logging

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FileLogConfig struct {
	Filename string //日志文件名
	Maxsize  int    //每个日志文件最大大小，达到后轮转文件，单位为字节，0为不限制
	Maxfiles int    //最多保留的文件个数，写满后会清除最旧的文件
	Rotate   bool   //日志文件轮转的总开关
}

// 实现了os.Writer接口
type FileLogWriter struct {
	*log.Logger
	*FileLogConfig

	mw              *MuxWriter
	daily_opendate  int
	startLock       sync.Mutex
}

//实现io.Writer接口，供go-logging调用
func (w FileLogWriter) Write(b []byte) (int, error) {
	w.rotate()
	return w.mw.Write(b)
}

var onlyOneConfig *FileLogConfig

type MuxWriter struct {
	sync.Mutex
	fd *os.File
}

func (l *MuxWriter) SetFd(fd *os.File) {
	if l.fd != nil {
		l.fd.Close()
	}
	l.fd = fd
}

//write to os.File
func (l *MuxWriter) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.fd.Write(b)
}

func (w *FileLogWriter) Init(conf *FileLogConfig) (err error) {
	w.FileLogConfig = conf
	onlyOneConfig = conf

	if len(w.Filename) == 0 {
		return errors.New("jsonconfig must have filename")
	}

	err = w.startLogger()
	return err
}

func (w *FileLogWriter) startLogger() error {
	fd, err := w.createLogFile()
	if err != nil {
		return err
	}
	w.mw.SetFd(fd)
	err = w.initFd()
	if err != nil {
		return err
	}
	return nil
}

//加锁，判断是否需要轮转，如需要则执行轮转
func (w *FileLogWriter) rotate() {
	w.startLock.Lock()
	defer w.startLock.Unlock()
	if w.needRotate() {
		err := w.DoRotate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
			return
		}
	}
}

//判断是否需要轮转到下一个文件
func (w *FileLogWriter) needRotate() bool {
	return w.Rotate && w.isExceedMaxsize()
}

func (w *FileLogWriter) createLogFile() (*os.File, error) {
	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	return fd, err
}

func (w *FileLogWriter) initFd() error {
	w.daily_opendate = time.Now().Day()
	return nil
}

//判断是否超出配置的最大文件大小
func (w *FileLogWriter) isExceedMaxsize() bool {
	info, err := w.mw.fd.Stat()
	if err == nil {
		return w.Maxsize > 0 && int(info.Size()) >= w.Maxsize
	} else {
		return false
	}
}

//进行文件轮转
//旧文件会被重命名为 <文件名>.<序号>，001为轮换的起始文件
func (w *FileLogWriter) DoRotate() interface{} {
	_, err := os.Lstat(w.Filename)
	if err == nil {
		//寻找下一个可用序号
		nextAvailableNum := 1
		newFilename := ""
		for {
			newFilename = w.Filename + fmt.Sprintf(".%03d", nextAvailableNum)
			_, err = os.Lstat(newFilename)

			if err == nil && nextAvailableNum <= w.Maxfiles {
				nextAvailableNum++
			} else {
				break
			}
		}

		//上锁，防止其他线程并发进行文件轮转
		w.mw.Lock()
		defer w.mw.Unlock()

		//当最大文件数还未满
		if actualFileNum := nextAvailableNum - 1; actualFileNum < w.Maxfiles {
			//重命名旧文件
			err = os.Rename(w.Filename, newFilename)
			if err != nil {
				return fmt.Errorf("Rotate: %s\n", err)
			}
		} else {
			for i := 1; i < w.Maxfiles; i++ {
				fname := w.Filename + fmt.Sprintf(".%03d", i)
				_, err = os.Lstat(fname)
				if err == nil {
					w.processOldFile(fname)
				}
			}
			newFilename = w.Filename + fmt.Sprintf(".%03d", w.Maxfiles)
			err = os.Rename(w.Filename, newFilename)
			if err != nil {
				return fmt.Errorf("Rotate: %s\n", err)
			}
		}

		//关闭旧文件句柄
		w.mw.fd.Close()

		//重启日志器，过程中会重新建立新日志文件
		err = w.startLogger()
		if err != nil {
			return fmt.Errorf("Rotate StartLogger: %s\n", err)
		}
	}
	return nil
}

func (w *FileLogWriter) processOldFile(path string) {
	oldExt := filepath.Ext(path)
	oldExt = oldExt[1:]
	oldLogNum, _ := strconv.Atoi(oldExt)
	if oldLogNum == 1 {
		os.Remove(path)
	} else {
		newLogNum := oldLogNum - 1

		newExt := fmt.Sprintf("%03d", newLogNum)
		newPath := strings.Replace(path, oldExt, newExt, -1)

		os.Rename(path, newPath)
	}
}

func getDefaultConfig() FileLogConfig {
	config := FileLogConfig{
		Filename: getDefaultLogPath(),
		Rotate:   true,
		Maxsize:  100 * 1024 * 1024,
		Maxfiles: 10,
	}
	return config
}

func getConfigFromViper() FileLogConfig {
	config := getDefaultConfig()

	if viper.IsSet("blockchain.log.filename") {
		filename := viper.GetString("blockchain.log.filename")

		//测试配置的日志路径是否可写
		_, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)

		if err == nil {
			config.Filename = filename
		} else {
			config.Filename = getDefaultLogPath()
		}
	}

	if viper.IsSet("blockchain.log.rotate") {
		config.Rotate = viper.GetBool("blockchain.log.rotate")
	}

	if viper.IsSet("blockchain.log.max-size") {
		config.Maxsize = viper.GetInt("blockchain.log.max-size")
	}

	if viper.IsSet("blockchain.log.max-files") {
		config.Maxfiles = viper.GetInt("blockchain.log.max-files")
	}

	return config
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getDefaultLogPath() string {
	dir := getCurrentDirectory()
	return filepath.Join(dir, "mpc.log")
}

//创建新的FileLogWriter对象
func NewFileWriter() FileLogWriter {
	config := getConfigFromViper()

	w := FileLogWriter{}
	// use MuxWriter instead direct use os.File for lock write when rotate
	w.mw = new(MuxWriter)
	err := w.Init(&config)
	if err != nil {
		log.Panicf("Error: Fail to initialize logging: %v\n", err)
	}

	return w
}
