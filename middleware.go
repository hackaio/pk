package pk

type Middleware func(keeper PasswordKeeper) PasswordKeeper

func InitWithMiddlewares(
	hasher Hasher, store PasswordStore,
	tokenizer Tokenizer, es EncoderSigner, middlewares []Middleware) PasswordKeeper {

	var keeper = NewPasswordKeeper(hasher, store, tokenizer, es)

	for _, middleware := range middlewares {
		keeper = middleware(keeper)
	}

	return keeper
}

func WireInMiddlewares(keeper PasswordKeeper, middlewares []Middleware) PasswordKeeper {
	for _, middleware := range middlewares {
		keeper = middleware(keeper)
	}
	return keeper
}
