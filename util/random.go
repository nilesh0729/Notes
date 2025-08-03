package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const Alphabets = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ1234567890"

const PasswordCharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"

var domains = []string{
	"gmail.com",
	"outlook.com",
	"outlook.com",
	"example.com",
	"protonmail.com",
}

// Random Integers
func RandomInts(min, max int64) int64 {
	return (min + rand.Int63n(max-min+1))
}

// Random Strings
func RandomString(n int) string {
	var sb strings.Builder
	k := len(Alphabets)

	for i := 0; i < n; i++ {
		c := Alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// NOW WE USE RANDOM INTEGERS AND RANDOM STRINGS TO GENERATE RANDOM USERS, EMAIL, PASSWORD

func RandomUser() string {
	return RandomString(6)
}

func RandomPassword(length int) string {
	Password := make([]byte, length)
	for i := range Password {
		Password[i] = PasswordCharset[rand.Intn(len(PasswordCharset))]
	}
	return string(Password)
}

func RandomEmail() string {
	username := RandomString(8) // adjust length as needed
	domain := domains[rand.Intn(len(domains))]
	return fmt.Sprintf("%s@%s", username, domain)
}
