/*
 * Copyright © 2021 PIUS ALFRED me.pius1102@gmail.com
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pk

type Middleware func(keeper PasswordKeeper) PasswordKeeper

func New(hasher Hasher, store PasswordStore, tokenizer Tokenizer,
	es EncoderSigner, middlewares []Middleware) PasswordKeeper {

	var keeper = NewPasswordKeeper(hasher, store, tokenizer, es)

	for _, middleware := range middlewares {
		keeper = middleware(keeper)
	}

	return keeper
}

func AddMiddlewares(keeper PasswordKeeper, middlewares []Middleware) PasswordKeeper {
	for _, middleware := range middlewares {
		keeper = middleware(keeper)
	}
	return keeper
}
