package main
import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "time"
)

func main() {
    // JWT yaratish uchun kalit so'z
    var mySigningKey = []byte("secret")

    // Token yaratish
    token := jwt.New(jwt.SigningMethodHS256)

    // Ma'lumotlar qo'shish
    claims := token.Claims.(jwt.MapClaims)
    claims["authorized"] = true
    claims["user"] = "John Doe"
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    // Tokenni imzolash
    tokenString, err := token.SignedString(mySigningKey)

    if err != nil {
        fmt.Println("Xatolik yuz berdi:", err)
    }

    // Tokenni ekranga chiqarish
    fmt.Println("Token:", tokenString)
}