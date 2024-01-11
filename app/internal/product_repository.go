package internal

type ProductRepository interface {
	GetAllProducts() []Product
	GetProductByID(id int) Product
	GetProductsByPriceGreaterThan(price float64) []Product
	AddProduct(product Product)
	UpdateProduct(product Product) (Product, error)
}
