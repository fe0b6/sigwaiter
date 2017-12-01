// Завершатель работы корректно
package sigwaiter

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	exitChan chan bool
)

func Run(waitTime int, chans ...chan bool) {
	exitChan = make(chan bool)

	log.Println("[info]", "Перехват сигналов инициализирован")

	// Перехват сигналов
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, os.Interrupt)

	waitExit(c)

	go func() {
		time.Sleep(time.Duration(waitTime) * time.Second)
		log.Println("[error]", "Неудалось завершить работу корректно")

		os.Exit(2)
	}()

	for _, ch := range chans {
		ch <- true
		_ = <-ch
	}

	log.Println("[info]", "Работа завершена корректно")

	os.Exit(0)
}

func waitExit(c chan os.Signal) {
	for {
		select {
		case s := <-c:
			log.Println("[info]", "Получен сигнал: ", s)
			return
		case <-exitChan:
			log.Println("[info]", "Самоинициализированный выход")
			return
		}
	}
}

func Exit() {
	exitChan <- true
}
