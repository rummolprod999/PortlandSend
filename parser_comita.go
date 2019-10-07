package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"time"
)

type ParserComita struct {
	addDoc    int
	sendDoc   int
	purchases []Puchase44
}

func (t *ParserComita) run() {
	t.purchases = make([]Puchase44, 0)
	for _, kw := range KeyWords {
		t.parsingComita(kw)
	}
	t.SendMessage()
}
func (t *ParserComita) SendMessage() {
	if len(t.purchases) > 0 {
		SendPurchaseInfo(EmailUs, t.purchases)
		t.sendDoc++
	}
}

func (t *ParserComita) parsingComita(kw string) {
	defer SaveStack()
	r := DownloadPageComita("https://etp.comita.ru/rest/site/procedures/223", kw)
	if r != "" {
		t.parsingTenderList(r)
	} else {
		Logging("Got empty string")
	}
}

func (t *ParserComita) parsingTenderList(p string) {
	defer SaveStack()
	_, err := jsonparser.ArrayEach([]byte(p), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			Logging(err, "in callback parsingTenderList")
			return
		}
		idByte, _, _, err := jsonparser.Get(value, "id")
		if err != nil {
			Logging(err, string(value), "id")
			return
		}
		id := string(idByte)
		if id == "" {
			Logging("id not found")
			return
		}
		typeByte, _, _, err := jsonparser.Get(value, "type", "id")
		if err != nil {
			Logging(err, string(value), "typeByte")
			return
		}
		typeAuct := string(typeByte)
		href := fmt.Sprintf("https://etp.comita.ru/openProcedure/%s/%s", typeAuct, id)
		purNumByte, _, _, err := jsonparser.Get(value, "number")
		if err != nil {
			Logging(err, href, "purNumByte")
			return
		}
		purNum := string(purNumByte)
		pubDateByteUnix, err := jsonparser.GetInt(value, "publishDate")
		if err != nil {
			pubDateByteUnix, err = jsonparser.GetInt(value, "created")
			if err != nil {
				Logging(err, href, "pubDateByteUnix")
				return
			}
		}
		pubDateUnix := time.Unix(0, pubDateByteUnix*int64(time.Millisecond))
		pubDate := pubDateUnix.String()
		nmckFloat, err := jsonparser.GetFloat(value, "contractSumNoNDS")
		if err != nil {
			Logging(err, href, "nmckFloat")
			nmckFloat = 0
		}
		if nmckFloat < 4*1000000 {
			return
		}
		nmck := fmt.Sprintf("%.2f RUB", nmckFloat)

		nameByte, _, _, err := jsonparser.Get(value, "name")
		if err != nil {
			Logging(err, string(value), "nameByte")
			return
		}
		purName := string(nameByte)

		purch := Puchase44{Href: href, PubDate: pubDate, PurName: purName, PurNum: purNum, Nmck: nmck}
		if purch.CheckPurchase() {
			t.checkPurchase(purch)
		}
	}, "items")
	if err != nil {
		Logging(err, "in parsingTenderList")
		return
	}

}

func (t *ParserComita) checkPurchase(p Puchase44) {
	db, err := DbConnection()
	if err != nil {
		Logging(err)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT id FROM purchase WHERE purchase_num=$1", p.PurNum)
	if err != nil {
		Logging(err)
		return
	}
	if rows.Next() {
		rows.Close()
		return
	}
	rows.Close()
	t.addDoc++
	t.purchases = append(t.purchases, p)
}
