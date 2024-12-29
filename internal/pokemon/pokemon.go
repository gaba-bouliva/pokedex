package pokemon

import (
	"fmt"
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
}

func NewPokemon() Pokemon {
	return Pokemon{}
}

func (p *Pokemon) Catch() bool{ 
	oddsGettingCaught := (float64(rand.Intn(p.BaseExperience))) / float64(p.BaseExperience)
	oddsGettingCaught = math.Floor(oddsGettingCaught * 100)/100
	fmt.Printf("Odds of pokemon %s getting caught: %f\n", p.Name, oddsGettingCaught)
	return oddsGettingCaught >= 0.80 
}