	ws := new(restful.WebService)
	ws.Path("/{scope}/systems").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_XML, restful.MIME_JSON).
		Param(ws.PathParameter("scope", "organization name to group system and connections").DataType("string"))
	idParam := ws.PathParameter("id", "identifier of the system").DataType("string")
	ws.Route(ws.GET("").To(s.getAll).
		Doc("list all known systems").
		Writes(model.System{})) // to the response ,TODO must be slice

	ws.Route(ws.GET("/{id}").To(s.get).
		Doc("get the system using its id").
		Param(idParam).
		Writes(model.System{})) // to the response

	ws.Route(ws.PUT("/{id}").To(s.put).
		Doc("create the system using its id").
		Param(idParam).
		Reads(model.System{})) // from the request

	ws.Route(ws.POST("").To(s.post).
		Doc("update the system using its id").
		Param(idParam).
		Reads(model.System{})) // from the request

	ws.Route(ws.DELETE("/{id}").To(s.delete).
		Doc("delete the system using its id").
		Param(idParam))

	restful.Add(ws)