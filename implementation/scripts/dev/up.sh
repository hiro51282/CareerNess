#!/bin/bash
set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
IMPL_ROOT="$REPO_ROOT/implementation"

echo "🚀 Starting CareerNess development servers..."
echo ""

# Start Go API server
echo "📡 Starting Go API server (:8080)..."
(cd "$IMPL_ROOT/apps/api" && go run ./cmd/server) &
API_PID=$!

# Start React dev server
echo "⚛️  Starting React dev server (:5173)..."
(cd "$IMPL_ROOT" && ~/.local/bin/pnpm dev:web) &
WEB_PID=$!

echo ""
echo "✅ Both servers running:"
echo "   • API  → http://localhost:8080"
echo "   • Web  → http://localhost:5173"
echo ""
echo "Press Ctrl+C to stop all servers"
echo ""

# Handle SIGINT to stop both processes
trap "kill $API_PID $WEB_PID 2>/dev/null; exit 0" SIGINT

# Wait for both processes
wait
