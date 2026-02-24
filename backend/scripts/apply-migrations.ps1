param(
    [string]$DatabaseUrl = $env:DATABASE_URL,
    [string]$MigrationsPath = "backend/migrations",
    [string]$PsqlPath = $env:PSQL_PATH
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

function Resolve-PsqlCommand {
    param([string]$ExplicitPath)

    if (![string]::IsNullOrWhiteSpace($ExplicitPath)) {
        if (!(Test-Path $ExplicitPath)) {
            throw "PSQL_PATH points to missing file: $ExplicitPath"
        }
        return $ExplicitPath
    }

    $cmd = Get-Command psql -ErrorAction SilentlyContinue
    if ($null -ne $cmd) {
        return $cmd.Source
    }

    throw "psql is not installed or not in PATH. Install PostgreSQL client or set PSQL_PATH."
}

function Invoke-Migration {
    param(
        [string]$PsqlCommand,
        [string]$Url,
        [string]$FilePath
    )
    Write-Host "Applying migration: $FilePath"
    & "$PsqlCommand" "$Url" -v ON_ERROR_STOP=1 -f "$FilePath"
}

Test-DatabaseUrl -Value $DatabaseUrl
$psql = Resolve-PsqlCommand -ExplicitPath $PsqlPath
$files = Get-MigrationFiles -Path $MigrationsPath

if ($files.Count -eq 0) {
    Write-Host "No migrations found in $MigrationsPath"
    exit 0
}

foreach ($file in $files) {
    Invoke-Migration -PsqlCommand $psql -Url $DatabaseUrl -FilePath $file.FullName
}

Write-Host "Migrations applied successfully"
