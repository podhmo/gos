# SPDX-FileCopyrightText: 2023-present podhmo <ababjam61+github@gmail.com>
#
# SPDX-License-Identifier: MIT
import typing as t
import os
import sys
import pathlib
from dictknife.jsonknife.resolver import get_resolver, ExternalFileResolver
from dictknife.jsonknife.bundler import Scanner, CachedItemAccessor
from prestring.go.codeobject import Module, goname

# todo: nested structure
# todo: use Reference instead of ReferenceName
# todo: todo document


def run(src: str) -> None:
    resolver = get_resolver(src)
    accessor = CachedItemAccessor(resolver)
    scanner = Scanner(accessor, {})
    conflicted = scanner.scan(resolver.doc)
    if conflicted:
        print("conflicted: ", conflicted)

    m = Module()
    m.package("M")

    oapicodegen = m.import_("github.com/podhmo/gos/oapicodegen")

    subresolver_dict: t.Dict[str, ExternalFileResolver] = resolver.cache
    for subresolver in subresolver_dict.values():
        components = subresolver.doc.get("components")
        if components is None:
            continue

        filename = pathlib.PurePath(subresolver.filename).name
        print(f"--{filename.ljust(20, '-')}------------------", file=sys.stderr)

        m.stmt(f"// from: {filename}")
        m.stmt("var (")
        with m.scope():
            for schema_name, schema in (components.get("schemas") or {}).items():
                if "$ref" in schema:
                    continue
                if "allOf" in schema:
                    continue
                if "oneOf" in schema:
                    continue
                if "anyOf" in schema:
                    continue

                type_name = goname(schema_name)
                underlying = schema.get("type", "object")

                if underlying == "array":
                    m.stmt(f'{type_name} = {oapicodegen}.Define("{type_name}", {to_type(schema)})')
                elif underlying == "object":
                    m.stmt(f'{type_name} = {oapicodegen}.Define("{type_name}", b.Object(')
                    with m.scope():
                        for name, prop in (schema.get("properties") or {}).items():
                            m.stmt(f'b.Field("{name}", {to_type(prop)}),')
                    m.stmt("))")
                    pass
                elif underlying == "string":
                    m.stmt(f'{type_name} = {oapicodegen}.Define("{type_name}", {to_type(schema)})')
                elif underlying == "integer":
                    m.stmt(f'{type_name} = {oapicodegen}.Define("{type_name}", {to_type(schema)})')
                elif underlying == "number":
                    m.stmt(f'{type_name} = {oapicodegen}.Define("{type_name}", {to_type(schema)})')
                elif underlying == "boolean":
                    m.stmt(f'{type_name} = {oapicodegen}.Define("{type_name}", {to_type(schema)})')
        m.stmt(")")
        m.sep()
    print(m)


def to_type(schema: t.Dict[str, t.Any]) -> str:
    ref = schema.get("$ref")
    if ref is not None:
        name = goname(ref.rsplit("/")[-1])
        return f'b.ReferenceName("{name}")'  # todo: name conflict
    underlying = schema.get("type", "object")

    # todo: map

    if underlying == "array":
        return f"b.Array({to_type(schema['items'])})"
    elif underlying == "object":
        return "b.Object(...)"
    elif underlying == "string":
        return "b.String()"
    elif underlying == "integer":
        return "b.Int()"
    elif underlying == "number":
        return "b.Float()"
    elif underlying == "boolean":
        return "b.Boolean()"


def main(argv: t.Optional[t.List[str]] = None) -> t.Any:
    import argparse

    # import logging
    # logging.basicConfig(level=logging.DEBUG)
    parser = argparse.ArgumentParser(
        prog=run.__name__,
        description=run.__doc__,
        formatter_class=type(
            '_HelpFormatter', (argparse.ArgumentDefaultsHelpFormatter, argparse.RawTextHelpFormatter), {}
        ),
    )
    parser.print_usage = parser.print_help  # type: ignore
    parser.add_argument('src', help='-')
    args = parser.parse_args(argv)
    params = vars(args).copy()
    action = run
    if bool(os.getenv("FAKE_CALL")):
        from inspect import getcallargs
        from functools import partial

        action = partial(getcallargs, action)  # type: ignore
    return action(**params)


if __name__ == '__main__':
    main()
