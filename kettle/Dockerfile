FROM gcr.io/k8s-testimages/kettle:v20180111-dac5561d0

ADD buckets.yaml /

# CMD cat buckets.yaml

CMD ["/kettle/runner.sh", "--buckets=/buckets.yaml", "--bq_instance=", "--project=istio-testing"]

# TODO (chx): override CMD to call runner.py with istio flags
# TODO (chx): create bq instance
# TODO (chx): after instument the flags, also update the tests and make sure they pass