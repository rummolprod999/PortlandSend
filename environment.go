package main

import (
	"database/sql"
	"fmt"
	"github.com/buger/jsonparser"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Filelog string

var DirLog = "log_portland"
var DirTemp = "temp_portland"
var SetFile = "settings.json"
var FileLog Filelog
var mutex sync.Mutex
var EmailUs string
var FileDB = "bd_purchase.sqlite"
var StartUrl = ""
var Urls = []string{"http://zakupki.gov.ru/epz/order/quicksearch/search.html?searchString=цемент&morphology=on&pageNumber=1&sortDirection=false&recordsPerPage=_50&showLotsInfoHidden=false&fz44=on&fz223=on&af=on&priceFrom=4000000&currencyId=-1&regionDeleted=false&sortBy=UPDATE_DATE", "http://zakupki.gov.ru/epz/order/quicksearch/search.html?searchString=Портландцемент&morphology=on&pageNumber=1&sortDirection=false&recordsPerPage=_50&showLotsInfoHidden=false&fz44=on&fz223=on&af=on&priceFrom=4000000&currencyId=-1&regionDeleted=false&sortBy=UPDATE_DATE", "http://zakupki.gov.ru/epz/order/quicksearch/search.html?searchString=Портланд-цемент&morphology=on&pageNumber=1&sortDirection=false&recordsPerPage=_50&showLotsInfoHidden=false&fz44=on&fz223=on&af=on&priceFrom=4000000&currencyId=-1&regionDeleted=false&sortBy=UPDATE_DATE", "http://zakupki.gov.ru/epz/order/quicksearch/search.html?searchString=Портланд+цемент&morphology=on&pageNumber=1&sortDirection=false&recordsPerPage=_50&showLotsInfoHidden=false&fz44=on&fz223=on&af=on&priceFrom=4000000&currencyId=-1&regionDeleted=false&sortBy=UPDATE_DATE"}

var UrlsNew = []string{}

var KeyWords = []string{"Цемент", "Портландцемент", "Портланд-цемент", "Портланд цемент", "цемент"}

func DbConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal_mode=OFF&_synchronous=OFF", FileDB))
	return db, err
}

func CreateLogFile() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirlog := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, DirLog))
	if _, err := os.Stat(dirlog); os.IsNotExist(err) {
		err := os.MkdirAll(dirlog, 0711)

		if err != nil {
			fmt.Println("Не могу создать папку для лога")
			os.Exit(1)
		}
	}
	t := time.Now()
	ft := t.Format("2006-01-02")
	FileLog = Filelog(filepath.FromSlash(fmt.Sprintf("%s/log_portland_%v.log", dirlog, ft)))
}

func CreateTempDir() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirtemp := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, DirTemp))
	if _, err := os.Stat(dirtemp); os.IsNotExist(err) {
		err := os.MkdirAll(dirtemp, 0711)

		if err != nil {
			fmt.Println("Не могу создать папку для временных файлов")
			os.Exit(1)
		}
	} else {
		err = os.RemoveAll(dirtemp)
		if err != nil {
			fmt.Println("Не могу удалить папку для временных файлов")
		}
		err := os.MkdirAll(dirtemp, 0711)
		if err != nil {
			fmt.Println("Не могу создать папку для временных файлов")
			os.Exit(1)
		}
	}
}

func Logging(args ...interface{}) {
	mutex.Lock()
	file, err := os.OpenFile(string(FileLog), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println("Ошибка записи в файл лога", err)
		return
	}
	fmt.Fprintf(file, "%v  ", time.Now())
	for _, v := range args {

		fmt.Fprintf(file, " %v", v)
	}
	//fmt.Fprintf(file, " %s", UrlXml)
	fmt.Fprintln(file, "")
	mutex.Unlock()
}

func ReadSetting() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	filetemp := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, SetFile))
	file, err := os.Open(filetemp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	EmailUs, err = jsonparser.GetString(b, "email")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if EmailUs == "" {
		fmt.Println("Check file with settings")
		os.Exit(1)
	}
	smtpHost, err = jsonparser.GetString(b, "host")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if smtpHost == "" {
		fmt.Println("Check file with settings")
		os.Exit(1)
	}

	smtpPort, err = jsonparser.GetInt(b, "port")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if smtpPort == 0 {
		fmt.Println("Check file with settings")
		os.Exit(1)
	}

	smtpUser, err = jsonparser.GetString(b, "user")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if smtpUser == "" {
		fmt.Println("Check file with settings")
		os.Exit(1)
	}

	smtpPass, err = jsonparser.GetString(b, "pass")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if smtpPass == "" {
		fmt.Println("Check file with settings")
		os.Exit(1)
	}
}
func CreateNewDB() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDB := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, FileDB))
	if _, err := os.Stat(fileDB); os.IsNotExist(err) {
		Logging(err)
		f, err := os.Create(fileDB)
		if err != nil {
			Logging(err)
			panic(err)
		}
		err = f.Chmod(0777)
		if err != nil {
			Logging(err)
			//panic(err)
		}
		err = f.Close()
		if err != nil {
			Logging(err)
			panic(err)
		}
		db, err := DbConnection()
		if err != nil {
			Logging(err)
			panic(err)
		}
		defer db.Close()
		_, err = db.Exec(`CREATE TABLE "purchase" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"purchase_num"	TEXT
)`)
		if err != nil {
			Logging(err)
			panic(err)
		}
		_, err = db.Exec(`CREATE INDEX "pur_index" ON "purchase" (
	"purchase_num"
)`)
		if err != nil {
			Logging(err)
			panic(err)
		}
	}
}
func CreateEnv() {
	ReadSetting()
	CreateLogFile()
	CreateTempDir()
	CreateNewDB()
}
