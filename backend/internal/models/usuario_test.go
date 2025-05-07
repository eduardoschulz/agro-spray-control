package models

import (
    "testing"
    "golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T){

    string_teste := "senhaforte123"

    hash, err := gerar_senhahash(string_teste)

    if err != nil {
        t.Fatalf("Erro ao gerar hash %s\n", err)
    }
    
    comp_hash ,_ := bcrypt.GenerateFromPassword([]byte(string_teste), 10)

    if string(comp_hash) == string(hash) { //nao tenho certeza se isso ta certo
        t.Fatalf("Hash da senha difere")
    }


}
