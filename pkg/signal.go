package pkg

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// NewShutdownSignal 新建一个监听中断信号的channel
func NewShutdownSignal() chan os.Signal {
	c := make(chan os.Signal)
	// SIGHUP: terminal closed
	// SIGINT: Ctrl+C
	// SIGTERM: program exit
	// SIGQUIT: Ctrl+/
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return c
}

// WaitExit 监听中断信号并优雅退出程序
func WaitExit(c chan os.Signal, exit func()) {
	for i := range c {
		switch i {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logrus.Info("收到中断信号 ", i.String(), ",准备退出程序...")
			exit()
			os.Exit(0)
		}
	}
}
