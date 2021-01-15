package pk

type Middleware func(service Service) Service

/*type loggingMiddleware struct {
	logger log.Logger
	next   PkService
}
*/
