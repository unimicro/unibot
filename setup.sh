#! /bin/bash
set -e

SERVER=${1:-unibot}

CMD='tar xzf - -C /tmp/
 && echo ≫ Making directory
 && sudo mkdir -p /srv/unibot/
 && echo ≫ Copying service
 && sudo cp /tmp/unibot.service /srv/unibot/
 && echo ≫ Enabling service
 && sudo systemctl enable /srv/unibot/unibot.service
'

echo 'Running command on "'${SERVER}'":' $CMD
tar czf - unibot.service |ssh $SERVER $CMD
