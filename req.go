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



// RegisterRequest collects the request parameters for the Register method.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse collects the response parameters for the Register method.
type RegisterResponse struct {
	Err error `json:"err"`
}

// Failed implements Failer.
func (r RegisterResponse) Failed() error {
	return r.Err
}

// LoginRequest collects the request parameters for the Login method.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse collects the response parameters for the Login method.
type LoginResponse struct {
	Token string `json:"token"`
	Err   error  `json:"err"`
}

// Failed implements Failer.
func (r LoginResponse) Failed() error {
	return r.Err
}

// AddRequest collects the request parameters for the Add method.
type AddRequest struct {
	Token   string     `json:"token"`
	Account  Account `json:"account"`
}

// AddResponse collects the response parameters for the Add method.
type AddResponse struct {
	Err error `json:"err"`
}

// Failed implements Failer.
func (r AddResponse) Failed() error {
	return r.Err
}

// GetRequest collects the request parameters for the Get method.
type GetRequest struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// GetResponse collects the response parameters for the Get method.
type GetResponse struct {
	Account  Account `json:"account"`
	Err     error      `json:"err"`
}

// Failed implements Failer.
func (r GetResponse) Failed() error {
	return r.Err
}

// DeleteRequest collects the request parameters for the Delete method.
type DeleteRequest struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// DeleteResponse collects the response parameters for the Delete method.
type DeleteResponse struct {
	Err error `json:"err"`
}

// Failed implements Failer.
func (r DeleteResponse) Failed() error {
	return r.Err
}

// ListRequest collects the request parameters for the List method.
type ListRequest struct {
	Token string                 `json:"token"`
	Args  map[string]interface{} `json:"args"`
}

// ListResponse collects the response parameters for the List method.
type ListResponse struct {
	Accounts [] Account `json:"accounts"`
	Err      error        `json:"err"`
}

// Failed implements Failer.
func (r ListResponse) Failed() error {
	return r.Err
}

// UpdateRequest collects the request parameters for the Update method.
type UpdateRequest struct {
	Token     Account `json:"token"`
	Name      Account `json:"name"`
	Username  Account `json:"username"`
	Account   Account `json:"account"`
}

// UpdateResponse collects the response parameters for the Update method.
type UpdateResponse struct {
	Acc  Account `json:"acc"`
	Err error      `json:"err"`
}

// Failed implements Failer.
func (r UpdateResponse) Failed() error {
	return r.Err
}

// AddAllRequest collects the request parameters for the AddAll method.
type AddAllRequest struct {
	Token    string       `json:"token"`
	Accounts [] Account `json:"accounts"`
}

// AddAllResponse collects the response parameters for the AddAll method.
type AddAllResponse struct {
	Err error `json:"err"`
}

// Failed implements Failer.
func (r AddAllResponse) Failed() error {
	return r.Err
}

// DeleteAllRequest collects the request parameters for the DeleteAll method.
type DeleteAllRequest struct {
	Token string                 `json:"token"`
	Args  map[string]interface{} `json:"args"`
}

// DeleteAllResponse collects the response parameters for the DeleteAll method.
type DeleteAllResponse struct {
	Err error `json:"err"`
}

// Failed implements Failer.
func (r DeleteAllResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}
