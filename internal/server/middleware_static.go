// Copyright 2025 The Toodofun Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"embed"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Static(prefix string, fs *StaticFileSystem) gin.HandlerFunc {
	fileServer := http.FileServer(fs)
	if prefix != "" {
		fileServer = http.StripPrefix(prefix, fileServer)
	}
	return func(ctx *gin.Context) {
		logrus.Debugf("prefix: %s, path: %s", prefix, ctx.Request.URL.Path)
		if fs.Exists(prefix, ctx.Request.URL.Path) {
			fileServer.ServeHTTP(ctx.Writer, ctx.Request)
			ctx.Abort()
		} else {
			p := ctx.Request.URL.Path
			pathHasAPI := strings.Contains(p, "/api")
			if pathHasAPI {
				return
			} else {
				adminFile, err := fs.Open("index.html")
				if err != nil {
					fmt.Println("file not found", ctx.Request.URL.Path)
					return
				}
				defer adminFile.Close()
				// 把文件返回
				http.ServeContent(ctx.Writer, ctx.Request, "index.html", time.Now(), adminFile)
				ctx.Abort()
			}
		}
	}
}

type StaticFileSystem struct {
	fs   http.FileSystem
	root string
}

func (s *StaticFileSystem) Open(name string) (http.File, error) {
	openPath := path.Join(s.root, name)
	logrus.Debugf("openPath: %s", openPath)
	return s.fs.Open(openPath)
}

func (s *StaticFileSystem) Exists(prefix string, filepath string) bool {
	logrus.Debugf("filepath: %s, prefix: %s", filepath, prefix)
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		var name string
		if p == "" {
			name = path.Join(s.root, p, "index.html")
		} else {
			name = path.Join(s.root, p)
		}
		if _, err := s.fs.Open(name); err != nil {
			return false
		}
		return true
	}
	return false
}

func NewStaticFileSystem(data embed.FS, root string) *StaticFileSystem {
	return &StaticFileSystem{
		fs:   http.FS(data),
		root: root,
	}
}
