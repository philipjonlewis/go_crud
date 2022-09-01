package main

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID 			string `json:"id"`
	Title 		string `json:"title"`
	Author 		string `json:"author"`
	Quantity 	int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book,err := getBookById(id)

	if err != nil { 
		c.IndentedJSON(http.StatusNotFound,gin.H{"message":"Book Not Found"})
		return 
	}

	c.IndentedJSON(http.StatusOK,book)
}

func getBookById(id string) (*book,error) {
	for i,b := range books {
		if b.ID == id {
			return &books[i],nil
		}
	}

	return nil,errors.New("book not Found")
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book ;

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books,newBook)

	c.IndentedJSON(http.StatusCreated,newBook)


}

func checkoutBook(c *gin.Context) {
	id,ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Missing Id Query parameter"})
	}

	book ,err := getBookById(id)
	
	if err != nil {
		c.IndentedJSON(http.StatusNotFound,gin.H{"message":"Book Not Found"})
		return 
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Book Not Available"})
		return 
	}

	book.Quantity -= 1

	c.IndentedJSON(http.StatusOK,book)

}

func returnBook(c *gin.Context) {
	id,ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Missing Id Query parameter"})
	}

	book ,err := getBookById(id)
	
	if err != nil {
		c.IndentedJSON(http.StatusNotFound,gin.H{"message":"Book Not Found"})
		return 
	}

	
	book.Quantity += 1

	c.IndentedJSON(http.StatusOK,book)

}

func main () {

	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
		})
	})

	router.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.html", gin.H{
			"content": "This is an about page and the content is from the main file",
		})
	})


	router.GET("/books", getBooks)

	router.GET("/books/:id", bookById)

	router.PATCH("/checkout",checkoutBook)

	router.PATCH("/return",returnBook)

	router.POST("/books", createBook)

	router.Run("localhost:8080")



}



