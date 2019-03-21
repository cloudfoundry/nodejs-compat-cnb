#!/usr/bin/env bash
set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

if [[ ! -d integration ]]; then
    echo -e "\n\033[0;31m** WARNING  No Integration tests **\033[0m"
    exit 0
fi

PACK_VERSION=${PACK_VERSION:-""}
source ./scripts/install_tools.sh $PACK_VERSION

export CNB_BUILD_IMAGE=${CNB_BUILD_IMAGE:-cfbuildpacks/cflinuxfs3-cnb-experimental:build}
export CNB_RUN_IMAGE=${CNB_RUN_IMAGE:-cfbuildpacks/cflinuxfs3-cnb-experimental:run}

# Always pull latest images
# Most helpful for local testing consistency with CI (which would already pull the latest)
docker pull $CNB_BUILD_IMAGE
docker pull $CNB_RUN_IMAGE

echo "Run Buildpack Runtime Integration Tests"
set +e
go test ./integration/... -timeout 0 -v -run Integration
exit_code=$?

if [ "$exit_code" != "0" ]; then
    echo -e "\n\033[0;31m** GO Test Failed **\033[0m"
else
    echo -e "\n\033[0;32m** GO Test Succeeded **\033[0m"
fi

exit $exit_code
