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

// NewCloseGoRoutineChannel 新建一个负责关闭协程的channel
func NewCloseGoRoutineChannel() chan bool {
	c := make(chan bool, 1)
	return c
}

// SendSignalToCloseChannel 向负责关闭协程的channel发送信号
func SendSignalToCloseChannel(c *chan bool) {
	logrus.Info("进入关闭channel的函数里面了")
	close(*c)
}

// CloseAllRoutines 关闭所有协程
func CloseAllRoutines(dbRoutine *chan bool,
	irisPortRoutine *chan bool) {
	logrus.Info("进入关闭协程函数里面了")
	SendSignalToCloseChannel(dbRoutine)
	SendSignalToCloseChannel(irisPortRoutine)
}

// WaitExit 监听中断信号并优雅退出程序
func WaitExit(wg *sync.WaitGroup, exit func(), dbRoutine *chan bool,
	irisPortRoutine *chan bool) {
	defer wg.Done()
	c := NewShutdownSignal()
	// 等待中断信号
	interruptSignal := <-c
	logrus.Info("收到中断信号: ", interruptSignal.String(), ",准备退出程序...")
	CloseAllRoutines(dbRoutine, irisPortRoutine)
	exit()
	os.Exit(0)
}
