package main

import "math"

type Positioned interface {
	Position() (int, int)
}

type Hex struct {
	row, col int
}

func (h *Hex) Position() (int, int) {
	return h.row, h.col
}

func (h *Hex) Cartesian() (float64, float64) {
	x := float64(h.col)
	if h.row % 2 == 1 {
		x = x + 0.5
	}

	y := float64(h.row) * math.Sqrt(3) / 2.0

	return x, y
}
	
func (h *Hex) Distance(h2 Hex) float64 {
	x1, y1 := h.Cartesian()
	x2, y2 := h2.Cartesian()

	return math.Sqrt(math.Pow(x1 - x2, 2) + math.Pow(y1 - y2, 2))
}

func (h *Hex) Adjacent(h2 Hex) bool {
	dist := h.Distance(h2)
	return dist > 0.0 && dist < 1.1
}

func (h *Hex) Closer(h2 Hex) {
	min := h.Distance(h2)
	closest := h
	
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			cand := Hex{h.row + dr, h.col + dc}

			if cand.Adjacent(*h) && h2.Distance(cand) < min {
				min = h2.Distance(cand)
				closest = &cand
			}
		}
	}

	h.row = closest.row
	h.col = closest.col
}

func (h *Hex) Copy(h2 Hex) {
	h.row = h2.row
	h.col = h2.col
}

