package preview

import (
	"fmt"
	"log"

	"github.com/samokw/foundre-webhook/internal/domain"
)

func Handle(req domain.PreviewRequest) error {
	log.Printf("preview action=%s key=%s safe=%s sha=%s",
		req.Action, req.Key(), req.SafeName(), req.SHA)

	switch req.Action {
	case domain.Create, domain.Update:
		return EnsureNamespace(req.SafeName())
	case domain.Delete:
		return DeleteNamespace(req.SafeName())
	default:
		return fmt.Errorf("unknown preview action: %q", req.Action)
	}
}
