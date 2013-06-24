	ws.Route(ws.GET("/customers/{customer_id}/orders").To(getOrdersForCustomer).
				
		Doc("return the orders of a customer").
		Param(ws.PathParameter("customer_id", "identifier of the customer").DataType("string")).
		Writes( []Order{} )) // on the response