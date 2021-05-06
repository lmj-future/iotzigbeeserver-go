# Zigbee Server

## 支持的平台

* iotprivate（私有云）
* iotware（物联网平台2.0）
* iotedge（边缘网关C530）

## 组件支持

* Redis
* Mongo
* Postgres
* RabbitMQ
* MQTT
* Kafka

### 不同平台目前用到的组件

* iotprivate
  * Redis（多实例所需）
  * Mongo
  * RabbitMQ
  * MQTT
  * Kafka

* iotware
  * Redis（多实例所需）
  * Postgres
  * MQTT

* iotedge
  * Redis（多实例所需）
  * Postgres
  * MQTT

## Prometheus

* curl -v "http://127.0.0.1/iot/iotzigbeeserver/prometheus"
