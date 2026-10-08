package restful

import "testing"

var parameterSink *Parameter

func BenchmarkParameterCreation_KubernetesRoutes(b *testing.B) {
	ws := new(WebService)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for gv := 0; gv < 60; gv++ {
			parameterSink = ws.PathParameter("namespace", "object name and auth scope, such as for teams and projects")
			parameterSink = ws.PathParameter("name", "name of the resource")
			parameterSink = ws.QueryParameter("pretty", "If 'true', then the output is pretty printed.")
			parameterSink = ws.QueryParameter("fieldSelector", "A selector to restrict the list of returned objects by their fields.")
			parameterSink = ws.QueryParameter("labelSelector", "A selector to restrict the list of returned objects by their labels.")
			parameterSink = ws.QueryParameter("resourceVersion", "resourceVersion sets a constraint on what resource versions a request may be served from.")
			parameterSink = ws.QueryParameter("timeoutSeconds", "Timeout for the list/watch call.")
			parameterSink = ws.QueryParameter("limit", "limit is a maximum number of responses to return for a list call.")
			parameterSink = ws.QueryParameter("continue", "The continue option should be set when retrieving more results from the server.")
			parameterSink = ws.BodyParameter("body", "The resource payload").DataType("v1.ConfigMap")
			parameterSink = ws.HeaderParameter("Accept", "Accept header")
			parameterSink = ws.FormParameter("fieldManager", "Name of the manager used to track field ownership.")
			parameterSink = ws.MultiPartFormParameter("file", "Multipart upload")
		}
	}
}
