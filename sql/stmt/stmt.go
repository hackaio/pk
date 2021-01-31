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

package stmt

//name
//username
//email
//password
//created

const (
	ADD_OWNER = "INSERT INTO masters (name, username,email,password, created) VALUES ($1, $2, $3, $4, $5);"
	GET_OWNER = "SELECT * FROM masters WHERE name = $1 AND username = $2;"
	ADD = "INSERT INTO accounts (name, username,email,hash,encoded,digest,signature,created) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"
	GET = "SELECT * FROM accounts WHERE name = $1 AND username = $2;"
	LIST = "SELECT * FROM accounts;"
)
