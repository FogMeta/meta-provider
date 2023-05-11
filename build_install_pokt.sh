#!/bin/bash

make clean all

BINARY_NAME=meta-provider
if [ -f "./build/swan-provider" ]; then
    echo "${BINARY_NAME} build success"
    mv ./build/swan-provider ./${BINARY_NAME}
    chmod +x ./${BINARY_NAME}
else
  echo "${BINARY_NAME} build failed"
  exit -1
fi


CONF_FILE_DIR=${HOME}/.swan/provider
mkdir -p ${CONF_FILE_DIR}

POKT_CONFIG_FILE=${CONF_FILE_DIR}/config-pokt.toml
if [ -f "${POKT_CONFIG_FILE}" ]; then
    echo "${POKT_CONFIG_FILE} exists"
else
    cp ./config/pokt/config-pokt.toml.example ${CONF_FILE_DIR}/config-pokt.toml
    echo "${CONF_FILE_DIR}/config-pokt.toml created"
fi

POKT_CHAINS_FILE=${CONF_FILE_DIR}/chains.json
if [ -f "${POKT_CHAINS_FILE}" ]; then
    echo "${POKT_CHAINS_FILE} exists"
else
    cp ./config/pokt/chains.json ${CONF_FILE_DIR}/chains.json
    echo "${CONF_FILE_DIR}/chains.json created"
fi

POKT_GENESIS_FILE=${CONF_FILE_DIR}/genesis.json
if [ -f "${POKT_GENESIS_FILE}" ]; then
    echo "${POKT_GENESIS_FILE} exists"
else
    cp ./config/pokt/mainnet-genesis.json ${CONF_FILE_DIR}/genesis.json
    echo "${CONF_FILE_DIR}/genesis.json created"
fi




