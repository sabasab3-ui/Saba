"""Advanced Analytics Engine for SABA"""

import json
from typing import Dict, List, Optional, Any
from datetime import datetime, timedelta
from collections import defaultdict
import statistics


class AnalyticsEngine:
    """Track and analyze SABA platform usage and performance"""

    def __init__(self):
        self.events: List[Dict] = []
        self.metrics: Dict[str, List[float]] = defaultdict(list)
        self.sessions: Dict[str, Dict] = {}

    def record_event(self, event_type: str, data: Dict[str, Any]) -> None:
        """Record an analytics event"""
        event = {
            "type": event_type,
            "timestamp": datetime.now().isoformat(),
            "data": data,
        }
        self.events.append(event)

    def record_agent_execution(self, agent_type: str, duration_ms: float, 
                               success: bool, country: str = None) -> None:
        """Record agent execution metrics"""
        self.record_event("agent_execution", {
            "agent_type": agent_type,
            "duration_ms": duration_ms,
            "success": success,
            "country": country,
        })
        self.metrics[f"agent_{agent_type}_duration"].append(duration_ms)

    def record_automation_run(self, automation_id: str, duration_ms: float,
                             success: bool, items_processed: int = 0) -> None:
        """Record automation execution"""
        self.record_event("automation_run", {
            "automation_id": automation_id,
            "duration_ms": duration_ms,
            "success": success,
            "items_processed": items_processed,
        })
        self.metrics[f"automation_{automation_id}_duration"].append(duration_ms)

    def get_agent_stats(self, agent_type: str) -> Dict[str, Any]:
        """Get statistics for an agent type"""
        key = f"agent_{agent_type}_duration"
        values = self.metrics.get(key, [])
        
        if not values:
            return {"no_data": True}
        
        return {
            "agent_type": agent_type,
            "execution_count": len(values),
            "avg_duration_ms": sum(values) / len(values),
            "min_duration_ms": min(values),
            "max_duration_ms": max(values),
            "std_dev": statistics.stdev(values) if len(values) > 1 else 0,
            "median_duration_ms": statistics.median(values),
        }

    def get_platform_health(self) -> Dict[str, Any]:
        """Get overall platform health metrics"""
        total_events = len(self.events)
        successful_events = sum(1 for e in self.events if e.get("data", {}).get("success", True))
        
        return {
            "total_executions": total_events,
            "success_rate": round((successful_events / total_events * 100) if total_events > 0 else 0, 2),
            "uptime_status": "operational",
            "agents_active": 6,
            "automations_running": len([e for e in self.events if e["type"] == "automation_run"]),
            "avg_response_time_ms": self._calculate_avg_response_time(),
        }

    def get_usage_by_country(self) -> Dict[str, int]:
        """Get usage statistics by country"""
        country_usage = defaultdict(int)
        for event in self.events:
            country = event.get("data", {}).get("country")
            if country:
                country_usage[country] += 1
        return dict(country_usage)

    def get_usage_by_agent(self) -> Dict[str, int]:
        """Get usage statistics by agent type"""
        agent_usage = defaultdict(int)
        for event in self.events:
            if event["type"] == "agent_execution":
                agent_type = event.get("data", {}).get("agent_type")
                if agent_type:
                    agent_usage[agent_type] += 1
        return dict(agent_usage)

    def _calculate_avg_response_time(self) -> float:
        """Calculate average response time across all executions"""
        durations = []
        for event in self.events:
            duration = event.get("data", {}).get("duration_ms")
            if duration:
                durations.append(duration)
        return sum(durations) / len(durations) if durations else 0
