# Nexus Ops (Wails Edition)

## 🚀 Nexus Hub - Cyberpunk Command Center

**Nexus Ops** is een moderne desktop applicatie, herbouwd met **Wails** (Go + React + TailwindCSS).
Het combineert de kracht van Go's concurrency met de visuele vrijheid van moderne webtechnologie.

![Nexus Hub Banner](assets/banner.png)

## ✨ Features (Wails Upgrade)

### 1. Modern UI (React + Tailwind)
-   **Glassmorphism**: Echte transparantie en blur effecten (Mac/Windows).
-   **Frameless Window**: Custom Cyberpunk titlebar.
-   **Event-Driven**: Real-time updates via Wails Runtime Events (geen polling).

### 2. Core Modules
-   **Forge**: Image Processing met streaming progress.
    -   *Non-blocking UI*: Scan duizenden bestanden zonder haperingen.
-   **Radar**: Live taak monitoring.
    -   *Hybrid Fetching*: Directe load + Event updates.
-   **Sidebar**: Navigatie naar Vault, Siphon, Terminal (Placeholders).

## 🛠 Installatie & Start

Vereisten: Go 1.20+, Node 16+

1.  **Installeer Wails**:
    ```bash
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    ```

2.  **Start Development Mode**:
    ```bash
    wails dev
    ```
    *Dit start de app en een browser-venster met hot-reload.*

3.  **Build voor Productie**:
    ```bash
    wails build
    ```
    *De binary verschijnt in `build/bin/`.*

## 🏗 Architectuur

-   `main.go`: Entry point & Wails configuratie.
-   `app.go`: Bridge tussen Go en Javascript.
-   `internal/sys`: Core state en metrics (Go).
-   `internal/forge`: Worker pool logic (Go).
-   `frontend/`: React applicatie (UI).
    -   `src/components`: UI Componenten (Forge, Radar, Sidebar).
    -   `wailsjs`: Automatisch gegenereerde bindings.

---
*Gemigreerd van Fyne naar Wails voor superieure aesthetics.*
