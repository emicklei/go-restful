//   GET /customers/{customer_id}/orders
//
func getOrdersForCustomer(req *restful.Request, resp *restful.Response) {

	orders, err := selectOrdersByCustomerId(req.PathParam("customer_id"))
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(orders)
}