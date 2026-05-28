document.addEventListener('DOMContentLoaded', function() {
    const backBtn = document.getElementById('backBtn');
    const analyticsContainer = document.getElementById('analyticsContainer');

    backBtn.addEventListener('click', function() {
        window.location.href = 'http://localhost:8080/static/index.html';
    });

    fetchAggregate();

    function fetchAggregate() {
        fetch(`http://localhost:8080/analytics/user_agent/`)
            .then(response => response.json())
            .then(data => {
                displayAnalytics(data);
            })
            .catch(error => {
                console.error('Error fetching aggregate analytics:', error);
                analyticsContainer.innerHTML = '<p>Failed to fetch analytics.</p>';
            });
    }

    function displayAnalytics(data) {
        if (!Array.isArray(data) || data.length === 0) {
            analyticsContainer.innerHTML = '<p>No analytics data available.</p>';
            return;
        }

        data.forEach((item, index) => {
            const card = document.createElement('div');
            card.className = 'analytics-card';
            card.style.animationDelay = `${index * 0.1}s`; // Stagger animations

            const shortUrl = item.short_url || 'Unknown';
            const userAgents = item.user_agent || [];
            const redirectCount = item.redirect_count || 0;

            card.innerHTML = `
                <h3>${shortUrl}</h3>
                <p><strong>Redirect Count:</strong> ${redirectCount}</p>
                <p><strong>User Agents:</strong></p>
                <ul>
                    ${userAgents.map(ua => `<li>${ua}</li>`).join('')}
                </ul>
            `;

            analyticsContainer.appendChild(card);
        });
    }
});