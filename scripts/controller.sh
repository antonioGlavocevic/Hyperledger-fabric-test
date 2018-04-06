#!/bin/bash

DIR="$(dirname "$0")"
source $DIR/helperFunctions.sh

POSITIONAL=()
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
  -c|--command)
  COMMAND="$2"
  shift # past argument
  shift # past value
  ;;
  -o|--organisation)
  ORGANISATION="$2"
  shift # past argument
  shift # past value
  ;;
  -p|--peer)
  PEER="$2"
  shift # past argument
  shift # past value
  ;;
  -i|--invoke)
  INVOKE="$2"
  shift # past argument
  shift # past value
  ;;
  -q|--query)
  QUERY="$2"
  shift # past argument
  shift # past value
  ;;
  *)    # unknown option
  POSITIONAL+=("$1") # save it in an array for later
  shift # past argument
  ;;
esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

if [ ! -z "$INVOKE" ]; then
  echo "executing: invoke $ORGANISATION $PEER $INVOKE"
  invoke $ORGANISATION $PEER $INVOKE
elif [ ! -z "$QUERY" ]; then
  echo "executing: query $ORGANISATION $PEER $QUERY"
  query $ORGANISATION $PEER $QUERY
elif [ ! -z "$COMMAND" ]; then
  echo "executing: $COMMAND $ORGANISATION $PEER"
  $COMMAND $ORGANISATION $PEER
else
  echo "usage: -o organisation -p peer [-c command | -i invokeChaincode | -q queryChaincode]"
fi
