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

![Screenshot](https://github.com/IwnuplyNotTyan/Hera/blob/main/.github/assets/screenshot.png?raw=true)

---

## 🪭 Tree 
-  [Build](https://github.com/IwnuplyNotTyan/Hera?tab=readme-ov-file#-build)
-  [Author](https://github.com/IwnuplyNotTyan/Hera?tab=readme-ov-file#%E2%80%8D-author)
-  [Libs](https://github.com/IwnuplyNotTyan/Hera?tab=readme-ov-file#%EF%B8%8F-libraries-used)
-  [License](https://github.com/IwnuplyNotTyan/Hera?tab=readme-ov-file#-license)

---

## 🪻 Build 
```sh
git clone https://github.com/IwnuplyNotTyan/Hera && cd Hera
go mod download
go build -o ./bin/hera ./cmd/hera/main.go
```

**Supported tags:**

| Tag | Desc                 |
|-----|----------------------|
| eng | Only English Locales |


**Supported flags:**

| Short | Long     | Desc            |
|-------|----------|-----------------|
| -h    | help     | Help            |
| -l    | --lang   | Change language |
| -t    | --theme  | Change themes   |
| -v    | --version| Version         |


---

## 👩‍💻 Author
- [iwnuplynottyan](https://github.com/iwnuplynottyan)

---

## 🛠️ Libraries Used
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework, core architecture
    - [Bubblezone](https://github.com/lrstanley/bubblezone) — Mouse support
    - [Bubbletint](https://github.com/lrstanley/bubbletint) — Themes
    - [Bubbles](https://github.com/charmbracelet/bubbles) — Modular widgets/components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [Log](https://github.com/charmbracelet/log) —  Pretty logs
- [Testify](https://github.com/stretchr/testify) —  Enchaned testing
- [Cobra](https://github.com/spf13/cobra) — Powerfull flags
    - [Fang](https://github.com/charmbracelet/fang) —  Make it pretty

---

## 📄 License
[MIT](https://github.com/IwnuplyNotTyan/Hera/blob/main/LICENSE).

<div align="center">
  <h1>Made with ❤️ </h1>
</div>
