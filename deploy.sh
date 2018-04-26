#! /bin/bash
set -e

rsync -r dist/* unibot:/tmp/unibot-tmp/
#scp -r dist/* unibot:/tmp/unibot-tmp/
ssh unibot 'sudo mv /tmp/unibot-tmp/* /srv/unibot/ 2>/dev/null; sudo mv /srv/unibot/unibot{,.old} ; sudo mv /srv/unibot/unibot{.new,} && sudo service unibot restart'
