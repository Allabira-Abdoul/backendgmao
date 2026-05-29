# 🚀 Déploiement GMAO Backend — Windows Server 2019 Datacenter (IIS + HTTPS, sans Docker)

> **Architecture** : IIS agit comme reverse proxy HTTPS devant le API Gateway Go (port 8080),
> qui route les requêtes vers les microservices via Consul Service Discovery.

```
┌─────────────────────────────────────────────────────────┐
│                 Windows Server 2019 DC                  │
│                                                         │
│  Client ──► IIS (443/HTTPS) ──► api-gateway (8080)      │
│                                    │                    │
│                              ┌─────┼─────┐              │
│                              ▼     ▼     ▼              │
│                          Consul (8500)                  │
│                              │                          │
│         ┌────────┬───────┬───┴───┬────────┬──────────┐  │
│         ▼        ▼       ▼       ▼        ▼          ▼  │
│     user-svc auth-svc asset-svc maint-svc analytics predict│
│     :8081    :8082    :8083     :8084     :8085    :8086 │
│                                                         │
│              PostgreSQL 17 (5432)                        │
└─────────────────────────────────────────────────────────┘
```

---

## Table des matières

1. [Prérequis](#1-prérequis)
2. [Installer PostgreSQL 17](#2-installer-postgresql-17)
3. [Installer Consul](#3-installer-consul)
4. [Installer Go et compiler les services](#4-installer-go-et-compiler-les-services)
5. [Configurer les variables d'environnement](#5-configurer-les-variables-denvironnement)
6. [Créer les services Windows (NSSM)](#6-créer-les-services-windows-nssm)
7. [Configurer IIS comme reverse proxy HTTPS](#7-configurer-iis-comme-reverse-proxy-https)
8. [Obtenir un certificat SSL](#8-obtenir-un-certificat-ssl)
9. [Pare-feu Windows](#9-pare-feu-windows)
10. [Vérification](#10-vérification)
11. [Maintenance et dépannage](#11-maintenance-et-dépannage)

---

## 1. Prérequis

| Composant | Version | Téléchargement |
|---|---|---|
| Windows Server 2019 Datacenter | — | Déjà installé |
| Go | 1.26+ | https://go.dev/dl/ |
| PostgreSQL | 17 | https://www.postgresql.org/download/windows/ |
| Consul | 1.20+ | https://developer.hashicorp.com/consul/install |
| NSSM | Latest | https://nssm.cc/download |
| IIS + ARR + URL Rewrite | — | Via Server Manager (voir section 7) |

---

## 2. Installer PostgreSQL 17

1. Télécharger et exécuter l'installeur PostgreSQL 17 pour Windows.
2. Pendant l'installation :
   - Définir le mot de passe du superutilisateur `postgres`.
   - Port par défaut : **5432**.
   - Locale : **French, France** ou **Default**.
3. Ouvrir **pgAdmin** ou **psql** et créer la base de données et l'utilisateur :

```sql
-- Créer l'utilisateur applicatif
CREATE USER gmao_user WITH PASSWORD 'VotreMotDePasseSecurise123!';

-- Créer la base de données
CREATE DATABASE gmao_db OWNER gmao_user;

-- Accorder les privilèges
GRANT ALL PRIVILEGES ON DATABASE gmao_db TO gmao_user;
```

4. Vérifier la connexion :

```powershell
psql -h 127.0.0.1 -U gmao_user -d gmao_db
```

---

## 3. Installer Consul

### 3.1 Téléchargement

```powershell
# Créer le répertoire d'installation
New-Item -ItemType Directory -Path "C:\consul" -Force
New-Item -ItemType Directory -Path "C:\consul\data" -Force
New-Item -ItemType Directory -Path "C:\consul\config" -Force

# Télécharger et dézipper Consul dans C:\consul\
# Lien : https://developer.hashicorp.com/consul/install
# Placer consul.exe dans C:\consul\
```

### 3.2 Configuration

Créer `C:\consul\config\server.json` :

```json
{
  "server": true,
  "bootstrap_expect": 1,
  "datacenter": "dc1",
  "data_dir": "C:\\consul\\data",
  "ui_config": {
    "enabled": true
  },
  "client_addr": "0.0.0.0",
  "bind_addr": "127.0.0.1",
  "log_level": "INFO"
}
```

### 3.3 Enregistrer Consul comme service Windows

```powershell
# Télécharger NSSM : https://nssm.cc/download
# Placer nssm.exe dans C:\nssm\

C:\nssm\nssm.exe install Consul "C:\consul\consul.exe" "agent -config-dir=C:\consul\config"
C:\nssm\nssm.exe set Consul AppDirectory "C:\consul"
C:\nssm\nssm.exe set Consul Start SERVICE_AUTO_START
C:\nssm\nssm.exe set Consul AppStdout "C:\consul\consul-stdout.log"
C:\nssm\nssm.exe set Consul AppStderr "C:\consul\consul-stderr.log"

# Démarrer le service
Start-Service Consul
```

### 3.4 Vérifier

```powershell
# Interface Web
Start-Process "http://127.0.0.1:8500/ui"

# CLI
C:\consul\consul.exe members
```

---

## 4. Installer Go et compiler les services

### 4.1 Installer Go

1. Télécharger Go 1.26+ depuis https://go.dev/dl/
2. Installer avec les options par défaut.
3. Vérifier :

```powershell
go version
# go version go1.26.2 windows/amd64
```

### 4.2 Cloner le projet

```powershell
# Créer le répertoire de déploiement
New-Item -ItemType Directory -Path "C:\gmao" -Force
New-Item -ItemType Directory -Path "C:\gmao\bin" -Force
New-Item -ItemType Directory -Path "C:\gmao\logs" -Force
New-Item -ItemType Directory -Path "C:\gmao\src" -Force

# Cloner le code source
cd C:\gmao\src
git clone https://github.com/Allabira-Abdoul/gmao-backend.git .
```

### 4.3 Compiler tous les services

Créer et exécuter le script `C:\gmao\build-all.ps1` :

```powershell
# build-all.ps1 — Compile tous les microservices pour Windows
$ErrorActionPreference = "Stop"
$services = @(
    "api-gateway",
    "user-service",
    "authentication-service",
    "asset-service",
    "maintenance-service",
    "analytics-service",
    "prediction-service"
)

Set-Location "C:\gmao\src"

foreach ($svc in $services) {
    Write-Host "=== Building $svc ===" -ForegroundColor Cyan
    $env:CGO_ENABLED = "0"
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    go build -ldflags="-s -w" -o "C:\gmao\bin\$svc.exe" "./apps/$svc/cmd/api"

    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERREUR: Compilation de $svc échouée!" -ForegroundColor Red
        exit 1
    }
    Write-Host "$svc.exe compilé avec succès" -ForegroundColor Green
}

Write-Host "`n=== Tous les services compilés ===" -ForegroundColor Green
Get-ChildItem "C:\gmao\bin\*.exe" | Format-Table Name, Length
```

```powershell
# Exécuter la compilation
powershell -ExecutionPolicy Bypass -File C:\gmao\build-all.ps1
```

---

## 5. Configurer les variables d'environnement

Créer `C:\gmao\env.ps1` (fichier de configuration centralisé) :

```powershell
# env.ps1 — Variables d'environnement partagées par tous les services

# PostgreSQL
$env:POSTGRES_USER     = "gmao_user"
$env:POSTGRES_PASSWORD = "VotreMotDePasseSecurise123!"
$env:POSTGRES_DB       = "gmao_db"
$env:DATABASE_URL      = "host=127.0.0.1 user=$env:POSTGRES_USER password=$env:POSTGRES_PASSWORD dbname=$env:POSTGRES_DB port=5432 sslmode=disable"

# Consul
$env:CONSUL_HOST = "127.0.0.1"
$env:CONSUL_PORT = "8500"

# Service Host (localhost car tout est sur la même machine)
$env:SERVICE_HOST = "127.0.0.1"

# CORS
$env:ALLOWED_ORIGINS = "https://gmao-frontend.vercel.app,https://votre-domaine.com"

# JWT
$env:JWT_SECRET         = "votre-secret-jwt-production-tres-long-et-unique"
$env:JWT_ACCESS_EXPIRY  = "15m"
$env:JWT_REFRESH_EXPIRY = "168h"
```

> [!CAUTION]
> **Ne jamais commiter ce fichier.** Ajoutez `env.ps1` à `.gitignore`.

---

## 6. Créer les services Windows (NSSM)

Chaque microservice sera enregistré comme un service Windows pour un démarrage automatique et un redémarrage en cas de crash.

### 6.1 Script d'installation de tous les services

Créer `C:\gmao\install-services.ps1` :

```powershell
# install-services.ps1 — Enregistre chaque microservice comme service Windows via NSSM
$ErrorActionPreference = "Stop"
$nssm = "C:\nssm\nssm.exe"

# Tableau : Nom du service, port, exécutable
$services = @(
    @{ Name="GMAO-ApiGateway";    Exe="api-gateway.exe";            Port="8080" },
    @{ Name="GMAO-UserService";   Exe="user-service.exe";           Port="8081" },
    @{ Name="GMAO-AuthService";   Exe="authentication-service.exe"; Port="8082" },
    @{ Name="GMAO-AssetService";  Exe="asset-service.exe";          Port="8083" },
    @{ Name="GMAO-MaintService";  Exe="maintenance-service.exe";    Port="8084" },
    @{ Name="GMAO-AnalyticsSvc";  Exe="analytics-service.exe";      Port="8085" },
    @{ Name="GMAO-PredictSvc";    Exe="prediction-service.exe";     Port="8086" }
)

foreach ($svc in $services) {
    $name = $svc.Name
    $exe  = "C:\gmao\bin\$($svc.Exe)"
    $port = $svc.Port

    Write-Host "Installing $name (port $port)..." -ForegroundColor Cyan

    & $nssm install $name $exe
    & $nssm set $name AppDirectory "C:\gmao\bin"
    & $nssm set $name Start SERVICE_AUTO_START

    # Logs
    & $nssm set $name AppStdout "C:\gmao\logs\$name-stdout.log"
    & $nssm set $name AppStderr "C:\gmao\logs\$name-stderr.log"
    & $nssm set $name AppRotateFiles 1
    & $nssm set $name AppRotateBytes 10485760  # 10 MB

    # Variables d'environnement
    $envVars  = "PORT=$port"
    $envVars += "`nCONSUL_HOST=127.0.0.1"
    $envVars += "`nCONSUL_PORT=8500"
    $envVars += "`nSERVICE_HOST=127.0.0.1"
    $envVars += "`nDATABASE_URL=host=127.0.0.1 user=gmao_user password=VotreMotDePasseSecurise123! dbname=gmao_db port=5432 sslmode=disable"
    $envVars += "`nJWT_SECRET=votre-secret-jwt-production-tres-long-et-unique"
    $envVars += "`nJWT_ACCESS_EXPIRY=15m"
    $envVars += "`nJWT_REFRESH_EXPIRY=168h"

    # Pour l'API Gateway uniquement
    if ($name -eq "GMAO-ApiGateway") {
        $envVars += "`nALLOWED_ORIGINS=https://gmao-frontend.vercel.app,https://votre-domaine.com"
    }

    & $nssm set $name AppEnvironmentExtra $envVars

    # Redémarrage automatique en cas de crash
    & $nssm set $name AppExit Default Restart
    & $nssm set $name AppRestartDelay 5000

    Write-Host "$name installé avec succès`n" -ForegroundColor Green
}

Write-Host "=== Tous les services installés ===" -ForegroundColor Green
```

### 6.2 Démarrer tous les services

```powershell
# Exécuter l'installation (en tant qu'Administrateur)
powershell -ExecutionPolicy Bypass -File C:\gmao\install-services.ps1

# Démarrer dans l'ordre
Start-Service Consul
Start-Sleep -Seconds 3

$services = @(
    "GMAO-UserService", "GMAO-AuthService", "GMAO-AssetService",
    "GMAO-MaintService", "GMAO-AnalyticsSvc", "GMAO-PredictSvc"
)
foreach ($svc in $services) {
    Start-Service $svc
    Write-Host "$svc démarré" -ForegroundColor Green
}
Start-Sleep -Seconds 2
Start-Service GMAO-ApiGateway
Write-Host "API Gateway démarré" -ForegroundColor Cyan
```

### 6.3 Vérifier l'état des services

```powershell
Get-Service GMAO-* | Format-Table Name, Status -AutoSize
Get-Service Consul  | Format-Table Name, Status -AutoSize
```

---

## 7. Configurer IIS comme reverse proxy HTTPS

### 7.1 Installer les rôles IIS

Ouvrir **PowerShell en Administrateur** :

```powershell
# Installer IIS avec les modules nécessaires
Install-WindowsFeature -Name Web-Server -IncludeManagementTools
Install-WindowsFeature -Name Web-Http-Redirect
Install-WindowsFeature -Name Web-WebSockets
```

### 7.2 Installer les extensions requises

Télécharger et installer **dans cet ordre** :

| Extension | Lien |
|---|---|
| **URL Rewrite 2.1** | https://www.iis.net/downloads/microsoft/url-rewrite |
| **Application Request Routing (ARR) 3.0** | https://www.iis.net/downloads/microsoft/application-request-routing |

> [!IMPORTANT]
> Redémarrer IIS après l'installation : `iisreset`

### 7.3 Activer le Proxy dans ARR

```powershell
# Activer le proxy ARR via appcmd
C:\Windows\System32\inetsrv\appcmd.exe set config -section:system.webServer/proxy /enabled:"True" /commit:apphost
C:\Windows\System32\inetsrv\appcmd.exe set config -section:system.webServer/proxy /preserveHostHeader:"True" /commit:apphost
```

### 7.4 Créer le site IIS

1. Ouvrir **IIS Manager** (`inetmgr`).
2. Supprimer ou arrêter le **Default Web Site**.
3. Clic droit sur **Sites** → **Add Website** :
   - **Site name** : `GMAO-API`
   - **Physical path** : `C:\gmao\www` (créer ce dossier vide)
   - **Binding** : HTTPS, port 443, choisir le certificat SSL (voir section 8)
   - Ajouter un deuxième binding : HTTP, port 80 (pour la redirection)

### 7.5 Configurer les règles de réécriture

Créer `C:\gmao\www\web.config` :

```xml
<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <system.webServer>

    <!-- === Redirection HTTP → HTTPS === -->
    <rewrite>
      <rules>
        <rule name="HTTP to HTTPS" stopProcessing="true">
          <match url="(.*)" />
          <conditions>
            <add input="{HTTPS}" pattern="^OFF$" />
          </conditions>
          <action type="Redirect" url="https://{HTTP_HOST}/{R:1}" redirectType="Permanent" />
        </rule>

        <!-- === Reverse Proxy vers API Gateway === -->
        <rule name="ReverseProxy to API Gateway" stopProcessing="true">
          <match url="(.*)" />
          <conditions>
            <add input="{HTTPS}" pattern="^ON$" />
          </conditions>
          <action type="Rewrite" url="http://127.0.0.1:8080/{R:1}" />
          <serverVariables>
            <set name="HTTP_X_FORWARDED_PROTO" value="https" />
            <set name="HTTP_X_FORWARDED_FOR" value="{REMOTE_ADDR}" />
            <set name="HTTP_X_REAL_IP" value="{REMOTE_ADDR}" />
          </serverVariables>
        </rule>
      </rules>
    </rewrite>

    <!-- === Autoriser les en-têtes personnalisés === -->
    <httpProtocol>
      <customHeaders>
        <add name="X-Content-Type-Options" value="nosniff" />
        <add name="X-Frame-Options" value="DENY" />
        <add name="Strict-Transport-Security" value="max-age=31536000; includeSubDomains" />
      </customHeaders>
    </httpProtocol>

    <!-- === WebSocket Support === -->
    <webSocket enabled="true" />

  </system.webServer>
</configuration>
```

### 7.6 Autoriser les Server Variables

Dans **IIS Manager** → Site `GMAO-API` → **URL Rewrite** → **View Server Variables** → Ajouter :

- `HTTP_X_FORWARDED_PROTO`
- `HTTP_X_FORWARDED_FOR`
- `HTTP_X_REAL_IP`

### 7.7 Appliquer les changements

```powershell
iisreset
```

---

## 8. Obtenir un certificat SSL

### Option A — Certificat auto-signé (développement / test)

```powershell
# Générer un certificat auto-signé valable 2 ans
$cert = New-SelfSignedCertificate `
    -DnsName "votre-domaine.com","localhost" `
    -CertStoreLocation "Cert:\LocalMachine\My" `
    -NotAfter (Get-Date).AddYears(2) `
    -FriendlyName "GMAO Backend SSL"

Write-Host "Thumbprint: $($cert.Thumbprint)"
```

Puis dans IIS Manager, modifier le binding HTTPS du site pour utiliser ce certificat.

### Option B — Certificat Let's Encrypt (production)

Utiliser **win-acme** (client ACME gratuit pour Windows) :

```powershell
# Télécharger win-acme : https://www.win-acme.com/
# Extraire dans C:\win-acme\

cd C:\win-acme
.\wacs.exe --target iis --siteid 1 --installation iis --accepttos --emailaddress admin@votre-domaine.com
```

> [!TIP]
> **win-acme** configure automatiquement le renouvellement via le Planificateur de tâches Windows.

### Option C — Certificat commercial

1. Générer un CSR (Certificate Signing Request) via IIS Manager.
2. Soumettre le CSR à votre autorité de certification (DigiCert, GlobalSign, etc.).
3. Importer le certificat `.pfx` dans IIS.

---

## 9. Pare-feu Windows

```powershell
# Autoriser HTTPS (port 443) depuis l'extérieur
New-NetFirewallRule -DisplayName "GMAO HTTPS" -Direction Inbound -Protocol TCP -LocalPort 443 -Action Allow

# Autoriser HTTP (port 80) pour la redirection
New-NetFirewallRule -DisplayName "GMAO HTTP Redirect" -Direction Inbound -Protocol TCP -LocalPort 80 -Action Allow

# BLOQUER l'accès direct aux ports internes depuis l'extérieur
# Les ports 8080-8086 et 8500 ne doivent PAS être exposés
New-NetFirewallRule -DisplayName "Block Direct Backend" -Direction Inbound -Protocol TCP -LocalPort 8080-8086 -Action Block -Profile Domain,Public
New-NetFirewallRule -DisplayName "Block Direct Consul" -Direction Inbound -Protocol TCP -LocalPort 8500 -Action Block -Profile Domain,Public
```

---

## 10. Vérification

### 10.1 Consul — Services enregistrés

```powershell
# Via CLI
C:\consul\consul.exe catalog services

# Résultat attendu :
# api-gateway
# user-service
# authentication-service
# asset-service
# maintenance-service
# analytics-service
# prediction-service
```

### 10.2 Health checks

```powershell
# API Gateway (direct, depuis le serveur)
Invoke-RestMethod -Uri http://127.0.0.1:8080/health
# Attendu : { "status": "UP", "service": "api-gateway" }

# Via IIS / HTTPS (depuis n'importe où)
Invoke-RestMethod -Uri https://votre-domaine.com/health -SkipCertificateCheck
# Même résultat attendu

# Microservices individuels
Invoke-RestMethod -Uri http://127.0.0.1:8081/health   # user-service
Invoke-RestMethod -Uri http://127.0.0.1:8082/health   # authentication-service
Invoke-RestMethod -Uri http://127.0.0.1:8083/health   # asset-service
```

### 10.3 Test d'un endpoint API

```powershell
# Login
$body = @{ email = "admin@gmao.com"; password = "admin123" } | ConvertTo-Json
Invoke-RestMethod -Uri https://votre-domaine.com/api/authentication/login -Method POST -Body $body -ContentType "application/json" -SkipCertificateCheck
```

---

## 11. Maintenance et dépannage

### Commandes utiles

```powershell
# Voir l'état de tous les services GMAO
Get-Service GMAO-*, Consul | Format-Table Name, Status

# Redémarrer un service spécifique
Restart-Service GMAO-UserService

# Redémarrer tous les services
Get-Service GMAO-* | Restart-Service

# Consulter les logs d'un service
Get-Content "C:\gmao\logs\GMAO-UserService-stderr.log" -Tail 50

# Suivre les logs en temps réel
Get-Content "C:\gmao\logs\GMAO-ApiGateway-stderr.log" -Tail 20 -Wait
```

### Mettre à jour les services

```powershell
# 1. Récupérer le nouveau code
cd C:\gmao\src
git pull origin main

# 2. Recompiler
powershell -ExecutionPolicy Bypass -File C:\gmao\build-all.ps1

# 3. Redémarrer les services
Get-Service GMAO-* | Stop-Service
Start-Sleep -Seconds 2
Get-Service GMAO-* | Start-Service
```

### Désinstaller un service

```powershell
Stop-Service GMAO-UserService
C:\nssm\nssm.exe remove GMAO-UserService confirm
```

---

## Résumé de l'architecture déployée

| Composant | Port | Rôle |
|---|---|---|
| **IIS (ARR)** | 80, **443** | Reverse proxy HTTPS, point d'entrée unique |
| **api-gateway.exe** | 8080 | Routage vers les microservices via Consul |
| **user-service.exe** | 8081 | Gestion des utilisateurs et rôles |
| **authentication-service.exe** | 8082 | Authentification JWT |
| **asset-service.exe** | 8083 | Gestion des actifs / équipements |
| **maintenance-service.exe** | 8084 | Gestion de la maintenance |
| **analytics-service.exe** | 8085 | Analyses et rapports |
| **prediction-service.exe** | 8086 | Prédiction de pannes (ML) |
| **Consul** | 8500 | Service Discovery |
| **PostgreSQL** | 5432 | Base de données |

> [!NOTE]
> Tous les ports internes (8080–8086, 8500, 5432) sont bloqués par le pare-feu
> pour le trafic externe. Seuls les ports **80** et **443** sont exposés via IIS.
