from __future__ import annotations

import os
from dataclasses import dataclass, field
from pathlib import Path

import yaml

DEFAULT_GROUP = "default"
DEFAULT_ICON_COLOR = "black"

COLOR_VARIANTS = {
    "primary",
    "secondary",
    "accent",
    "info",
    "success",
    "warning",
    "error",
    "neutral",
}

DEFAULT_CONFIG_PATH = "config/config.yml"
CONFIG_PATH_ENV_VAR = "DASHBOARD_CONFIG_PATH"

_DOLLAR_SENTINEL = "__DASHBOARD_LITERAL_DOLLAR_SIGN__"


class ConfigError(Exception):
    """Raised when a configuration is missing, unreadable, or invalid."""


@dataclass
class Icon:
    name: str = ""
    color: str = DEFAULT_ICON_COLOR

    @classmethod
    def from_dict(cls, data: dict | None) -> "Icon":
        data = data or {}
        color = data.get("color") or DEFAULT_ICON_COLOR
        return cls(name=data.get("name", ""), color=color)


@dataclass
class Tag:
    name: str
    color: str

    @classmethod
    def from_dict(cls, data: dict) -> "Tag":
        return cls(name=data.get("name", ""), color=data.get("color", ""))


@dataclass
class Service:
    link: str
    title: str
    icon: Icon = field(default_factory=Icon)
    group: str = DEFAULT_GROUP
    tags: list[str] = field(default_factory=list)

    @classmethod
    def from_dict(cls, data: dict) -> "Service":
        # Unknown keys (e.g. the legacy "description") are ignored.
        return cls(
            link=data.get("link", ""),
            title=data.get("title", ""),
            icon=Icon.from_dict(data.get("icon")),
            group=data.get("group") or DEFAULT_GROUP,
            tags=list(data.get("tags") or []),
        )


@dataclass
class Config:
    services: list[Service] = field(default_factory=list)
    title: str = ""
    tags: list[Tag] = field(default_factory=list)
    favicon: str = ""

    @classmethod
    def from_dict(cls, data: dict) -> "Config":
        data = data or {}
        return cls(
            title=data.get("title", "") or "",
            tags=[Tag.from_dict(t) for t in (data.get("tags") or [])],
            services=[Service.from_dict(s) for s in (data.get("services") or [])],
            favicon=data.get("favicon", "") or "",
        )


def _expand_env(text: str) -> str:
    """Expand ``$VAR``/``${VAR}`` while preserving literal ``$$`` as ``$``"""
    text = text.replace("$$", _DOLLAR_SENTINEL)
    text = os.path.expandvars(text)
    return text.replace(_DOLLAR_SENTINEL, "$")


def _resolve_path(config_path: str | None) -> Path:
    """Resolve the config path: explicit arg -> env var -> default"""
    for candidate in (config_path, os.getenv(CONFIG_PATH_ENV_VAR), DEFAULT_CONFIG_PATH):
        if candidate:
            p = Path(candidate)
            if p.exists():
                return p
    raise ConfigError("configuration file not found")


def _read_yaml_text(path: Path) -> str:
    """Read raw YAML text from a file or a directory of .yml/.yaml files.
    Directory files are concatenated (title last-wins, services/tags accumulate)
    """
    if path.is_dir():
        parts: list[str] = []
        for child in sorted(path.iterdir()):
            if ".." in child.name:
                continue
            if child.suffix in (".yml", ".yaml") and child.is_file():
                parts.append(child.read_text(encoding="utf-8"))
        return "\n".join(parts)
    return path.read_text(encoding="utf-8")


def _merge_docs(text: str) -> dict:
    """Parse one-or-more concatenated YAML docs and merge them.
    ``title`` last-wins; ``tags`` and ``services`` lists accumulate.
    """
    merged: dict = {}
    for doc in yaml.safe_load_all(text):
        if doc is None:
            continue
        if not isinstance(doc, dict):
            raise ConfigError("configuration must be a mapping")
        for key, value in doc.items():
            if key in ("tags", "services") and isinstance(value, list):
                merged.setdefault(key, [])
                merged[key].extend(value)
            else:
                merged[key] = value
    return merged


def load(config_path: str | None = None) -> Config:
    """Load and construct a :class:`Config` (without validating it)."""
    path = _resolve_path(config_path)
    text = _read_yaml_text(path)
    text = _expand_env(text)
    if not text.strip():
        raise ConfigError("configuration file not found")
    data = _merge_docs(text)
    return Config.from_dict(data)


def validate(config: Config) -> None:
    """Validate a config.
    Raises :class:`ConfigError` on the first problem found.
    """
    if not config.services:
        raise ConfigError("configuration should contain at least one Service")

    for service in config.services:
        if not service.link:
            raise ConfigError(
                f"invalid service configuration for {service.title!r}: missing link"
            )
        if not service.title:
            raise ConfigError("invalid service configuration: missing title")
        seen: set[str] = set()
        for tag_name in service.tags:
            if tag_name in seen:
                raise ConfigError(
                    f"invalid service configuration for {service.title!r}: "
                    f"duplicate tag {tag_name!r}"
                )
            seen.add(tag_name)

    tag_names: set[str] = set()
    for tag in config.tags:
        if tag.name in tag_names:
            raise ConfigError(f"invalid tag configuration: duplicate tag {tag.name!r}")
        tag_names.add(tag.name)
        if tag.color not in COLOR_VARIANTS:
            raise ConfigError(
                f"invalid tag color variant for tag {tag.name}: {tag.color}"
            )


def load_and_validate(config_path: str | None = None) -> Config:
    config = load(config_path)
    validate(config)
    return config
