package model

type GetProductsParams struct {
	Name *string
}

func StringToPointer(s string) *string {
	return &s
}
