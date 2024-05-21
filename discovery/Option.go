package discovery

import "context"

type Options func(*ProberService)

func addTimeout(Option Options) func(service *ProberService) {
	return func(Prob *ProberService) {
		Prob.Ctx = context.Background()
	}
}
