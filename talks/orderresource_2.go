func getOrdersForCustomer(req *restful.Request, resp *restful.Response) {

	customer_id := req.PathParam("customer_id")
	orders := dao.selectOrdersByCustomerId(customer_id)
	resp.WriteEntity(orders)
}