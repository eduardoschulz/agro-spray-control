package models

import (
	"log"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type Usuario struct {
	CPF           string    `json:"cpf" db:"cpf"`
	Email         string    `json:"email" db:"email"`
	Nome          string    `json:"nome" db:"nome"`
	PasswordHash  []byte    `json:"-" db:"password_hash"` 
	NivelPermissao int8       `json:"nivel_permissao" db:"nivel_permissao"`
	CriadoEm      time.Time `json:"criado_em" db:"criado_em"`
//	AtualizadoEm  time.Time `json:"atualizado_em" db:"atualizado_em"`
}


func GerarSenha(senha string) ([]byte, error) {
    return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost);
}

func (u *Usuario) VerificarSenha(senha string) error{
    //retorna 0 ou 1 erros 
    return bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(senha))
}


func NovoUsuario(cpf, email, nome, senha string, nivel_permissao int8) (*Usuario, error) {
    
    hash, err := GerarSenha(senha)

    if err != nil {
        log.Printf("Erro ao gerar senha %s\n", err)
        return nil, err
    }
    
  u := &Usuario{
    CPF: cpf,
    Email: email, 
    Nome: nome,
    PasswordHash: hash,
    NivelPermissao: nivel_permissao,
    CriadoEm: time.Now(),
  }

  return u, nil
 }

