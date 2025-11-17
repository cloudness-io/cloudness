package store

import "errors"

var (
	ErrResourceNotFound           = errors.New("resource not found")
	ErrDuplicate                  = errors.New("resource is a duplicate")
	ErrVersionConflict            = errors.New("resource version conflict")
	ErrPathTooLong                = errors.New("the path is too long")
	ErrPrimaryPathAlreadyExists   = errors.New("primary path already exists for resource")
	ErrPrimaryPathRequired        = errors.New("path has to be primary")
	ErrAliasPathRequired          = errors.New("path has to be an alias")
	ErrPrimaryPathCantBeDeleted   = errors.New("primary path can't be deleted")
	ErrNoChangeInRequestedMove    = errors.New("the requested move doesn't change anything")
	ErrIllegalMoveCyclicHierarchy = errors.New("the requested move is not permitted as it would cause a " +
		"cyclic depdency")
	ErrSpaceWithChildsCantBeDeleted = errors.New("the space can't be deleted as it still contains " +
		"spaces or repos")
	ErrPreConditionFailed = errors.New("precondition failed")
)
