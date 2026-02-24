param(
    [string]$DatabaseUrl = $env:DATABASE_URL,
    [string]$MigrationsPath = "backend/migrations"
)

$ErrorActionPreference = "Stop"

function Test-DatabaseUrl {
    param([string]$Value)
    if ([string]::IsNullOrWhiteSpace($Value)) {
        throw "DATABASE_URL is empty. Pass -DatabaseUrl or set env:DATABASE_URL"
    }
}

function Get-MigrationFiles {
    param([string]$Path)
    if (!(Test-Path $Path)) {
        throw "Migrations path not found: $Path"
    }
    return Get-ChildItem -Path $Path -Filter *.sql | Sort-Object Name
}

function Invoke-Migration {
    param(
        [string]$Url,
        [string]$FilePath
    )
    Write-Host "Applying migration: $FilePath"
    psql "$Url" -v ON_ERROR_STOP=1 -f "$FilePath"
}

Test-DatabaseUrl -Value $DatabaseUrl
$files = Get-MigrationFiles -Path $MigrationsPath

if ($files.Count -eq 0) {
    Write-Host "No migrations found in $MigrationsPath"
    exit 0
}

foreach ($file in $files) {
    Invoke-Migration -Url $DatabaseUrl -FilePath $file.FullName
}

Write-Host "Migrations applied successfully"
