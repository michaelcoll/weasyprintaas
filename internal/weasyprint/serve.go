/*
 * Copyright (c) 2022 Michaël COLL.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package weasyprint

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"

	"github.com/michaelcoll/weasyprintaas/internal/weasyprint/model"
)

const apiPort = ":8080"

type Controller struct {
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Serve() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.POST("/", c.convert)

	// Listen and serve on 0.0.0.0:<portNumber>
	fmt.Printf("Listening API on 0.0.0.0%s\n", color.GreenString(apiPort))
	err := router.Run(apiPort)
	if err != nil {
		log.Fatalf("Error starting server : %v", err)
	}
}

func (c *Controller) convert(ctx *gin.Context) {
	var request = &model.Request{}

	err := ctx.BindJSON(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Body is malformed : %v", err)})
	}

	fmt.Printf("Converting url %s\n", color.CyanString(request.Url))
	outFilename := fmt.Sprintf("/tmp/out-%d.pdf", rand.Intn(1000000))

	stderr := new(strings.Builder)
	cmd := exec.CommandContext(ctx.Request.Context(), "weasyprint", request.Url, outFilename)
	cmd.Stderr = stderr

	err = cmd.Run()
	defer removeFile(outFilename)

	if err != nil {
		fmt.Printf("Can't execute command : \n%v\n", color.RedString(stderr.String()))
	}

	ctx.File(outFilename)
}

func removeFile(name string) {
	err := os.Remove(name)
	if err != nil {
		fmt.Printf("Can't remove output file : %v\n", err)
	}
}
