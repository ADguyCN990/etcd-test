package start

import (
	"etcd-test/interval/dao"
	"etcd-test/interval/model/entity"
	"etcd-test/interval/webServer"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

var ConfigData Config

type Init struct {
}

type Config struct {
	Server     Server     `json:"server"`
	DataSource DataSource `json:"dataSource"`
}

type Server struct {
	Port string `yaml:"port"`
}

type DataSource struct {
	DriverName string `yaml:"driverName"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Database   string `yaml:"database"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Charset    string `yaml:"charset"`
}

// InitConfig 读取配置文件
func (s *Init) InitConfig() error {
	// 获取当前工作目录的绝对路径
	wd, err := os.Getwd()
	if err != nil {
		logrus.WithError(err).Error("无法打开配置文件")
		return err
	}

	// 拼接配置文件路径
	configPath := filepath.Join(wd, "config", "application.yaml")

	// 读取配置文件
	configFile, err := os.Open(configPath)
	if err != nil {
		logrus.WithError(err).Error("无法读取配置文件")
		return err
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			logrus.WithError(err).Error("无法断开与配置文件的连接")
		}
	}(configFile)

	// 解析配置文件
	err = yaml.NewDecoder(configFile).Decode(&ConfigData)
	if err != nil {
		logrus.WithError(err).Error("无法解析配置文件")
		return err
	}
	return nil
}

// Database 连接数据库并启动gorm
func (s *Init) Database() (*gorm.DB, error) {
	//从配置文件中获取数据库连接信息
	driverName := ConfigData.DataSource.DriverName
	host := ConfigData.DataSource.Host
	port := ConfigData.DataSource.Port
	database := ConfigData.DataSource.Database
	username := ConfigData.DataSource.Username
	password := ConfigData.DataSource.Password
	charset := ConfigData.DataSource.Charset

	//构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username, password, host, port, database, charset)
	// 连接数据库
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: driverName,
		DSN:        dsn,
	}), &gorm.Config{})
	if err != nil {
		logrus.WithError(err).Error("无法连接数据库")
		return nil, err
	}

	err = db.AutoMigrate(&entity.TbTask{})
	if err != nil {
		logrus.WithError(err).Error("无法启用AutoMigrate")
		return nil, err
	}

	return db, nil
}

func (s *Init) Exit() {
	// 关闭与数据库的连接
	var da = &dao.Dao{}
	err := da.CloseDB()
	if err != nil {
		logrus.WithError(err).Error("无法断开与数据库的连接")
		return
	}
	logrus.Info("成功断开与数据库的连接")

	// 关闭iris server
	var se = &webServer.Server{}
	err = se.ShutdownServer()
	if err != nil {
		logrus.WithError(err).Error("无法关闭iris")
		return
	}
	logrus.Info("成功关闭iris")

}
