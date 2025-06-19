from pyspark.sql import SparkSession
import pyspark.sql.functions as F
from pyspark.sql.types import *
import os


os.environ['PYSPARK_SUBMIT_ARGS'] = '--packages org.apache.spark:spark-streaming-kafka-0-10_2.13:4.0.0,org.apache.spark:spark-sql-kafka-0-10_2.13:4.0.0 pyspark-shell'

print("Creating Spark session!")

spark = SparkSession.\
    builder.\
    appName("Consumer-raw-app").\
    master("local[*]").\
    getOrCreate()


# kafka konfig


print("Creating kafka input config!")

kafka_input_config = {
    "kafka.bootstrap.servers": "kafka:9092",
    "subscribe": "raw-topic",
    "startingOffsets": "latest",
    "failOnDataLoss": "false"
}


# I guess this will later be processed-topic ( Which we first need to create ), but now this should be commented out until we check kafka is setup correctly
# kafka_output_config = {
#     "kafka.bootstrap.servers": "kafka:9092",
#     "topic": "processed-topic",
#     "checkpointLocation": "./check.txt" # ?
# }


print("Creating df schema!")


# Input schema
df_schema = StructType([
    StructField(name="id",          dataType=IntegerType(),         nullable=False),
    StructField(name="text",        dataType=StringType(),          nullable=False),
    StructField(name="author_id",   dataType=IntegerType(),         nullable=False),
    StructField(name="retweets",    dataType=IntegerType(),         nullable=False),
    StructField(name="likes",       dataType=IntegerType(),         nullable=False),
    StructField(name="created_at",  dataType=TimestampNTZType(),    nullable=False),  # Check if ntz or some else
    StructField(name="edited",      dataType=BooleanType(),         nullable=False)
])  

print("Creating Read stream!")

# Read stream
df = spark.readStream.\
    format("kafka").\
    options(**kafka_input_config).\
    load().\
    select(
        F.from_json(
            F.col("value").cast("string"), 
            df_schema
        )
        .alias("json_data")
    ).select("json_data.*")


query = df.writeStream.\
    format("console")\
    .trigger(processingTime="100 milliseconds")\
    .start()

query.awaitTermination()


