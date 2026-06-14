import os
from dotenv import load_dotenv

load_dotenv()

REDIS_URL = os.getenv("GLIMPSE_REDIS_URL", "redis://localhost:6379/0")
QUEUE_NAME = os.getenv("GLIMPSE_QUEUE_NAME", "ml:jobs")
STREAM_NAME = os.getenv("STREAM_NAME", "image-tasks")