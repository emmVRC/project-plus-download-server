package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"os"
	"strings"
)

// idc this is bad
var BlacklistedWorlds = []string{
	"wrld_7d19d792-e57c-438a-a56e-d4f1e2daa1d5",
	"wrld_4473217e-823d-49d1-909a-ca3424f7795f",
	"wrld_860bdea8-0499-4c16-8170-5683a3a3ab8a",
	"wrld_02cb84c8-446a-484b-b52b-2c1777e55f8d",
	"wrld_588c1ab1-1ddc-4378-9533-a3aab06c0ee2",
	"wrld_41c22a73-30a4-438f-914b-3d3627f0fe8a",
	"wrld_9b5bd9e0-fdfd-453d-96e1-5292f6790874",
	"wrld_856e5ef2-69f8-4e9e-863b-51cba9c3fd20",
	"wrld_c8e3ba0c-592f-499c-a06f-9ddf997e1b24",
	"wrld_4a65ba53-d8df-40a7-b67d-30c63bff0e95",
	"wrld_7487d91a-3ef4-44c6-ad6d-9bdc7dee5efd",
	"wrld_7e10376a-29b6-43af-ac5d-6eb72732e90c",
	"wrld_f1242b62-023a-4ee1-8e0b-2946274e194a",
	"wrld_8194a6d5-3aa9-425a-8282-86e0af69690b",
	"wrld_3be14f03-125f-4e4f-acf2-0a565738dd95",
	"wrld_d2cb61e6-b941-45fa-9d47-aa7b8713e309",
	"wrld_0d8a98c9-891a-4665-a682-97e6998db7af",
	"wrld_a5bcb226-b040-4c56-adff-b1dc2f0189c8",
	"wrld_de3400c4-491f-4506-b8e8-71f06b4660bb",
	"wrld_3ff0ac0f-18e5-4c3d-80ad-c0a11dea0ecb",
	"wrld_cff7a309-4cd0-4ce1-a5aa-fb4992cd1f09",
	"wrld_e5c30b56-efa8-42d5-a8d4-a2cca2bf3403",
	"wrld_f35460ce-5d6d-40ed-8828-41e590deee58",
	"wrld_aa5ef7a9-2a9a-4b9f-a1ef-daf3565c5628",
	"wrld_9c72e56b-d2b0-4c9b-b816-07a857f6ae4e",
	"wrld_953eff93-20c9-457b-8ef0-1ac2130d5b8a",
	"wrld_c6417518-6469-4296-a5fb-4cc51c5c8a58",
	"wrld_b155ff30-0bbf-486c-afc5-e7a05145387b",
	"wrld_1b0090ad-9c39-4b6a-bfdf-9c53bfd5988b",
	"wrld_6fc49682-46ee-4ffb-a869-75c17c70ad1d",
	"wrld_0a6542e0-2bab-435a-9d53-f8f1f5aa4b4d",
	"wrld_11f3457f-da08-4817-99d2-416b1c9db6f3",
	"wrld_d319c58a-dcec-47de-b5fc-21200116462c",
}

var WhitelistedWorlds = []string{
	"wrld_f8474f94-684c-47be-84fd-a7ce90a1b5ce",
	"wrld_fae3fa95-bc18-46f0-af57-f0c97c0ca90a",
	"wrld_6a5f5dfd-6146-4d18-b2ac-8bc2193c365c",
	"wrld_5d1c606f-d940-4725-a05f-ff2dfe9db942",
	"wrld_2e647e8d-8430-4b40-8bd7-914f1fa817b2",
	"wrld_1b482eca-bede-4de8-88a8-bbb6ca7e24cd",
	"wrld_1256060a-6fe6-4b5f-b6f6-619c99ddeb1c",
	"wrld_831d080e-5577-4234-b436-0c1dec75b145",
	"wrld_5c6565f0-50b8-49ef-8a12-0fdf4d1f8c2a",
	"wrld_6d641df0-3d65-4e7a-bce0-b97b87d50d2c",
	"wrld_f106374e-1744-4708-a980-fb91c9afd6d0",
	"wrld_ed55ea0a-851c-4eda-b292-62d648a717de",
	"wrld_b8f44a58-4355-48ba-8f6c-b1434a19b467",
	"wrld_d6e55ea7-5fdc-42ac-93cd-f285fc22cd04",
	"wrld_ad969e26-f327-450a-9a04-8c0b1db35382",
	"wrld_ea2aa8b2-d78a-41ee-af64-cfb2598d440a",
	"wrld_2147b034-c6e9-4919-b74a-a0c81b596d95",
	"wrld_84b4fe0a-efb3-4f9d-b8e8-34f11d2c4a54",
	"wrld_042dffc5-9a6c-4ae4-a655-f64516a64bcb",
	"wrld_c183ed0c-87d3-4124-97ba-4898df5d3daf",
	"wrld_6e77e43e-3fdd-447f-9c49-f13351febc0a",
	"wrld_e68dd0cb-c176-4f18-9ea9-e6a189747ffd",
	"wrld_ffda6963-447c-44a6-983c-b0c014a3eced",
	"wrld_8c9b5e18-2711-4dfc-b896-f47e75fed7b9",
	"wrld_1c63aa36-9180-419e-b8be-e27e219a8139",
	"wrld_848eaf0d-3c34-4253-bbfe-ac6307d2ca12",
	"wrld_bd179c6a-be43-4cb2-b1ee-0f250b5cb77a",
	"wrld_e06ef2d3-34b1-4f8d-835a-4ca58c9bd1c1",
	"wrld_899f0635-9335-4c0d-a7a5-04b4117bc0dc",
	"wrld_39403e77-5132-4c76-b3d5-6b94027b0b15",
	"wrld_9d428ade-60f1-4c34-898a-c6dd8aa13c53",
	"wrld_c5303565-4b52-45b1-b2e7-5da1f99de2a6",
	"wrld_9190a89b-c0c8-4ad4-9e35-b50549e686cd",
	"wrld_bd8aad50-b6ce-40e8-8b04-3221411e6bdd",
	"wrld_b8d07d37-a02a-429b-a508-46bf88497106",
	"wrld_e075192e-99a3-41c5-ad4e-59c74ec4f5cb",
	"wrld_a3d06073-d1a0-4640-b719-8bfce671df0d",
	"wrld_31329115-4296-4d89-b47d-09dd96577486",
	"wrld_62a0710e-26e4-4339-b0ae-a33a5ed5148f",
	"wrld_47237d61-4fd9-4066-bf91-41c482a67f0a",
	"wrld_bcf249ec-fcef-4b9f-9037-55521f0f1ab7",
	"wrld_306779d1-d027-4f2f-bbd7-68bf97bf045f",
	"wrld_12b8d9d8-d719-4925-bb07-5da0888f198d",
	"wrld_65ce4a4b-f66f-4844-adfe-072663394514",
	"wrld_ae2eaf16-e704-4414-97c3-5fddd022c1f9",
	"wrld_953eff93-20c9-457b-8ef0-1ac2130d5b8a",
	"wrld_e9a31011-8401-4b72-af0f-0d7595328c0c",
	"wrld_383c3fd1-ef79-489f-a971-a025b53aadf3",
	"wrld_02cb84c8-446a-484b-b52b-2c1777e55f8d",
	"wrld_78869f89-fe4a-41c7-bbef-06836eb0176c",
	"wrld_791ebf58-54ce-4d3a-a0a0-39f10e1b20b2",
	"wrld_9e45e97e-c629-4e08-8a9a-6b04060b1485",
}

var DllPath string
var DllData []byte
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
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Get("/should_load", shouldLoad)
	app.Get("/:hash", downloadMod)

	app.Get("/risky_func/:worldId", isAllowed)

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
	c.Set("surrogate-key", "mod-resource")

	if strings.ToLower(c.Params("hash")) == DllHash {
		return c.SendStatus(204)
	}

	c.Set("content-type", "application/octet-stream")

	return c.Send(DllData)
}

func isAllowed(c *fiber.Ctx) error {
	c.Set("surrogate-key", "risky-func")

	worldId := c.Params("worldId")

	if slices.Contains(WhitelistedWorlds, worldId) {
		return c.SendString("allowed")
	}

	if slices.Contains(BlacklistedWorlds, worldId) {
		return c.SendString("denied")
	}

	return c.SendStatus(200)
}
