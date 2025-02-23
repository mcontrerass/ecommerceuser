package main

import (
	"context"
	"ecommerceuser/awsgo"
	"ecommerceuser/bd"
	"ecommerceuser/models"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(EjecutarLambda)
}

func EjecutarLambda(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	awsgo.InicializarAWS()

	if !ValidarParametros() {
		fmt.Println("Error en los parámetros, debe enviar 'SecretName'")
		err := errors.New("Error en los parámetros, debe enviar SecretName")
		return event, err
	}

	var datos models.SignUp

	for row, att := range event.Request.UserAttributes {
		switch row {
		case "email":
			datos.UserEmail = att
			fmt.Println("Email: " + datos.UserEmail)

		case "sub":
			datos.UserUUID = att
			fmt.Println("Sub: " + datos.UserUUID)
		}
	}

	err := bd.ReadSecret()

	if err != nil {
		fmt.Println("Error al leer el Secret " + err.Error())
		return event, err
	}

	err = bd.SignUp(datos)
	return event, err
}

func ValidarParametros() bool {
	var traeParametro bool
	_, traeParametro = os.LookupEnv("SecretName")
	return traeParametro
}
