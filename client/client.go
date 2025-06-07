package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Valor string `json:"bid"`
}

func main() {
	cotacao, err := obterCotacao()

	if err != nil {
		panic(err)
	}

	err = criaArquivo(cotacao)

	if err != nil {
		panic(err)
	}

	println("Arquivo criado com sucesso com o valor da cotação: " + cotacao.Valor)
}

func criaArquivo(cotacao *Cotacao) error {
	f, err := os.Create("cotacao.txt")

	if err != nil {
		return err
	}

	defer f.Close()

	if err != nil {
		return err
	}

	_, err = f.WriteString("Dólar: " + cotacao.Valor + "\n")

	if err != nil {
		return err
	}

	return nil
}

func obterCotacao() (*Cotacao, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		if err == context.DeadlineExceeded {
			log.Println("Request timed out")
		}
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro ao obter cotação: %s", resp.Status)
		return nil, errors.New("erro ao obter cotação: " + resp.Status)
	}

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}
	var cotacao Cotacao

	err = json.NewDecoder(resp.Body).Decode(&cotacao)

	return &cotacao, err
}
