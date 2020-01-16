package main

import (
	"Gin/VinIDRentCar/model"
	"fmt"
	"strings"

	myJwt "Gin/basicGinMiddlewareJwt/authentication.jwt"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
)

type testHeader struct {
	Rate   int    `header:"Rate"`
	Domain string `header:"Domain"`
	Tocken string `header:"Authorization"`
}

type Authen struct {
	Id    int    `json:"Id"`
	Types string `json:"Types"`
	Exp   int    `json:"exp"`
}

func main() {
	fmt.Println("hello world")
	router := gin.Default()
	router.Use(CORSMiddleware())
	private := router.Group("/private/api")
	private.Use(middlerWare4Private)
	private.Use(jwt.Auth(myJwt.GetSecretKey()))

	{
		private.GET("/ping", func(c *gin.Context) {
			privateRes := ""
			if value, existed := c.Get("ID"); existed {
				privateRes += "ID: " + fmt.Sprintf("%v", value)
			}

			if value, existed := c.Get("Type"); existed {
				privateRes += ", Type: " + fmt.Sprintf("%v", value)
			}

			if value, existed := c.Get("Exp"); existed {
				privateRes += ", Exp: " + fmt.Sprintf("%v", value)
			}

			c.JSON(200, model.MakeRespond(privateRes, 200, "private"))
		})
	}

	public := router.Group("/public/api")

	{
		public.GET("/ping", func(c *gin.Context) {
			successRes := make(map[string]interface{})
			successRes["tocken"] = myJwt.CreateTocken(myJwt.GetSecretKey(), 1, "customer")
			c.JSON(200, model.MakeRespond(successRes, 200, "public"))
		})
	}

	router.Run(":2234")
}

func middlerWare4Private(c *gin.Context) {
	header := &testHeader{}

	if err := c.BindHeader(&header); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		// split header
		payloadArr := strings.Split(header.Tocken, " ")
		payload := payloadArr[1]

		// split payload header
		payloadArr = strings.Split(header.Tocken, ".")
		payload = payloadArr[1]
		fmt.Println("payload = ", payload)
		fmt.Println(len(payload))

		// add "=" string to end base64 endcode
		if coutEqualStr := len(payload) % 4; coutEqualStr != 0 {
			fmt.Println(4 - coutEqualStr)
			for i := 0; i < (4 - coutEqualStr); i++ {
				payload += "="
			}
		}

		// decode base 64
		payloadByte, err := b64.StdEncoding.DecodeString(payload)

		if err != nil {
			fmt.Println("base64 decode error: ", err.Error())
			return
		}

		// json Unmarshal
		payloadStruct := Authen{}
		if err := json.Unmarshal(payloadByte, &payloadStruct); err != nil {
			return
		} else {
			fmt.Println(payloadStruct)
			c.Set("ID", payloadStruct.Id)
			c.Set("Type", payloadStruct.Types)
			c.Set("Exp", payloadStruct.Exp)
		}
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
