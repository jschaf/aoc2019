package combinations

type Perm struct {
	orig  []int
	state []int
}

func NewPermuter(orig []int) *Perm {
	state := make([]int, len(orig))
	return &Perm{orig, state}
}

func (p *Perm) HasNext() bool {
	return p.state[0] < len(p.state)
}

func (p *Perm) Next() {
	for i := len(p.state) - 1; i >= 0; i-- {
		if i == 0 || p.state[i] < len(p.state)-i-1 {
			p.state[i]++
			return
		}
		p.state[i] = 0
	}
}

func (p *Perm) Get() []int {
	result := make([]int, len(p.orig))
	copy(result, p.orig)
	for i, v := range p.state {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}
