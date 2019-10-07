package main

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

}
