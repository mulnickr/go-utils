package serve

type Middleware func(Handler) HandlerFunc
