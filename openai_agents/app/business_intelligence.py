"""Advanced Business Intelligence Module for SABA"""

import json
from typing import Dict, List, Optional, Any
from datetime import datetime, timedelta
from dataclasses import dataclass, asdict
import asyncio


@dataclass
class BusinessInsight:
    """Represents a business intelligence insight"""
    insight_type: str  # market, customer, competitor, trend
    title: str
    description: str
    confidence: float  # 0-1
    data_points: List[str]
    recommendations: List[str]
    created_at: datetime
    expires_at: Optional[datetime] = None

    def to_dict(self) -> Dict:
        return asdict(self)


class BusinessIntelligenceEngine:
    """Advanced BI engine for market and business analysis"""

    def __init__(self):
        self.insights: Dict[str, BusinessInsight] = {}
        self.market_data: Dict[str, Dict] = {}
        self.trends: List[Dict] = []

    async def analyze_market(self, country: str, industry: str) -> Dict[str, Any]:
        """Analyze market conditions for a specific country and industry"""
        analysis = {
            "market_size": self._estimate_market_size(country, industry),
            "growth_rate": self._estimate_growth_rate(country, industry),
            "competition_level": self._assess_competition(country, industry),
            "entry_barriers": self._identify_barriers(country, industry),
            "opportunities": await self._find_opportunities(country, industry),
            "risks": self._identify_risks(country, industry),
        }
        return analysis

    async def analyze_customer_segment(self, segment: str, metrics: Dict) -> Dict[str, Any]:
        """Analyze a specific customer segment"""
        return {
            "segment": segment,
            "size_estimate": metrics.get("size", 0),
            "purchasing_power": self._assess_purchasing_power(segment),
            "preferences": self._identify_preferences(segment),
            "growth_potential": self._estimate_growth(segment),
            "acquisition_cost": self._estimate_acq_cost(segment),
        }

    async def project_roi(self, investment: float, timeline_months: int, 
                         industry: str) -> Dict[str, Any]:
        """Project ROI for an automation investment"""
        monthly_savings = self._estimate_monthly_savings(investment, industry)
        payback_months = investment / monthly_savings if monthly_savings > 0 else 0
        
        return {
            "initial_investment": investment,
            "monthly_savings": monthly_savings,
            "payback_period_months": round(payback_months, 1),
            "total_savings_12m": monthly_savings * 12,
            "roi_percentage": round((monthly_savings * 12 / investment * 100) if investment > 0 else 0, 1),
            "break_even_date": (datetime.now() + timedelta(days=payback_months * 30)).isoformat(),
        }

    def _estimate_market_size(self, country: str, industry: str) -> str:
        """Estimate market size"""
        # Simplified estimation
        market_sizes = {
            ("Uganda", "Retail"): "$2.5B",
            ("Uganda", "Agriculture"): "$8.5B",
            ("Kenya", "Fintech"): "$5.2B",
            ("Rwanda", "Tech"): "$1.8B",
        }
        return market_sizes.get((country, industry), "$Unknown")

    def _estimate_growth_rate(self, country: str, industry: str) -> str:
        """Estimate market growth rate"""
        return "12-18% annually"

    def _assess_competition(self, country: str, industry: str) -> str:
        """Assess competition level"""
        return "Medium-High"

    def _identify_barriers(self, country: str, industry: str) -> List[str]:
        """Identify market entry barriers"""
        return [
            "Regulatory compliance requirements",
            "Capital investment needs",
            "Existing market players",
            "Infrastructure limitations",
        ]

    async def _find_opportunities(self, country: str, industry: str) -> List[str]:
        """Find market opportunities"""
        return [
            f"Digital transformation gap in {industry}",
            f"Under-served SMEs in {country}",
            "Automation adoption accelerating",
            "Growing middle class purchasing power",
        ]

    def _identify_risks(self, country: str, industry: str) -> List[str]:
        """Identify market risks"""
        return [
            "Currency fluctuation",
            "Regulatory changes",
            "Economic volatility",
            "Competition from international players",
        ]

    def _assess_purchasing_power(self, segment: str) -> str:
        """Assess customer purchasing power"""
        return "Moderate to High"

    def _identify_preferences(self, segment: str) -> List[str]:
        """Identify customer preferences"""
        return [
            "Cost-effective solutions",
            "Mobile-first platforms",
            "Local language support",
            "Flexible payment options",
        ]

    def _estimate_growth(self, segment: str) -> str:
        """Estimate segment growth"""
        return "15-25% annually"

    def _estimate_acq_cost(self, segment: str) -> str:
        """Estimate customer acquisition cost"""
        return "$50-200"

    def _estimate_monthly_savings(self, investment: float, industry: str) -> float:
        """Estimate monthly savings from automation"""
        # Typical ROI: 30-50% of investment per year
        annual_roi = investment * 0.4
        return annual_roi / 12
