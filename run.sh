#!/bin/bash

if [ -n "$VERBOSE" ]; then
  set -x
fi

# Store the process IDs
declare -a pids

# Kill all background processes and their children when the script exits
cleanup() {
  echo "🛑 Cleaning up processes..."
  # Kill all background processes and their children
  for pid in "${pids[@]}"; do
    pkill -P "$pid" 2>/dev/null
    kill "$pid" 2>/dev/null
  done
}

# Set up trap for various signals
trap cleanup EXIT

autobuild() {
  # Run the file watcher in background
  if ! command -v watchexec >/dev/null 2>&1; then
    echo "Error: watchexec not found."
    echo "Please install watchexec:"
    echo "  Mac: brew install watchexec"
    echo "  Linux: cargo install watchexec-cli"
    echo "  or download from: https://github.com/watchexec/watchexec/releases"
    exit 1
  fi

  echo "👀 Watching for changes..."
  # Watch for changes with pretty output
  watchexec \
    --on-busy-update=do-nothing \
    --quiet \
    -w resources/ \
    -w persist/ \
    -w lib/ \
    -w game \
    -- \
    'echo "🔄 Rebuilding..." && ./build.sh'
}

serve() {
  python3 -m http.server 8080 --directory out/
}

main() {
  serve &
  pids+=($!)

  if [ -n "$AUTOBUILD" ]; then
    autobuild &
    pids+=($!)
  fi

  wait
}

main
