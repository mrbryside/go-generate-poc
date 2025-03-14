package dto

// this code below is generated by go-generate do not edit manually
// -----------------------------------------------------------------

type CreateProductsRequest struct {
	Type  string                     `json:"type" validate:"required"`
	Name  CreateProductsNameRequest  `json:"name"`
	Price CreateProductsPriceRequest `json:"price"`
} // @name CreateProductsRequest

type CreateProductsNameRequest struct {
	Type string `json:"type"`
} // @name CreateProductsNameRequest

type CreateProductsPriceRequest struct {
	Type int `json:"type"`
} // @name CreateProductsPriceRequest

type CreateProductsOKResponse struct {
	Type string `json:"type"`
} // @name CreateProductsOKResponse

type CreateProductsBadRequestResponse struct {
	Data CreateProductsBadRequestDataResponse `json:"data"`
} // @name CreateProductsBadRequestResponse

type CreateProductsBadRequestDataResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Manbank string `json:"manbank"`
} // @name CreateProductsBadRequestDataResponse
