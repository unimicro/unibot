#! /bin/bash
set -e

rsync -r dist/* unidiv:/tmp/unibot-tmp/
#scp -r dist/* unidiv:/tmp/unibot-tmp/
ssh unidiv 'sudo mv /tmp/unibot-tmp/* /srv/unibot/ 2>/dev/null; sudo mv /srv/unibot/unibot{,.old} ; sudo mv /srv/unibot/unibot{.new,} && sudo service unibot restart'
