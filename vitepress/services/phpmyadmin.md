# phpMyAdmin

MageBox supports [phpMyAdmin](https://www.phpmyadmin.net/), a web-based administration tool for MySQL and MariaDB databases.

## Overview

phpMyAdmin provides a visual interface for:

- **Browsing databases**, tables, and records
- **Running SQL queries** directly in the browser
- **Importing and exporting** data (SQL, CSV, etc.)
- **Managing users**, permissions, and privileges
- **Viewing and editing** table structure and data

## Enabling phpMyAdmin

### Via Command

```bash
magebox phpmyadmin enable
```

### Via Global Config

```bash
magebox config set phpmyadmin true
```

Then start services:

```bash
magebox global start
```

### Via Project Config

Add to your `.magebox.yaml`:

```yaml
services:
  mysql: "8.0"
  phpmyadmin: true
```

To use a custom port:

```yaml
services:
  mysql: "8.0"
  phpmyadmin:
    port: 9090
```

## Connection Details

| Setting | Value |
|---------|-------|
| Web UI | `http://localhost:8036` |
| Username | `root` |
| Password | `magebox` |

## Getting Started

1. Enable phpMyAdmin:
   ```bash
   magebox phpmyadmin enable
   ```

2. Open **http://localhost:8036** in your browser

3. phpMyAdmin auto-connects to your default MySQL/MariaDB instance with root credentials

## Connecting to Multiple Databases

phpMyAdmin runs with **arbitrary server mode** enabled, which means you can connect to any MySQL or MariaDB instance from the login screen. Use the container name as the server host:

| Database | Server Host |
|----------|-------------|
| MySQL 5.7 | `magebox-mysql-5.7` |
| MySQL 8.0 | `magebox-mysql-8.0` |
| MySQL 8.4 | `magebox-mysql-8.4` |
| MariaDB 10.4 | `magebox-mariadb-10.4` |
| MariaDB 10.6 | `magebox-mariadb-10.6` |
| MariaDB 11.4 | `magebox-mariadb-11.4` |

## MageBox Commands

### Open in Browser

```bash
magebox phpmyadmin open
```

Opens the phpMyAdmin web UI in your default browser.

### Check Status

```bash
magebox phpmyadmin status
```

Shows whether phpMyAdmin is enabled and running.

### Disable

```bash
magebox phpmyadmin disable
```

Stops the container and removes it from docker-compose.

## Docker Container

phpMyAdmin runs as a Docker container (`magebox-phpmyadmin`) with:

- **Image**: `phpmyadmin:latest`
- **Port**: 8036 (Web UI)
- **Network**: magebox

## Troubleshooting

### Cannot Connect to Database

Ensure your database service is running:

```bash
magebox global status
```

If not, start services:

```bash
magebox global start
```

### Port 8036 Already in Use

If another service is using port 8036, stop it before enabling phpMyAdmin:

```bash
lsof -i :8036
```
