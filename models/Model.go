package models

import (
	"gorm.io/gorm"
)

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name"`
	Lastname    *string `json:"lastname"`
	Email       string  `gorm:"unique;not null" json:"email"`
	Password    string  `gorm:"not null" json:"-"`
	Role        Role    `gorm:"type:varchar(10);default:'user'" json:"role"`
	Phonenumber *string `json:"phonenumber"`

	Orders    []Order    `gorm:"foreignKey:OrderByID;constraint:OnDelete:CASCADE;" json:"orders"`
	Addresses []Address  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"addresses"`
	Carts     []Cart     `gorm:"foreignKey:OrderByID;constraint:OnDelete:CASCADE;" json:"carts"`
	Favorites []Favorite `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"favorites"`
}

type Product struct {
	gorm.Model
	Name  string  `gorm:"not null" json:"name"`
	Price float64 `gorm:"not null" json:"price"`
	// Categories      []Category       `gorm:"many2many:product_categories;" json:"categories"`
	Categories      []Category       `gorm:"many2many:product_categories;constraint:OnDelete:CASCADE;" json:"categories"`
	Images          []Images         `gorm:"many2many:product_images;constraint:OnDelete:CASCADE;" json:"images"`
	Sold            int              `json:"sold"`
	Quantity        int              `json:"quantity"`
	Favorites       []Favorite       `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"favorites"`
	ProductsOnCart  []ProductOnCart  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"products_on_cart"`
	ProductsOnOrder []ProductOnOrder `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"products_on_order"`
}

type Images struct {
	gorm.Model
	URL     string    `gorm:"not null" json:"url"`
	Product []Product `gorm:"many2many:product_images;constraint:OnDelete:CASCADE;" json:"products"`
}

type Category struct {
	gorm.Model
	Name     string    `gorm:"not null" json:"name"`
	Icon     string    `gorm:"not null" json:"icon"`
	Products []Product `gorm:"many2many:product_categories;constraint:OnDelete:CASCADE;" json:"products"`
}

type Order struct {
	gorm.Model
	CartTotal float64          `gorm:"not null" json:"cart_total"`
	OrderByID uint             `gorm:"not null" json:"order_by_id"`
	OrderBy   User             `gorm:"foreignKey:OrderByID" json:"order_by"`
	Amount    int              `gorm:"not null" json:"amount"`
	Status    string           `gorm:"not null" json:"status"`
	Currency  *string          `json:"currency"`
	AddressID *uint            `json:"address_id"`
	Address   *Address         `gorm:"foreignKey:AddressID" json:"address"`
	Products  []ProductOnOrder `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"products"`
}

type Address struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name"`
	Phonenumber int     `gorm:"not null" json:"phonenumber"`
	Address     string  `gorm:"not null" json:"address"`
	Province    string  `gorm:"not null" json:"province"`
	Amphure     string  `gorm:"not null" json:"amphure"`
	Tambon      string  `gorm:"not null" json:"tambon"`
	Zipcode     int     `gorm:"not null" json:"zipcode"`
	UserID      uint    `gorm:"not null" json:"user_id"`
	User        User    `gorm:"foreignKey:UserID" json:"user"`
	Orders      []Order `gorm:"foreignKey:AddressID;constraint:OnDelete:CASCADE;" json:"orders"`
}

type Cart struct {
	gorm.Model
	CartTotal float64         `gorm:"not null" json:"cart_total"`
	OrderByID uint            `gorm:"not null" json:"order_by_id"`
	OrderBy   User            `gorm:"foreignKey:OrderByID" json:"order_by"`
	Products  []ProductOnCart `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;" json:"products"`
}

type ProductOnCart struct {
	gorm.Model
	CartID    uint    `gorm:"not null" json:"cart_id"`
	Cart      Cart    `gorm:"foreignKey:CartID" json:"cart"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Count     int     `gorm:"not null" json:"count"`
	Price     float64 `gorm:"not null" json:"price"`
}

type Favorite struct {
	gorm.Model
	UserID    uint    `gorm:"not null" json:"user_id"`
	User      User    `gorm:"foreignKey:UserID" json:"user"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
}

type ProductOnOrder struct {
	gorm.Model
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	OrderID   uint    `gorm:"not null" json:"order_id"`
	Order     Order   `gorm:"foreignKey:OrderID" json:"order"`
	Count     int     `gorm:"not null" json:"count"`
	Price     float64 `gorm:"not null" json:"price"`
}
