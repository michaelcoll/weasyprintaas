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
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"

	"github.com/michaelcoll/weasyprintaas/internal/weasyprint/model"
)

type Controller struct {
	multithreading bool
	mu             sync.Mutex
}

func New(multithreading bool) *Controller {
	return &Controller{multithreading: multithreading}
}

func (c *Controller) Serve(port uint16) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.POST("/", c.convert)
	router.GET("/health", c.health)

	// Listen and serve on 0.0.0.0:<portNumber>
	fmt.Printf("Listening API on 0.0.0.0:%s\n", color.GreenString(strconv.Itoa(int(port))))
	err := router.Run(fmt.Sprintf(":%d", port))
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

	if !c.multithreading {
		c.mu.Lock()
	}
	err = cmd.Run()
	if !c.multithreading {
		c.mu.Unlock()
	}
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

func (c *Controller) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}
