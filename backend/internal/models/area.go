package models

type Area struct {
	Cod        string `json:"cod" db:"cod"`
	Tamanho    int    `json:"tamanho" db:"tamanho"`
	CodFazenda string `json:"cod_fazenda" db:"cod_fazenda"`
}

func NovaArea(cod, codfazenda string, tamanho int) *Area {

	a := &Area{
		Cod:        cod,
		Tamanho:    tamanho,
		CodFazenda: codfazenda,
	}

	return a
}
