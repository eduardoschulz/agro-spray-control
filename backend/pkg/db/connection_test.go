package db

import (
    "testing"
)

func TestConnect(t *testing.T) {
    /*string de conexao ao banco de dados */
    connstr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

    /* Chama a funcao Connect do connection.go */
    db := Connect(connstr)

    // Verifica se a conexão não é nula
    if db == nil {
        t.Fatalf("Esperava uma conexão válida, mas a conexão é nil")
    }

    // Testa se o banco está acessível (ping).
    err := db.Ping()
    if err != nil {
        t.Fatalf("Erro ao tentar pingar o banco: %v", err)
    }
}

