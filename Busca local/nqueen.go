package main

import (
	"fmt"
	"math"
	"math/rand"
	"flag"
)

const (
	// Constante que define a proporção do tabuleiro
	STATE_SIZE int = 100
)

const (
	ALG_HILL_CLIMB = iota
	ALG_SIM_ANNELING
)

// Estrutura de dados que representa um estado do tabuleiro. Mantendo as 
// informações sobre a posição da rainha em cada colouna, o tamanho do
// estado e a força comparada a um estado objetivo.
type NQueenState struct {
	Positions []int
	Size int
	Fitness int
}

// Retorna um ponteiro de um estado qualquer definido pelo argumento passado à
// função.
func NewNQueenState ( state []int ) *NQueenState {
	s := new(NQueenState)
	s.Positions = state
	s.Size = len(state)
	return s
}

// Retorna true se a rainha, denotada pelo argumento queen, não for atacada na
// linha definida pelo argumento row.
func ( Me *NQueenState ) IsSafe ( queen int, row int ) bool {
	for i := 0; i < queen; i += 1 {
		val := Me.Positions[i]

		if val == row ||                // Mesma Linha
		   val == row - (queen - i) ||  // Mesma Diagonal
		   val == row + (queen - i) {   // Mesma Diagonal
			return false
		}
	}
	return true
}

// Retorna true se o estado for um dos estados objetivos. Verificando se nenhuma
// das rainhas no tabuleiro esta atacando ou sofrendo ataque em suas respectivas
// posições.
func ( Me *NQueenState ) IsSolution () bool {
	for i, val := range Me.Positions {
		if !Me.IsSafe(i, val) {
			return false
		}
	}
	return true
}

// Retorna os estados vizinhos no espaço de estados da busca local.
func ( Me *NQueenState ) GetNeighbors () []*NQueenState {
	neighbors := make([]*NQueenState, 0)

	for i := 0; i < Me.Size; i += 1 {
		for j := 1; j <= Me.Size; j += 1 {
			state := make([]int, Me.Size)
			copy(state, Me.Positions)
			state[i] = j
			neighbors = append(neighbors, NewNQueenState(state))
		}
	}
	return neighbors
}

// Retorna a força do estado em relação a uma possivel solução.
func ( Me *NQueenState ) GetFitness () int {
	if Me.Fitness == 0 {
		for i, val := range Me.Positions {
			if Me.IsSafe(i, val) {
				Me.Fitness++
			}
		}
	}
	return Me.Fitness
}

// Retorna a maior força de um conjunto de estados.
func GetMaxFitness ( states []*NQueenState ) int {
	max := states[0].Fitness
	for _, val := range states {
		if val.GetFitness() > max {
			max = val.Fitness
		}
	}
	return max
}

// Retorna os melhores estados de um conjunto de estados.
func GetBestStates ( states []*NQueenState ) []*NQueenState {
	max := GetMaxFitness(states)
	bestStates := make([]*NQueenState, 0)
	for _, val := range states {
		if val.Fitness == max {
			bestStates = append(bestStates, val)
		}
	}
	return bestStates
}

// Algoritmo de escalada com reínicio aleatório. Sobe colinas na geografia do
// espaço de estados.
func HillClimbAlgorithm ( initialState []int ) *NQueenState {
	current := NewNQueenState(initialState)
	seed := int64(0)

	for {
		neighbors := current.GetNeighbors()
		bestStates := GetBestStates(neighbors)

		r := rand.New(rand.NewSource(seed))

		// Escolhe um dos melhores entre os estados vizinhos
		index := r.Intn(len(bestStates))
		neighbor := bestStates[index]

		if neighbor.IsSolution() {
			return neighbor
		}

		seed += 1
		current = neighbor
	}
}

// Algoritmo de tempera simulada com reinicio aleatório. Agita o ponteiro entre
// o espaço de estados até achar uma solução ou a temperatura atingir o máximo.
func SimulatedAnnelingAlgorithm ( initialState []int  ) *NQueenState {
	current := NewNQueenState(initialState)
	best := current
	k := 0
	// O número máximo de simulações.
	kmax := 10 * current.Size

	for !current.IsSolution() && k < kmax {
		// O valor da temperatura atual.
		temp := current.GetFitness() / kmax
		neighbors := current.GetNeighbors()
		bestStates := GetBestStates(neighbors)

		r := rand.New(rand.NewSource(int64(k)))

		index := r.Intn(len(bestStates))
		neighbor := bestStates[index]

		if current.Fitness <= neighbor.Fitness {
			current = neighbor
			if best.Fitness <= current.Fitness {
				best = current
			}
		} else {
			// Se a diferença sobre a tempratura for maior que um 
			// número aleatório. Altere o foco da busca para outro 
			// ponto do espaço de estado. 
			aux := float64(current.Fitness - neighbor.Fitness / temp)
			if math.Exp(aux) > r.Float64()  {
				current = neighbor
			}
		}

		k += 1
	}
	return best
}

func main () {
	stateSize := flag.Int("size", STATE_SIZE, "O tamanho de cada estado")
	alg := flag.Int("alg", ALG_HILL_CLIMB, "O algoritmo: \n\t\t 0 - HillClimbAlgorithm \n\t\t 1 - SimulatedAnnelingAlgorithm")

	flag.Parse()

	st := make([]int, *stateSize)

	for i, _ := range st {
		st[i] = 1
	}

	var res *NQueenState

	if *alg == ALG_HILL_CLIMB {
		res = HillClimbAlgorithm(st)
	} else if *alg == ALG_SIM_ANNELING {
		res = SimulatedAnnelingAlgorithm(st)
	}

	for i, val := range res.Positions {
		if i < len(res.Positions) - 1 {
			fmt.Printf("%d,", val)
		} else {
			fmt.Printf("%d", val)
		}
	}
}
