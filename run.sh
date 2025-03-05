#!/bin/sh
echo "Starting Tailwind CSS..."
mkdir -p public/styles
tailwindcss -i tailwind.css -o public/styles/app.css --watch &
TAILWIND_PID=$!
echo "Tailwind started with PID $TAILWIND_PID"

echo "Starting Air..."
air -c .air.toml

# This ensures Tailwind is properly terminated when Air exits
kill $TAILWIND_PID