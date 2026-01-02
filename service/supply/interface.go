package supply

import "context"

type Service interface {
	GetSupplyDetails(ctx context.Context, supplyId int) (*SupplyResponse, error)
}
