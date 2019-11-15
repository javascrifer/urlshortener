package shortener

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFound = errors.New("redirect not found")
	ErrRedirectInvalid  = errors.New("redirect is invalid")
)

type redirectService struct {
	repository RedirectRepository
}

func (r redirectService) Find(code string) (*Redirect, error) {
	return r.repository.Find(code)
}

func (r redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.repository.Store(redirect)
}

func NewRedirectService(r RedirectRepository) RedirectService {
	return &redirectService{repository: r}
}
