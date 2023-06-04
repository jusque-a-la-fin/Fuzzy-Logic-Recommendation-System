package adapters

type Context interface {
	Bind(i interface{}) error
	BindJSON(i interface{}) error
	HTML(code int, fileName string, i interface{})
	JSON(code int, i interface{})
	PostForm(key string) string
	Param(key string) string
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	Cookie(name string) (string, error)
}
