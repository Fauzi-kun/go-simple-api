package main

type Note struct{
	ID uint `gorm:"primaryKey" json:"id"`
	Title string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}