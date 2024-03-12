package domain

type Book struct {
	CommonModel
	Title string `json:"title"`
}
