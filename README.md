# hatrack

[![ci](https://github.com/baodq97/hatrack/actions/workflows/ci.yml/badge.svg)](https://github.com/baodq97/hatrack/actions/workflows/ci.yml)
[![release](https://img.shields.io/github/v/release/baodq97/hatrack)](https://github.com/baodq97/hatrack/releases/latest)
[![license](https://img.shields.io/badge/license-Apache--2.0-blue)](LICENSE)

Switch the account Claude Code is signed in to, without signing out and back in.

If you have more than one Claude subscription (personal and work, or two plans), hatrack keeps each
login and swaps the one Claude uses with one click in the Windows tray, or one command on macOS,
Linux and WSL.
It works for the `claude` CLI and the Claude Code extension for VS Code, since both read the same login.

- **Windows:** a tray app (click the hat, pick an account) plus the `hat` CLI.
- **macOS:** a menu bar app with a dashboard window, plus the `hat` CLI.
- **Linux and WSL:** the `hat` CLI.

```
$ hat
* personal@gmail.com             personal@gmail.com  personal@gmail.com's Organization
  work@company.com               work@company.com    Company
$ hat use work@company.com
now using work@company.com. If an open Claude session still shows the old account, restart it (VS Code: Reload Window).
```

## Contents

- [Install](#install)
- [Getting started](#getting-started)
- [Commands](#commands)
- [Windows tray app](#windows-tray-app)
- [macOS menu bar app](#macos-menu-bar-app)
- [Adding an account](#adding-an-account)
- [VS Code](#vs-code)
- [WSL and Windows](#wsl-and-windows)
- [macOS](#macos)
- [How it works](#how-it-works)
- [Security and privacy](#security-and-privacy)
- [Limits](#limits)
- [Troubleshooting](#troubleshooting)
- [Uninstall](#uninstall)
- [Build from source](#build-from-source)
- [Disclaimer](#disclaimer)
- [License](#license)

## Install

### Windows

In PowerShell:

```powershell
irm https://raw.githubusercontent.com/baodq97/hatrack/main/install.ps1 | iex
```

This installs `hat.exe` and `hatrack-tray.exe` to `%LOCALAPPDATA%\hatrack`, adds that folder to your
user `PATH`, starts the tray app, and adds it to `shell:startup` so it runs when you sign in to Windows.
Nothing needs admin rights.

Prefer to look first? Download `hat-windows-amd64.exe` and `hatrack-tray-windows-amd64.exe` from the
[latest release](https://github.com/baodq97/hatrack/releases/latest), rename them to `hat.exe` and
`hatrack-tray.exe`, and keep them in the same folder. Read [install.ps1](install.ps1) to see exactly what
the script does.

### macOS, Linux and WSL

```sh
curl -fsSL https://raw.githubusercontent.com/baodq97/hatrack/main/install.sh | sh
```

This puts `hat` in `~/.local/bin` (override with `HAT_INSTALL_DIR`). Builds exist for amd64 and arm64
(Apple silicon and Intel Macs).

### macOS menu bar app

Download `Hatrack-macos.zip` from the [latest release](https://github.com/baodq97/hatrack/releases/latest),
unzip it and move `Hatrack.app` to Applications. It runs on Apple silicon and Intel Macs (macOS 11
or later). The app is not notarized by Apple, so macOS blocks the first launch. Either run

```sh
xattr -d com.apple.quarantine /Applications/Hatrack.app
```

or open it once, then go to System Settings → Privacy & Security and click **Open Anyway**. (On
macOS 14 and earlier, right-click the app and choose **Open** also works.)

Or build it yourself (needs Go and the Xcode command line tools): `./build-mac-app.sh`.

The `hat` CLI from the command above works alongside the app; both use the same saved accounts.

### Verify a download

Every release lists SHA-256 sums in `checksums.txt`, and each binary has a signed build provenance
attestation that proves it was built by this repository's release workflow:

```sh
gh attestation verify hat-linux-amd64 --repo baodq97/hatrack
```

The macOS app ships as `Hatrack-macos.zip`, with its own `Hatrack-macos.zip.sha256` and attestation.

### Update

Run the install command again. On Windows it stops the running tray app first, since Windows locks a
running exe. For the macOS app, quit Hatrack from the menu bar and replace `Hatrack.app` with the one
from the new release.

## Getting started

1. Sign in to Claude Code as usual, if you have not already.
2. Keep that login: `hat save` (tray: **Save current account**).
3. Add each other account: `hat add` (tray: **Add account…**). See [Adding an account](#adding-an-account).
4. Switch: `hat use <name>` (tray: click the account).

Profile names default to the account's email. Pass a name to pick your own: `hat save personal`.

## Commands

| Command | What it does |
|---|---|
| `hat` | List saved accounts. `*` marks the one Claude uses now. |
| `hat use <name>` | Switch Claude to that account. |
| `hat add [name]` | Sign in another account. Your current login is not touched. |
| `hat save [name]` | Save the account Claude is signed in to now. |
| `hat rm <name>` | Forget a saved account. Does not sign anything out. |
| `hat version` | Print the version. |


## Windows tray app

Click the hat icon in the notification area. Windows 11 puts new icons behind the `^` overflow
arrow at first; drag it onto the taskbar to keep it visible.

```
✓ personal@gmail.com  (personal@gmail.com's Organization)
  work@company.com  (Company)
─────────────
Add account…
Save current account
Remove ▸
─────────────
Quit
```

- **Click an account** to switch to it. The check mark shows the active one; hover the icon to see it too.
- **Add account…** opens a console that walks you through signing in another account.
- **Save current account** keeps the login Claude has now.
- **Remove** forgets a saved account.
- If an action fails, the error shows at the top of the menu. The app also logs to `%USERPROFILE%\.hatrack\tray.log`.

The menu refreshes every few seconds, so changes made with `hat` show up on their own. Quitting the
tray app changes nothing: Claude keeps using the last account you picked.

## macOS menu bar app

Hatrack.app puts a hat in the menu bar, with no Dock icon:

- **Click an account** to switch to it. The check mark shows the active one.
- **Open Full Dashboard** opens a window to switch, save, rename and remove accounts. It opens as an
  app window if Chrome is installed, otherwise in your default browser.
- **Add Account…** opens Terminal and runs `hat add` there, so you can pick the browser to sign in
  with and paste the code if asked. The new account shows up in the menu when sign-in finishes.
- **Save Current Account** keeps the login Claude has now.

The dashboard is served only on `127.0.0.1`. Each launch makes a random token that only the window
Hatrack opens receives, and every action must carry it, so other web pages cannot switch or remove
accounts through it. If the dashboard says it has expired, open it again from the menu bar. To start
Hatrack when you log in, add it in System Settings → General → Login Items.

## Adding an account

`hat add` (or **Add account…** in the tray) asks how you want to sign in:

```
Add a Claude account

Sign-in happens in a separate folder; the account you use now stays signed in.

  [1] Open the sign-in page in my default browser
  [2] Show me the link - I'll open it in the right browser profile myself

Choose 1 or 2:
```

- **[1]** runs `claude auth login` as usual. Pick this if your default browser is signed in to
  claude.ai with the account you want to add.
- **[2]** prints the sign-in link instead of opening a browser. Open it in the browser profile (or a
  private window) that is signed in to the account you want, approve, then paste the code it shows
  back into the console. On Windows the code shows as `********` and its last four characters.

Do **not** add accounts with Claude's `/login`, `claude auth login`, or the Sign in button in VS Code
while another account is active: those replace the current login without saving it, so you would
have to sign that account in again. If it happens anyway, `hat save` right after keeps the new one.

## VS Code

The Claude Code extension and the CLI share one login, so `hat use` switches both. A conversation
that is already open may keep the old account until you run **Developer: Reload Window**.

- **Remote - WSL / SSH / containers:** the extension runs on the remote side and uses that side's
  login. Switch with `hat` inside WSL (or on the remote), not with the Windows tray.
- **`CLAUDE_CONFIG_DIR` in `claudeCode.environmentVariables`:** the extension then reads that folder.
  `hat` only sees `CLAUDE_CONFIG_DIR` from its own environment, so set it there too.

## WSL and Windows

Windows and each WSL distro have their own home folder, so each has its own Claude login and its own
saved accounts. Set up accounts on each side you use; switching on one side does not affect the other.

## macOS

On macOS Claude Code keeps its tokens in the login **Keychain**, not in a file, and hatrack follows
the same order Claude does:

1. **Keychain first.** Claude reads the generic password `Claude Code-credentials` (account: your
   user name). If you set `CLAUDE_CONFIG_DIR`, the item is `Claude Code-credentials-<8 hex>` instead,
   where the suffix is the start of the SHA-256 of that folder path.
2. **`~/.claude/.credentials.json` only as a fallback**, when there is no Keychain item (for example
   if the Keychain was locked or unavailable when Claude signed in).

So on a Mac:

- An old `.credentials.json` next to a Keychain item is ignored by Claude, and by hatrack too: it
  saves and switches the Keychain item. Do not treat that file as your current login.
- Switching always writes the Keychain, never the file. If the login was only in the file, hatrack
  moves it into the Keychain and removes the file, the same thing Claude does when its Keychain
  write succeeds. hatrack never falls back to writing tokens to a plain file on macOS; if the
  Keychain write fails, the switch fails and says so.
- `hat add` signs in with a scratch `CLAUDE_CONFIG_DIR`, so Claude puts the new tokens in their own
  Keychain item. hatrack copies them out and deletes that scratch item.
- hatrack calls `/usr/bin/security`, the same tool Claude uses, so macOS does not ask for Keychain
  access. If a prompt does appear, it names `security`; allow it.
- Saved logins in `~/.hatrack` are files (mode `0600`), not Keychain items.

## How it works

A Claude Code login is two things:

- the OAuth tokens: `~/.claude/.credentials.json`, or on macOS the Keychain (see [macOS](#macos)).
- the `oauthAccount` block in `~/.claude.json`: who the tokens belong to.

(Both live under `$CLAUDE_CONFIG_DIR` if you set it.)

hatrack keeps a copy of both for each account in `~/.hatrack/claude/<name>/` and swaps them in.
Everything else, like settings, history, projects, MCP servers and plugins, is shared between accounts
and never touched. Only the `oauthAccount` key of `~/.claude.json` changes; every other key is written
back as it was.

- **Switching** first writes the current login back to its saved copy, because Claude refreshes tokens
  while you work and an old copy would stop working. A login that was never saved is saved under its
  email instead of being lost.
- **Adding** runs `claude auth login` with `CLAUDE_CONFIG_DIR` pointing at a scratch folder, saves the
  result, and deletes the folder. Your active login is never replaced.
- **One copy per account:** saving the same account under a second name is refused, since only one
  copy can stay current.
- Every write goes to a temporary file first and is renamed into place, so a crash never leaves half a
  login file.
- `CLAUDE_SECURESTORAGE_CONFIG_DIR`, which moves only where Claude keeps tokens, is honored too, and is
  left out of the environment of the sign-in that `hat add` runs.

## Security and privacy

- hatrack makes **no network requests**. Tokens never leave your machine; only Claude itself talks to Anthropic.
- Saved logins are plain files, like Claude's own `.credentials.json`, readable only by your user
  (mode `0600` on Linux and macOS; your user profile's permissions on Windows). Anyone who can read
  your home folder can read them, just as with Claude's own login. On macOS this means a saved login
  is less protected than the active one, which Claude keeps in the Keychain.
- hatrack only reads and writes Claude's login files (on macOS, Claude's Keychain items) and
  `~/.hatrack`. It does not touch Claude's traffic, telemetry or anything else.
- No telemetry, no accounts, no server. The macOS dashboard is a local page on `127.0.0.1` that
  needs a per-launch token, loads nothing from the internet, and stops when you quit Hatrack.

Anthropic's terms do not allow sharing one subscription between people. hatrack is for switching
between accounts that are yours. Found a security problem? See [SECURITY.md](SECURITY.md).

## Limits

- **Claude Code only** for now. Codex and other tools may come later.
- **macOS app is not notarized**, so Gatekeeper asks before the first launch.
- **Open sessions** may keep using the previous account until restarted.
- **Saved logins are not refreshed** while not in use. One left unused long enough expires; add it again with `hat add`.
- Claude Code may change how it stores logins, which can break hatrack until it is updated.

## Troubleshooting

| Problem | Fix |
|---|---|
| `claude is not signed in` | Sign in to Claude Code first, then `hat save`. |
| `profile "x" already holds …` | That name belongs to another account. Pick another name. |
| `… is already saved as "x"` | That account is saved under another name. Use that one, or `hat rm` it first. |
| A switched account says it is logged out | Its saved tokens expired while unused, or another tool signed it out. `hat add` it again. |
| VS Code still shows the old account | Run **Developer: Reload Window**. |
| No tray icon | Check the `^` overflow arrow. If the app is not running, start `hatrack-tray.exe` and check `%USERPROFILE%\.hatrack\tray.log`. |
| `keychain …` error on macOS | Unlock the login Keychain (Keychain Access, or `security unlock-keychain`) and retry. Over SSH the Keychain is usually locked. |
| `hat` not found right after installing on Windows | Open a new terminal so it picks up the updated `PATH`. |
| Windows SmartScreen warns about the exe | The binaries are not code-signed. Verify them as shown in [Verify a download](#verify-a-download). |

## Uninstall

hatrack never changes Claude's login on its own, so after uninstalling, Claude keeps using whichever
account was active.

**Windows:** quit the tray app, then in PowerShell:

```powershell
Remove-Item "$([Environment]::GetFolderPath('Startup'))\hatrack.lnk", "$env:LOCALAPPDATA\hatrack" -Recurse
Remove-Item "$HOME\.hatrack" -Recurse   # your saved logins
```

and remove `%LOCALAPPDATA%\hatrack` from your user `PATH` (Settings → System → About → Advanced system
settings → Environment Variables).

**macOS, Linux and WSL:** on macOS, quit Hatrack from the menu bar first.

```sh
rm ~/.local/bin/hat
rm -rf /Applications/Hatrack.app ~/Library/Caches/hatrack   # macOS app
rm -r ~/.hatrack   # your saved logins
```

## Build from source

Needs Go (see `go.mod` for the version). No C compiler is needed, except for the macOS app.

```sh
go test ./...
go build ./cmd/hat                                                         # CLI for this OS
GOOS=windows go build -ldflags -H=windowsgui -o hatrack-tray.exe ./cmd/hatrack-tray  # tray, from any OS
./build-mac-app.sh                                                         # Hatrack.app, on a Mac
```

Layout:

```
cmd/hat/            CLI
cmd/hatrack-tray/   Windows tray app
cmd/hatrack-gui/    macOS menu bar app
internal/claude/    reading, saving and swapping Claude logins
internal/gui/       the macOS dashboard: a local page and its API
internal/tray/      the macOS menu bar
```

Releases are built by [.github/workflows/release.yml](.github/workflows/release.yml) when a `v*` tag is
pushed. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Disclaimer

hatrack is an independent, unofficial project. It is **not affiliated with, endorsed by, or supported
by Anthropic**. "Claude", "Claude Code" and related names are trademarks of Anthropic, PBC, used here
only to say what hatrack works with.

hatrack works by reading and rewriting Claude Code's local login files. Their format is not a public
interface and can change without notice. The software is provided **"as is", without warranty of any
kind**, as set out in sections 7 and 8 of the [license](LICENSE). In particular:

- **Use it at your own risk.** A bug, or a change in Claude Code, could sign you out or make a saved
  login unusable. The worst case should be signing in again, but back up `~/.claude/.credentials.json`
  and `~/.claude.json` if losing them would hurt.
- **You are responsible for how you use your accounts.** Using hatrack must comply with
  [Anthropic's terms](https://www.anthropic.com/legal) and your organization's policies. Do not use it
  to share a subscription between people or to get around usage limits or other restrictions.
- The authors are not liable for lost access, suspended accounts, lost work, or any other damage
  arising from using hatrack.

## License

Copyright 2026 baodq97. Licensed under the [Apache License, Version 2.0](LICENSE).
