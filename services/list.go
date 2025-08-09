package services

import (
	"tasks/errorcodes"
	"tasks/store"
)

type List struct {
	ID   uint64
	Name string
}

type CreateListParams struct {
	Name string
}

type UpdateListParams struct {
	Name *string
}

func CreateList(ctx ServiceContext, args CreateListParams) (uint64, *ServiceError) {

	var listID uint64
	var serr *store.StoreError
	listID, serr = ST.CreateList(ctx.Ctx, store.CreateListParams{
		Name: args.Name,
	})

	if serr != nil {

		return 0, &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	return listID, nil
}

func GetList(ctx ServiceContext, id uint64) (*List, *ServiceError) {

	var list *store.List
	var serr *store.StoreError
	list, serr = ST.GetList(ctx.Ctx, id)

	if serr != nil {

		return nil, &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	return &List{
		ID:   list.ID,
		Name: list.Name,
	}, nil
}

func UpdateList(ctx ServiceContext, id uint64, args UpdateListParams) *ServiceError {

	var serr *store.StoreError = ST.UpdateList(ctx.Ctx, id, store.UpdateListParams{
		Name: args.Name,
	})

	if serr != nil {

		return &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	return nil
}

func DeleteList(ctx ServiceContext, id uint64) *ServiceError {

	var serr *store.StoreError = ST.DeleteList(ctx.Ctx, id)
	if serr != nil {

		return &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	return nil
}
