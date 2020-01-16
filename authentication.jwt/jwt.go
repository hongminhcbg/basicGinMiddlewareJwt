package jwt

import (
	"os"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"
)

func CreateTocken(secretKey string, id interface{}, types interface{}) string {
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))

	// Truyền dữ liệu vào phần Claim của token
	// Dữ liệu có kiểu map[string]interface{} mô phỏng một cấu trúc dạng JSON
	token.Claims = jwt_lib.MapClaims{
		"Id":    id,
		"Types": types,
		"exp":   time.Now().Add(time.Second * 15).Unix(),
	}

	// Tạo Signature cho token
	// Signature = HS256(Header, Claim, mysupersecretpassword)
	// Sử dụng mysupersecretpassword như một input đầu vào
	// để thuật toán HS256 tạo ra chuỗi signature
	if tokenString, err := token.SignedString([]byte(secretKey)); err != nil {
		return "create error"
	} else {
		return tokenString
	}

}

// GetSecretKey secret key is system enviroment
func GetSecretKey() string {
	return os.Getenv("RENT_CAR_SECRET_KEY")
}
