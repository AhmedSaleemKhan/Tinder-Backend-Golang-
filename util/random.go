package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomName generates a random owner name
func RandomName() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomPhoneNumber generates a random phone number
func RandomPhoneNumber() string {
	return fmt.Sprintf("+%v", RandomInt(9999999999, 100000000000))
}

// RandomUID generates a random uid
func RandomUID(n int) string {
	chars := "abcdefghijklmnopqrstuvwxy" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"
	random := ""
	for i := 0; i < n; i++ {
		random += string([]rune(chars)[rand.Intn(len(chars))])
	}
	return random
}

// Pick random gender from slice of gerders
func PickRandomGender() string {
	gerders := []string{"male", "female"}
	randomIndex := rand.Intn(len(gerders))
	return gerders[randomIndex]
}

// Give slice of ethnicity
func PickRandomEthnicity() []string {
	ethnicities := []string{"American Indian", "Black/African Descent", "East Asian", "Hispanic/Latino", "Middle Easter", "Pacific Islander", "South Asian", "Southeast Asian", "White/Caucasian"}
	return ethnicities
}

// Pick random type from slice
func PickRandomType() string {
	types := []string{"chat", "audio", "video", "payin", "payout"}
	randomIndex := rand.Intn(len(types))
	return types[randomIndex]
}

// RandomBool generate random boolean
func RandomBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}
