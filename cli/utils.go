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

package cli

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	prettyjson "github.com/hokaccha/go-prettyjson"
)

func logJSON(iList ...interface{}) {
	for _, i := range iList {
		m, err := json.Marshal(i)
		if err != nil {
			logError(err)
			return
		}

		pj, err := prettyjson.Format(m)
		if err != nil {
			logError(err)
			return
		}

		fmt.Printf("\n%s\n\n", string(pj))
	}
}

func logUsage(u string) {
	fmt.Printf(color.YellowString("\nusage: %s\n\n"), u)
}

func logMessage(key, msg string) {
	fmt.Printf(color.YellowString("\n%s: %s\n\n"),key, msg)
}

func logError(err error) {
	boldRed := color.New(color.FgRed, color.Bold)
	boldRed.Print("\nerror: ")

	fmt.Printf("%s\n\n", color.RedString(err.Error()))
}

func logOK() {
	fmt.Printf("\n%s\n\n", color.BlueString("ok"))
}

func logCreated(e string) {

	fmt.Printf(color.BlueString("\ncreated: %s\n\n"), e)

}

