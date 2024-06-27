#!/bin/bash
#
# Updates the bundled prebuilt librdkafka libraries to specified version.
#

set -e


usage() {
    echo "Usage: $0 librdkafka-static-bundle-<VERSION>.tgz"
    echo ""
    echo "This tool must be run from the TOPDIR/kafka/librdkafka_vendor directory"
    exit 1
}


# Parse dynamic libraries from linker command line.
# Will print a list matching -lfoo and -framework X..
parse_dynlibs() {
    local libs=
    while [[ $# -gt 0 ]]; do
        if [[ $1 == -l* ]]; then
            libs="${libs} $1"
        elif [[ $1 == -framework ]]; then
            libs="${libs} $1 $2"
            shift # remove one (extra) arg
        fi
        shift # remove one arg
    done

    echo "$libs"
}

# Parse dynamic library dependecies from pkg-config file and print
# them to stdout.
parse_pc_dynlibs() {
    local pc=$1
    parse_dynlibs $(sed -n 's/^Libs: \(..*\)/\1/p' "$pc")
}

setup_build() {
    # Copies static library from the temp directory into final location,
    # extracts dynamic lib list from the pkg-config file,
    # and generates the build_..go file
    local btype=$1
    local apath=$2
    local pc=$3
    local srcinfo=$4
    local build_tag=
    local gpath="../build_${btype}.go"
    local dpath="librdkafka_${btype}.a"

    if [[ $btype =~ ^glibc_linux.*$ ]]; then
        build_tag="// +build !musl"
    elif [[ $btype =~ ^musl_linux.*$ ]]; then
        build_tag="// +build musl"
    fi

    local dynlibs=$(parse_pc_dynlibs $pc)

    echo "Copying $apath to $dpath"
    cp "$apath" "$dpath"

    echo "Generating $gpath (extra build tag: $build_tag)"

    cat >$gpath <<EOF
// +build !dynamic
$build_tag

// This file was auto-generated by librdkafka_vendor/bundle-import.sh, DO NOT EDIT.

package kafka

// #cgo CFLAGS: -DUSE_VENDORED_LIBRDKAFKA -DLIBRDKAFKA_STATICLIB
// #cgo LDFLAGS: \${SRCDIR}/librdkafka_vendor/${dpath} $dynlibs
import "C"

// LibrdkafkaLinkInfo explains how librdkafka was linked to the Go client
const LibrdkafkaLinkInfo = "static ${btype} from ${srcinfo}"
EOF

    git add "$dpath" "$gpath"

}


bundle="$1"
[[ -f $bundle ]] || usage

bundlename=$(basename "$bundle")

bdir=$(mktemp -d tmpXXXXXX)

echo "Extracting bundle $bundle:"
tar -xzvf "$bundle" -C "$bdir/"

echo "Copying librdkafka files"
for f in rdkafka.h LICENSES.txt ; do
    cp $bdir/$f . || true
    git add "$f"
done


for btype in glibc_linux_amd64 \
                 glibc_linux_arm64 \
                 musl_linux_amd64 \
                 musl_linux_arm64 \
                 darwin_amd64 \
                 darwin_arm64 \
                 windows ; do
    lib=$bdir/librdkafka_${btype}.a
    pc=${lib/%.a/.pc}
    [[ -f $lib ]] || (echo "Expected file $lib missing" ; exit 1)
    [[ -f $pc ]] || (echo "Expected file $pc missing" ; exit 1)

    setup_build $btype $lib $pc $bundlename
done

rm -rf "$bdir"

echo "All done"
