package models

type Area struct {
	Cod             string    `json:"cod" db:"cod"`
	Tamanho     int    `json:"tamanho" db:"tamanho"`
	CodFazenda string    `json:"codfazenda" db:"codfazenda"`
}


func NovaArea(cod, codfazenda string, tamanho int) *Area {

    
  a := &Area{
    Cod: cod,
    Tamanho: tamanho,
    CodFazenda: codfazenda,
  }

  return a
 }

