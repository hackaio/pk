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

type Signer interface {
	Sign(string) ([]byte, []byte, error)

	Verify(password string, dbDigest []byte, dbSignature []byte) (err error)
}

type Encoder interface {
	Encode(password string) ([]byte, error)
	Decode(encoded []byte) (string, error)
}

type EncoderSigner interface {
	Encoder
	Signer
}
