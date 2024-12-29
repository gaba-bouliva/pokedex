package pokemon

import (
	"math"
	"math/rand"
)

type Pokemon struct {
	ID 							int     			`json:"id"` 
	Name 						string   			`json:"name"`
	Order						int 					`json:"order"`
	Height 					int						`json:"height"`
	Weight    			int						`json:"weight"`
	BaseExperience 	int						`json:"base_experience"`
	Stats 					[]Stat        `json:"stats"`
	Types 					[]TypeDetail	`json:"types"`
}

type Stat struct { 
	Value     		int 			`json:"base_stat"`			
	Detail 									`json:"stat"` 
}

type TypeDetail struct { 
	Detail 			`json:"type"`
}

type Detail struct { 
	Name      string `json:"name"`
	URL 			string `json:"url"`
}

func NewPokemon() Pokemon {
	return Pokemon{}
}

func (p *Pokemon) Catch() bool{ 
	oddsGettingCaught := (float64(rand.Intn(p.BaseExperience))) / float64(p.BaseExperience)
	oddsGettingCaught = math.Floor(oddsGettingCaught * 100)/100
	return oddsGettingCaught >= 0.74 
}