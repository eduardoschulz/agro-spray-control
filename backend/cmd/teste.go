package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
)

const baseURL = "http://localhost:30000"

func main() {
	for {
		prompt := promptui.Select{
			Label: "Escolha uma ação",
			Items: []string{
				"Adicionar Usuário",
				"Listar Usuários",
				"Adicionar Fazenda",
				"Listar Fazendas",
				"Adicionar Área",
				"Listar Áreas",
				"Sair",
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println("Erro ao selecionar:", err)
			os.Exit(1)
		}

		switch result {
		case "Adicionar Usuário":
			addUsuario()
		case "Listar Usuários":
			listar("usuarios")
		case "Adicionar Fazenda":
			addFazenda()
		case "Listar Fazendas":
			listar("fazendas")
		case "Adicionar Área":
			addArea()
		case "Listar Áreas":
			listar("areas")
		case "Sair":
			return
		}
	}
}

func input(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Println("Erro:", err)
		return ""
	}
	return result
}

func addUsuario() {
	cpf := input("CPF")
	nome := input("Nome")
	senha := input("Senha")

	payload := map[string]string{
		"cpf":    cpf,
		"nome":   nome,
		"senha":  senha,
        "nivelPermissao": "1",
	}
	post("usuarios", payload)
}

func addFazenda() {
	nome := input("Nome da Fazenda")
	localizacao := input("Localização")

	payload := map[string]string{
		"nome":        nome,
		"localizacao": localizacao,
	}
	post("fazendas", payload)
}

func addArea() {
	nome := input("Nome da Área")
	tamanho := input("Tamanho")
	codFazenda := input("Código da Fazenda")

	payload := map[string]string{
		"nome":       nome,
		"tamanho":    tamanho,
		"codFazenda": codFazenda,
	}
	post("areas", payload)
}

func post(endpoint string, data map[string]string) {
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(fmt.Sprintf("%s/%s", baseURL, endpoint), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erro ao enviar requisição:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Resposta:", string(body))
}

func listar(endpoint string) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseURL, endpoint))
	if err != nil {
		fmt.Println("Erro ao buscar dados:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Resultado:", string(body))
}

