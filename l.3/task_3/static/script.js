const API_BASE = '';

document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('load-comments-btn').addEventListener('click', loadComments);
    document.getElementById('search-btn').addEventListener('click', searchComments);
    document.getElementById('clear-search-btn').addEventListener('click', loadComments);
    document.getElementById('create-comment-btn').addEventListener('click', createComment);
});

async function loadComments() {
    try {
        const response = await fetch(`${API_BASE}/comments/all`);
        if (!response.ok) throw new Error('Failed to load comments');
        const comments = await response.json();
        const container = document.getElementById('comments-container');
        container.innerHTML = ''; // Clear container before displaying
        displayComments(comments, container, 0);
    } catch (error) {
        alert('Error loading comments: ' + error.message);
    }
}

function displayComments(comments, container, level) {
    if (!comments || comments.length === 0) {
        container.innerHTML = '<p>No comments found.</p>';
        return;
    }

    comments.forEach(comment => {
        const commentDiv = document.createElement('div');
        commentDiv.className = 'comment';
        commentDiv.style.marginLeft = `${level * 20}px`;

        commentDiv.innerHTML = `
            <div class="text">${escapeHtml(comment.text)}</div>
            <div class="meta">ID: ${comment.id} | Created: ${new Date(comment.created_at).toLocaleString()}</div>
            <div class="actions">
                <button class="reply-btn" data-id="${comment.id}">Reply</button>
                <button class="delete-btn" data-id="${comment.id}">Delete</button>
            </div>
            <div class="reply-form" id="reply-form-${comment.id}">
                <textarea placeholder="Write your reply..."></textarea>
                <button class="submit-reply-btn" data-parent-id="${comment.id}">Submit Reply</button>
                <button class="cancel-reply-btn" data-id="${comment.id}">Cancel</button>
            </div>
        `;

        container.appendChild(commentDiv);

        // Add event listeners
        commentDiv.querySelector('.reply-btn').addEventListener('click', toggleReplyForm);
        commentDiv.querySelector('.delete-btn').addEventListener('click', deleteComment);
        commentDiv.querySelector('.submit-reply-btn').addEventListener('click', submitReply);
        commentDiv.querySelector('.cancel-reply-btn').addEventListener('click', cancelReply);

        // Recursively display children
        if (comment.children && comment.children.length > 0) {
            displayComments(comment.children, container, level + 1);
        }
    });
}

function toggleReplyForm(event) {
    const id = event.target.dataset.id;
    const form = document.getElementById(`reply-form-${id}`);
    form.style.display = form.style.display === 'block' ? 'none' : 'block';
}

async function submitReply(event) {
    const parentId = event.target.dataset.parentId;
    const textarea = event.target.previousElementSibling;
    const text = textarea.value.trim();
    if (!text) return alert('Please enter a reply.');

    try {
        const response = await fetch(`${API_BASE}/comments`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ parent_id: parseInt(parentId), text })
        });
        if (!response.ok) throw new Error('Failed to create reply');
        textarea.value = '';
        document.getElementById(`reply-form-${parentId}`).style.display = 'none';
        loadComments();
    } catch (error) {
        alert('Error creating reply: ' + error.message);
    }
}

function cancelReply(event) {
    const id = event.target.dataset.id;
    document.getElementById(`reply-form-${id}`).style.display = 'none';
}

async function deleteComment(event) {
    const id = event.target.dataset.id;
    if (!confirm('Are you sure you want to delete this comment?')) return;

    try {
        const response = await fetch(`${API_BASE}/comments/${id}`, { method: 'DELETE' });
        if (!response.ok) throw new Error('Failed to delete comment');
        loadComments();
    } catch (error) {
        alert('Error deleting comment: ' + error.message);
    }
}

async function createComment() {
    const text = document.getElementById('new-comment-text').value.trim();
    const parentIdStr = document.getElementById('new-comment-parent-id').value.trim();
    if (!text) return alert('Please enter a comment.');

    const body = { text };
    if (parentIdStr) {
        const parentId = parseInt(parentIdStr);
        if (isNaN(parentId)) return alert('Parent ID must be a number.');
        body.parent_id = parentId;
    }

    try {
        const response = await fetch(`${API_BASE}/comments`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });
        if (!response.ok) throw new Error('Failed to create comment');
        document.getElementById('new-comment-text').value = '';
        document.getElementById('new-comment-parent-id').value = '';
        loadComments();
    } catch (error) {
        alert('Error creating comment: ' + error.message);
    }
}

async function searchComments() {
    const query = document.getElementById('search-input').value.trim();
    if (!query) return alert('Please enter search text.');

    try {
        const response = await fetch(`${API_BASE}/comments/search`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ text: query })
        });
        if (!response.ok) throw new Error('Failed to search comments');
        const comments = await response.json();
        const container = document.getElementById('comments-container');
        container.innerHTML = ''; // Clear container before displaying
        displayComments(comments, container, 0);
    } catch (error) {
        alert('Error searching comments: ' + error.message);
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}