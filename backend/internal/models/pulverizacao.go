package models

import "time"

type Pulverizacao struct {
	Cod             string    `json:"cod" db:"cod"`
    DtAplicacao time.Time `json:"dtaplicacao" db:"dtaplicacao"`
	Cultura     string    `json:"cultura" db:"cultura"`
	CodLote            string    `json:"codlote" db:"codlote"`
	CpfResponsavel string    `json:"cpfresponsavel" db:"cpfresponsavel"`
}


func NovaPulverizacao(dtaplicacao time.Time, cod, cultura, codlote, cpf_responsavel string) *Pulverizacao {

  p := &Pulverizacao{
    Cod: cod,
    DtAplicacao: dtaplicacao,
    Cultura: cultura,
    CpfResponsavel: cpf_responsavel,
  }

  return p
 }

