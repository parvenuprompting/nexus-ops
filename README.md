# Nexus Ops (Wails Edition)

## 🚀 Nexus Hub - Cyberpunk Command Center

**Nexus Ops** is een moderne desktop applicatie, herbouwd met **Wails** (Go + React + TailwindCSS).
Het combineert de brute kracht van Go's concurrency met de visuele elegantie van moderne webtechnologie en Glassmorphism.

![Nexus Hub Banner](assets/nexus_banner.png)

## ✨ Highlights (Wails Upgrade)

### 🎨 Modern UI (React + Tailwind)
-   **Glassmorphism**: Prachtige transparantie en `backdrop-blur` effecten (geoptimaliseerd voor Mac/Windows).
-   **Frameless Window**: Een custom Cyberpunk titlebar voor een echte desktop ervaring.
-   **Event-Driven Telemetry**: Real-time updates via het Wails event systeem (geen polling nodig).

### 🛠 Core Modules
-   **Forge**: Geavanceerde Image Processing met streaming progress.
    -   *NIEUW*: Keuze uit formaten (JPG/PNG) en Aspect Ratios (16:9, 1:1, etc.).
    -   *Non-blocking UI*: Scan en verwerk duizenden bestanden zonder haperingen in de interface.
-   **Radar**: Live monitoring van alle systeemtaken met kleurgecodeerde badges en real-time duur tracking.
-   **Terminal**: Directe uitvoering van systeemcommando's voor power-users.
-   **Vault**: Veilige, lokale opslag voor API keys en configuratie geheimen.
-   **Siphon**: Netwerk diagnostiek voor het testen van API latency en bereikbaarheid.

### 🛡 Stabiliteit
-   **Panic Recovery**: Automatische detectie en logging van crashes in achtergrondprocessen zonder de UI te bevriezen.

## 🚀 Snel Starten

Vereisten: Go 1.25+, Node 18+

1.  **Installeer Wails CLI**:
    ```bash
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    ```

2.  **Start Development Mode**:
    ```bash
    wails dev
    ```
    *Met hot-reload voor zowel Go als React wijzigingen.*

3.  **Bouw voor Productie**:
    ```bash
    wails build
    ```
    *De binary verschijnt direct in `build/bin/`.*

## 🏗 Architectuur

-   `main.go`: Applicatie entry point & venster configuratie.
-   `app.go`: De bridge die Go functies en events blootstelt aan de frontend.
-   `internal/sys`: Beheert de wereldwijde state, metrics en taakbeheer (Go).
-   `internal/workers`: High-performance worker pools voor parallelle taken.
-   `frontend/`: Moderne React applicatie met TailwindCSS v3 styling.

---
*Nexus Ops: Kracht ontmoet Esthetiek.*
*© 2026 Tiëndo Welles*
