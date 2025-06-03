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
     if len(senha) < 8 {
        return nil, errors.New("senha deve ter pelo menos 8 caracteres")
    }
    return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

func (u *Usuario) VerificarSenha(senha string) error{
    //retorna 0 ou 1 erros 
    return bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(senha))
}

func NovoUsuario(cpf, email, nome, senha string, nivel_permissao int8) (*Usuario, error) {
    // Validação básica
    if !validarCPF(cpf) {
        return nil, errors.New("CPF inválido")
    }
    
    if !validarEmail(email) {
        return nil, errors.New("e-mail inválido")
    }
    
    hash, err := GerarSenha(senha)
    if err != nil {
        return nil, err
    }
    
    return &Usuario{
        CPF:           cpf,
        Email:         email, 
        Nome:          nome,
        PasswordHash:  hash,
        NivelPermissao: nivel_permissao,
        CriadoEm:      time.Now(),
    }, nil
}

func validarCPF(cpf string) bool {
    // Implementar validação real de CPF
    return len(cpf) == 11
}

func validarEmail(email string) bool {
    // Regex simples para validação de e-mail
    return regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(email)
}
