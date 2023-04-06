package category

import (
	"context"
	"net/http"

	"github.com/Invan2/invan_kanban_service/models"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func (c *CategoryService) Create(ctx context.Context, event cloudevents.Event) models.Response {

	return models.Response{

		StatusCode: http.StatusOK,
	}
}
