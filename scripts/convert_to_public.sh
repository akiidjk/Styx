#!/bin/bash

set -euo pipefail

GENERATED_DIR="./internal/ebpf/generated/"
TARGET_EXTENSION="*.go"

GREEN="\033[32m"
YELLOW="\033[33m"
RED="\033[31m"
RESET="\033[0m"

log() {
    echo -e "${YELLOW}[INFO]${RESET} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${RESET} $1"
}

error() {
    echo -e "${RED}[ERROR]${RESET} $1" >&2
}

if [ ! -d "$GENERATED_DIR" ] || [ -z "$(find "$GENERATED_DIR" -name "$TARGET_EXTENSION")" ]; then
    error "Nessun file generato trovato in '$GENERATED_DIR'."
    exit 1
fi

process_file() {
    local file="$1"
    log "Modifica in corso: $file"

    local functions=($(grep -oP '\bfunc \K[a-z][a-zA-Z0-9_]*' "$file"))
    local interfaces=($(grep -oP '\btype \K[a-z][a-zA-Z0-9_]*' "$file"))

    for func in "${functions[@]}"; do
        local new_func="$(tr '[:lower:]' '[:upper:]' <<< ${func:0:1})${func:1}"
        log "Rendendo pubblica la funzione: $func -> $new_func"
        sed -i "s/\b$func\b/$new_func/g" "$file"
    done

    for interface in "${interfaces[@]}"; do
        local new_interface="$(tr '[:lower:]' '[:upper:]' <<< ${interface:0:1})${interface:1}"
        log "Rendendo pubblica l'interfaccia: $interface -> $new_interface"
        sed -i "s/\b$interface\b/$new_interface/g" "$file"
    done
}

for file in $(find "$GENERATED_DIR" -name "$TARGET_EXTENSION"); do
    process_file "$file"
done

success "Tutte le funzioni generate sono state rese pubbliche con successo."
