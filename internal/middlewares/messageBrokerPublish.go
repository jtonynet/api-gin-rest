package middlewares

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/api-gin-rest/config"
	"github.com/tidwall/gjson"

	"github.com/jtonynet/api-gin-rest/internal/interfaces"
)

// Tratamento adequedo de Middlewares para obter "Separation of Concerns"
// https://gin-gonic.com/docs/examples/custom-middleware/

func MessageBrokerPublishPostRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := c.MustGet("cfg").(config.API)

		c.Next()

		if cfg.FeatureFlags.PostAlunoAsMessageEnabled {
			messageBroker := c.MustGet("messageBroker").(interfaces.MessageBroker)

			model, modelExists := c.Get("model")
			if modelExists {

				modelJSON, err := json.Marshal(model)
				if err != nil {
					slog.Error("middlewares:MessageBrokerPublishPostRequest:json.Marshal error: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Erro ao serializar request body",
					})
					c.Abort()
				}

				UUIDField := gjson.Get(string(modelJSON), "uuid")
				if !UUIDField.Exists() {
					slog.Error("middlewares:MessageBrokerPublishPostRequest:UUIDField.Exists error: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Erro ao obter UUID",
					})
					c.Abort()
				}

				err = messageBroker.Publish(string(modelJSON))
				if err != nil {
					slog.Error("controllers:CriaNovoAluno:messageBroker.Publish error: %v", err)
				}
			}
			c.Abort()
		}
	}
}
