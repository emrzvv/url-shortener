package service

import (
	storage "github.com/emrzvv/url-shortener/internal/app/db"
	"math/rand/v2"
	"regexp"
)

type URLShortenerService interface {
	GenerateShortURL(url string) (string, error)
	GetOriginURLByID(id string) (string, error)
}

type URLShortenerServiceImpl struct {
	db storage.Storage
}

func NewURLShortenerService(db storage.Storage) URLShortenerService {
	return &URLShortenerServiceImpl{db: db}
}

func (s *URLShortenerServiceImpl) GenerateShortURL(url string) (string, error) {
	if isValid, _ := regexp.MatchString("^https?://(.*)\\.(.*)$", url); !isValid {
		return "", &InvalidURLError{value: url}
	}
	shorten := generate(6)
	for _, ok := s.db.Get(shorten); ok; shorten = generate(6) { // TODO: придумать другую стратегию с однозначной генерацией
	}
	s.db.Set(shorten, url)
	return shorten, nil
}

func (s *URLShortenerServiceImpl) GetOriginURLByID(id string) (string, error) {
	if isValid, _ := regexp.MatchString("^[0-9a-zA-Z]{6}$", id); !isValid {
		return "", &InvalidIDError{value: id}
	}
	value, ok := s.db.Get(id)
	if !ok {
		return "", &InvalidIDError{value: id}
	}
	return value, nil
}

func randFromRange(min, max int) int {
	return rand.IntN(max-min+1) + min
}

func generate(length int) string {
	var result = make([]byte, length)
	var ranges = [3][2]int{
		{48, 57},  // 0-9
		{65, 90},  // A-Z
		{97, 122}, // a-z
	}
	for i := 0; i < length; i++ {
		var r = randFromRange(0, 2)
		result[i] = byte(randFromRange(ranges[r][0], ranges[r][1]))
	}
	return string(result)
}
