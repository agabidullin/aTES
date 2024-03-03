# Kafka
## Create topic

curl \
  -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic <KEY>" \
  https://pkc-75m1o.europe-west3.gcp.confluent.cloud:443/kafka/v3/clusters/lkc-y6d7po/topics \
  -d '{"topic_name":"<TOPIC_NAME>"}'