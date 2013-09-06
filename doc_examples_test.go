package restful

func ExampleOPTIONSFilter() {
	Filter(OPTIONSFilter())
}
func ExampleContainer_OPTIONSFilter() {
	myContainer := new(Container)
	myContainer.Filter(myContainer.OPTIONSFilter)
}
