#!python3
import json
import os
import sys
import ast

"""
# environment setting example

export GENERATOR_NAMESPACE='kube-namespace'
export GENERATOR_TARGET_IP=localhost
export GENERATOR_TARGET_PORT=8972
export GENERATOR_HOT_RING_CAPACITY=8
export GENERATOR_COLD_RING_CAPACITY=32
export GENERATOR_COLD_TIMEOUT_MILLIS=3000
export GENERATOR_POLLING_INTERVAL_MILLIS=1000
export GENERATOR_HOT_RING_THRESHOLD=0
export GENERATOR_COLD_RING_THRESHOLD=0
export GENERATOR_COLD_SEND_THRESHOLD_BYTES=4096
export GENERATOR_FILES='[
  {"filename":"test1.txt", "hotFilter":["error","critical"]},
  {"filename":"test2.txt", "hotFilter":["critical","warn"]}
]'
"""


def get_value_from_environment(env: str) -> str:
    try:
        val = os.getenv(env).strip()
    except Exception as e:
        print(e)
        return None
    if len(val) == 0:
        return None
    return val


def assign_config_contents(d: map, k: str, v):
    if v is None:
        return
    if k in [
        "hotRingCapacity",
        "coldRingCapacity",
        "coldTimeoutMilli",
        "hotRingThreshold",
        "coldRingThreshold",
        "coldSendThresholdBytes",
        "pollingIntervalMilli",
    ]:
        d[k] = int(v)
    elif k in ["files"]:
        d[k] = ast.literal_eval(v)
    else:
        d[k] = v


if __name__ == "__main__":
    configContents = {}
    assign_config_contents(
        configContents, "namespace", get_value_from_environment("GENERATOR_NAMESPACE")
    )
    assign_config_contents(
        configContents, "targetIp", get_value_from_environment("GENERATOR_TARGET_IP")
    )
    assign_config_contents(
        configContents,
        "targetPort",
        get_value_from_environment("GENERATOR_TARGET_PORT"),
    )
    assign_config_contents(
        configContents,
        "hotRingCapacity",
        get_value_from_environment("GENERATOR_HOT_RING_CAPACITY"),
    )
    assign_config_contents(
        configContents,
        "coldRingCapacity",
        get_value_from_environment("GENERATOR_COLD_RING_CAPACITY"),
    )
    assign_config_contents(
        configContents,
        "coldTimeoutMilli",
        get_value_from_environment("GENERATOR_COLD_TIMEOUT_MILLIS"),
    )
    assign_config_contents(
        configContents,
        "hotRingThreshold",
        get_value_from_environment("GENERATOR_HOT_RING_THRESHOLD"),
    )
    assign_config_contents(
        configContents,
        "coldRingThreshold",
        get_value_from_environment("GENERATOR_COLD_RING_THRESHOLD"),
    )
    assign_config_contents(
        configContents,
        "coldSendThresholdBytes",
        get_value_from_environment("GENERATOR_COLD_SEND_THRESHOLD_BYTES"),
    )
    assign_config_contents(
        configContents,
        "pollingIntervalMilli",
        get_value_from_environment("GENERATOR_POLLING_INTERVAL_MILLIS"),
    )
    assign_config_contents(
        configContents, "files", get_value_from_environment("GENERATOR_FILES")
    )

    fp = open("config.json", "w")
    fp.write(json.dumps(configContents, indent=4))
    fp.close()
