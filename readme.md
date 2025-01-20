![Snaptrack Logo](./logo.png)

# Snaptrack

> **Snaptrack** is a site-snapshot and change-tracking tool. It captures the HTML of any given site (or set of pages), stores snapshots in a local SQLite database, and highlights differences between consecutive snapshots.

## Table of Contents

1. [Overview](#overview)  
2. [Features](#features)  
3. [Getting Started](#getting-started)  
   1. [Prerequisites](#prerequisites)  
   2. [Installation](#installation)  
   3. [Building from Source](#building-from-source)  
4. [Usage](#usage)  
   1. [Crawling a Domain](#crawling-a-domain)  
   2. [Checking or Diffing a Single URL](#checking-or-diffing-a-single-url)  
   3. [TUI (Interactive Mode)](#tui-interactive-mode)  
5. [Use Cases](#use-cases)  
   1. [Site Owners & Content Managers](#site-owners--content-managers)  
   2. [Security Teams & Professionals](#security-teams--professionals)  
   3. [Testers & QA Engineers](#testers--qa-engineers)  
   4. [SEO & Marketing Analysts](#seo--marketing-analysts)  
6. [Configuration & Customization](#configuration--customization)  
7. [Contributing](#contributing)  
8. [License](#license)

---

## Overview

**Snaptrack** is a Go-based application designed to **monitor websites** for changes over time. It does this by:

1. **Fetching** the raw HTML from a URL (or recursively crawling an entire domain).  
2. **Storing** each snapshot in a local SQLite database.  
3. **Comparing** each new snapshot to the previous version for that page and presenting a **unified diff** of what changed.

The tool can be run in **CLI mode** for batch usage (`crawl`, `check`, etc.) or in a **TUI (Text-based User Interface)** for interactive exploration of snapshots.

---

## Features

- **Recursive Crawl**: Optionally follow links within the same domain to capture snapshots of multiple pages.  
- **SQLite Database**: Stores snapshots locally—simple, lightweight, no external DB required.  
- **Diff Highlights**: Compares new HTML to the previous version for each page, generating a **unified diff** (with optional color).  
- **TUI Interface**: A text-based interface lets you browse tracked URLs, see changes, and re-check pages on demand.  
- **Raw HTTP** Approach (No Headless Browser)**: Faster, simpler for static or server-rendered pages.  
  - *(If needed, revert to a headless browser approach for JavaScript-heavy sites.)*

---

## Getting Started

### 1. Prerequisites

- **Go** (version 1.18+ recommended).  
- A **SQLite** driver (e.g., `github.com/mattn/go-sqlite3`) automatically installed via `go mod tidy`.  
- (Optional) A color-supporting terminal for color-coded diffs.

### 2. Installation

1. **Clone** this repository:
   ```bash
   git clone https://github.com/copyleftdev/snaptrack.git
   ```
2. **Change** to the directory:
   ```bash
   cd snaptrack
   ```
3. **Install** dependencies:
   ```bash
   go mod tidy
   ```

### 3. Building from Source

You can build via our **Makefile**:

```bash
make build
make run
```

Or manually:

```bash
go build -o bin/snapstack ./cmd/snapstack
```

The executable `snapstack` is placed in `./bin/`.

---

## Usage

Snaptrack can be invoked via CLI subcommands or launched in a TUI if no arguments are provided.

### 1. Crawling a Domain

```bash
./bin/snapstack crawl https://example.com --max-depth=2
```

- **Crawl** the specified domain (`example.com`) recursively up to 2 levels.  
- Store HTML snapshots in `snapshots.db`.  
- Show diff logs if changes are detected on subsequent crawls.

### 2. Checking or Diffing a Single URL

```bash
./bin/snapstack check https://example.com
```
*(If implemented—example usage. Checks a single page.)*

```bash
./bin/snapstack diff https://example.com
```
*(If implemented—example usage. Shows a unified diff for that page’s last two snapshots.)*

### 3. TUI (Interactive Mode)

```bash
./bin/snapstack
```

- Launches a text-based interface to:
  - **List** all tracked URLs in your DB.  
  - **Select** a URL to see if it changed.  
  - Press <kbd>d</kbd> for diff output, <kbd>r</kbd> to recapture, etc.  
  - Press <kbd>q</kbd> or <kbd>Esc</kbd> to quit.

---

## Use Cases

### 1. Site Owners & Content Managers
- Maintain a **historical record** of content changes over time.  
- Quickly identify any **unapproved modifications** or mistakes in text or layout.  

### 2. Security Teams & Professionals
- **Monitor** a public site for **unexpected** or **malicious** insertions (like scripts or hidden form fields).  
- **Diff** after each deployment or scheduled check to confirm the site hasn’t been tampered with.  
- Helps detect **defacement** or **backdoors** if an attacker alters the HTML.  
- Local snapshots ensure an **audit trail** for post-incident analysis.

### 3. Testers & QA Engineers
- **Compare** staging and production pages by capturing snapshots from each environment.  
- Confirm no undesired changes slipped into a new release.  
- **Record** each build’s output so you can see exactly what changed from one version to the next.  
- Useful for **UI regression** checks on static or server-rendered sites.

### 4. SEO & Marketing Analysts
- Track how **metadata**, headings, or content changes might affect SEO.  
- Keep a historical log of **keyword** or **content** modifications.

---

## Configuration & Customization

- **Database Path**: Defaults to `snapshots.db` in the current directory. Change in `main.go` or environment variables as desired.  
- **Crawl Depth & Concurrency**: `--max-depth` plus internal concurrency settings let you control the scope and speed of crawling.  
- **Timeout**: Each HTTP request uses a default of ~15 seconds. Adjust in `capture.go` if needed.  
- **Unified Diff**: We currently produce a standard “unified diff.” For color highlighting, ensure your terminal supports ANSI or integrate with external tools.

---

## Contributing

We welcome contributions! Please:

1. **Fork** this repo & create a feature branch.  
2. **Submit** a pull request when ready.  
3. **Open** an issue to discuss features, request improvements, or report bugs.

---

## License

Snaptrack is licensed under the [MIT License](./LICENSE). See the [LICENSE](./LICENSE) file for more info.

