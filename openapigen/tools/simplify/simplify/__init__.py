# SPDX-FileCopyrightText: 2023-present podhmo <ababjam61+github@gmail.com>
#
# SPDX-License-Identifier: MIT
"""
simplify openapi doc
"""
import typing as t
from collections import defaultdict
from dictknife import loading
from dictknife import DictWalker


def simplify(src: t.Optional[str], *, format: t.Literal["yaml", "json"], output: t.Optional[str]) -> None:
    src = loading.loadfile(src)
    dst = Transformer().transform(src)
    loading.dumpfile(dst, output)


class Transformer:
    def transform(self, doc: t.Dict[str, t.Any]) -> t.Dict[str, t.Any]:
        self._omit_toplevel_parameters(doc)
        self._dref_parameters(doc)
        return doc

    def _dref_parameters(self, doc: t.Dict[str, t.Any]) -> None:
        ref_used = defaultdict(list)
        for path, sd in DictWalker(["$ref"]).walk(doc):
            ref = sd["$ref"]
            if "#/components/schemas" in ref:
                continue
            ref_used[ref].append(sd)

        components = doc.get("components") or {}
        parameters = components.get("parameters")
        if parameters is not None:
            for name, definition in parameters.items():
                ref_key = f"#/components/parameters/{name}"
                for sd in ref_used.get(ref_key) or []:
                    sd.pop("$ref")
                    sd.update(definition)
            components.pop("parameters")

    def _omit_toplevel_parameters(self, doc: t.Dict[str, t.Any]) -> None:
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
