COMPANY_NAME=webdevelop-pro
SERVICE_NAME=go-common

case $1 in

run)
  GIT_COMMIT=$(git rev-parse --short HEAD)
  BUILD_DATE=$(date "+%Y%m%d")
  build && ./app
  ;;

audit)
  echo "running gosec"
  gosec ./...
  ;;

*)
  BRANCH_NAME=`git rev-parse --abbrev-ref HEAD`
  GIT_COMMIT=`git rev-parse --short HEAD`
  echo $BRANCH_NAME, $GIT_COMMIT
  podman build --platform linux/amd64 -t cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT-amd64 -t cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev-amd64 .
  podman build --platform linux/arm64 -t cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT-arm64 -t cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev-arm64 .
  podman manifest create cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT
  podman manifest add cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT-amd64 
  podman manifest add cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT-arm64 
  podman manifest push --all cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT
  podman manifest create cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev
  podman manifest add cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev-amd64 
  podman manifest add cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev-arm64 
  podman manifest push --all cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev
  # docker buildx build --platform linux/amd64,linux/arm64 -t docker.io/webdeveloppro/$SERVICE_NAME:$GIT_COMMIT -t docker.io/webdeveloppro/$SERVICE_NAME:latest-dev -t cr.webdevelop.biz/$COMPANY_NAME/$SERVICE_NAME:latest-dev --platform=linux/amd64 .
  # docker push cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT
  # docker push cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev
  # docker push docker.io/webdeveloppro/$SERVICE_NAME:$GIT_COMMIT
  # docker push webdeveloppro/$SERVICE_NAME:latest-dev
  ;;

esac

