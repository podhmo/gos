# SPDX-FileCopyrightText: 2023-present podhmo <ababjam61+github@gmail.com>
#
# SPDX-License-Identifier: MIT
"""
simplify openapi doc
"""
import typing as t
from dictknife import loading


def simplify(src: t.Optional[str], *, format: t.Literal["yaml", "json"], output: t.Optional[str]) -> None:
    src = loading.loadfile(src)
    dst = _simplify_transform(src)
    loading.dumpfile(dst, output)


def _simplify_transform(doc: t.Dict[str, t.Any]) -> t.Dict[str, t.Any]:
    for path, path_item in (doc.get("paths") or {}).items():
        toplevel_parameters = path_item.pop("parameters", None)
        if toplevel_parameters is None:
            continue
        for method, op in path_item.items():
            if method.upper() not in ("GET", "PUT", "POST", "DELETE", "OPTIONS", "HEAD", "PATCH", "TRACE"):
                continue
            if "parameters" not in op:
                op["parameters"] = toplevel_parameters[:]
            else:
                op["parameters"].extend(toplevel_parameters[:])
    return doc


def main(argv: t.Optional[t.List[str]] = None) -> t.Any:
    import argparse
    import os
    import sys

    parser = argparse.ArgumentParser(
        prog=simplify.__name__,
        description=simplify.__doc__,
        formatter_class=type(
            '_HelpFormatter', (argparse.ArgumentDefaultsHelpFormatter, argparse.RawTextHelpFormatter), {}
        ),
        epilog=sys.modules[__name__].__doc__,
    )
    parser.print_usage = parser.print_help  # type: ignore
    parser.add_argument('src', help='source file')
    parser.add_argument("--format", choices=["yaml", "json"], default="json")
    parser.add_argument(
        "-o",
        "--output",
    )
    args = parser.parse_args(argv)
    params = vars(args).copy()
    action = simplify
    if bool(os.getenv("FAKE_CALL")):
        from inspect import getcallargs
        from functools import partial

        action = partial(getcallargs, action)  # type: ignore
    return action(**params)


if __name__ == '__main__':
    main()
