package main

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const smmRelease = "https://github.com/satisfactorymodding/SatisfactoryModManager/releases/latest"
const cliRelease = "https://github.com/satisfactorymodding/ficsit-cli/releases/latest"

var smm = map[string]map[string]string{
	"windows": {
		"amd64": "https://github.com/satisfactorymodding/SatisfactoryModManager/releases/latest/download/Satisfactory-Mod-Manager-Setup.exe",
	},
	"linux": {
		"amd64": "https://github.com/satisfactorymodding/SatisfactoryModManager/releases/latest/download/Satisfactory-Mod-Manager.AppImage",
	},
}

var cli = map[string]map[string]map[string]string{
	"windows": {
		"amd64": {
			"binary": "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_windows_amd64.exe",
		},
		"386": {
			"binary": "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_windows_386.exe",
		},
		"arm64": {
			"binary": "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_windows_arm64.exe",
		},
		"armv7": {
			"binary": "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_windows_armv7.exe",
		},
	},
	"linux": {
		"amd64": {
			"binary": "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_amd64",
			"deb":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_amd64.deb",
			"rpm":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_amd64.rpm",
			"apk":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_amd64.apk",
		},
		"386": {
			"binary": "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_386",
			"deb":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_386.deb",
			"rpm":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_386.rpm",
			"apk":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_386.apk",
		},
		"arm64": {
			"binary": "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_arm64",
			"deb":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_arm64.deb",
			"rpm":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_arm64.rpm",
			"apk":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_arm64.apk",
		},
		"armv7": {
			"binary": "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_armv7",
			"deb":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_armv7.deb",
			"rpm":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_armv7.rpm",
			"apk":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_armv7.apk",
		},
		"ppc64le": {
			"binary": "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_ppc64le",
			"deb":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_ppc64le.deb",
			"rpm":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_ppc64le.rpm",
			"apk":    "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_linux_ppc64le.apk",
		},
	},
	"darwin": {
		"amd64": {
			"binary": "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_darwin_all",
		},
		"386": {
			"binary": "https://github.com/satisfactorymodding/ficsit-cli/releases/latest/download/ficsit_darwin_all",
		},
	},
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/:platform/:arch/:packaging", handleRequest)
	e.GET("/:platform/:arch", handleRequest)
	e.GET("/:platform", handleRequest)
	e.GET("/", handleRequest)

	e.Use(middleware.Logger())

	_ = e.Start(":8080")
}

type Platform string

const (
	Windows = Platform("windows")
	Linux   = Platform("linux")
	Darwin  = Platform("darwin")
)

func identifyPlatform(userAgent string) Platform {
	lower := strings.ToLower(userAgent)

	if strings.Contains(lower, "macintosh") ||
		strings.Contains(lower, "darwin") {
		return Darwin
	}

	if strings.Contains(lower, "linux") ||
		strings.Contains(lower, "ubuntu") ||
		strings.Contains(lower, "debian") {
		return Linux
	}

	return Windows
}

func handleRequest(context echo.Context) error {
	host := context.Request().Host

	if strings.HasPrefix(host, "smm.") {
		return handleSMM(context)
	}

	if strings.HasPrefix(host, "cli.") {
		return handleCLI(context)
	}

	return context.NoContent(404)
}

func handleSMM(context echo.Context) error {
	if url := resolveMap(smm, context); url != nil {
		return context.Redirect(301, *url)
	}
	return context.Redirect(301, smmRelease)
}

func handleCLI(context echo.Context) error {
	archData := resolveMap(cli, context)
	if archData == nil {
		return context.Redirect(301, cliRelease)
	}

	packaging := context.Param("packaging")
	if packaging == "" {
		packaging = "binary"
	}

	url, ok := (*archData)[packaging]
	if !ok {
		return context.Redirect(301, cliRelease)
	}

	return context.Redirect(301, url)
}

func resolveMap[T any](data map[string]map[string]T, context echo.Context) *T {
	platform := Platform(context.Param("platform"))
	if platform == "" {
		platform = identifyPlatform(context.Request().UserAgent())
	}

	platformData, ok := data[string(platform)]
	if !ok {
		return nil
	}

	arch := context.Param("arch")
	if arch == "" {
		arch = "amd64"
	}

	archData, ok := platformData[arch]
	if !ok {
		return nil
	}

	return &archData
}
