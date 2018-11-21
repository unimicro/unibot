#! /bin/bash
set -e

SERVER=${1:-unibot}

CMD='
tar xzf - -C /srv/unibot/
 && echo ≫ Replacing executable...
 && test -f /srv/unibot/unibot
 && mv /srv/unibot/unibot{,.old} || true
 && mv /srv/unibot/unibot{.new,}
 && echo ≫ Restarting service...
 && sudo systemctl daemon-reload
 && sudo service unibot restart
 && echo ≫ Checking status...
 && sleep 1
 && sudo service unibot status
 && echo ≫ Done
'

echo '≫ Transmitting executable "unibot.new" to "'${SERVER}'" remote...'
tar czf - unibot.new tokens.json unibot.service | ssh $SERVER $CMD
