package main

import "sync"

type Filelog string

var DirLog = "log_portland"
var DirTemp = "temp_portland"
var SetFile = "settings.json"
var FileLog Filelog
var mutex sync.Mutex
var BotToken string
var ChannelId int64
var FileDB = "bd_purchase.sqlite"
var StartUrl = ""
