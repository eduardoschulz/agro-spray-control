package models

import  ( 
    "time"
    "golang.org/x/crypto/bcrypt"
)

type Usuario struct {
    cpf string
    email string
    nome string
    password_hash []byte
    creation_time time.Time
}

/*
Isso provavelmente não está correto
TODO procurar uma maneira melhor de tratar isso aqui
*/
func gerar_senhahash(senha string) ([]byte, error) {
    
    hash, err := bcrypt.GenerateFromPassword([]byte(senha), 10);
    return hash, err

}


func novo_usuario(cpf, email, nome string, hash []byte) *Usuario {

    
  u := &Usuario{
    cpf: cpf,
    email: email, 
    nome: nome,
    password_hash: hash,
    creation_time: time.Now(),
  }

  return u
 }

