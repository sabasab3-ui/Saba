#!/usr/bin/env bash
set -euo pipefail

echo "SABA OpenAI Agents installer"
echo "Python: $(python --version 2>&1)"

python -m pip install --upgrade pip
python -m pip install -r requirements.txt

echo
echo "Installation complete."
echo "Next:"
echo "  cp .env.example .env"
echo "  nano .env"
echo "  ./run.sh"
