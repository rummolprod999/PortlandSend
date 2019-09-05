package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type ParserEis struct {
	addDoc    int
	sendDoc   int
	purchases []Puchase44
}

func (t *ParserEis) run() {
	t.purchases = make([]Puchase44, 0)
	for _, url := range Urls {
		t.parsingPage(url)
	}
	t.SendMessage()
}
func (t *ParserEis) SendMessage() {
	if len(t.purchases) > 0 {
		SendPurchaseInfo(EmailUs, t.purchases)
	}

}

func (t *ParserEis) parsingPage(p string) {
	defer SaveStack()
	r := DownloadPage(p)
	if r != "" {
		t.parsingTenderList(r, p)
	} else {
		Logging("Got empty string", p)
	}
}

func (t *ParserEis) parsingTenderList(p string, url string) {
	defer SaveStack()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(p))
	if err != nil {
		Logging(err)
		return
	}
	doc.Find("div.parametrs div.registerBox.registerBoxBank.margBtm20").Each(func(i int, s *goquery.Selection) {
		t.parsingTenderFromList(s, url)

	})

}

func (t *ParserEis) parsingTenderFromList(p *goquery.Selection, url string) {
	defer SaveStack()
	purNum := strings.TrimSpace(p.Find("td.descriptTenderTd dl dt a").First().Text())
	purNum = strings.Replace(purNum, "№ ", "", -1)
	purName := strings.TrimSpace(p.Find("td.descriptTenderTd dl dd:nth-of-type(2)").First().Text())
	pubDate := strings.TrimSpace(p.Find("td.amountTenderTd ul li:nth-of-type(1)").First().Text())
	pubDate = strings.TrimSpace(strings.Replace(pubDate, "Размещено:", "", -1))
	nmck := strings.TrimSpace(p.Find("span:contains('Начальная цена') ~ strong").First().Text())
	nmck = delallwhitespace(nmck)
	hrefT := p.Find("td.descriptTenderTd dl dt a")
	href, exist := hrefT.Attr("href")
	if !exist {
		Logging("The element have no Href attribute", hrefT.Text())
		return
	}
	if !strings.Contains(href, "http://zakupki.gov.ru") {
		href = fmt.Sprintf("http://zakupki.gov.ru%s", href)
	}
	purch := Puchase44{Href: href, PubDate: pubDate, PurName: purName, PurNum: purNum, Nmck: nmck}
	if purch.CheckPurchase() {
		t.checkPurchase(purch)
	}

}

func (t *ParserEis) checkPurchase(p Puchase44) {
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
	_, err = db.Exec("INSERT INTO purchase (id, purchase_num) VALUES (NULL, $1)", p.PurNum)
	if err != nil {
		Logging(err)
		return
	}
	t.addDoc++
	t.purchases = append(t.purchases, p)
}
