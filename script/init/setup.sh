#!/usr/bin/env bash


echo "Starting setup for all package"
is_formula_installed() {
    brew list "$1" &>/dev/null
}

install_with_brew() {
    brew install "$1"
}

if ! command -v brew &>/dev/null; then
    echo "Homebrew is not installed. Please install it first: https://brew.sh/"
    exit 1
fi

while IFS= read -r formula; do
    if ! is_formula_installed "$formula"; then
        read -rp "The formula '$formula' is not installed. Do you want to install it? (y/n): " answer
        if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
            install_with_brew "$formula"
        fi
    fi
done < formulas.txt

if ! is_formula_installed "ruby"; then
    install_with_brew "ruby"
fi

if  is_formula_installed "go"; then
    go install github.com/swaggo/swag/cmd/swag@latest
fi

while IFS= read -r gem; do
    gem install "$gem"
done < rubygem.txt
