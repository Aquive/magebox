# Expose / Share

MageBox can expose your local project to the internet via [Cloudflare Tunnels](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/), giving you a temporary public URL. No Cloudflare account required.

## Use Cases

- **Client demos** - Share a working preview without deploying
- **Webhook testing** - Receive callbacks from payment providers, shipping APIs, etc.
- **Mobile testing** - Test on real devices over the network
- **Collaborative debugging** - Let a colleague see the issue in their browser

## Requirements

Install `cloudflared`:

```bash
brew install cloudflared
```

## Usage

### Start a tunnel

```bash
magebox expose
```

This will:
1. Start a Cloudflare quick tunnel
2. Add the tunnel domain to `.magebox.yaml` and reload nginx
3. Update all Magento base URLs (base, media, static) across all scopes
4. Flush Magento cache

Once ready, the public URL is displayed:

```
Public URL: https://random-words-here.trycloudflare.com
```

### Expose a specific domain

If your project has multiple domains, specify which one to expose:

```bash
magebox expose store.mystore.test
```

### Stop the tunnel

Press **Ctrl+C** in the terminal where expose is running, or from another terminal:

```bash
magebox expose stop
```

This reverts all changes:
- Restores original base URLs in the database
- Restores `env.php` from backup
- Removes the tunnel domain from `.magebox.yaml`
- Regenerates nginx vhosts and reloads nginx
- Flushes Magento cache

### Check tunnel status

```bash
magebox expose status
```

## How It Works

```
Browser → Cloudflare Edge → cloudflared → nginx (HTTP) → PHP-FPM
                                            ↑
                              Tunnel hostname in .magebox.yaml
                              so nginx routes it correctly
```

1. `cloudflared` creates a quick tunnel and provides a `*.trycloudflare.com` URL
2. The tunnel hostname is added to `.magebox.yaml` as a domain (SSL disabled, Cloudflare handles SSL)
3. Nginx vhosts are regenerated so the tunnel hostname is accepted as a `server_name`
4. Magento base URLs are updated in `core_config_data` (all scopes) and locked in `env.php` so they override any `config.php` settings
5. `app:config:import` runs to sync the config changes
6. Cache is flushed so pages render with the tunnel URLs

## What Gets Modified

| Resource | On expose | On stop |
|----------|-----------|---------|
| `.magebox.yaml` | Tunnel domain added | Tunnel domain removed |
| `core_config_data` | All base URLs updated | Restored from backup |
| `app/etc/env.php` | Base URLs locked | Restored from backup |
| nginx vhosts | Regenerated | Regenerated |

All changes are fully reverted on stop or Ctrl+C.

## Limitations

- **Random URLs** - Quick tunnels generate random `*.trycloudflare.com` hostnames that change each session. Custom domains require a Cloudflare account and named tunnels.
- **No SLA** - Quick tunnels are best-effort, not for production use
- **Rate limits** - Maximum 200 concurrent requests through the tunnel

## Troubleshooting

### cloudflared not found

```bash
brew install cloudflared
```

### Redirect loop

If the tunnel URL redirects in a loop, there may be stale tunnel URLs from a previous session. Run:

```bash
magebox expose stop
```

This cleans up any leftover state even if no tunnel is running.

### Media/static assets showing localhost

Ensure the tunnel was started after the project's services are running (`magebox start` first). The expose command updates all 6 base URL paths (base, media, static x secure/unsecure) across all database scopes and locks them in `env.php`.

### Stale tunnel domain in .magebox.yaml

If a previous session wasn't cleanly stopped, stale `*.trycloudflare.com` entries may remain in `.magebox.yaml`. Running `magebox expose` again automatically cleans these up before adding the new tunnel domain.
