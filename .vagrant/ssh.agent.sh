#!/bin/bash
ssh-add -l >/dev/null
if [ "$?" == "1" ]; then
  ssh-add
fi
