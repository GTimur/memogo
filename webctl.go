// Реализует http-сервер с возможностью корректного завершения.
package memogo

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type WebCtl struct {
	host     net.IP
	port     uint16
	islisten bool
}

// Типизируем страницу для передачи данных в шаблон
type Page struct {
	Title   string
	Body    template.HTML
	LnkHome string
	DateNow template.HTML
}

var (
	NeedExit bool                // Флаг для завершения работы сервера
	Quit     = make(chan int, 1) // Канал для завершения сервера HTTP
)

var (
	HTMLDOC = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{ .Title }}</title>
</head>
<body>
	{{ .Body }}
</body>
</html>`
)

/*Сервер*/
//Запускает goroutine http.Server
func (w *WebCtl) StartServe() (err error) {
	//signal.Notify(Quit, os.Interrupt)
	srv := &http.Server{Addr: w.connString(),
		Handler:      http.DefaultServeMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// для отдачи сервером статичных файлов из папки /static
	fs := http.FileServer(http.Dir("./static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))

	cssFileServer := http.StripPrefix("/static/", fs)
	http.Handle("/static/", cssFileServer)
	http.HandleFunc("/", urlhome)       //Домашняя страница
	http.HandleFunc("/queue", urlqueue) //Task queue

	go func() {
		log.Println("Starting HTTP-server...")
		log.Fatalln("WebCtl error:", srv.ListenAndServe())
	}()

	go func() {
		<-Quit
		fmt.Println("Shutting down HTTP-server...")
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalln("HTTP Shutdown error:", err)
		}
	}()
	w.islisten = true
	return err
}

//SetHost - set Host
func (w *WebCtl) SetHost(host net.IP) {
	w.host = host
}

//SetPort - set Port
func (w *WebCtl) SetPort(port uint16) {
	w.port = port
}

/**/
func (w WebCtl) connString() string {
	return fmt.Sprintf("%s:%d", w.host.String(), w.port)
}

func (c *Config) SetManagerSrv(addr string, port uint16) {
	c.MgrSrv = ManagerSrv{
		Addr: addr,
		Port: port,
	}
}

func (c *Config) ManagerSrvAddr() string {
	return c.MgrSrv.Addr
}

func (c *Config) ManagerSrvPort() uint16 {
	return c.MgrSrv.Port
}

/****/
// Обработчик запросов для home - пример
func urlhome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	srvURI := "http://" + GlobalConfig.ManagerSrvAddr() + ":" + strconv.Itoa(int(GlobalConfig.ManagerSrvPort()))

	body := "<center><h1>" + BANNER + " " + VERSION + "</h1>" +
		"<p>" +
		"<ul>" +
		"<li><h4><a href=" + srvURI + "/queue" + ">Очередь заданий на рассылку уведомлений</a></h4></li>" +
		"</ul></center> </p>"

	main := HTMLDOC

	page := Page{Title: "MEMO GO reminder utility",
		Body:    template.HTML(body),
		LnkHome: "none",
		DateNow: "",
	}

	homeTemplate := template.Must(template.New("main").Parse(main))

	if r.Method == "GET" {
		if err := homeTemplate.ExecuteTemplate(w, "main", page); err != nil {
			fmt.Sprintln("Homepage handling error:", err.Error())
		}
		//fmt.Println("Homepage: GET request.")
	} else {
		//fmt.Println("Homepage: POST request.")
	}

}

// urlqueue - /queue page
func urlqueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	body := GlobalQueue.StringByID()
	body += "<br>"
	body += "InitTasks countdown (minutes):" + string(GlobalInitDelay/60) + "<br>"

	main := HTMLDOC

	page := Page{Title: "MEMO TASKS QUEUE",
		Body:    template.HTML(body),
		LnkHome: "none",
		DateNow: "",
	}

	homeTemplate := template.Must(template.New("main").Parse(main))

	if r.Method == "GET" {
		if err := homeTemplate.ExecuteTemplate(w, "main", page); err != nil {
			fmt.Sprintln("Queue page handling error:", err.Error())
		}
		//fmt.Println("Queue page: GET request.")
	} else {
		//fmt.Println("Queue page: POST request.")
	}
}

/****/
