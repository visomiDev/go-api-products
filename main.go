package main

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
)

// Product model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, `Name`, `Code`, `Price`
type Product struct {
	gorm.Model
	Name  string
	Code  string
	Price uint
}

// ProductRequest definition, including fields `Name`, `Code`, `Price`
type ProductRequest struct {
	Name  string `json:"name"`
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

// BaseResponse definition, including fields `Message`, `Success`
type BaseResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// ProductCreatedResponse definition, including fields `Data`
type ProductCreatedResponse struct {
	BaseResponse
	Data *Product `json:"data"`
}

// ProductsGetAllResponse definition include fields `Data`
type ProductsGetAllResponse struct {
	BaseResponse
	Data []*Product `json:"data"`
}

// ProductGetOneResponse definition include fields `Data`
type ProductGetOneResponse struct {
	BaseResponse
	Data *Product `json:"data"`
}

func main() {
	e := echo.New()

	db, err := gorm.Open("postgres", "postgres://tcpgqyol:t59UI8naIHw-FuQhLCsRtkfuf_RvwALH@ruby.db.elephantsql.com:5432/tcpgqyol")

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&Product{})

	e.POST("/products", func(c echo.Context) error {
		pr := new(ProductRequest)

		if err = c.Bind(pr); err != nil {
			return err
		}

		p := new(Product)

		p.Name = pr.Name
		p.Code = pr.Code
		p.Price = pr.Price

		db.Create(&p)

		return c.JSON(http.StatusCreated, &ProductCreatedResponse{
			BaseResponse: BaseResponse{
				Message: "The product was created",
				Success: true,
			},
			Data: p,
		})
	})

	e.GET("/products", func(c echo.Context) error {
		products := []*Product{}

		db.Find(&products)

		return c.JSON(http.StatusOK, &ProductsGetAllResponse{
			BaseResponse: BaseResponse{
				Message: "Products found",
				Success: true,
			},
			Data: products,
		})
	})

	e.GET("/products/:id", func(c echo.Context) error {
		id := c.Param("id")

		i, err := strconv.Atoi(id)

		if err != nil {
			return err
		}

		p := new(Product)

		db.First(p, i)

		if p.ID == 0 {
			return c.JSON(http.StatusNotFound, &BaseResponse{
				Message: "The product not exists",
				Success: false,
			})
		}

		return c.JSON(http.StatusOK, &ProductGetOneResponse{
			BaseResponse: BaseResponse{
				Message: "Product found",
				Success: true,
			},
			Data: p,
		})
	})

	e.PUT("/products/:id", func(c echo.Context) error {
		pr := new(ProductRequest)

		if err := c.Bind(pr); err != nil {
			return err
		}

		id := c.Param("id")

		i, err := strconv.Atoi(id)

		if err != nil {
			return err
		}

		p := new(Product)

		db.First(p, i)

		p.Code = pr.Code
		p.Name = pr.Name
		p.Price = pr.Price

		db.Save(p)

		return c.JSON(http.StatusOK, &ProductGetOneResponse{
			BaseResponse: BaseResponse{
				Message: "Product updated",
				Success: true,
			},
			Data: p,
		})
	})

	e.DELETE("/products/:id", func(c echo.Context) error {
		id := c.Param("id")

		i, err := strconv.Atoi(id)

		if err != nil {
			return err
		}

		p := new(Product)

		db.First(p, i)
		db.Delete(p)

		return c.JSON(http.StatusOK, &BaseResponse{
			Message: "Product Deleted",
			Success: true,
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
