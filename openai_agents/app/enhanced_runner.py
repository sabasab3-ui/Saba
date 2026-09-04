"""Enhanced Runner with Advanced Features"""

import asyncio
from typing import Dict, Any, Optional
from datetime import datetime

from .agents import build_agents
from .config import settings
from .business_intelligence import BusinessIntelligenceEngine
from .analytics_engine import AnalyticsEngine
from .connectors import ConnectorFactory


class EnhancedRunner:
    """Advanced task runner with BI, analytics, and connectors"""

    def __init__(self):
        self.agents = build_agents()
        self.bi_engine = BusinessIntelligenceEngine()
        self.analytics = AnalyticsEngine()
        self.connectors = {}

    async def run_task(
        self,
        task: str,
        mode: str = "auto",
        input_data: Optional[str] = None,
        country: str = None,
        industry: str = None,
    ) -> Dict[str, Any]:
        """Execute a task with enhanced capabilities"""
        start_time = datetime.now()

        try:
            # Select appropriate agent
            agent = self.agents.get(mode)
            if not agent:
                return {"error": f"Unknown mode: {mode}"}

            # For business analysis tasks, use BI engine
            if mode == "business" and country and industry:
                market_analysis = await self.bi_engine.analyze_market(country, industry)
                result = {
                    "agent": agent.name,
                    "output": market_analysis,
                    "analysis_type": "market",
                }
            else:
                # Standard agent execution
                result = await self._execute_agent(agent, task, input_data)

            # Record analytics
            duration_ms = (datetime.now() - start_time).total_seconds() * 1000
            self.analytics.record_agent_execution(
                agent_type=mode,
                duration_ms=duration_ms,
                success=True,
                country=country,
            )

            return result

        except Exception as e:
            duration_ms = (datetime.now() - start_time).total_seconds() * 1000
            self.analytics.record_agent_execution(
                agent_type=mode,
                duration_ms=duration_ms,
                success=False,
                country=country,
            )
            return {"error": str(e), "agent": mode}

    async def _execute_agent(self, agent, task: str, input_data: Optional[str]) -> Dict[str, Any]:
        """Execute agent with the given task"""
        # Simulate agent execution
        await asyncio.sleep(0.1)
        return {
            "agent": agent.name,
            "output": f"Executed {task}",
            "status": "completed",
        }

    async def run_automation(
        self,
        automation_id: str,
        automation_config: Dict[str, Any],
    ) -> Dict[str, Any]:
        """Run an automation workflow"""
        start_time = datetime.now()
        items_processed = 0

        try:
            # Get connector if needed
            connector_type = automation_config.get("connector")
            if connector_type:
                await self._setup_connector(connector_type, automation_config)

            # Execute automation steps
            for step in automation_config.get("steps", []):
                items_processed += await self._execute_step(step)

            # Record analytics
            duration_ms = (datetime.now() - start_time).total_seconds() * 1000
            self.analytics.record_automation_run(
                automation_id=automation_id,
                duration_ms=duration_ms,
                success=True,
                items_processed=items_processed,
            )

            return {
                "automation_id": automation_id,
                "status": "completed",
                "items_processed": items_processed,
                "duration_ms": duration_ms,
            }

        except Exception as e:
            duration_ms = (datetime.now() - start_time).total_seconds() * 1000
            self.analytics.record_automation_run(
                automation_id=automation_id,
                duration_ms=duration_ms,
                success=False,
                items_processed=items_processed,
            )
            return {"error": str(e), "automation_id": automation_id}

    async def _setup_connector(self, connector_type: str, config: Dict) -> None:
        """Setup external system connector"""
        if connector_type not in self.connectors:
            connector = ConnectorFactory.create_connector(connector_type, config)
            if connector:
                await connector.connect()
                self.connectors[connector_type] = connector

    async def _execute_step(self, step: Dict[str, Any]) -> int:
        """Execute a single automation step"""
        await asyncio.sleep(0.05)  # Simulate work
        return step.get("items", 1)

    def get_platform_stats(self) -> Dict[str, Any]:
        """Get comprehensive platform statistics"""
        return {
            "health": self.analytics.get_platform_health(),
            "usage_by_agent": self.analytics.get_usage_by_agent(),
            "usage_by_country": self.analytics.get_usage_by_country(),
            "agents_available": list(self.agents.keys()),
        }
