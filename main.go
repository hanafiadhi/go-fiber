package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main()  {
	app := fiber.New(fiber.Config{
		IdleTimeout: time.Second * 5,
		WriteTimeout :time.Second * 5,
		ReadTimeout: time.Second * 5,
		Prefork: true,
	})

	app.Use(func (ctx *fiber.Ctx) error {
		fmt.Println("Iam middleware before processing request")
		err:= ctx.Next()
		fmt.Println("Iam middleware after processing request")
		return err
	})
	if fiber.IsChild() {
		fmt.Println("I'am a child")
	}else{
		fmt.Println("I'am a parent")
	}
	app.Get("/",func(c *fiber.Ctx) error {
		return c.SendString("Hallo golang fiber")
	})
	err := app.Listen("127.0.0.1:8080")

	if err != nil {
		panic(err)
	}
}
