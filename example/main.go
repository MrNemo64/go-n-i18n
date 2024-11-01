package main

import (
	"fmt"

	"github.com/MrNemo64/go-n-i18n/example/lang"
)

func main() {
	bundle := lang.MessagesForMust("en-EN")

	fmt.Println(bundle.WhereAmI())                       // Assume this json is in the file "en-EN.json"
	fmt.Println(bundle.NestedMessages().Parametrized(4)) // This message has an amout parameter of type int: 4
	fmt.Println(bundle.ConditionalMessages(100))
	/*
		This is the "else" branch
		This multiline message is used
		And shows the amount: 100
	*/
	fmt.Println(bundle.MultilineMessage("MrNemo64", 13.1267))
	/*
		Hello MrNemo64!
		Messages can be multiline
		And each one can have parameters
		This one has a float formated with 2 decimals! 13.13
	*/
}
