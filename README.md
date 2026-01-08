# Nexus Ops

Nexus Ops is een krachtige desktopapplicatie gebouwd met Go en Fyne v2. De applicatie demonstreert geavanceerde concurrency patronen en systeemobservatie.

## Functionaliteiten

### 1. Forge (De Smederij)
Een krachtige worker pool voor het verwerken van afbeeldingen.
- **Image Processing**: Selecteer een map met afbeeldingen en laat Nexus Ops ze automatisch verkleinen.
- **Concurrency**: Gebruikt meerdere workers (afhankelijk van je CPU cores) om taken parallel uit te voeren.
- **Live Voortgang**: Volg de status via een voortgangsbalk.

### 2. Radar (Het Systeem)
Real-time inzicht in wat er onder de motorkap gebeurt.
- **Visualisatie**: Zie actieve taken en processen live op het scherm.
- **Metrics**: Monitor het aantal actieve goroutines, voltooide taken en eventuele fouten.

## Installatie en Starten

Zorg dat je Go 1.25+ geïnstalleerd hebt.

```bash
# Clone de repository
git clone https://github.com/parvenuprompting/nexus-ops.git
cd nexus-ops

# Afhankelijkheden installeren
go mod tidy

# Starten
go run ./cmd/nexus
```

## Architectuur

De code is modulair opgebouwd:
- `cmd/nexus`: Het entrypoint van de applicatie.
- `internal/forge`: De logica voor de worker pool en image processing.
- `internal/sys`: Systeem-brede types, state en metrics.
- `internal/ui`: De gebruikersinterface (Fyne).

## Auteur

Gebouwd met ❤️ door Antigravity.
