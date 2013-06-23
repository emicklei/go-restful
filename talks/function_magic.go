func getOrdersForCustomer(customer_id int) []Orders {

	return dao.selectOrdersByCustomerId(customer_id)
}