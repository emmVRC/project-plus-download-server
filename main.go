package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"os"
	"strings"
)

var DllPath string
var DllData []byte
var DllBase64Data string
var DllHash string

var ShouldLoad bool

func init() {
	shouldLoad := flag.Bool("should-load", true, "")
	dllPath := flag.String("dll-path", "emmvrc.dll", "")
	flag.Parse()

	ShouldLoad = *shouldLoad
	DllPath = *dllPath

	open, err := os.Open(DllPath)
	if err != nil {
		log.Fatalf("failed to open dll: %s", err)
	}

	DllData, err = io.ReadAll(open)
	if err != nil {
		log.Fatalf("failed to read dll: %s", err)
	}

	DllHash = fmt.Sprintf("%x", sha256.Sum256(DllData))
	DllBase64Data = base64.StdEncoding.EncodeToString(DllData)
}

func main() {
	app := fiber.New()

	app.Get("/:hash", downloadMod)

	log.Fatal(app.Listen(":3000"))
}

func shouldLoad(c *fiber.Ctx) error {
	c.Set("surrogate-key", "mod-resource")
	
	if ShouldLoad {
		return c.SendStatus(200)
	}

	return c.SendStatus(403)
}

func downloadMod(c *fiber.Ctx) error {
	if strings.ToLower(c.Params("hash")) == DllHash {
		return c.SendStatus(204)
	}

	c.Set("content-type", "application/octet-stream")
	c.Set("surrogate-key", "mod-resource")

	return c.SendString(DllBase64Data)
}
