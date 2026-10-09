package service

import (
	"crypto/rand"
	"errors"
	"math/big"
	"net/url"

	"github.com/agatml/url-shortener/internal/repository"
)

const (
	alphabet   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	codeLength = 7
	maxRetries = 5
)

var ErrInvalidURL = errors.New("url inválida")

// Repository define o que o service precisa. Quem implementa é a camada de repository.
type Repository interface {
	Save(code, url string) error
	Find(code string) (string, error)
}

type Shortener struct {
	repo Repository
}

func NewShortener(repo Repository) *Shortener {
	return &Shortener{repo: repo}
}

// Shorten valida a URL, gera um código único e salva.
func (s *Shortener) Shorten(rawURL string) (string, error) {
	if err := validateURL(rawURL); err != nil {
		return "", err
	}

	for i := 0; i < maxRetries; i++ {
		code, err := generateCode()
		if err != nil {
			return "", err
		}

		err = s.repo.Save(code, rawURL)
		if err == nil {
			return code, nil
		}
		if !errors.Is(err, repository.ErrAlreadyExists) {
			return "", err
		}
		// colisão de código: tenta de novo com outro
	}

	return "", errors.New("não foi possível gerar um código único")
}

// Resolve devolve a URL original a partir do código.
func (s *Shortener) Resolve(code string) (string, error) {
	return s.repo.Find(code)
}

func validateURL(raw string) error {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return ErrInvalidURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrInvalidURL
	}
	if u.Host == "" {
		return ErrInvalidURL
	}
	return nil
}

func generateCode() (string, error) {
	b := make([]byte, codeLength)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		b[i] = alphabet[n.Int64()]
	}
	return string(b), nil
}
