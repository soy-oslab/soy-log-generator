#!python3
import json
import os
import sys
import ast

"""
# environment setting example

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
    val = os.getenv(env).strip()
    if len(val) == 0:
        print(f"{env} must be specified", file=sys.stderr)
        sys.exit(-1)
    return val


if __name__ == "__main__":
    ip = get_value_from_environment("GENERATOR_TARGET_IP")
    port = get_value_from_environment("GENERATOR_TARGET_PORT")
    hot_ring_capacity = get_value_from_environment("GENERATOR_HOT_RING_CAPACITY")
    cold_ring_capacity = get_value_from_environment("GENERATOR_COLD_RING_CAPACITY")
    cold_timeout_millis = get_value_from_environment("GENERATOR_COLD_TIMEOUT_MILLIS")
    hot_ring_threshold = get_value_from_environment("GENERATOR_HOT_RING_THRESHOLD")
    cold_ring_threshold = get_value_from_environment("GENERATOR_COLD_RING_THRESHOLD")
    cold_send_threshold = get_value_from_environment("GENERATOR_COLD_SEND_THRESHOLD_BYTES")
    polling_interval_millis = get_value_from_environment(
        "GENERATOR_POLLING_INTERVAL_MILLIS"
    )
    files = get_value_from_environment("GENERATOR_FILES")

    configContents = {
        "targetIp": ip,
        "targetPort": port,
        "hotRingCapacity": int(hot_ring_capacity),
        "coldRingCapacity": int(cold_ring_capacity),
        "coldTimeoutMilli": int(cold_timeout_millis),
        "hotRingThreshold": int(hot_ring_threshold),
        "coldRingThreshold": int(cold_ring_threshold),
        "coldSendThresholdBytes": int(cold_send_threshold),
        "pollingIntervalMilli": int(polling_interval_millis),
        "files": ast.literal_eval(files),
    }

    fp = open("config.json", "w")
    fp.write(json.dumps(configContents, indent=4))
    fp.close()
