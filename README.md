# Nexus Ops

![Nexus Hub Banner](assets/banner.png)

**Nexus Ops** is een high-performance desktop applicatie gebouwd met **Go** en **Fyne v2**. Het fungeert als een "Command Center" voor diverse systeemtaken, met een focus op concurrency, image processing en real-time monitoring.

## 🚀 Features

### 1. Nexus Hub Architectuur
De applicatie maakt gebruik van een moderne **Sidebar + Content** layout ("Glassmorphism" stijl).
-   **Forge**: Een krachtige Image Processing module voor batch resizing.
    -   *Streaming Scanner*: Verwerkt duizenden bestanden zonder UI freeze.
    -   *Smart Workers*: Gebruikt `runtime.NumCPU()` workers en semaphores voor geheugenveiligheid.
-   **Radar**: Real-time taak monitoring.
    -   *Event-Driven*: Directe updates bij start/stop van taken.
    -   *Visuals*: Kleurgecodeerde badges (CPU, I/O, IMG).
-   **The Vault** (Placeholder): Secure storage voor encrypted assets.
-   **Siphon** (Placeholder): Data stream dashboard.
-   **Terminal** (Placeholder): In-app console emulator.

### 2. Technologie & Performance
-   **Concurrency**: Gebouwd op Go's concurrency primitieven (Goroutines, Channels, WaitGroups).
-   **Type-Safe State**: Robuust state management met RWMutex beschermde maps.
-   **Memory Safety**: Automatische limitering van zware operaties om OOM crashes te voorkomen.
-   **Fyne v2 UI**: Cross-platform GUI met een custom Cyberpunk thema.

## 🛠 Installatie & Start

Vereisten: Go 1.20 of hoger.

```bash
# Clone de repository
git clone https://github.com/parvenuprompting/nexus-ops.git
cd nexus-ops

# Dependencies installeren
go mod tidy

# Applicatie starten
go run ./cmd/nexus
```

## 🏗 Architectuur

De codebase volgt een strikte scheiding van verantwoordelijkheden:

-   `cmd/nexus`: Entry point en setup.
-   `internal/sys`: Core state, metrics en task definities.
-   `internal/ui`: Fyne UI layout, tabbladen en custom widgets.
-   `internal/forge`: Business logic voor image processing en worker pools.

## 📷 Screenshots

Zie de banner hierboven voor de actuele "Glassmorphism" interface.

---
*Ontwikkeld als demonstratie van Advanced Agentic Coding met Go en Fyne.*
