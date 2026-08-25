from typing import Literal
from pydantic import BaseModel, Field

Mode = Literal["auto", "research", "analysis", "reasoning", "business", "coding"]

class RunRequest(BaseModel):
    task: str = Field(min_length=1, max_length=12000)
    mode: Mode = "auto"
    input: str | None = Field(default=None, max_length=12000)
    session_id: str | None = Field(default=None, max_length=256)

class RunResponse(BaseModel):
    status: str
    mode: str
    agent: str
    output: str
    session_id: str | None = None
