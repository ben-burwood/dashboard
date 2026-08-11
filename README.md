# Dashboard

Simple Dashboard (<3MB RAM) written in go/vue.

Inspired stylistically by [mafl](https://github.com/hywax/mafl) - great project but I wanted something much simpler and smaller.

## Stack

The tech stack is inpired by [gatus](https://github.com/TwiN/gatus) - another reat project with a simple go configuration file approach to Uptime Monitoring.

- go Http Server
- VueJS
    - Tailwind CSS
    - Iconfiy Icons

## Deployment

Container Deployed to Github Container Repository: `ghcr.io/ben-burwood/dashboard`

### Docker Compose

```yml
  dashboard:
    image: ghcr.io/ben-burwood/dashboard:latest
    ports:
      - 8080:8080
    environment:
      - DASHBOARD_CONFIG_PATH=/config
    volumes:
      - ./config:/config:ro
```

## Configuration

The Configuration for the Homepage is defined in the config/ Yaml Files (Configurations can be split into multiple files and will be merged).

### Configuration

| Field     | Type                    | YAML Key      | Description                       |
|-----------|------------------------|---------------|-----------------------------------|
| Title     | string                 | `title`       | Title of the dashboard            |
| Tags      | []string (tag.Tags)    | `tags`        | List of tags                      |
| Services  | []*Service             | `services`    | List of service objects           |
| Favicon   | string                 | `favicon`     | Page icon (optional)              |

The `favicon` accepts an Iconify name (e.g. `mdi:view-dashboard`, served from the Iconify SVG API), or a URL / path / `data:` URI. 


### Tags

Custom Tags can be defined (see example). The color for each tag must be one of: `primary`, `secondary`, `accent`, `info`, `success`, `warning`, `error`, `neutral`.

### Services

| Field       | Type         | YAML Key      | Description                       |
|-------------|--------------|---------------|-----------------------------------|
| Link        | string       | `link`        | Service link (URL)                |
| Icon        | Icon         | `icon`        | Icon - iconify (name, color)      |
| Title       | string       | `title`       | Service title                     |
| Description | string       | `description` | Service description (optional)    |
| Group       | Group        | `group`       | Service group (optional)          |
| Tags        | []string     | `tags`        | List of tags (optional)           |

## Static build (Python)

As an alternative to running the Go server + Vue frontend, a small Python
builder ([`dashboard-builder`](pyproject.toml)) reads the **same** `config/`
YAML and renders a single self-contained `build/index.html` you can open
directly (`file://`) or drop on any static host. It reproduces the app's look
(icons via the Iconify CDN, dark-mode toggle), so no server or runtime fetch is
needed. This lives alongside the Go/Vue app and changes nothing about it.

Run it with [uv](https://docs.astral.sh/uv/):

```sh
# from the repo (config path defaults to config/config.yml)
uvx --from . dashboard-builder                          # -> build/index.html
# pass a config path explicitly, and/or choose the output file
uv run dashboard-builder config/config.yml --out build/index.html
```

To rebuild after editing the config, run the command again. Tests:
`uv run --group dev pytest`.
