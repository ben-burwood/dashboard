# Dashboard

Simple Dashboard (<3MB RAM) written in go/vue.

Inspired stylistically by [mafl](https://github.com/hywax/mafl) - great project but I wanted something much simpler and smaller.

## Stack

The tech stack is inpired by [gatus](https://github.com/TwiN/gatus) - another reat project with a simple go configuration file approach to Uptime Monitoring.

- go Http Server
- VueJS
    - Tailwind CSS
    - DaisyUI Component Library
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


### Tags

Custom Tags can be defined (see example). The color for these follows the DaisyUI CSS Color names (https://daisyui.com/docs/colors)

### Services

| Field       | Type         | YAML Key      | Description                       |
|-------------|--------------|---------------|-----------------------------------|
| Link        | string       | `link`        | Service link (URL)                |
| Icon        | Icon         | `icon`        | Icon - iconify (name, color)      |
| Title       | string       | `title`       | Service title                     |
| Description | string       | `description` | Service description (optional)    |
| Group       | Group        | `group`       | Service group (optional)          |
| Tags        | []string     | `tags`        | List of tags (optional)           |
