#!/usr/bin/env bash
export SPARKLE_DEPLOYMENT_TARGET=$SPARKLE_DEPLOYMENT_TARGET
if [ -z "$SPARKLE_DEPLOYMENT_TARGET" ]
then
      echo "\$SPARKLE_DEPLOYMENT_TARGET is empty"
      exit 2
fi
make build
echo $SPARKLE_DEPLOYMENT_TARGET
rsync -avz ./dist $SPARKLE_DEPLOYMENT_TARGET:~/
echo "Restart service"
ssh -o StrictHostKeyChecking=no $SPARKLE_DEPLOYMENT_TARGET -t 'sudo service sparkle stop'
ssh -o StrictHostKeyChecking=no $SPARKLE_DEPLOYMENT_TARGET -t 'sudo service sparkle start'
