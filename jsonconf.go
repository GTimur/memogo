package memogo

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// CONFIGFILE - constant path to json config file
const CONFIGFILE = "./config.json"

// Config define initialization parameters
type Config struct {
	Root    string `json:"TASKS_ROOT_DIR_PATH"`
	SMTPSrv SrvSMTP
	MgrSrv  ManagerSrv
}

// SrvSMTP setup smtp credentials
type SrvSMTP struct {
	Addr     string
	Port     uint
	Account  string //smtp account from whom sending email
	Password string
	From     string //sender address (FROM)
	FromName string //Name of sender (example "Memo GO sender")
	UseTLS   bool   //auth: use TLS or plain/text
}

//ManagerSrv  - Web-server address:port
type ManagerSrv struct {
	Addr string
	Port uint16
}

//Read from Json-config file
func (c *Config) readJSON() (err error) {
	file, err := os.Open(CONFIGFILE)
	defer file.Close()
	if err != nil {
		log.Printf("Config: error reading JSON file: %v\n", err)
		return err
	}

	//Готовим для импорта структуру JSON
	var jsonConfig Config

	//Читаем файл JSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonConfig)
	if err != nil {
		log.Printf("JSON decoder error: %v", err)
		return err
	}

	var aeserr error
	encStr := strings.Split(jsonConfig.SMTPSrv.Password, "*")

	var encBytes []byte
	b := 0
	for _, elm := range encStr {
		if len(elm) == 0 {
			continue
		}
		b, err = strconv.Atoi(elm)
		if err != nil {
			break
		}
		encBytes = append(encBytes, byte(b))
	}

	jsonConfig.SMTPSrv.Password, aeserr = AesDecript(encBytes)
	if aeserr != nil {
		log.Println("readJSON decript error:", aeserr)
		return aeserr
	}

	return err
}

// ReadJSON - reads JSON file into struct
func (c *Config) ReadJSON() (err error) {
	err = c.readJSON()
	return
}

//Write config into JSON-file
//File must be exist and accessible
func (c *Config) writeJSON() (err error) {
	file, err := os.OpenFile(CONFIGFILE, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Printf("Config writer: error writing JSON file: %v\n", err)
		return err
	}

	//Готовим данные JSON (конвертируем в экспортируемый вид)
	jsonConfig := c

	var aeserr error
	// Зашифруем строку пароля
	var encBytes []byte
	encBytes, aeserr = AesEncrypt(jsonConfig.SMTPSrv.Password)
	if aeserr != nil {
		return aeserr
	}

	str := ""
	for i := range encBytes {
		str += fmt.Sprintf("%d*", encBytes[i])
	}

	jsonConfig.SMTPSrv.Password = str

	// пишем в файл
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(&jsonConfig)
	if err != nil {
		log.Printf("JSON encoder error: %v", err)
	}
	return err
}

// WriteJSON - writes struct to JSON file
// filepath is set in CONFIGFILE constant
func (c *Config) WriteJSON() (err error) {
	err = c.writeJSON()
	return
}

// MakeConfig - creates config file if it not exists
func (c *Config) MakeConfig() (err error) {
	if _, err = os.Stat(CONFIGFILE); err == nil {
		//File exist and will not be rewrited
		//fmt.Println(os.IsExist(err),err)
		return err
	}
	file, err := os.Create(CONFIGFILE)
	defer file.Close()
	return err
}
