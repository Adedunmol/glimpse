import os
from dotenv import load_dotenv

load_dotenv()

REDIS_URL = os.getenv("GLIMPSE_REDIS_URL", "redis://localhost:6379/0")
QUEUE_NAME = os.getenv("GLIMPSE_QUEUE_NAME", "ml:jobs")
STREAM_NAME = os.getenv("STREAM_NAME", "image-tasks")

DATABASE_URL = os.getenv("GLIMPSE_DATABASE_URL", "postgresql+asyncpg://postgres:postgres@postgres:5432/glimpse")

AWS_ACCESS_KEY = os.getenv("GLIMPSE_AWS.ACCESS_KEY_ID")
AWS_SECRET_KEY = os.getenv("GLIMPSE_AWS.SECRET_ACCESS_KEY")
AWS_REGION = os.getenv("GLIMPSE_AWS.REGION", "us-east-1")
S3_BUCKET = os.getenv("GLIMPSE_AWS.UPLOAD_BUCKET")
ENDPOINT_URL=os.getenv("GLIMPSE_AWS.ENDPOINT_URL")
MODEL_NAME = os.getenv("MODEL_NAME", "buffalo_sc")