package models

import (
    "testing"
    "golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T){

    string_teste := "senhaforte123"

    hash, err := GerarSenha(string_teste)

    if err != nil {
        t.Fatalf("Erro ao gerar hash %s\n", err)
    }
    
    // Teste se a senha pode ser verificada com o hash
    err = bcrypt.CompareHashAndPassword(hash, []byte(string_teste))
    if err != nil {
        t.Fatalf("Hash não corresponde à senha original: %v", err)
    }
    
    // Teste com senha errada
    err = bcrypt.CompareHashAndPassword(hash, []byte("senhaerrada"))
    if err == nil {
        t.Fatal("Hash não deveria corresponder a senha errada")
    }
  

}
