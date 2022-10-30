package schemes

type Paginator struct {
	Page  int32
	Limit int32
}

func (p *Paginator) Offset() int32 {
	return (p.Page - 1) * p.Limit
}
