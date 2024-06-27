package service

import (
	"math/rand/v2"
	"regexp"
)

func randFromRange(min, max int) int {
	return rand.IntN(max-min+1) + min
}

func GenerateShortenURL(length int) string {
	var result = make([]byte, length)
	var ranges = [3][2]int{
		{48, 57},
		{65, 90},
		{97, 122},
	}
	for i := 0; i < length; i++ {
		var r = randFromRange(0, 2)
		result[i] = byte(randFromRange(ranges[r][0], ranges[r][1]))
	}
	return string(result)
}

func IsURLValid(url string) bool {
	_, err := regexp.MatchString("^https?://(.*)\\.(.*)$", url)
	return err == nil
}

func IsIDValid(id string) bool {
	_, err := regexp.MatchString("^[0-9a-zA-Z]{6}$", id)
	return err == nil
}
