package models

type Produto struct {
	Cod         string `json:"cod" db:"cod"`
	Descricao   string `json:"descricao" db:"descricao"`
	Fabricante  string `json:"fabricante" db:"fabricante"`
	CompQuimica string `json:"compquimica" db:"compquimica"`
}

func NovoProduto(cod, descricao, fabricante, compquimica string) *Produto {

	p := &Produto{
		Cod:         cod,
		Descricao:   descricao,
		Fabricante:  fabricante,
		CompQuimica: compquimica,
	}
	return p
}
