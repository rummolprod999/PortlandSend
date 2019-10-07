package main

import "fmt"

func init() {
	CreateEnv()
}

func main() {
	defer SaveStack()
	Logging("Start work")
	p := ParserEis{}
	p.run()
	c := ParserComita{}
	c.run()
	Logging(fmt.Sprintf("Add purchases %d", p.addDoc))
	Logging(fmt.Sprintf("Send purchases %d", p.sendDoc))
	Logging(fmt.Sprintf("Add purchases comita %d", c.addDoc))
	Logging(fmt.Sprintf("Send purchases comita %d", c.sendDoc))
	Logging("End work")
}
