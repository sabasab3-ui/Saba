import os
from dataclasses import dataclass
from dotenv import load_dotenv

load_dotenv()

@dataclass(frozen=True)
class Settings:
    openai_api_key: str = os.getenv("OPENAI_API_KEY", "")
    model: str | None = os.getenv("OPENAI_MODEL") or None
    host: str = os.getenv("HOST", "127.0.0.1")
    port: int = int(os.getenv("PORT", "8090"))
    saba_gateway_url: str = os.getenv(
        "SABA_GATEWAY_URL",
        "http://127.0.0.1:8080",
    )
    max_turns: int = int(os.getenv("MAX_TURNS", "12"))

settings = Settings()
