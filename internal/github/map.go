package github

import "github.com/samokw/foundre-webhook/internal/domain"

func MapPRAction(a string) (domain.PreviewAction, bool) {
	switch a {
	case "opened", "reopened":
		return domain.Create, true
	case "synchronize":
		return domain.Update, true
	case "closed":
		return domain.Delete, true
	default:
		return "", false
	}
}
