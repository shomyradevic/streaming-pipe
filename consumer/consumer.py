from pyspark.sql import SparkSession
import pyspark.sql.functions as F
from pyspark.sql.types import *
import os


os.environ['PYSPARK_SUBMIT_ARGS'] = '--packages org.apache.spark:spark-streaming-kafka-0-10_2.13:4.0.0,org.apache.spark:spark-sql-kafka-0-10_2.13:4.0.0 pyspark-shell'

print("Creating Spark session ...")
spark = SparkSession.\
    builder.\
    appName("Consumer-raw-app").\
    master("local[*]").\
    getOrCreate()


print("Creating kafka config ...")
kafka_config = {
    "raw": {
        "kafka.bootstrap.servers": "kafka:9092",
        "subscribe": "raw-topic",
        "startingOffsets": "latest",
        "failOnDataLoss": "false"
    },
    "processed": {
        "kafka.bootstrap_servers": "kafka:9092",
        "topic": "processed-topic",
        "checkpointLocation": "./check.txt" # ?
    }  # pyspark.errors.exceptions.captured.AnalysisException: checkpointLocation must be specified either through option("checkpointLocation", ...) or SparkSession.conf.set("spark.sql.streaming.checkpointLocation", ...).
}

print("Creating df schema ...")
df_schema = StructType([
    StructField(name="id",          dataType=IntegerType(),         nullable=False),
    StructField(name="text",        dataType=StringType(),          nullable=False),
    StructField(name="author_id",   dataType=IntegerType(),         nullable=False),
    StructField(name="retweets",    dataType=IntegerType(),         nullable=False),
    StructField(name="likes",       dataType=IntegerType(),         nullable=False),
    StructField(name="created_at",  dataType=TimestampNTZType(),    nullable=False),  # Check if ntz or some else
    StructField(name="edited",      dataType=BooleanType(),         nullable=False)
])  

print("Creating Read stream ...")
df = spark.readStream.\
    format("kafka").\
    options(**kafka_config["raw"]).\
    load().\
    select(
        F.from_json(
            F.col("value").cast("string"), 
            df_schema
        )
        .alias("json_data")
    ).select("json_data.*")


print("Performing filter ...")
# Filter down tweets that were created from monday to thursday and have > 10K likes and > 1K retweets
df = df.filter(
    (
        (F.dayofweek(F.col("created_at")).between(2, 5)) & 
        (F.col("likes") > 10000) & 
        (F.col("retweets") > 1000)
    )
).select(
    F.col("author_id"),
    F.col("created_at"),
    F.col("likes"),
    F.col("retweets"),
    F.col("text")
)

print("Creating output dataframe ...")
output_df = df.select(F.to_json(F.struct(*df.columns)).alias("value")) 
print(output_df)


print("Creating Write stream ...")
query = output_df.writeStream.\
        format("kafka").\
        options(**kafka_config["processed"]).\
        start()

query.awaitTermination()

# Message examples:
# {
#   "author_id": 687004,
#   "created_at": "2020-07-10T21:51:25Z",
#   "edited": false,
#   "id": 1280432,
#   "likes": 764,
#   "retweets": 93420,
#   "text": "Go maps are typed, so the value type must be consistent or use interface{} to support multiple types."
# }���{
#   "author_id": 695898,
#   "created_at": "2018-11-09T09:37:49Z",
#   "edited": false,
#   "id": 5364929,
#   "likes": 926,
#   "retweets": 45404,
#   "text": "For more accurate and detailed performance analysis, use Go's built-in benchmarking tools with the testing package."
# }