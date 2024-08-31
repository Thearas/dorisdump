#!/bin/bash

set -ex

antlr4 -Dlanguage=Go -package parser *.g4
