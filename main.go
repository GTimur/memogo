package main

func main() {
	memogo.GlobalConfig = memogo.Config{
		Root: "./root/",
		SMTPSrv: memogo.SrvSMTP{
			Addr:     "10.20.20.6",
			Port:     25,
			Account:  "noti",
			Password: "Bank999",
			From:     "noti@ymkbank.ru",
			FromName: "Memo GO",
			UseTLS:   false,
		},
		MgrSrv: memogo.ManagerSrv{
			Addr: "127.0.0.1",
			Port: 8000,
		},
	}
	

	//var files map[string]string

	//files, err := memogo.FindAllFiles(globalconfig.Root, []string{"*.*"})
	//if err != nil {
	//	log.Fatalf("Main(): FindFiles error: %v", err)
	//}

	//groupname := path.Dir(strings.Replace("./root/test2/rule01.json", globalconfig.Root, "", -1))
	//fmt.Println("GROUPS FOUND:", groupname)
	//fmt.Println("FILES FOUND:", files)

	/*err = GlobalConfig.MakeConfig()
	if err != nil {
		panic(err)
	}

	err = GlobalConfig.WriteJSON()
	if err != nil {
		panic(err)
	}

	err = memogo.TestJSON()
	if err != nil {
		log.Fatalf("TestJSON() error: %v", err)
	}*/

}
