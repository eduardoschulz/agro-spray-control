package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "http://localhost:30000"

func main() {
	fmt.Println("CLI - Operações com Área e Pulverização")

	for {
		fmt.Println("\nEscolha uma opção:")
		fmt.Println("1. Cadastrar área")
		fmt.Println("2. Listar áreas")
		fmt.Println("3. Cadastrar pulverização")
		fmt.Println("4. Listar pulverizações")
		fmt.Println("5. Associar área a pulverização")
		fmt.Println("6. Sair")

		var escolha int
		fmt.Print("Opção: ")
		fmt.Scanln(&escolha)

		switch escolha {
		case 1:
			cadastrar("areas", map[string]interface{}{
				"tamanho":     prompt("Tamanho (ex: 10.5)"),
				"cod_fazenda": prompt("Código da fazenda"),
			})
		case 2:
			listar("areas")
		case 3:
			cadastrar("pulverizacoes", map[string]interface{}{
				"dt_aplicacao":    prompt("Data de aplicação (ex: 2025-06-10T00:00:00Z)"),
				"cultura":         prompt("Cultura"),
				"cod_lote":        prompt("Código do lote"),
				"cpf_responsavel": prompt("CPF do responsável"),
			})
		case 4:
			listar("pulverizacoes")
		case 5:
			cadastrar("associar", map[string]interface{}{
				"cod_pulv": prompt("Código da pulverização"),
				"cod_area": prompt("Código da área"),
			})
		case 6:
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("Opção inválida")
		}
	}
}

func cadastrar(endpoint string, dados map[string]interface{}) {
	body, err := json.Marshal(dados)
	if err != nil {
		fmt.Println("Erro ao gerar JSON:", err)
		return
	}

	resp, err := http.Post(baseURL+"/"+endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Erro na requisição:", err)
		return
	}
	defer resp.Body.Close()
	res, _ := io.ReadAll(resp.Body)
	fmt.Println("Resposta:", string(res))
}

func listar(endpoint string) {
	resp, err := http.Get(baseURL + "/" + endpoint)
	if err != nil {
		fmt.Println("Erro ao listar:", err)
		return
	}
	defer resp.Body.Close()
	res, _ := io.ReadAll(resp.Body)
	fmt.Println(string(res))
}

func prompt(label string) string {
	fmt.Print(label + ": ")
	var input string
	fmt.Scanln(&input)
	return input
}
