package services

import (
	"tasks/errorcodes"
	"tasks/store"
)

type ListAccessType = store.ListAccessType

type ListAccess struct {
	UserID string
	ListID uint64
	Access []ListAccessType
}

type GetListAccessParams struct {
	UserID string
	ListID uint64
}

type AddListAccessParams struct {
	UserID string
	ListID uint64
	Access []ListAccessType
}

type RemoveListAccessParams struct {
	UserID *string
	ListID *uint64
	Access *[]ListAccessType
}

func GetListAccess(ctx ServiceContext, args GetListAccessParams) (*ListAccess, *ServiceError) {

	var listAccess *store.ListAccess
	var serr *store.StoreError
	listAccess, serr = ST.GetListAccess(ctx.Ctx, store.GetListAccessParams{
		UserID: args.UserID,
		ListID: args.ListID,
	})

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound,
			errorcodes.ListNotFound:
			return nil, &ServiceError{
				Code: serr.Code,
			}

		default:
			return nil, &ServiceError{
				Code: errorcodes.Internal,
			}
		}
	}

	return &ListAccess{
		UserID: listAccess.UserID,
		ListID: listAccess.ListID,
		Access: listAccess.Access,
	}, nil
}

func AddListAccess(ctx ServiceContext, args AddListAccessParams) *ServiceError {

	// var listAccess *store.ListAccess
	var serr *store.StoreError = ST.AddListAccess(ctx.Ctx, store.AddListAccessParams{
		UserID: args.UserID,
		ListID: args.ListID,
		Access: args.Access,
	})

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound,
			errorcodes.ListNotFound:
			return &ServiceError{
				Code: serr.Code,
			}

		default:
			return &ServiceError{
				Code: errorcodes.Internal,
			}
		}
	}

	return nil
}

func RemoveListAccess(ctx ServiceContext, args RemoveListAccessParams) *ServiceError {

	var serr *store.StoreError = ST.RemoveListAccess(ctx.Ctx, store.RemoveListAccessParams{
		UserID: args.UserID,
		ListID: args.ListID,
		Access: args.Access,
	})

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound,
			errorcodes.ListNotFound:
			return &ServiceError{
				Code: serr.Code,
			}

		default:
			return &ServiceError{
				Code: errorcodes.Internal,
			}
		}
	}

	return nil
}
