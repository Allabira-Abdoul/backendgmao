# ==============================================================================
# GMAO Backend - Local Development Launch Script (Windows PowerShell)
# Usage: .\run_all.ps1
#
# Port Mapping (chosen to avoid Hyper-V exclusions 7886-8085, 9227-9726):
#   user-service        -> 8100
#   analytics-service   -> 8101
#   asset-service       -> 8102
#   auth-service        -> 8103
#   maintenance-service -> 8104
#   prediction-service  -> 8105
#   audit-service       -> 8106
#   api-gateway         -> 8200
# ==============================================================================

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "Starting GMAO Backend Environment for Local Testing..." -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# 1. Start Docker infrastructure (Postgres and Consul)
Write-Host "Spin up Postgres, Consul, and RabbitMQ containers..." -ForegroundColor Yellow
docker compose -f deploy/docker-compose.yml up -d consul postgres rabbitmq

# Wait for Postgres health check
Write-Host "Waiting 6 seconds for database to initialize..." -ForegroundColor Yellow
Start-Sleep -Seconds 6

# 2. Create helper function to spawn microservice in a dedicated window
# IMPORTANT: No spaces before/after & in CMD set commands —
# CMD includes trailing spaces in variable values which breaks URL parsing.
function Start-ServiceProcess {
    param (
        [string]$Name,
        [string]$Path,
        [int]$Port
    )
    Write-Host "Starting $Name on port $Port..." -ForegroundColor Green
    $root = (Get-Location).Path
    $argList = "/k title $Name&cd /d `"$root\$Path`"&set PORT=$Port&set CONSUL_HOST=127.0.0.1&set CONSUL_PORT=8500&set JWT_SECRET=gmao-dev-secret-change-in-production&set `"DATABASE_URL=host=127.0.0.1 user=gmao_user password=gmao_password dbname=gmao_db port=5432 sslmode=disable`"&go run ."
    Start-Process cmd.exe -ArgumentList $argList
}

# 3. Start all microservices (each in its own window for log visibility)
Start-ServiceProcess -Name "user-service"        -Path "apps\user-service\cmd\api"        -Port 8100
Start-ServiceProcess -Name "analytics-service"   -Path "apps\analytics-service\cmd\api"   -Port 8101
Start-ServiceProcess -Name "asset-service"       -Path "apps\asset-service\cmd\api"       -Port 8102
Start-ServiceProcess -Name "auth-service"        -Path "apps\auth-service\cmd\api"        -Port 8103
Start-ServiceProcess -Name "maintenance-service" -Path "apps\maintenance-service\cmd\api" -Port 8104
Start-ServiceProcess -Name "prediction-service"  -Path "apps\prediction-service\cmd\api"  -Port 8105
Start-ServiceProcess -Name "audit-service"       -Path "apps\audit-service\cmd\api"       -Port 8106

# Wait before launching Gateway so services can register with Consul first
Start-Sleep -Seconds 3

# 4. Launch the API Gateway
Write-Host "Starting API Gateway on http://127.0.0.1:8200..." -ForegroundColor Green
$root = (Get-Location).Path
$gatewayArg = "/k title api-gateway&cd /d `"$root\apps\api-gateway\cmd\api`"&set PORT=8200&set CONSUL_HOST=127.0.0.1&set CONSUL_PORT=8500&set JWT_SECRET=gmao-dev-secret-change-in-production&go run ."
Start-Process cmd.exe -ArgumentList $gatewayArg

Write-Host ""
Write-Host "Environment launched! Close individual service windows to stop them." -ForegroundColor Cyan
Write-Host "Gateway Endpoint  : http://127.0.0.1:8200" -ForegroundColor Cyan
Write-Host "Consul Dashboard  : http://127.0.0.1:8500" -ForegroundColor Cyan
Write-Host ""
Write-Host "Service Port Map:" -ForegroundColor DarkCyan
Write-Host "  user-service        -> http://127.0.0.1:8100" -ForegroundColor DarkCyan
Write-Host "  analytics-service   -> http://127.0.0.1:8101" -ForegroundColor DarkCyan
Write-Host "  asset-service       -> http://127.0.0.1:8102" -ForegroundColor DarkCyan
Write-Host "  auth-service        -> http://127.0.0.1:8103" -ForegroundColor DarkCyan
Write-Host "  maintenance-service -> http://127.0.0.1:8104" -ForegroundColor DarkCyan
Write-Host "  prediction-service  -> http://127.0.0.1:8105" -ForegroundColor DarkCyan
Write-Host "  audit-service       -> http://127.0.0.1:8106" -ForegroundColor DarkCyan
