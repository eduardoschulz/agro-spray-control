package models

type Fazenda struct {
	Cod             string `json:"cod" db:"cod"`
	Localizacao     string `json:"localizacao" db:"localizacao"`
	CpfProprietario string `json:"cpf_proprietario" db:"cpf_proprietario"`
}

func NovaFazenda(cod, localizacao, cpf_proprietario string) *Fazenda {

	f := &Fazenda{
		Cod:             cod,
		Localizacao:     localizacao,
		CpfProprietario: cpf_proprietario,
	}

	return f
}
