# soy-log-generator

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/94b16cd6d8fa4cf99eb108e4d4e1c922)](https://www.codacy.com/gh/soyoslab/soy_log_generator/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=soyoslab/soy_log_generator&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/94b16cd6d8fa4cf99eb108e4d4e1c922)](https://www.codacy.com/gh/soyoslab/soy_log_generator/dashboard?utm_source=github.com&utm_medium=referral&utm_content=soyoslab/soy_log_generator&utm_campaign=Badge_Coverage)
![log-generator-build](https://github.com/soyoslab/soy_log_generator/actions/workflows/log-generator-build.yml/badge.svg)
![dockerize](https://github.com/soyoslab/soy_log_generator/actions/workflows/dockerize.yml/badge.svg)

# Introduction

This project sends the messages got from the log files to soy\_log\_collector.

# Build

Prepare the build environment. we use the go with 1.16 version.

```bash
sudo apt update -y && sudo apt install -y python3
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.7.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

Build this program by using the below commands.

```bash
git clone --recurse-submodules https://github.com/soyoslab/soy_log_generator.git
cd soy_log_generator
make test && make
```

# Usage

Set the valid environment variables.

```bash
export GENERATOR_NAMESPACE=test-server
export GENERATOR_TARGET_IP=${LOG_COLLCETOR_IP}
export GENERATOR_TARGET_PORT=${LOG_COLLCETOR_PORT}
export GENERATOR_HOT_RING_CAPACITY=8
export GENERATOR_COLD_RING_CAPACITY=32
export GENERATOR_HOT_RING_THRESHOLD=2
export GENERATOR_COLD_RING_THRESHOLD=2
export GENERATOR_COLD_TIMEOUT_MILLIS=3000
export GENERATOR_COLD_SEND_THRESHOLD_BYTES=4096
export GENERATOR_POLLING_INTERVAL_MILLIS=1000
export GENERATOR_FILES='[{"filename":"/var/log/*log","hotFilter":["error","failed","critical"]},]'
```

Generate the configuration file by using the configuration generator script.

```bash
sudo python3 ./scripts/config-file-generator.py
```

Now you can run the program.

```bash
sudo make generator-run
```

# Docker

You can build the docker image.

```bash
sudo pip install git-archive-all
sudo make gen-src-archive
sudo mkdir -p /tmp/dockerize/generator
sudo tar xvzf build/generator.src.tar.gz --strip 1 -C /tmp/dockerize/generator
sudo docker build -t generator:latest -f scripts/Dockerfile /tmp/dockerize
```

Now you need the docker environment file named `.env` in your localhost machine. Its contents are like below.

```
GENERATOR_NAMESPACE=test-server
GENERATOR_TARGET_IP=${LOG_COLLCETOR_IP}
GENERATOR_TARGET_PORT=${LOG_COLLCETOR_PORT}
GENERATOR_HOT_RING_CAPACITY=8
GENERATOR_COLD_RING_CAPACITY=32
GENERATOR_HOT_RING_THRESHOLD=2
GENERATOR_COLD_RING_THRESHOLD=2
GENERATOR_COLD_TIMEOUT_MILLIS=3000
GENERATOR_COLD_SEND_THRESHOLD_BYTES=4096
GENERATOR_POLLING_INTERVAL_MILLIS=1000
GENERATOR_FILES=[{"filename":"/var/log/*log","hotFilter":["error","failed","critical"]},]
```

Run the docker image by using the below command.

```bash
sudo docker run --env-file ./.env --name generator-test generator
```

If you want to change the config setting in the container, you can restart by sending the interrupt signal to process

```bash
sudo kill -s SIGINT $PID
```

However, you must not give a relative path in the container configuration file. It will cause unexpected behavior in the analysis.
