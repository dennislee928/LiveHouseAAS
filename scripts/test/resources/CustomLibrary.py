"""Custom Robot Framework library for Redis, NATS, and MinIO connectivity checks."""

import asyncio
import redis as redis_client
import boto3
from botocore.exceptions import ClientError, EndpointResolutionError
from robot.api.deco import keyword
from robot.api import logger


class CustomLibrary:
    ROBOT_LIBRARY_SCOPE = "SUITE"

    # ── Redis ──────────────────────────────────────────────────────────────

    @keyword("Connect To Redis")
    def connect_to_redis(self, host, port=6379, password=None, db=0):
        kwargs = {"host": host, "port": int(port), "db": int(db),
                  "socket_connect_timeout": 5, "decode_responses": True}
        if password:
            kwargs["password"] = password
        self._redis = redis_client.Redis(**kwargs)
        self._redis.ping()
        logger.info(f"Connected to Redis at {host}:{port}")

    @keyword("Redis Set And Get")
    def redis_set_and_get(self, key, value):
        self._redis.set(key, value, ex=30)
        result = self._redis.get(key)
        assert result == value, f"Expected '{value}', got '{result}'"
        logger.info(f"Redis SET/GET round-trip OK: {key}={value}")
        return result

    @keyword("Redis Ping")
    def redis_ping(self):
        response = self._redis.ping()
        assert response, "Redis PING returned False"
        logger.info("Redis PING OK")

    @keyword("Redis Delete Key")
    def redis_delete_key(self, key):
        self._redis.delete(key)

    @keyword("Disconnect From Redis")
    def disconnect_from_redis(self):
        if hasattr(self, "_redis"):
            self._redis.close()

    # ── MinIO / S3 ─────────────────────────────────────────────────────────

    @keyword("Connect To MinIO")
    def connect_to_minio(self, endpoint, access_key, secret_key):
        self._minio = boto3.client(
            "s3",
            endpoint_url=endpoint,
            aws_access_key_id=access_key,
            aws_secret_access_key=secret_key,
            region_name="us-east-1",
        )
        # Verify connectivity by listing buckets
        self._minio.list_buckets()
        logger.info(f"Connected to MinIO at {endpoint}")

    @keyword("MinIO Create Bucket If Not Exists")
    def minio_create_bucket_if_not_exists(self, bucket):
        try:
            self._minio.head_bucket(Bucket=bucket)
            logger.info(f"Bucket '{bucket}' already exists")
        except ClientError as e:
            if e.response["Error"]["Code"] == "404":
                self._minio.create_bucket(Bucket=bucket)
                logger.info(f"Created bucket '{bucket}'")
            else:
                raise

    @keyword("MinIO Upload And Download")
    def minio_upload_and_download(self, bucket, key, content):
        import io
        data = content.encode()
        self._minio.put_object(Bucket=bucket, Key=key, Body=data)
        obj = self._minio.get_object(Bucket=bucket, Key=key)
        result = obj["Body"].read().decode()
        assert result == content, f"Expected '{content}', got '{result}'"
        logger.info(f"MinIO PUT/GET round-trip OK: s3://{bucket}/{key}")
        return result

    @keyword("MinIO Delete Object")
    def minio_delete_object(self, bucket, key):
        self._minio.delete_object(Bucket=bucket, Key=key)

    @keyword("MinIO List Buckets")
    def minio_list_buckets(self):
        response = self._minio.list_buckets()
        buckets = [b["Name"] for b in response.get("Buckets", [])]
        logger.info(f"MinIO buckets: {buckets}")
        return buckets

    # ── NATS (HTTP monitoring) ─────────────────────────────────────────────

    @keyword("NATS Monitoring Healthz")
    def nats_monitoring_healthz(self, monitor_url):
        import urllib.request
        url = monitor_url.rstrip("/") + "/healthz"
        with urllib.request.urlopen(url, timeout=5) as resp:
            body = resp.read().decode()
            assert resp.status == 200, f"NATS /healthz returned {resp.status}"
            logger.info(f"NATS /healthz OK: {body}")
            return body

    @keyword("NATS Monitoring Varz")
    def nats_monitoring_varz(self, monitor_url):
        import urllib.request, json
        url = monitor_url.rstrip("/") + "/varz"
        with urllib.request.urlopen(url, timeout=5) as resp:
            data = json.loads(resp.read())
            logger.info(f"NATS server: {data.get('server_name','?')} "
                        f"v{data.get('version','?')}")
            return data
