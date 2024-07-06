package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	City         string `json:"city"`
	State        string `json:"state"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
}

type ViaCepAPIResponse struct {
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
	Cep        string `json:"cep"`
}

func main() {
	cep := os.Args[1]

	viacepCh := make(chan string)
	brasilApiCh := make(chan string)

	go getAddressFromBrasilAPI(cep, brasilApiCh)
	go getAddressFromViaCepAPI(cep, viacepCh)

	select {
	case address := <-viacepCh:
		fmt.Printf("Received from ViaCep: \n%s", address)
	case address := <-brasilApiCh:
		fmt.Printf("Received from Brasil API: \n%s", address)
	case <-time.After(1 * time.Second):
		log.Fatal("Timeout")
	}
}

func getAddressFromBrasilAPI(cep string, ch chan<- string) {
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var jsonResponse BrasilAPIResponse
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return
	}

	result := fmt.Sprintf(
		"\tCEP: %s\n\tCity: %s\n\tState: %s\n\tStreet: %s\n\tNeighborhood: %s\n",
		jsonResponse.Cep,
		jsonResponse.City,
		jsonResponse.State,
		jsonResponse.Street,
		jsonResponse.Neighborhood,
	)

	ch <- string(result)
}

func getAddressFromViaCepAPI(cep string, ch chan<- string) {
	resp, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var jsonResponse ViaCepAPIResponse
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return
	}

	result := fmt.Sprintf(
		"\tCEP: %s\n\tCity: %s\n\tState: %s\n\tStreet: %s\n\tNeighborhood: %s\n",
		jsonResponse.Cep,
		jsonResponse.Localidade,
		jsonResponse.Uf,
		jsonResponse.Logradouro,
		jsonResponse.Bairro,
	)

	ch <- string(result)
}
