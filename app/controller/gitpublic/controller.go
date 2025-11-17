package gitpublic

import "github.com/cloudness-io/cloudness/app/services/gitpublic"

type Controller struct {
	gitPublicSvc *gitpublic.Service
}

func NewController(gitPublicSvc *gitpublic.Service) *Controller {
	return &Controller{
		gitPublicSvc: gitPublicSvc,
	}
}
