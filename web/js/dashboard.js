// SABA Dashboard - Main Application Logic

class SabaApp {
    constructor() {
        this.baseUrl = window.location.origin;
        this.apiUrl = '/api';
        this.stats = {};
        this.init();
    }

    async init() {
        console.log('🤖 SABA Dashboard Initializing...');
        await this.loadStats();
        await this.setupEventListeners();
        this.startHealthCheck();
    }

    async loadStats() {
        try {
            const response = await fetch('/health');
            const data = await response.json();
            this.updateUI(data);
        } catch (error) {
            console.error('Failed to load stats:', error);
            this.showAlert('Failed to connect to SABA service', 'danger');
        }
    }

    updateUI(data) {
        const resultElement = document.getElementById('result');
        if (resultElement) {
            resultElement.textContent = JSON.stringify(data, null, 2);
        }
    }

    setupEventListeners() {
        const checkHealthBtn = document.getElementById('checkHealthBtn');
        if (checkHealthBtn) {
            checkHealthBtn.addEventListener('click', () => this.checkHealth());
        }

        const runInventoryBtn = document.getElementById('runInventoryBtn');
        if (runInventoryBtn) {
            runInventoryBtn.addEventListener('click', () => this.runAgent('inventory'));
        }
    }

    async checkHealth() {
        const resultEl = document.getElementById('result');
        if (!resultEl) return;

        resultEl.innerHTML = '<div class="loading"></div> Checking SABA health...';

        try {
            const response = await fetch('/health');
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const data = await response.json();
            resultEl.textContent = JSON.stringify(data, null, 2);
            this.showAlert('✓ SABA is healthy and operational', 'success');
        } catch (error) {
            resultEl.textContent = `Error: ${error.message}`;
            this.showAlert('✗ Failed to reach SABA service', 'danger');
        }
    }

    async runAgent(task) {
        const resultEl = document.getElementById('result');
        if (!resultEl) return;

        resultEl.innerHTML = '<div class="loading"></div> Running agent task...';

        try {
            const response = await fetch('/api/agent', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    task: task,
                    params: { action: 'check' }
                })
            });

            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const data = await response.json();
            resultEl.textContent = JSON.stringify(data, null, 2);
            this.showAlert('✓ Agent task completed successfully', 'success');
        } catch (error) {
            resultEl.textContent = `Error: ${error.message}`;
            this.showAlert('✗ Agent task failed', 'danger');
        }
    }

    startHealthCheck() {
        // Check health every 30 seconds
        setInterval(() => this.loadStats(), 30000);
    }

    showAlert(message, type = 'info') {
        const alertContainer = document.getElementById('alertContainer');
        if (!alertContainer) return;

        const alert = document.createElement('div');
        alert.className = `alert ${type}`;
        alert.textContent = message;
        alertContainer.appendChild(alert);

        // Auto-remove after 5 seconds
        setTimeout(() => alert.remove(), 5000);
    }
}

// Initialize app when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => new SabaApp());
} else {
    new SabaApp();
}
