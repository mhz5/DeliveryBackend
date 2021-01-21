#!/usr/bin/env bash

# Process bash arguments: https://stackoverflow.com/questions/192249/how-do-i-parse-command-line-arguments-in-bash
#UPDATE_PROTOS=0
#while [[ $# -gt 0 ]]
#do
#    key="$1"
#    case ${key} in
#        -p|--update-protos)
#        UPDATE_PROTOS=1
#        shift
#        ;;
#    esac
#done


cd ..
# Generate BUILD files with Gazelle.
bazel run //:gazelle
# Build entire project with Bazel.
bazel build ...

if [[ $? -ne 0 ]]; then
    exit 1
fi

# Delete old compiled protos
sudo rm proto/*/*.pb.go

# Move newly compiled proto files to /proto directory, so IDE can find them.
for proto_dir in $(ls -d proto/*/)
do
    proto_pkg=$(echo ${proto_dir} | sed 's/proto//g' | sed 's/\///g')
#    TODO: Copy files without sudo
    sudo cp -r bazel-bin/proto/${proto_pkg}/${proto_pkg}_go_proto_/fda/proto/${proto_pkg}/* \
        proto/${proto_pkg}
done

