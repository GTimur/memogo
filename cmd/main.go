package main

import (
	"memogo/pkg/memogo"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	var err error

	makeConfig := flag.Int("makeConfig", 0, "1 - config file template will be created")
	initDelay := flag.Int("initDelay", 0, "InitTasks loop delay in minutes, default=60")
	flag.Parse()

	if *makeConfig == 1 {
		memogo.GlobalConfig = memogo.Config{
			Root: "root\\",
			SMTPSrv: memogo.SrvSMTP{
				Addr:     "10.20.20.6",
				Port:     25,
				Account:  "noti",
				Password: "Password111",
				From:     "noti@ymkbank.ru",
				FromName: "Memo GO",
				UseTLS:   false,
			},
			MgrSrv: memogo.ManagerSrv{
				Addr: "127.0.0.1",
				Port: 8800,
			},
			LogFile: memogo.LogFile{
				Filename:   memogo.LOGFILE, // path to file
				CleanStale: false,          // if Stale flag is set - data will be rewrited after TTLDays
				TTLDays:    180,            // TTL in days if Stale is set
			},
		}
		err = memogo.GlobalConfig.MakeConfig()
		if err != nil {
			log.Fatalln("Creation error: ", memogo.CONFIGFILE, "\n Program terminated")
		}
		err = memogo.GlobalConfig.WriteJSON()
		if err != nil {
			log.Fatalln("Write error: ", memogo.CONFIGFILE, "\n Program terminated")
		}
		os.Exit(0)
	}

	if *initDelay <= 0 || *initDelay >= 168 {
		memogo.GlobalInitDelay = 60 * 60 //1 hour
	} else {
		memogo.GlobalInitDelay = *initDelay * 60 //to minutes
	}

	memogo.Banner()

	memogo.InitConfig()
	memogo.InitLog()
	memogo.InitEvents()

	// Инициализация web-сервера
	memogo.NeedExit = false // флаг для завершения работы
	var web memogo.WebCtl
	web.SetHost(net.ParseIP(memogo.GlobalConfig.ManagerSrvAddr()))
	web.SetPort(memogo.GlobalConfig.ManagerSrvPort())

	fmt.Println("\nWeb control configured: " + "http://" + memogo.GlobalConfig.ManagerSrvAddr() + ":" + strconv.Itoa(int(memogo.GlobalConfig.ManagerSrvPort())))

	/* Запускаем сервер обслуживания WebCtl */
	err = web.StartServe()
	if err != nil {
		log.Println("HTTP сервер: Ошибка. ", err)
		os.Exit(1)
	}

	/* Перехват CTRL+C для завершения приложения */
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("\nReceived %v, shutdown procedure initiated.\n\n", sig)
			memogo.Quit <- 1
			memogo.NeedExit = true
		}
	}()

	// Цикл с таймером для ожидания команды завершения
	ticker := time.NewTicker(time.Second * 1) // Запускаем обработчик каждую секунду

	// Зациклимся с таймером посекундно пока не получим команду завершения работы.
	memogo.DelayCounter = 0
	for range ticker.C {
		err := memogo.MemoSvc(memogo.GlobalQueue)
		if err != nil {
			log.Println("MemoSvc error:", err)
		}

		// Опрос папки /root раз в час + переинициализация списков
		if memogo.DelayCounter >= memogo.GlobalInitDelay {
			// read tasks from disk and rebuild GlobalTask
			// read GlobalTask array and rebuild GlobalTimeMap
			// read GlobalTimeMap and build GlobalQueue
			memogo.InitEvents()
			memogo.GlobalConfig.LogFile.Add("Tasks reloaded. Time maps rebuilded. Delay=" + strconv.Itoa(memogo.GlobalInitDelay/60) + "minutes.")
			memogo.DelayCounter = 0
		}
		memogo.DelayCounter++

		if !memogo.NeedExit {
			continue
		}
		break
	}

	ticker.Stop()
}
