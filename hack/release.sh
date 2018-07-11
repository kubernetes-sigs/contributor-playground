#!/bin/bash

RELEASE=${RELEASE:=0.1.0}

IMAGE=cce-cloud-controller-manager

SHA=${SHA:=$(git rev-parse --short=8 HEAD)}

function do_release() {
    if git rev-parse "$RELEASE" >/dev/null 2>&1; then
        echo "Tag $RELEASE already exists. Doing nothing."
        exit 1
    fi

    echo "Creating new release $RELEASE for SHA $SHA"

    git tag -a "$RELEASE" -m "Release version: $RELEASE"
    git push --tags

    docker pull $IMAGE:$SHA
    docker tag $IMAGE:$SHA $IMAGE:$RELEASE
    docker push $IMAGE:$RELEASE
}

read -r -p "Are you sure you want to release ${SHA} as ${RELEASE}? [y/N] " response
case "$response" in
    [yY][eE][sS]|[yY])
        do_release
        ;;
    *)
        exit
        ;;
esac