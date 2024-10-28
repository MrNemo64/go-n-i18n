package main

import (
	"fmt"

	"github.com/MrNemo64/go-n-i18n/example/lang"
)

func main() {
	fmt.Println(lang.MessagesForMust("en-EN").Cmds().Multiline(5, "juan"))
	fmt.Println()
	fmt.Println(lang.MessagesForMust("es-ES").Cmds().Multiline(5, "juan"))
}
