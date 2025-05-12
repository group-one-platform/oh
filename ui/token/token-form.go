package token

import (
	"github.com/charmbracelet/huh"
	"log"
)

func RequestToken() (token string, err error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("API Token").
				EchoMode(huh.EchoModePassword).
				Value(&token),
		),
	)

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	return
}
