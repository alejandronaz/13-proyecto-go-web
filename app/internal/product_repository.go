package internal

type ProductRepository interface {
	GetAllProducts() []Product
	GetProductByID(id int) Product
	GetProductsByPriceGreaterThan(price float64) []Product
	AddProduct(product Product) Product
	UpdateProduct(product Product) (Product, error)
	DeleteProduct(id int) error
}
