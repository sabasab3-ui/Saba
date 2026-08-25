from typing import Literal, Optional
from pydantic import BaseModel, Field

Mode = Literal["auto", "research", "analysis", "reasoning", "business", "coding"]

class RunRequest(BaseModel):
    task: str = Field(min_length=1, max_length=12000)
    mode: Mode = "auto"
    input: Optional[str] = None
    session_id: Optional[str] = None

class RunResponse(BaseModel):
    status: str
    mode: str
    agent: str
    output: str
    session_id: Optional[str] = None
