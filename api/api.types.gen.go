// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package api

// Defines values for Currency.
const (
	CurrencyRUB Currency = "RUB"
)

// Defines values for OrderStatus.
const (
	OrderStatusCanceled OrderStatus = "canceled"

	OrderStatusCreated OrderStatus = "created"

	OrderStatusDone OrderStatus = "done"

	OrderStatusProcessing OrderStatus = "processing"
)

// Defines values for UpdatedOrderByBotStatus.
const (
	UpdatedOrderByBotStatusCanceled UpdatedOrderByBotStatus = "canceled"

	UpdatedOrderByBotStatusCreated UpdatedOrderByBotStatus = "created"

	UpdatedOrderByBotStatusDone UpdatedOrderByBotStatus = "done"

	UpdatedOrderByBotStatusProcessing UpdatedOrderByBotStatus = "processing"
)

// Category defines model for category.
type Category struct {
	Id        string `json:"id"`
	ShownName string `json:"shownName"`
}

// Currency defines model for currency.
type Currency string

// Связь с `category` по `categoryId`
type Dish struct {
	Calories        *int   `json:"calories,omitempty"`
	Carbohydrates   *int   `json:"carbohydrates,omitempty"`
	CategoryId      string `json:"categoryId"`
	Description     string `json:"description"`
	DishId          string `json:"dishId"`
	Fats            *int   `json:"fats,omitempty"`
	ImageUrl        string `json:"imageUrl"`
	PreviewImageUrl string `json:"previewImageUrl"`
	Price           struct {
		Amount   *int      `json:"amount,omitempty"`
		Currency *Currency `json:"currency,omitempty"`
	} `json:"price"`
	Proteins  *int   `json:"proteins,omitempty"`
	ShownName string `json:"shownName"`

	// Текст, который будет отображаться в случае отсутствия блюда в меню в данный момент
	UnavailableLabel *UnavailableLabel `json:"unavailableLabel,omitempty"`
	Weight           *int              `json:"weight,omitempty"`
}

// NewOrder defines model for newOrder.
type NewOrder struct {
	Positions []OrderPosition `json:"positions"`
}

// Order defines model for order.
type Order struct {
	Comment      string          `json:"comment"`
	OrderId      string          `json:"orderId"`
	Positions    []OrderPosition `json:"positions"`
	RestaurantId string          `json:"restaurantId"`
	Status       OrderStatus     `json:"status"`
}

// OrderStatus defines model for Order.Status.
type OrderStatus string

// OrderPosition defines model for orderPosition.
type OrderPosition struct {
	// Связь с `category` по `categoryId`
	Dish     Dish `json:"dish"`
	Quantity int  `json:"quantity"`
}

// Payment option
type PaymentOption struct {
	// Описание способа оплаты
	Description string `json:"description"`

	// Ссылка на изображение способа оплаты
	ImageUrl string `json:"imageUrl"`

	// Ссылка на переход на сайт платежной системы
	PaymentForwardUrl *string `json:"paymentForwardUrl,omitempty"`

	// Отображаемое название способа оплаты
	ShownName string `json:"shownName"`
}

// RestaurantStyle defines model for restaurantStyle.
type RestaurantStyle struct {
	BackgroundColor     string `json:"backgroundColor"`
	IconUrl             string `json:"iconUrl"`
	Id                  string `json:"id"`
	RestaurantShownName string `json:"restaurantShownName"`
}

// Текст, который будет отображаться в случае отсутствия блюда в меню в данный момент
type UnavailableLabel struct {
	ShownText string `json:"shownText"`
}

// UpdatedOrderByBot defines model for updatedOrderByBot.
type UpdatedOrderByBot struct {
	Status *UpdatedOrderByBotStatus `json:"status,omitempty"`
}

// UpdatedOrderByBotStatus defines model for UpdatedOrderByBot.Status.
type UpdatedOrderByBotStatus string

// CreateOrderResp defines model for createOrderResp.
type CreateOrderResp Order

// GetOrderResp defines model for getOrderResp.
type GetOrderResp Order

// GetOrdersResp defines model for getOrdersResp.
type GetOrdersResp []Order

// GetPaymentOptionsResp defines model for getPaymentOptionsResp.
type GetPaymentOptionsResp []PaymentOption

// GetRestaurantResp defines model for getRestaurantResp.
type GetRestaurantResp struct {
	Categories []Category `json:"categories"`

	// Отсортированный список блюд
	Dishes []Dish          `json:"dishes"`
	Style  RestaurantStyle `json:"style"`
}

// UpdateOrderByBotJSONBody defines parameters for UpdateOrderByBot.
type UpdateOrderByBotJSONBody UpdatedOrderByBot

// CreateOrderJSONBody defines parameters for CreateOrder.
type CreateOrderJSONBody NewOrder

// UpdateOrderByBotJSONRequestBody defines body for UpdateOrderByBot for application/json ContentType.
type UpdateOrderByBotJSONRequestBody UpdateOrderByBotJSONBody

// CreateOrderJSONRequestBody defines body for CreateOrder for application/json ContentType.
type CreateOrderJSONRequestBody CreateOrderJSONBody
