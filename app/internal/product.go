package internal

type Product struct {
	ID          int
	Name        string
	Quantity    int
	CodeValue   string
	IsPublished bool
	Expiration  string
	Price       float64
}

func (p *Product) IsEmpty() bool {
	return p.ID == 0 && p.Name == "" && p.Quantity == 0 && p.CodeValue == "" && !p.IsPublished && p.Expiration == "" && p.Price == 0
}
