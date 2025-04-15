from __future__ import annotations

from pathlib import Path

from agent import REPO_ROOT
from agent.utils.config import _convert_path_to_abspath, _convert_paths_to_abspath


def test_convert_path_to_abspath():
    assert _convert_path_to_abspath("sadf") == REPO_ROOT / "sadf"
    assert _convert_path_to_abspath("/sadf") == Path("/sadf")


def test_convert_paths_to_abspath():
    assert _convert_paths_to_abspath([Path("sadf"), Path("/sadf")]) == [
        REPO_ROOT / "sadf",
        Path("/sadf"),
    ]
