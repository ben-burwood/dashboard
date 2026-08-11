"""Turn a validated :class:`~dashboard_builder.config.Config` into static HTML"""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path

from jinja2 import Environment, PackageLoader, select_autoescape

from .config import Config, Service, Tag


@dataclass
class ServiceView:
    """A service with its tag names already resolved to :class:`Tag` objects."""

    service: Service
    tags: list[Tag]


@dataclass
class GroupView:
    title: str
    services: list[ServiceView]


def _resolve_tags(service: Service, lookup: dict[str, Tag]) -> list[Tag]:
    """Resolve a service's tag names to Tag objects, dropping unknown names"""
    resolved: list[Tag] = []
    for name in service.tags:
        tag = lookup.get(name)
        if tag is not None:
            resolved.append(tag)
    return resolved


def group_services(config: Config) -> list[GroupView]:
    """Group services by ``group`` preserving first-occurrence order"""
    lookup = {tag.name: tag for tag in config.tags}
    groups: dict[str, GroupView] = {}
    for service in config.services:
        view = groups.get(service.group)
        if view is None:
            view = GroupView(title=service.group, services=[])
            groups[service.group] = view
        view.services.append(ServiceView(service, _resolve_tags(service, lookup)))
    return list(groups.values())


def favicon_href(value: str) -> str:
    """Resolve the configured page icon to an ``href`` for ``<link rel="icon">``.

    Accepts, in order of detection:
    - a URL, absolute/relative path, or ``data:`` URI -> used verbatim;
    - an Iconify name (``prefix:name``, as used for service icons) -> served
      self-contained from the Iconify SVG API (no local file needed);
    - anything else -> treated as a filename relative to the page.

    Returns ``""`` when unset (no favicon link is emitted).
    """
    value = (value or "").strip()
    if not value:
        return ""
    if "://" in value or value.startswith(("/", "./", "../", "data:")):
        return value
    if ":" in value:
        prefix, _, name = value.partition(":")
        return f"https://api.iconify.design/{prefix}/{name}.svg"
    return value


def _environment() -> Environment:
    return Environment(
        loader=PackageLoader("dashboard_builder", "templates"),
        autoescape=select_autoescape(["html", "xml", "html.j2"]),
    )


def render(config: Config) -> str:
    """Render the dashboard page to an HTML string."""
    env = _environment()
    template = env.get_template("index.html.j2")
    return template.render(
        title=config.title,
        groups=group_services(config),
        favicon=favicon_href(config.favicon),
    )


def build(config: Config, out_path: str | Path) -> Path:
    """Render the page and write it to ``out_path``"""
    out = Path(out_path)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(render(config), encoding="utf-8")
    return out
