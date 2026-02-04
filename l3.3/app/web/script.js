let currentPage = 1;
let totalPages = 1;
let totalComments = 0;
let limit = 10;
let currentTreeCommentId = null; 
let currentSearch = ''; 

document.addEventListener('DOMContentLoaded', () => {
    loadComments();
});

async function loadComments(page = 1) {
    try {
        showMessage('Загрузка...');
        
        let url = `/comments?page=${page}&limit=${limit}`;
        
        if (currentSearch) {
            url += `&query=${encodeURIComponent(currentSearch)}`;
        }
        
        const response = await fetch(url);
        
        if (!response.ok) {
            showError('Ошибка загрузки комментариев');
            return;
        }
        
        const data = await response.json();
        
        currentPage = data.page || page;
        totalPages = data.total_pages || 1;
        totalComments = data.total || 0;
        
        if (currentTreeCommentId) {
            await displayCommentsWithTree(data.comments);
        } else {
            displayComments(data.comments);
        }
        
        updatePagination();
        
    } catch (error) {
        showError('Ошибка при загрузке комментариев');
        console.error(error);
    }
}

async function loadCommentTree(commentId) {
    try {
        const response = await fetch(`/comments?parent=${commentId}&limit=1000`);
        
        if (!response.ok) {
            return [];
        }
        
        const data = await response.json();
        return data.comments || [];
        
    } catch (error) {
        console.error('Ошибка при загрузке дерева:', error);
        return [];
    }
}

async function displayCommentsWithTree(rootComments) {
    const container = document.getElementById('commentsTree');
    
    if (!rootComments || rootComments.length === 0) {
        container.innerHTML = '<div class="empty">Нет комментариев</div>';
        return;
    }
    
    const mainComment = rootComments.find(c => c.id === currentTreeCommentId);
    if (!mainComment) {
        currentTreeCommentId = null;
        loadComments();
        return;
    }
    
    const allTreeComments = await loadCommentTree(currentTreeCommentId);
    
    let html = '';
    
    rootComments.forEach(comment => {
        if (comment.id === currentTreeCommentId) {
            html += renderCommentWithFullTree(comment, allTreeComments);
        } else {
            html += renderComment(comment, 0, false);
        }
    });
    
    container.innerHTML = html;
}

function displayComments(comments) {
    const container = document.getElementById('commentsTree');
    
    if (!comments || comments.length === 0) {
        container.innerHTML = '<div class="empty">Нет комментариев</div>';
        return;
    }
    
    let html = '';
    
    comments.forEach(comment => {
        html += renderComment(comment, 0, false);
    });
    
    container.innerHTML = html;
}

function renderComment(comment, depth = 0, isChild = false) {
    const marginLeft = depth * 30;
    const isExpanded = comment.id === currentTreeCommentId;
    
    return `
        <div class="comment ${isChild ? 'child-comment' : ''}" 
             style="margin-left: ${marginLeft}px" 
             data-id="${comment.id}">
            <div class="comment-content">
                ${escapeHtml(comment.content)}
            </div>
            <div class="comment-meta">
                <span>${formatDate(comment.created_at)}</span>
                <span>#${comment.id.substring(0, 8)}</span>
            </div>
            <div class="comment-actions">
                ${!isChild ? `
                    <button class="tree-btn" onclick="toggleCommentTree('${comment.id}')">
                        ${isExpanded ? 'Скрыть дерево' : 'Показать дерево'}
                    </button>
                ` : ''}
                <button class="reply-btn" onclick="replyToComment('${comment.id}')">
                    Ответить
                </button>
                <button class="delete-btn" onclick="deleteComment('${comment.id}')">
                    Удалить
                </button>
            </div>
        </div>
    `;
}

function renderCommentWithFullTree(comment, allTreeComments) {
    let html = '';
    
    const sortedTreeComments = [...allTreeComments].sort((a, b) => a.path.localeCompare(b.path));
    
    const commentsById = {};
    sortedTreeComments.forEach(c => {
        commentsById[c.id] = c;
    });
    
    function renderTree(currentComment, currentDepth = 0) {
        let result = renderComment(currentComment, currentDepth, currentDepth > 0);
        
        const children = sortedTreeComments.filter(c => 
            c.parent_id === currentComment.id && c.id !== currentComment.id
        );
        
        children.sort((a, b) => new Date(a.created_at) - new Date(b.created_at));
        
        children.forEach(child => {
            result += renderTree(child, currentDepth + 1);
        });
        
        return result;
    }
    
    html = renderTree(comment, 0);
    
    return html;
}

async function toggleCommentTree(commentId) {
    if (currentTreeCommentId === commentId) {
        currentTreeCommentId = null;
        showMessage('Дерево скрыто');
    } else {
        currentTreeCommentId = commentId;
        showMessage('Загрузка дерева...');
    }
    
    loadComments(currentPage);
}

async function addComment() {
    const content = document.getElementById('commentText').value.trim();
    const parentId = document.getElementById('parentId').value.trim();
    
    if (!content) {
        showError('Введите текст комментария');
        return;
    }
    
    const requestBody = {
        content: content
    };
    
    if (parentId) {
        requestBody.parent_id = parentId;
    }
    
    try {
        showMessage('Отправка...');
        
        const response = await fetch('/comments', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody)
        });
        
        if (!response.ok) {
            const error = await response.json().catch(() => ({}));
            showError(error.error || 'Ошибка создания');
            return;
        }
        
        showSuccess('Комментарий добавлен');
        document.getElementById('commentText').value = '';
        document.getElementById('parentId').value = '';
        loadComments(currentPage);
        
    } catch (error) {
        showError('Ошибка при добавлении комментария');
        console.error(error);
    }
}

function replyToComment(commentId) {
    document.getElementById('parentId').value = commentId;
    document.getElementById('commentText').focus();
    showMessage('Введите ответ в поле ниже');
}

async function deleteComment(commentId) {
    try {
        showMessage('Удаление...');
        
        const response = await fetch(`/comments/${commentId}`, {
            method: 'DELETE'
        });
        
        if (!response.ok) {
            showError('Ошибка удаления');
            return;
        }
        
        showSuccess('Комментарий удален');
        loadComments(currentPage);
        
    } catch (error) {
        showError('Ошибка при удалении');
        console.error(error);
    }
}

async function searchComments() {
    const query = document.getElementById('searchInput').value.trim();
    
    if (!query) {
        showError('Введите поисковый запрос');
        return;
    }
    
    try {
        showMessage('Поиск...');
        currentSearch = query;
        currentTreeCommentId = null; 
        currentPage = 1;
        
        const response = await fetch(`/comments?query=${encodeURIComponent(query)}&limit=100`);
        
        if (!response.ok) {
            showError('Ошибка поиска');
            return;
        }
        
        const data = await response.json();
        
        currentPage = data.page || 1;
        totalPages = data.total_pages || 1;
        totalComments = data.total || 0;
        
        displayComments(data.comments);
        
        updatePagination();
        
        showSuccess(`Найдено: ${data.comments ? data.comments.length : 0} комментариев`);
        
    } catch (error) {
        showError('Ошибка поиска');
        console.error(error);
    }
}

function clearSearch() {
    document.getElementById('searchInput').value = '';
    currentSearch = '';
    currentTreeCommentId = null;
    currentPage = 1;
    loadComments();
}

function prevPage() {
    if (currentPage > 1) {
        loadComments(currentPage - 1);
    }
}

function nextPage() {
    if (currentPage < totalPages) {
        loadComments(currentPage + 1);
    }
}

function updatePagination() {
    const pagination = document.getElementById('pagination');
    const pageInfo = document.getElementById('pageInfo');
    const prevBtn = document.getElementById('prevBtn');
    const nextBtn = document.getElementById('nextBtn');
    
    if (!pagination || !pageInfo || !prevBtn || !nextBtn) return;
    
    if (totalPages > 1 && !currentSearch) {
        pagination.style.display = 'flex';
        pageInfo.textContent = `Страница ${currentPage} из ${totalPages}`;
        
        prevBtn.disabled = currentPage <= 1;
        nextBtn.disabled = currentPage >= totalPages;
    } else {
        pagination.style.display = 'none';
    }
}

function showMessage(text) {
    const messageDiv = document.getElementById('message');
    if (!messageDiv) return;
    
    messageDiv.textContent = text;
    messageDiv.className = 'message';
    messageDiv.style.display = 'block';
    
    setTimeout(() => {
        if (messageDiv.textContent === text) {
            messageDiv.style.display = 'none';
        }
    }, 2000);
}

function showSuccess(text) {
    const messageDiv = document.getElementById('message');
    if (!messageDiv) return;
    
    messageDiv.textContent = text;
    messageDiv.className = 'message success';
    messageDiv.style.display = 'block';
    
    setTimeout(() => {
        if (messageDiv.textContent === text) {
            messageDiv.style.display = 'none';
        }
    }, 2000);
}

function showError(text) {
    const messageDiv = document.getElementById('message');
    if (!messageDiv) return;
    
    messageDiv.textContent = text;
    messageDiv.className = 'message error';
    messageDiv.style.display = 'block';
    
    setTimeout(() => {
        if (messageDiv.textContent === text) {
            messageDiv.style.display = 'none';
        }
    }, 3000);
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    try {
        const date = new Date(dateString);
        return date.toLocaleDateString('ru-RU') + ', ' + 
               date.toLocaleTimeString('ru-RU', {hour: '2-digit', minute:'2-digit'});
    } catch (e) {
        return dateString;
    }
}

document.getElementById('searchInput').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') searchComments();
});

document.getElementById('commentText').addEventListener('keydown', (e) => {
    if (e.ctrlKey && e.key === 'Enter') addComment();
});