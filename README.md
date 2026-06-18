<div align="center">
  <h1>🐙 ~ Hera</h1>
  <p>Turn Based RogueLike?</p>
</div>

<p align="center">
  <a href="https://github.com/IwnuplyNotTyan/Hera/actions/workflows/ci.yml">
    <img src="https://img.shields.io/github/actions/workflow/status/IwnuplyNotTyan/Hera/ci.yml" alt="Build Status"/>
  </a>
  <img src="https://img.shields.io/github/license/IwnuplyNotTyan/Hera" alt="License"/>
  <img src="https://img.shields.io/github/stars/IwnuplyNotTyan/Hera" alt="Stars"/>
  <img src="https://img.shields.io/github/last-commit/IwnuplyNotTyan/Hera" alt="Last Commit"/>
</p>

<p align="center">
  <img src="https://github.com/IwnuplyNotTyan/Hera/blob/main/.github/assets/screenshot.png?raw=true" alt="Screenshot">
</p>

---

# 🪭 Tree 
-  [Build](https://github.com/IwnuplyNotTyan/Hera?tab=readme-ov-file#-build)
-  [Author](https://github.com/IwnuplyNotTyan/Hera?tab=readme-ov-file#%E2%80%8D-author)
-  [Libs](https://github.com/IwnuplyNotTyan/Hera?tab=readme-ov-file#%EF%B8%8F-libraries-used)
-  [License](https://github.com/IwnuplyNotTyan/Hera?tab=readme-ov-file#-license)

---

# 🪻 Install

## ❄️ Nix
``` bash
nix run github:iwnuplynottyan/koi
```


## ⛏️ Build from source
```sh
git clone https://github.com/IwnuplyNotTyan/Hera && cd Hera
go mod download
go build -o ./bin/hera ./cmd/hera/main.go
```

<details>
<summary><b>Flags and tags</b></summary>

**Tags:**

| Tag | Desc                 |
|-----|----------------------|
| eng | Only English Locales |
| notint| Disable themes     |


**LDFlags:**
| Flags | Example | Desc |
|-------|---------|------|
|main.version|-ldflags="-X main.version=1.0.0| --version number|
|main.commit|-ldflags="-X main.commit=$(git rev-parse HEAD)"| --version commit|


**Flags** can be find in help command

</details>

---

## 🐛 Developer Console

Opened by pressing `v` after enabling **Debug Mode** (`--debug`/`-d` flag or toggle in Settings).

Toggle: `v` cycles effects → skip turn → console (Debug Mode only).

| Command | Description |
|---------|-------------|
| `help` | Show this help |
| `status` | Show all entities with HP, position, effects |
| `clear` | Clear console output |
| `regen` | Regenerate walls, water; clear fire/smoke |
| `add player` | Spawn a new player on a free tile |
| `add enemy` | Spawn a new enemy on a free tile |
| `remove player <N>` | Remove player N (1‑based) |
| `remove enemy <N>` | Remove enemy N (1‑based) |
| `remove effect <target> <type>` | Remove an effect |
| `heal <target> <N>` | Heal target by N HP |
| `damage <target> <N>` | Damage target by N HP |
| `set hp <target> <N>` | Set exact HP |
| `effect <target> <type> [dur]` | Add effect (`fire`/`wet`/`smoke`), default duration 3 |

**Target notation:** `p1`, `p2`, … for players, `e1`, `e2`, … for enemies (1‑based).

---

## 👩‍💻 Author
- [iwnuplynottyan](https://github.com/iwnuplynottyan)

---

## 🛠️ Libraries Used
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework, core architecture
    - [Bubbletint](https://github.com/lrstanley/bubbletint) — Themes
    - [Bubbles](https://github.com/charmbracelet/bubbles) — Modular widgets/components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [Log](https://github.com/charmbracelet/log) —  Pretty logs
- [Testify](https://github.com/stretchr/testify) —  Enchaned testing
- [Cobra](https://github.com/spf13/cobra) — Powerfull flags
    - [Fang](https://github.com/charmbracelet/fang) —  Make it pretty

### 🌑 Assets

- [Tmplr](https://github.com/patorjk/figlet.js/blob/main/fonts/Tmplr.flf) —  ASCII font

---

## 📄 License
[MIT](https://github.com/IwnuplyNotTyan/Hera/blob/main/LICENSE).

<div align="center">
  <h1>Made with ❤️ </h1>
</div>
