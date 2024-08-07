#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
    echo "usage: $0 <package-name>"
    exit 1
fi

package_split=(${package//\// })
package_name=${package_split[-1]}

buildir="$PWD/build"
mkdir -p "$buildir"

platforms=("windows/amd64" "darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
        CGO_ENABLED=1
        CC=x86_64-w64-mingw32-gcc
    else
        CGO_ENABLED=0
        CC=
    fi

    env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=$CGO_ENABLED CC=$CC go build -o "$buildir"/$output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done