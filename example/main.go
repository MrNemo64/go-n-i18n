package main

import (
	"fmt"

	"github.com/MrNemo64/go-n-i18n/example/lang"
)

func main() {
	fmt.Println(lang.MessagesForMust("en-EN").In().Depeer())
	fmt.Println(lang.MessagesForMust("es-ES").In().Depeer())
}
