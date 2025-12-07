package web

// M is a shortcut for map[string]interface{}, similar to gin.H
// It provides a convenient way to create map responses without typing the full type signature.
//
// Example:
//
//	return web.Ok(web.M{
//	    "name": "John",
//	    "age": 30,
//	    "active": true,
//	})
type M map[string]interface{}
