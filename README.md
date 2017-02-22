# Requirements to run locally:

Start by creating a file name "slack_api.token" and put the slack api token into it,
it's already in .gitignore so you don't have to worry about it being committed.

# Installation on server:

Create the target folders(server):

    mkdir /tmp/unibot-tmp/
    sudo mkdir /srv/unibot/
    sudo chmod 755 /srv/unibot/

Run the build and deploy scripts (locally):

    ./build.sh && ./deploy.sh

Start by installing it as a service in a system that has systemd
installed(server):

    sudo systemctl enable /srv/unibot/unibot.service

And start the service(server):

    sudo systemctl start unibot.service

And check the logs that everything went ok(server):

    journalctl -u unibot.service # add -f to continuously print new entries
