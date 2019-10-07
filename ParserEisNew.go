package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type ParserEisNew struct {
	addDoc    int
	sendDoc   int
	purchases []Puchase44
}

func (t *ParserEisNew) run() {
	t.purchases = make([]Puchase44, 0)
	for _, url := range UrlsNew {
		t.parsingPage(url)
	}
	t.SendMessage()
}

func (t *ParserEisNew) SendMessage() {
	if len(t.purchases) > 0 {
		SendPurchaseInfo(EmailUs, t.purchases)
		t.sendDoc++
	}

}
func (t *ParserEisNew) parsingPage(p string) {
	defer SaveStack()
	r := DownloadPage(p)
	if r != "" {
		t.parsingTenderList(r, p)
	} else {
		Logging("Got empty string", p)
	}
}

func (t *ParserEisNew) parsingTenderList(p string, url string) {
	defer SaveStack()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(p))
	if err != nil {
		Logging(err)
		return
	}
	d := 0
	doc.Find("div.search-registry-entry-block > div.row").Each(func(i int, s *goquery.Selection) {
		d++
		t.parsingTenderFromList(s, url)

	})
	if d == 0 {
		Logging("We got no purchases", url)
	}
}

func (t *ParserEisNew) parsingTenderFromList(p *goquery.Selection, url string) {
	defer SaveStack()
	purNum := strings.TrimSpace(p.Find("div.registry-entry__header-top__number a").First().Text())
	purNum = strings.Replace(purNum, "№ ", "", -1)
	purName := strings.TrimSpace(p.Find("div:contains('Объект закупки') + div.registry-entry__body-value").First().Text())
	pubDate := strings.TrimSpace(p.Find("div.data-block > div:contains('Размещено') + div").First().Text())
	nmck := strings.TrimSpace(p.Find("div.price-block > div:contains('Начальная цена') + div").First().Text())
	//nmck = delallwhitespace(nmck)
	hrefT := p.Find("div.registry-entry__header-top__number a")
	href, exist := hrefT.Attr("href")
	if !exist {
		fmt.Println(p.Text())
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

func (t *ParserEisNew) checkPurchase(p Puchase44) {
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
