[CmdletBinding()]
param(
    [Parameter(Position = 0)]
    [string]$Start = "."
)

# Print the Go language version for the project rooted at the current working
# directory (or at -Start). Compatible with Windows PowerShell 5.1+.

$SkippedDirectories = @(".git", "vendor", "testdata", "node_modules")

function Get-GoDirective {
    param([string]$Path)

    foreach ($Line in Get-Content -LiteralPath $Path -ErrorAction SilentlyContinue) {
        if ($Line -match '^\s*go\s+([0-9]+\.[0-9]+(?:\.[0-9]+)?(?:[A-Za-z]+[0-9]*)?)') {
            return $Matches[1]
        }
    }
    return $null
}

function Get-ScannedVersions {
    param([string]$Root)

    foreach ($Item in Get-ChildItem -LiteralPath $Root -Force -ErrorAction SilentlyContinue) {
        if ($Item.PSIsContainer) {
            if ($SkippedDirectories -notcontains $Item.Name) {
                Get-ScannedVersions $Item.FullName
            }
        } elseif ($Item.Name -eq 'go.mod') {
            $Version = Get-GoDirective $Item.FullName
            if ($Version) {
                Write-Output $Version
            }
        }
    }
}

function Get-VersionParts {
    param([string]$Version)

    if ($Version -notmatch '^([0-9]+)\.([0-9]+)(?:\.([0-9]+))?') {
        return $null
    }
    return @([int]$Matches[1], [int]$Matches[2], $(if ($Matches[3]) { [int]$Matches[3] } else { 0 }))
}

try {
    $StartItem = Get-Item -LiteralPath $Start -ErrorAction Stop
    $Root = if ($StartItem.PSIsContainer) { $StartItem.FullName } else { $StartItem.Directory.FullName }
} catch {
    Write-Output unknown
    exit 0
}

# The nearest go.mod up the tree is the active module; its version wins.
$Directory = $Root
while ($true) {
    $ModFile = Join-Path $Directory 'go.mod'
    if (Test-Path -LiteralPath $ModFile -PathType Leaf) {
        $Version = Get-GoDirective $ModFile
        Write-Output $(if ($Version) { $Version } else { 'unknown' })
        exit 0
    }
    $Parent = [IO.Directory]::GetParent($Directory)
    if ($null -eq $Parent) {
        break
    }
    $Directory = $Parent.FullName
}

# No enclosing module: use the lowest version across modules below the root.
$Minimum = $null
$MinimumParts = $null
foreach ($Version in @(Get-ScannedVersions $Root)) {
    $Parts = Get-VersionParts $Version
    if ($null -eq $Parts) {
        continue
    }
    $IsLess = $null -eq $Minimum
    if (-not $IsLess) {
        for ($Index = 0; $Index -lt 3; $Index++) {
            if ($Parts[$Index] -lt $MinimumParts[$Index]) {
                $IsLess = $true
                break
            }
            if ($Parts[$Index] -gt $MinimumParts[$Index]) {
                break
            }
        }
    }
    if ($IsLess) {
        $Minimum = $Version
        $MinimumParts = $Parts
    }
}
Write-Output $(if ($Minimum) { $Minimum } else { 'unknown' })
