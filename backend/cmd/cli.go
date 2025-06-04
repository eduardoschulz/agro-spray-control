package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	baseURL  = "http://localhost:30000" // Deve corresponder à porta do seu servidor
	jwtToken string
	reader   = bufio.NewReader(os.Stdin)
)

func main() {
	clearScreen()
	displayWelcome()

	for {
		displayMainMenu()
		choice := readInput("Digite sua opção: ")

		switch choice {
		case "1":
			login()
		case "2":
			menuCRUD("áreas")
		case "3":
			menuCRUD("fazendas")
		case "4":
			menuCRUD("lotes")
		case "5":
			menuCRUD("pulverizações")
		case "6":
			menuCRUD("produtos")
		case "7":
			menuCRUD("usuários")
		case "8":
			menuAssociacoes()
		case "0":
			fmt.Println("Saindo...")
			return
		default:
			color.Red("Opção inválida! Tente novamente.")
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func displayWelcome() {
	color.Cyan("====================================")
	color.Cyan("    Agro Spray Control - CLI")
	color.Cyan("====================================")
	fmt.Println()
}

func displayMainMenu() {
	color.Yellow("\nMENU PRINCIPAL")
	fmt.Println("1. Login")
	fmt.Println("2. Áreas")
	fmt.Println("3. Fazendas")
	fmt.Println("4. Lotes")
	fmt.Println("5. Pulverizações")
	fmt.Println("6. Produtos")
	fmt.Println("7. Usuários")
	fmt.Println("8. Associações")
	fmt.Println("0. Sair")
}

func readInput(prompt string) string {
	color.Cyan(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func login() {
	clearScreen()
	color.Cyan("==== LOGIN ====")
	cpf := readInput("CPF: ")
	senha := readInput("Senha: ")

	loginData := map[string]string{
		"cpf":   cpf,
		"senha": senha,
	}

	resp, err := makeRequest("POST", "/login", loginData, "")
	if err != nil {
		color.Red("Erro ao fazer login: %v", err)
		return
	}

	var result map[string]string
	if err := json.Unmarshal(resp, &result); err != nil {
		color.Red("Erro ao decodificar resposta: %v", err)
		return
	}

	jwtToken = result["token"]
	color.Green("Login realizado com sucesso!")
}

func menuCRUD(entity string) {
	for {
		clearScreen()
		color.Yellow("\nMENU %s", strings.ToUpper(entity))
		fmt.Println("1. Listar")
		fmt.Println("2. Adicionar")
		fmt.Println("0. Voltar")

		choice := readInput("Digite sua opção: ")
		switch choice {
		case "1":
			listEntities(entity)
		case "2":
			addEntity(entity)
		case "0":
			return
		default:
			color.Red("Opção inválida!")
		}
	}
}

func listEntities(entity string) {
	clearScreen()
	color.Cyan("==== LISTAR %s ====", strings.ToUpper(entity))

	page := readInput("Página (Enter para padrão 1): ")
	if page == "" {
		page = "1"
	}
	limit := readInput("Itens por página (Enter para padrão 20): ")
	if limit == "" {
		limit = "20"
	}

	url := fmt.Sprintf("/%s?page=%s&limit=%s", entity, page, limit)
	resp, err := makeRequest("GET", url, nil, jwtToken)
	if err != nil {
		color.Red("Erro ao listar %s: %v", entity, err)
		readInput("\nPressione Enter para continuar...")
		return
	}

	prettyPrintJSON(resp)
	readInput("\nPressione Enter para continuar...")
}

func addEntity(entity string) {
	clearScreen()
	color.Cyan("==== ADICIONAR %s ====", strings.ToUpper(entity))

	var data map[string]interface{}
	switch entity {
	case "áreas":
		data = map[string]interface{}{
			"cod":         readInput("Código: "),
			"nome":        readInput("Nome: "),
			"tamanhoHa":   readFloatInput("Tamanho (ha): "),
			"fazendaCod":  readInput("Código da fazenda: "),
			"localizacao": readInput("Localização: "),
			"tipoSolo":    readInput("Tipo de solo: "),
		}
	case "fazendas":
		data = map[string]interface{}{
			"cod":  readInput("Código: "),
			"nome": readInput("Nome: "),
			"localizacao": readInput("Localização: "),
			"tamanhoTotal": readFloatInput("Tamanho total (ha): "),
		}
	case "lotes":
		data = map[string]interface{}{
			"cod":      readInput("Código: "),
			"nome":     readInput("Nome: "),
			"areaCod":  readInput("Código da área: "),
			"produtoCod": readInput("Código do produto: "),
			"quantidade": readFloatInput("Quantidade: "),
		}
	case "pulverizações":
		data = map[string]interface{}{
			"cod":        readInput("Código: "),
			"data":       readInput("Data (YYYY-MM-DD): "),
			"produtoCod": readInput("Código do produto: "),
			"quantidade": readFloatInput("Quantidade: "),
			"observacao": readInput("Observação: "),
		}
	case "produtos":
		data = map[string]interface{}{
			"cod":        readInput("Código: "),
			"nome":       readInput("Nome: "),
			"tipo":       readInput("Tipo: "),
			"fabricante": readInput("Fabricante: "),
			"descricao":  readInput("Descrição: "),
		}
	case "usuários":
		data = map[string]interface{}{
			"cpf":           readInput("CPF: "),
			"nome":          readInput("Nome: "),
			"email":         readInput("Email: "),
			"passwordHash":  readInput("Senha: "),
			"permissaoNivel": readInput("Nível de permissão (admin/operador): "),
		}
	}

	resp, err := makeRequest("POST", "/"+entity, data, jwtToken)
	if err != nil {
		color.Red("Erro ao adicionar %s: %v", entity, err)
		readInput("\nPressione Enter para continuar...")
		return
	}

	color.Green("%s adicionado com sucesso!", strings.Title(entity))
	prettyPrintJSON(resp)
	readInput("\nPressione Enter para continuar...")
}

func menuAssociacoes() {
	for {
		clearScreen()
		color.Yellow("\nMENU ASSOCIAÇÕES")
		fmt.Println("1. Listar áreas de uma pulverização")
		fmt.Println("2. Listar pulverizações de uma área")
		fmt.Println("3. Associar área a pulverização")
		fmt.Println("0. Voltar")

		choice := readInput("Digite sua opção: ")
		switch choice {
		case "1":
			listAreasPulverizacao()
		case "2":
			listPulverizacoesArea()
		case "3":
			associarAreaPulverizacao()
		case "0":
			return
		default:
			color.Red("Opção inválida!")
		}
	}
}

func listAreasPulverizacao() {
	clearScreen()
	color.Cyan("==== ÁREAS DE UMA PULVERIZAÇÃO ====")
	codPulv := readInput("Digite o código da pulverização: ")

	resp, err := makeRequest("GET", "/pulverizacoes/"+codPulv+"/areas", nil, jwtToken)
	if err != nil {
		color.Red("Erro ao buscar áreas: %v", err)
		readInput("\nPressione Enter para continuar...")
		return
	}

	prettyPrintJSON(resp)
	readInput("\nPressione Enter para continuar...")
}

func listPulverizacoesArea() {
	clearScreen()
	color.Cyan("==== PULVERIZAÇÕES DE UMA ÁREA ====")
	codArea := readInput("Digite o código da área: ")

	resp, err := makeRequest("GET", "/areas/"+codArea+"/pulverizacoes", nil, jwtToken)
	if err != nil {
		color.Red("Erro ao buscar pulverizações: %v", err)
		readInput("\nPressione Enter para continuar...")
		return
	}

	prettyPrintJSON(resp)
	readInput("\nPressione Enter para continuar...")
}

func associarAreaPulverizacao() {
	clearScreen()
	color.Cyan("==== ASSOCIAR ÁREA A PULVERIZAÇÃO ====")

	data := map[string]string{
		"codPulv": readInput("Código da pulverização: "),
		"codArea": readInput("Código da área: "),
	}

	resp, err := makeRequest("POST", "/associar", data, jwtToken)
	if err != nil {
		color.Red("Erro ao associar: %v", err)
		readInput("\nPressione Enter para continuar...")
		return
	}

	color.Green("Associação criada com sucesso!")
	prettyPrintJSON(resp)
	readInput("\nPressione Enter para continuar...")
}

func readFloatInput(prompt string) float64 {
	for {
		input := readInput(prompt)
		val, err := strconv.ParseFloat(input, 64)
		if err == nil {
			return val
		}
		color.Red("Valor inválido! Digite um número.")
	}
}

func makeRequest(method, path string, body interface{}, token string) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s - %s", resp.Status, string(respBody))
	}

	return respBody, nil
}

func prettyPrintJSON(data []byte) {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		fmt.Println(string(data))
		return
	}

	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println(string(data))
		return
	}

	fmt.Println(string(pretty))
}
