package start

import (
	"etcd-test/interval/model/entity"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

var config Config

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
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		logrus.WithError(err).Error("无法解析配置文件")
		return err
	}
	return nil
}

func (s *Init) Database() (*gorm.DB, error) {
	//从配置文件中获取数据库连接信息
	driverName := config.DataSource.DriverName
	host := config.DataSource.Host
	port := config.DataSource.Port
	database := config.DataSource.Database
	username := config.DataSource.Username
	password := config.DataSource.Password
	charset := config.DataSource.Charset

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

func (s *Init) Iris() (*iris.Application, error) {
	app := iris.New()
	return app, nil
}

func (s *Init) IrisListen(app *iris.Application) error {
	serverPort := config.Server.Port
	err := app.Run(iris.Addr(":" + serverPort))
	if err != nil {
		logrus.WithError(err).Error("Iris无法监听端口：", serverPort)
		return err
	}
	return nil
}
