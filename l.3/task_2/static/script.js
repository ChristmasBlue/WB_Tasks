document.addEventListener('DOMContentLoaded', function() {
    const shortUrlInput = document.getElementById('shortUrlInput');
    const redirectBtn = document.getElementById('redirectBtn');
    const analyticsBtn = document.getElementById('analyticsBtn');
    const userAgentBtn = document.getElementById('userAgentBtn');
    const dateBtn = document.getElementById('dateBtn');
    const monthBtn = document.getElementById('monthBtn');
    const resultsDiv = document.getElementById('results');

    redirectBtn.addEventListener('click', function() {
        const shortUrl = shortUrlInput.value.trim();
        if (shortUrl) {
            window.location.href = 'http://localhost:8080/s/' + shortUrl;
        } else {
            alert('Please enter a short URL');
        }
    });

    const createBtn = document.getElementById('createBtn');
    createBtn.addEventListener('click', function() {
        const url = shortUrlInput.value.trim();
        if (!url) {
            alert('Please enter a URL');
            return;
        }
        fetch('http://localhost:8080/shorten', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ url: url })
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                displayResults('Shortened URL', data);
            })
            .catch(error => {
                console.error('Error creating short URL:', error);
                displayResults('Error', 'Failed to create short URL');
            });
    });

    analyticsBtn.addEventListener('click', function() {
        const shortUrl = shortUrlInput.value.trim();
        if (shortUrl) {
            window.location.href = 'http://localhost:8080/static/analytics.html?shortUrl=' + encodeURIComponent(shortUrl);
        } else {
            alert('Please enter a short URL');
        }
    });

    userAgentBtn.addEventListener('click', function() {
        window.location.href = 'http://localhost:8080/static/analytics-useragent.html';
    });

    dateBtn.addEventListener('click', function() {
        window.location.href = 'http://localhost:8080/static/analytics-date.html';
    });

    monthBtn.addEventListener('click', function() {
        window.location.href = 'http://localhost:8080/static/analytics-month.html';
    });

    function fetchAnalytics(shortUrl) {
        fetch(`http://localhost:8080/analytics/${shortUrl}`)
            .then(response => response.json())
            .then(data => {
                displayResults('Analytics for ' + shortUrl, data);
            })
            .catch(error => {
                console.error('Error fetching analytics:', error);
                displayResults('Error', 'Failed to fetch analytics');
            });
    }


    function displayResults(title, data) {
        resultsDiv.innerHTML = `<h3>${title}</h3><pre>${JSON.stringify(data, null, 2)}</pre>`;
    }
});