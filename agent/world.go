package main

import (
	"math"
	"math/rand"
        "bytes"
)

// A World is the hex-grid world the simulation takes place in.
type World struct {
	width, height int
	agents []Agent
}

// Random returns a random location in the world.
func (w *World) Random() Hex {
	return Hex{rand.Intn(w.height), rand.Intn(w.width)}
}

// RandomBorder returns a random location on the border of the world.
func (w *World) RandomBorder() Hex {
	perimeter := (w.width + w.height) * 2 - 4;

	r := rand.Intn(perimeter)

	if r < w.width {
		// top edge ([0, 0] to [0, width - 1]
		return Hex{0, r}
	}

	r = r - w.width
	if r < w.height - 1 {
		// right edge (1, width - 1] to [height - 1, width - 1])
		return Hex{r + 1, w.width - 1}
	}

	r = r - (w.height - 1)
	if r < w.width - 1 {
		// bottom edge ([height - 1, 0] to [height - 1, width - 2]
		return Hex{w.height - 1, r}
	}
	
	// left egde ([1, 0] to [height - 1, 0]
	r = r - (w.width - 1)
	return Hex{r + 1, 0}
}

// Update updates all the agents in the world, initiates interactions between them,
// removes dead agents, and adds newly spawned agents.
func (w *World) Update() {
	// update all the agents
	for _, a := range w.agents {
		a.Update()
	}

	// check for interactions
	for _, a1 := range w.agents {
		for _, a2 := range w.agents {
			if a1 != a2 {
				r1, c1 := a1.Position()
				r2, c2 := a2.Position()

				if r1 == r2 && c1 == c2 {
					switch a2.(type) {
					case *Prey:
						a1.AcceptPrey(a2.(*Prey))
					case *Predator:
						a1.AcceptPredator(a2.(*Predator))
					case *Scent:
						a1.AcceptScent(a2.(*Scent))
					}

				}
			}
		}
	}

	// remove dead agents and spawn new ones
	next := make([]Agent, 0, 10)
	for _, a := range w.agents {
		if a.Alive() {
			next = append(next, nil)
			next[len(next) - 1] = a

			s := a.Spawn()
			if s != nil {
				next = append(next, nil)
				next[len(next) - 1] = s
			}
		}
	}

	w.agents = next
}
// String returns a printable representation of the world.
func (w *World) String() string {
	agentMap := make([][]byte, w.height)
	for i := 0; i < len(agentMap); i++ {
		agentMap[i] = make([]byte, w.width)
		for j := 0; j < len(agentMap[i]); j++ {
			agentMap[i][j] = '.'
		}
	}

	for _, a := range w.agents {
		r, c := a.Position()
		if r >= 0 && r < w.height && c >= 0 && c < w.width {
			agentMap[r][c] = a.String()[0]
		}
	}

	var buf bytes.Buffer
	for r := 0; r < len(agentMap); r++ {
		if r % 2 == 1 {
			buf.WriteString(" ")
		}

		for c := 0; c < len(agentMap[r]); c++ {
			buf.WriteByte(agentMap[r][c])
			buf.WriteString(" ")
		}
		buf.WriteString("\n")
	}

	return buf.String()
			
}

// NewWorld creates a new world with one predator, one food source, and one prey.
func NewWorld(w int) *World {
	world := World{w, int(float64(w) * 2 / math.Sqrt(3)), make([]Agent, 0, 10)}
	world.agents = append(world.agents, nil)
	world.agents[len(world.agents) - 1] = NewPredator(&world)

	world.agents = append(world.agents, nil)
	world.agents[len(world.agents) - 1] = NewFood(&world)

	world.agents = append(world.agents, nil)
	world.agents[len(world.agents) - 1] = NewPrey(&world)

	return &world
}
