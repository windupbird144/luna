#!/bin/env sh
# Helper script to call the /reminders endpoint every 5 seconds 11 times.
# Set it up in crontab to run at the full minute.
export LUNAPORT=12345
for i in {1..11}; do
    curl localhost:$LUNAPORT/reminders
    sleep 5
done