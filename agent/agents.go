package main

import "math/rand"

// An Agent is an agent in a simulation.
type Agent interface {
	Positioned
	Update()
	Spawn() Agent
	Alive() bool
        String() string
	World() *World
	AcceptScent(s *Scent)
	AcceptPredator(p *Predator)
	AcceptPrey(p *Prey)
}

// A Generic is a generic agent with a home world and a location.
type Generic struct {
	world *World
	Hex
}

// World returns the world a generic agent lives in.
func (g *Generic) World() *World {
	return g.world
}

// NewGeneric creates and returns a new generic at a random location in the given world.
func NewGeneric(world *World) Generic {
	return Generic{world, world.Random()}
}

// NoPredatorAction supplies a do-nothing AcceptPredator function.
type NoPredatorAction struct {
}

func (a *NoPredatorAction) AcceptPredator(p *Predator) {
}

// NoPreyAction supplies a do-nothing AcceptPrey function.
type NoPreyAction struct {
}

func (a *NoPreyAction) AcceptPrey(p *Prey) {
}

// NoScentAction supplies a do-nothing AcceptScent function.
type NoScentAction struct {
}

func (a *NoScentAction) AcceptScent(s *Scent) {
}

// NoActions supplies do-nothing functions for interations with predators, prey, and scents.
type NoActions struct {
     NoPredatorAction
     NoPreyAction
     NoScentAction
}

// EatenByPredator supplies a AcceptPredator method that sets a bool to false.
type EatenByPredator struct {
	alive *bool
}

func (e *EatenByPredator) AcceptPredator(p *Predator) {
	*(e.alive) = false
}

// EatenByPredator supplies a AcceptPrey method that sets a bool to false.
type EatenByPrey struct {
	alive *bool
}

func (e *EatenByPrey) AcceptPrey(p *Prey) {
	*(e.alive) = false
}

// A Stationary is a Generic agent that also supplies a do-nothing Update function.
type Stationary struct {
	Generic
}

func (s *Stationary) Update() {
}

// A Mobile is a Generic agent that also supplies an Update function that moves the
// agent one step closer to a goal.  Agents that embed Mobile must supply a ReachedGoal
// function that determines what action to take when the agent reaches its goal.
type Mobile struct {
	Generic
	goal *Hex
	ReachedGoal func()
}

// Update updates the location of a Mobile agent to be one step closer to its current goal.
// When it reaches the goal, it calls the ReachedGoal function to determine the action
// to take.
func (s *Mobile) Update() {
	s.Closer(*(s.goal))
	r1, c1 := s.Position()
	r2, c2 := s.goal.Position()
	if r1 == r2 && c1 == c2 {
		s.ReachedGoal()
	}
}

// Immortal supplies an always-alive Alive function.
type Immortal struct {
}

func (i *Immortal) Alive() bool {
	return true
}

// NonEmitter supplies a spawn-nothing Spawn function.
type NonEmitter struct {
}

func (n *NonEmitter) Spawn() Agent {
	return nil
}

// ScentFollower supplies an AcceptScent function that changes an externally allocated
// Hex to the origin of the Scent.
type ScentFollower struct {
     goal *Hex
}

// AcceptScent changes the external goal to the origin of the Scent.
func (a *ScentFollower) AcceptScent(s *Scent) {
	*(a.goal) = s.Origin()
}

// A Predator is a mobile, immortal agent that spawns nothing, doesn't react when acted on
// by other predators or preys, follows scents, and chooses a random location in the world
// as its new goal when it reaches its current goal.
type Predator struct {
	Mobile
	Immortal
	NonEmitter
	NoPredatorAction
	NoPreyAction
	ScentFollower
}

// String returns "P" as the printable representation of a Predator.
func (p *Predator) String() string {
	return "P"
}

// NewPredator creates and returns a new predator at a random location in the world
// with a randomly selected goal.
func NewPredator(w *World) *Predator {
	goal := w.Random()
	p := &Predator{
		Mobile{
			NewGeneric(w),
			&goal,
			func () { goal.Copy(w.Random()) }} ,
		Immortal{},
		NonEmitter{},
		NoPredatorAction{},
		NoPreyAction{},
		ScentFollower{&goal} }
	
	return p
}

// Mortal supplies an Alive function that returns the value of an externally allocated
// boolean.
type Mortal struct {
	alive *bool
}

// Alive returns the value of the external alive field.
func (m *Mortal) Alive() bool {
	return *(m.alive)
}

// An Origined is something that has a starting location.
type Origined interface {
	Origin() Hex
}

// Emitted implements Origined by supplying a field to hold the origin.
type Emitted struct {
	origin Hex
}

func (e *Emitted) Origin() Hex {
	return e.origin
}

// A Scent is an Origined, Mobile, Mortal Agent that is spawned by an Agent.  It dies when
// it reaches the goal it randomly selects upon creation.
type Scent struct {
	Mobile
	NonEmitter
	Mortal
	Emitted
	NoActions
}

// String returns "x" as the printable representation of a Scent.
func (s *Scent) String() string {
	return "x"
}

// NewScent returns a new Scent that starts at the location of the given Agent.
func NewScent(a Agent) *Scent {
	r, c := a.Position()
	origin := Hex{r, c}
	goal := a.World().RandomBorder()

	alive := true

	return &Scent{
		Mobile{
			Generic{a.World(), Hex{r, c}},
			&goal,
			func () { alive = false } },
		NonEmitter{},
		Mortal{&alive},
	        Emitted{origin},
		NoActions{} }
}

// Emitter supplies a Spawn function that spawns a new Agent at the rate given by the
// rate field.  The agent to be spawned is determined by the Emit function that embedding
// objects must supply.
type Emitter struct {
	rate int
	Emit func () Agent
}

// Spawn spawns an Agent at an average rate of once per whatever value is held in the rate
// field.  The Agent to spawn is determined by the Emit function.
func (e *Emitter) Spawn() Agent {
	if rand.Intn(e.rate) == 0 {
		return e.Emit()
	} else {
		return nil
	}
}

// Food is a Stationary, Mortal Agent that emits Scents on average once per 5 time steps
// and dies when it is eaten by a Prey.
type Food struct {
	Stationary
	Emitter
	Mortal
	NoPredatorAction
	EatenByPrey
	NoScentAction
}

// String returns "*" as the printable representation of Food.
func (f *Food) String() string {
	return "*"
}

// NewFood returns a new Food at a randomly chosen location in the world.
func NewFood(w *World) *Food {
	alive := true

	f := Food{
		Stationary{
			Generic{
				w,
				w.Random()} },
		Emitter{
			5,
			nil },
		Mortal{&alive},
		NoPredatorAction{},
		EatenByPrey{&alive},
		NoScentAction{} }
	f.Emit = func () Agent { return NewScent(&f) }

	return &f
}

// A Prey is a Mobile, Mortal agent that emits nothing, dies when eaten by a Predator,
// and follows Scents.
type Prey struct {
	Mobile
	NonEmitter
	Mortal
	EatenByPredator
	NoPreyAction
	ScentFollower
}

// String returns "p" as the printable representation of a Prey.
func (p *Prey) String() string {
	return "p"
}

// NewPrey creates and returns a new Prey at a randomly chosen location and with a randomly
// chosen goal.
func NewPrey(w *World) *Prey {
	goal := w.Random()
	alive := true
	p := &Prey{
		Mobile{
			NewGeneric(w),
			&goal,
			func () { goal.Copy(w.Random()) }} ,
		NonEmitter{},
	        Mortal{&alive},
		EatenByPredator{&alive},
		NoPreyAction{},
		ScentFollower{&goal} }

	return p
}
