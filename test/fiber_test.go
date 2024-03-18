package test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/stretchr/testify/assert"
)
var app = fiber.New(fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		c.Status(fiber.StatusInternalServerError)
		return c.SendString("Error " + err.Error())
	},
})
func TestRoutingHelloWord(t *testing.T)  {

	app.Get("/",func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	request := httptest.NewRequest("GET","/",nil)
	response,err := app.Test(request)

	assert.Nil(t,err)
	assert.Equal(t,200,response.StatusCode)

	byte, err:= io.ReadAll(response.Body)
	assert.Nil(t,err)
	assert.Equal(t,"Hello World", string(byte))
}

func TestCtx(t *testing.T)  {
	app.Get("/hello", func(c *fiber.Ctx) error {
		query := c.Query("name","Guest")
		return c.SendString("Hello " + query)
	})

	request := httptest.NewRequest("GET","/hello?name=Hanafi",nil)
	response,err:= app.Test(request)
	assert.Nil(t,err)
	assert.Equal(t,200,response.StatusCode)

	byte, err := io.ReadAll(response.Body)
	assert.Nil(t,err)
	assert.Equal(t,"Hello Hanafi", string(byte))

	request = httptest.NewRequest("GET","/hello",nil)
	response,err= app.Test(request)
	assert.Nil(t,err)
	assert.Equal(t,200,response.StatusCode)

	byte, err = io.ReadAll(response.Body)
	assert.Nil(t,err)
	assert.Equal(t,"Hello Guest", string(byte))
}

func TestHttpRequest(t *testing.T)  {
	app.Get("/request", func(c *fiber.Ctx) error {
		firstName := c.Get("firstname")
		lastName := c.Cookies("lastname")
		return c.SendString("Hello " + firstName + " " + lastName )
	})

	request := httptest.NewRequest("GET","/request",nil)
	request.Header.Set("firstname","hanafi")
	request.AddCookie(&http.Cookie{Name: "lastname", Value: "adhi"})


	response, err:= app.Test(request)
	assert.Nil(t,err)
	assert.Equal(t,200, response.StatusCode)

	byte, err:= io.ReadAll(response.Body)
	assert.Nil(t,err)
	assert.Equal(t,"Hello hanafi adhi", string(byte))
}

func TestRouteParams(t *testing.T)  {
	app.Get("user/:userId/order/:orderId",func(c *fiber.Ctx) error {
		userId := c.Params("userId")
		orderId := c.Params("orderId")
		return c.SendString("Get Order "+ orderId + " From User " + userId)
	})

	request := httptest.NewRequest("GET","/user/hanafi/order/123",nil)
	response, err:= app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t,200,response.StatusCode)

	byte,err:= io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t,"Get Order 123 From User hanafi",string(byte))
}

func TestFromRequest(t *testing.T)  {
	app.Post("/hello",func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		log.Debug(name)
		return c.SendString("Hello " + name)
	})
	body := strings.NewReader("name=hanafi")

	request := httptest.NewRequest("POST","/hello", body)
	request.Header.Set("Content-Type","application/x-www-form-urlencoded")
	response, err:= app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t,200,response.StatusCode)

	byte,err:= io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t,"Hello hanafi",string(byte))
}
//go:embed source/ancol.jpeg
var ancol []byte
func TestFormUpload(t *testing.T)  {
	app.Post("/upload",func(c *fiber.Ctx) error {
		file, err:= c.FormFile("image")
		if err != nil {
			return  err
		}
		err = c.SaveFile(file,"./target/"+ file.Filename)
		if err != nil {
			return  err
		}
		return c.SendString("Upload Success")
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, err:= writer.CreateFormFile("image","jalan-jalan.jpeg")
	assert.Nil(t,err)
	file.Write(ancol)
	writer.Close()


	request := httptest.NewRequest("POST","/upload",body)
	request.Header.Set("Content-Type",writer.FormDataContentType())
	response,err:= app.Test(request)

	assert.Nil(t,err)
	assert.Equal(t,200,response.StatusCode)


	byte,err:= io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t,"Upload Success",string(byte))
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
func TestLogin(t *testing.T)  {
	app.Post("/login",func(c *fiber.Ctx) error {
		body := c.Body()
		request := new(LoginRequest)

		err:= json.Unmarshal(body,request)
		if err != nil {
			return err
		}

		return c.SendString("Login Success "+ request.Username)
	})

	body := strings.NewReader(`{"username":"hanafi","password":"adhi"}`)

	request := httptest.NewRequest("POST","/login",body)
	request.Header.Set("Content-Type","application/json")

	response,err:= app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t,200,response.StatusCode)

	byte,err:= io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t,"Login Success hanafi",string(byte))
}

type RegisterRequest struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
	Name string `json:"name" xml:"name" form:"name"`
}
func TestBodyParserJson(t *testing.T)  {
	app.Post("/register",func(c *fiber.Ctx) error {
		request := new(RegisterRequest)
		err:= c.BodyParser(request)
		if err != nil {
			return  err
		}
		return c.SendString("Register Success "+ request.Name)
	})

	body := strings.NewReader(`{"username":"hanafi","password":"adhi","name":"Hanafi Adhi"}`)

	request := httptest.NewRequest("POST","/register",body)
	request.Header.Set("Content-Type","application/json")

	response,err:=app.Test(request)
	assert.Nil(t,err)
	assert.Equal(t,200,response.StatusCode)

	byte,err:= io.ReadAll(response.Body)
	assert.Nil(t,err)
	assert.Equal(t,"Register Success Hanafi Adhi",string(byte))
}

func TestBodyParserForm(t *testing.T)  {
	app.Post("/register",func(c *fiber.Ctx) error {
		request := new(RegisterRequest)
		err:= c.BodyParser(request)
		if err != nil {
			return  err
		}
		return c.SendString("Register Success "+ request.Name)
	})

	body := strings.NewReader(`username=hanafi&password=adhi&name=Hanafi+Adhi`)

	request := httptest.NewRequest("POST","/register",body)
	request.Header.Set("Content-Type","application/x-www-form-urlencoded")

	response,err:=app.Test(request)
	assert.Nil(t,err)
	assert.Equal(t,200,response.StatusCode)

	byte,err:= io.ReadAll(response.Body)
	assert.Nil(t,err)
	assert.Equal(t,"Register Success Hanafi Adhi",string(byte))
}

func TestResponseJson(t *testing.T)  {
	app.Get("/user",func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"username":"hanafi",
			"password":"adhi",
		})
	})

	request := httptest.NewRequest("GET","/user",nil)
	request.Header.Set("Accept","application/json")
	response ,err:= app.Test(request)
	assert.Nil(t,err)
	assert.Equal(t,200,response.StatusCode)

	byte, err:= io.ReadAll(response.Body)
	assert.Nil(t,err)
	assert.Equal(t,`{"password":"adhi","username":"hanafi"}`,string(byte))
}

func TestDownloadFile(t *testing.T)  {
	app.Get("/download",func(c *fiber.Ctx) error {
		return c.Download("./source/dahlah.txt","contoh.txt")
	})

	request := httptest.NewRequest("GET","/download",nil)
	response,err:= app.Test(request)
	assert.Nil(t,err)
	assert.Equal(t,200,response.StatusCode)
	assert.Equal(t,"attachment; filename=\"contoh.txt\"",response.Header.Get("Content-Disposition"))

	byte, err:= io.ReadAll(response.Body)
	assert.Nil(t,err)
	assert.Equal(t,"bismillah\n",string(byte))
}

func  TestRoutingGroup(t *testing.T)  {
	heloWorld:= func (ctx *fiber.Ctx) error {
		return ctx.SendString("Hello World")
	}

	api := app.Group("/api")

	api.Get("/hello",heloWorld)
	api.Get("/world",heloWorld)

	web := app.Group("/web")
	web.Get("/hello",heloWorld)
	web.Get("/world",heloWorld)

	request := httptest.NewRequest("GET","/api/hello",nil)
	response, err:= app.Test(request)
	assert.Nil(t,err)
	assert.Equal(t,200,response.StatusCode)

	byte, err:= io.ReadAll(response.Body)
	assert.Nil(t,err)
	assert.Equal(t,"Hello World",string(byte))
}

func  TestErrorHandler(t *testing.T)  {
	app.Get("/error", func(c *fiber.Ctx) error {
		return errors.New("Ups")
	})

	request := httptest.NewRequest("GET","/error",nil)
	response , err:= app.Test(request)
	assert.Nil(t,err)
	assert.Equal(t,500,response.StatusCode)

	byte, err:= io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t,"Error Ups", string(byte))
}
