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
  docker buildx build --platform linux/amd64,linux/arm64 -t cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT -t cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev -t docker.io/webdeveloppro/$SERVICE_NAME:$GIT_COMMIT -t docker.io/webdeveloppro/$SERVICE_NAME:latest-dev -t cr.webdevelop.biz/$COMPANY_NAME/$SERVICE_NAME:latest-dev --platform=linux/amd64 .
  docker push cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:$GIT_COMMIT
  docker push cr.webdevelop.pro/$COMPANY_NAME/$SERVICE_NAME:latest-dev
  docker push cr.webdevelop.biz/$COMPANY_NAME/$SERVICE_NAME:latest-dev
  docker push docker.io/webdeveloppro/$SERVICE_NAME:$GIT_COMMIT
  docker push webdeveloppro/$SERVICE_NAME:latest-dev
  ;;

esac

