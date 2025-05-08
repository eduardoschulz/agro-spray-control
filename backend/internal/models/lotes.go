package models

import "time"

type Lote struct {
	Cod             string    `json:"cod" db:"cod"`
	DtValidade time.Time     `json:"descricao" db:"descricao"`
    Quantidade int          `json:"quantidade" db:"quantidade"`
	CodProduto string    `json:"CodProduto" db:"CodProduto"`
}


func NovoLote(cod, codproduto string, quantidade int, dtvalidade time.Time) *Lote {

    
  l := &Lote{
      Cod: cod, 
      DtValidade: dtvalidade,
      Quantidade: quantidade,
      CodProduto: codproduto,
}
  return l
 }

