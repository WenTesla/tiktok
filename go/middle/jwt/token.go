package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"tiktok/go/config"
	"tiktok/go/model"
	"time"
)

// 提取密钥
var jwtSecret = []byte(config.Secret)

// 签名token
func SignToken(user model.User) string {
	// dev
	Id := user.Id
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// iss (issuer): 签发人
		"iss": "tiktok",
		// sub (subject): 主题
		"sub": "用户Id",
		// Id
		"id": Id,
		// nbf (Not Before): 生效时间
		"nbf": time.Now().Unix(),
		// exp (Expiration): 过期时间，由 JWT_EXPIRE_HOURS 控制，golang-jwt 会在解析时自动校验
		"exp": time.Now().Add(time.Duration(config.AppConfig.JWT.Expire) * time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtSecret)
	// 如果签名失败
	if err != nil {
		return "签名失败"
	}
	fmt.Println(tokenString, err)
	return tokenString
}

//  解析token 返回用户Id

func ParseToken(tokenString string) (interface{}, error) {

	// sample token string taken from the New example
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//fmt.Println(claims["foo"], claims["nbf"])
		return claims["id"], nil
	} else {
		fmt.Println(err)
		return nil, err
	}
	//return true, err
}

// 返回token中的签发日期

func ParseTokenTime(tokenString string) (any, error) {
	// sample token string taken from the New example
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//fmt.Println(claims["foo"], claims["nbf"])
		return claims["nbf"], nil
	} else {
		fmt.Println(err)
		return nil, err
	}
}
