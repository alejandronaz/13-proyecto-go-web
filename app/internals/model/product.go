package model

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func (p *Product) IsEmpty() bool {
	return p.ID == 0 && p.Name == "" && p.Quantity == 0 && p.CodeValue == "" && p.IsPublished == false && p.Expiration == "" && p.Price == 0
}
