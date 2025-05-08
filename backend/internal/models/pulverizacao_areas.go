package models

type PulverizacaoArea struct {
	CodPulv            string    `json:"codpulv" db:"codpulv"`
	CodArea string    `json:"codarea" db:"codarea"`
}


func NovaPulverizacaoArea (codpulv, codarea string) *PulverizacaoArea {

    
    
  p := &PulverizacaoArea{
      CodPulv: codpulv,
      CodArea: codarea,
  }

  return p
 }

