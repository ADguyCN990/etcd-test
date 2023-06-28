package pkg

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// NewShutdownSignal 新建一个监听中断信号的channel
func NewShutdownSignal() chan os.Signal {
	c := make(chan os.Signal, 1)
	// SIGHUP: terminal closed
	// SIGINT: Ctrl+C
	// SIGTERM: program exit
	// SIGQUIT: Ctrl+/
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return c
}

// WaitExit 监听中断信号并优雅退出程序
func WaitExit(wg *sync.WaitGroup, exit func()) {
	c := NewShutdownSignal()
	defer wg.Done()
	//等待中断信号
	interruptSignal := <-c
	logrus.Info("收到中断信号 ", interruptSignal.String(), ",准备退出程序...")
	exit()
	os.Exit(0)
}
