package spec

import (
	"github.com/cloudness-io/cloudness/app/services/githubapp"
	"github.com/cloudness-io/cloudness/app/services/gitpublic"
)

type Service struct {
	ghAppSvc     *githubapp.Service
	gitpublicSvc *gitpublic.Service
}

func NewService(
	ghAppSvc *githubapp.Service,
	gitpublicSvc *gitpublic.Service,
) *Service {
	return &Service{
		ghAppSvc:     ghAppSvc,
		gitpublicSvc: gitpublicSvc,
	}
}
