#!/bin/bash
# Parameters are mandatory for posting the message to slack.
webhook_url="$WEBHOOK_URL"
message_body="$MESSAGE_BODY"
channel_name="$CHANNEL_NAME" # Slack channel name to post the message to
username="$USERNAME" # Username to send message as
user_icon="$USER_ICON" # User icon like ':ghost:'

if [[ -z "$webhook_url" ]]
then
  echo "Webhook URL can not be empty"
elif [[ -z "$message_body" ]]
then
  echo "Message body can not be empty"
elif [[ -z "$channel_name" ]]
then
  echo "Channel name can not be empty"
elif [[ -z "$username" ]]
then
  echo "User name can not be empty"
elif [[ -z "$user_icon" ]]
then
  echo "User icon can not be empty"
fi

# Command for posting the message to slack.
curl -f -X POST --data-urlencode "payload={\"channel\": \"#${channel_name}\", \"username\": \"${username}\", \"text\": ${message_body} , \"icon_emoji\": \"${user_icon}\"}" ${webhook_url}
