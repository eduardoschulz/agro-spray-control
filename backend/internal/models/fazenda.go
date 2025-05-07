package models

/* import  ( 
    "time"
    "golang.org/x/crypto/bcrypt"
) */



type Fazenda struct{
    cod string
    localizacao string
    cpf_proprietario string
}


func novo_fazenda(cod, localizacao, cpf_proprietario string) *Fazenda {

    
  f := &Fazenda{
    cod: cod,
    localizacao: localizacao,
    cpf_proprietario: cpf_proprietario,
  }

  return f
 }

