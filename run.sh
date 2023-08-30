#!/bin/bash

# Display a help message
function show_help {
    cat << EOF
Usage: $0 [OPTION]
Run the server in different modes.

Options:
    --test                 Run the server in test mode.
    --fast                 Run tests without bringing up all services.
    --generate-goldens     Re-generate golden files when tests don't match.

No option will run the server in normal mode.
EOF
}

# Check for no arguments and run the server in normal mode
if [ "$#" -eq 0 ]; then
    docker compose up --build --remove-orphans -d
    exit 0
fi

# Check passed arguments
while [ "$#" -gt 0 ]; do
    case "$1" in
        --test)
            docker compose -f docker-compose.test.yml up --build --remove-orphans --abort-on-container-exit
            shift
            ;;
        --fast)
            docker compose -f docker-compose.test.yml up -d test_db --build --remove-orphans
            go test tests/highlevel_test.go -v
            docker compose -f docker-compose.test.yml down
            shift
            ;;
        --generate-goldens)
            # Implement logic to re-generate golden files here. Placeholder for now.
            echo "Re-generating golden files..."
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done
