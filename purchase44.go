package main

import "fmt"

type Puchase44 struct {
	PurNum  string
	PurName string
	PubDate string
	Nmck    string
	Href    string
}

func (p *Puchase44) CheckPurchase() bool {
	if p.PurNum == "" || p.PurName == "" || p.Href == "" {
		Logging(fmt.Sprintf("The purchase is bad %+v", p))
		return false
	}
	return true
}
