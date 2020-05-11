package utils

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
	"log"
)

func HandleDbError(e validator.FieldError) string {
	var handledError = ""
	fmt.Println(e)
	switch e.Field() {
	case "Author":
		handledError = "Autor do livro precisa ser informado!"
		break
	case "Title":
		handledError = "Titulo do livro precisa ser informado!"
		break
	case "Publisher":
		handledError = "Editora do livro precisa ser informado!"
		break
	case "Language":
		handledError = "Idioma do livro precisa ser informado!"
		break
	case "ISBN":
		handledError = "O Codigo ISBN do livro precisa ser informado!"
		break
	default:
		handledError = "Houve um error nao previsto durante operacao no Banco, favor tente mais tarde ou entre em contato"
	}

	return handledError
}

func HandleWriteError(err error) string {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return "O Codigo ISBN informado ja consta no banco de dados!"
			}
		}
	}
	log.Fatal(err)
	return "Unexpected error"
}
