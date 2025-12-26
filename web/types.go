package web

// M 是 map[string]interface{} 的简写，类似于 gin.H。
// 提供了一种方便的方式来创建 map 响应，而无需输入完整的类型签名。
//
// 示例：
//
//	return web.Ok(web.M{
//	    "name": "John",
//	    "age": 30,
//	    "active": true,
//	})
type M map[string]interface{}
