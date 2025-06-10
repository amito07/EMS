package randomfunction

import (
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

   func GetStudentId(length int) string {
       rand.Seed(time.Now().UnixNano()) // Seed with current time
       const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
       result := make([]byte, length)
       for i := 0; i < length; i++ {
           result[i] = charset[rand.Intn(len(charset))]
       }
       return string(result)
   }

   var secretKey = []byte("secret-key")
   func GenerateJwtToken(email, studentId string) (string, error) {
	    token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "email": email,
		"student_id": studentId,
        "exp": time.Now().Add(time.Hour * 24).Unix(), 
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
    return "", err
    }

 return tokenString, nil


   }