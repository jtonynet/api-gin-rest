package middlewares

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Fortemente baseado no artigo:
// https://articles.wesionary.team/make-pagination-easy-in-golang-using-pagination-middleware-and-gorm-scope-a5f6eb3bebaa
// TODO Aplicar a service DB e o gorm como sugerido no artigo

func PaginateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		per_page, err := strconv.ParseInt(c.Query("per_page"), 10, 0)
		if err != nil {
			per_page = 5
		}

		page, err := strconv.ParseInt(c.Query("page"), 10, 0)
		if err != nil {
			page = 0
		}

		c.Set("Limit", per_page)
		c.Set("Page", page)

		c.Next()
	}
}
