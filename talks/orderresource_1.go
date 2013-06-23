ws := new(restful.WebService)

ws.Route(ws.GET("/customers/{customer_id}/orders").To(getOrdersForCustomer))
ws.Route(ws.POST("/customers/{customer_id}/orders").To(addOrderForCustomer))
ws.Route(ws.DELETE("/customers/{customer_id}/orders/{order_id}").To(deleteOrderForCustomer))