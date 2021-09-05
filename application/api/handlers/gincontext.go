package handlers

type GinContext interface {
	JSON(code int, obj interface{})
}
