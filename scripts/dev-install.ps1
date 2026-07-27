# Developer-only helper: build this checkout into the tool's cache so any agent
# using the plugin runs your local changes instead of the released version.
$ErrorActionPreference = "Stop"

$binaryName = "go-modern-guidelines.exe"

if ($env:LOCALAPPDATA) {
    $cacheRoot = Join-Path $env:LOCALAPPDATA "go-modern-guidelines"
} else {
    Write-Error "go-modern-guidelines: LOCALAPPDATA must be set"
    exit 1
}

$devDir = Join-Path $cacheRoot "dev"
$devBinary = Join-Path $devDir $binaryName

$scriptDir = Split-Path -Parent $PSCommandPath
$moduleDir = (Resolve-Path -LiteralPath (Join-Path $scriptDir "..")).Path

$command = if ($args.Count -ge 1) { $args[0] } else { "install" }

switch ($command) {
    "install" {
        if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
            Write-Error "go-modern-guidelines: Go toolchain is required to build the dev binary"
            exit 1
        }
        if (-not (Test-Path -LiteralPath (Join-Path $moduleDir "go.mod") -PathType Leaf)) {
            Write-Error "go-modern-guidelines: no Go module found at $moduleDir"
            exit 1
        }
        New-Item -ItemType Directory -Path $devDir -Force | Out-Null
        $staged = "$devBinary.tmp.$PID"
        $previousGoFlags = $env:GOFLAGS
        $previousGoWork = $env:GOWORK
        $previousCgoEnabled = $env:CGO_ENABLED
        $env:GOFLAGS = ""
        $env:GOWORK = "off"
        $env:CGO_ENABLED = "0"
        Push-Location -LiteralPath $moduleDir
        try {
            go build -o $staged .
            if ($LASTEXITCODE -ne 0) {
                Write-Error "go-modern-guidelines: failed to build dev binary"
                exit 1
            }
            Move-Item -LiteralPath $staged -Destination $devBinary -Force
        } finally {
            Pop-Location
            $env:GOFLAGS = $previousGoFlags
            $env:GOWORK = $previousGoWork
            $env:CGO_ENABLED = $previousCgoEnabled
            Remove-Item -LiteralPath $staged -Force -ErrorAction SilentlyContinue
        }
        Write-Host "go-modern-guidelines: installed dev build to $devBinary" -ForegroundColor DarkGray
        Write-Host "go-modern-guidelines: set GO_MODERN_GUIDELINES_DEV=1 to use it" -ForegroundColor DarkGray
    }
    "uninstall" {
        Remove-Item -LiteralPath $devDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Host "go-modern-guidelines: removed dev build ($devBinary)" -ForegroundColor DarkGray
    }
    default {
        Write-Error "usage: pwsh scripts/dev-install.ps1 [install|uninstall]"
        exit 2
    }
}
