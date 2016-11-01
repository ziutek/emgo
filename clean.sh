#!/bin/bash

set -e

rm -rf egroot/pkg/* 
rm -rf egpath/pkg/*
rm -f $(find egroot/src egpath/src -name '__*.[ch]' -print)
rm -f $(find egpath/src/*/examples -name '*.elf' -o -name '*.bin' -print)