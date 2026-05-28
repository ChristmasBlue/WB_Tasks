document.addEventListener('DOMContentLoaded', function() {
    const backBtn = document.getElementById('backBtn');
    const analyticsContainer = document.getElementById('analyticsContainer');

    backBtn.addEventListener('click', function() {
        window.location.href = 'http://localhost:8080/static/index.html';
    });

    const urlParams = new URLSearchParams(window.location.search);
    const shortUrl = urlParams.get('shortUrl');

    if (shortUrl) {
        fetchAnalytics(shortUrl);
    } else {
        analyticsContainer.innerHTML = '<p>No short URL provided.</p>';
    }

    function fetchAnalytics(shortUrl) {
        fetch(`http://localhost:8080/analytics/${shortUrl}`)
            .then(response => response.json())
            .then(data => {
                displayAnalytics(data);
            })
            .catch(error => {
                console.error('Error fetching analytics:', error);
                analyticsContainer.innerHTML = '<p>Failed to fetch analytics.</p>';
            });
    }

    function displayAnalytics(data) {
        if (!Array.isArray(data) || data.length === 0) {
            analyticsContainer.innerHTML = '<p>No analytics data available.</p>';
            return;
        }

        // Flatten all request times from all items
        const allTimes = data.reduce((acc, item) => {
            if (Array.isArray(item.request_time)) {
                return acc.concat(item.request_time);
            }
            return acc;
        }, []);

        // Add timestamps list
        const timeDiv = document.createElement('div');
        timeDiv.innerHTML = '<h3>Time: </h3><ul>' + allTimes.map(time => `<li>${time}</li>`).join('') + '</ul>';
        analyticsContainer.appendChild(timeDiv);

        // Then the cards for user agents
        data.forEach((item, index) => {
            const card = document.createElement('div');
            card.className = 'analytics-card';
            card.style.animationDelay = `${index * 0.1}s`; // Stagger animations

            const userAgent = (item.user_agent && item.user_agent.length > 0) ? item.user_agent.join(', ') : 'Unknown';

            card.innerHTML = `
                <h3>Visit ${index + 1}</h3>
                <p><strong>User Agent:</strong> ${userAgent}</p>
            `;

            analyticsContainer.appendChild(card);
        });
    }
});