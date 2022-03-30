#!/usr/bin/env sh

rm -f internal/cgo/*.c internal/cgo/*.h
cp libimagequant/*.c libimagequant/*.h internal/cgo/
rm -f internal/cgo/example.c