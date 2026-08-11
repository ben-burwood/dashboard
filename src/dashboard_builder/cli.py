"""Command-line entry point for the dashboard static-site builder."""

from __future__ import annotations

import argparse
import sys

from .builder import build
from .config import ConfigError, load_and_validate

DEFAULT_OUT = "build/index.html"


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        prog="dashboard-builder",
        description=(
            "Build a self-contained static dashboard page from the same YAML "
            "config as the Go/Vue app."
        ),
    )
    parser.add_argument(
        "config",
        nargs="?",
        help=(
            "Path to a config file or directory of .yml/.yaml files. "
            "Falls back to $DASHBOARD_CONFIG_PATH, then config/config.yml."
        ),
    )
    parser.add_argument(
        "-o",
        "--out",
        default=DEFAULT_OUT,
        help=f"Output HTML file path (default: {DEFAULT_OUT}).",
    )
    args = parser.parse_args(argv)

    try:
        config = load_and_validate(args.config)
    except ConfigError as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 1

    out = build(config, args.out)
    print(f"Wrote {out}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
