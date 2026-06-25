# ── base ────────────────────────────────────────────────────────────────────
# Shared foundation: just the image + working directory.
# Both dependency stages inherit from here, keeping things DRY.
FROM python:3.13-slim AS base
WORKDIR /app


# ── runtime-deps ─────────────────────────────────────────────────────────────
# Installs only the system libraries needed at runtime.
# Runs in parallel with build-deps (no dependency between the two).
FROM base AS runtime-deps
RUN apt-get update && apt-get install -y --no-install-recommends \
    libglib2.0-0 \
    libxcb1 \
    libgl1 \
    && rm -rf /var/lib/apt/lists/*


# ── build-deps ───────────────────────────────────────────────────────────────
# Installs Python packages. Build tools (gcc, etc.) live only in this stage
# and never make it into the final image.
# Runs in parallel with runtime-deps.
FROM base AS build-deps
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
# Install into an isolated directory so the COPY into the final image is clean.
RUN pip install --no-cache-dir --retries 5 --timeout 120 \
    --target=/deps \
    -r requirements.txt


# ── final ────────────────────────────────────────────────────────────────────
# Starts from the runtime-deps layer (system libs already present),
# then pulls in the pre-built Python packages from build-deps.
FROM runtime-deps AS final

# Drop compiled packages into site-packages so Python finds them normally.
COPY --from=build-deps /deps /usr/local/lib/python3.13/site-packages

COPY src/ ./src/

CMD ["python", "src/main.py"]