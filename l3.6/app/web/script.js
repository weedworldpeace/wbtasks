let currentPage = 0;
const limit = 20;
let currentTransactions = [];
let editingId = null;

document.addEventListener('DOMContentLoaded', function() {
    const today = new Date().toISOString().split('T')[0];
    document.getElementById('fromDate').value = today;
    document.getElementById('toDate').value = today;
    
    loadData();
});

async function loadData() {
    currentPage = 0;
    document.getElementById('transactionsList').innerHTML = '';
    document.getElementById('loadMoreBtn').style.display = 'none';
    
    await Promise.all([
        loadAnalytics(),
        loadTransactions()
    ]);
}

async function loadAnalytics() {
    const from = document.getElementById('fromDate').value;
    const to = document.getElementById('toDate').value;
    
    try {
        const response = await fetch(`/api/v1/analytics?from=${from}T00:00:00Z&to=${to}T23:59:59Z`);
        const data = await response.json();
        
        document.getElementById('totalCount').textContent = data.total || 0;
        
        document.getElementById('sumIncome').textContent = data.sum?.income?.toFixed(2) || '0';
        document.getElementById('sumExpense').textContent = data.sum?.expense?.toFixed(2) || '0';
        document.getElementById('sumTotal').textContent = data.sum?.amount?.toFixed(2) || '0';
        
        document.getElementById('avgIncome').textContent = data.average?.income?.toFixed(2) || '0';
        document.getElementById('avgExpense').textContent = data.average?.expense?.toFixed(2) || '0';
        document.getElementById('avgTotal').textContent = data.average?.amount?.toFixed(2) || '0';
        
        document.getElementById('medianIncome').textContent = data.median?.income?.toFixed(2) || '0';
        document.getElementById('medianExpense').textContent = data.median?.expense?.toFixed(2) || '0';
        document.getElementById('medianTotal').textContent = data.median?.amount?.toFixed(2) || '0';
        
        document.getElementById('perc90Income').textContent = data.percentile_90?.income?.toFixed(2) || '0';
        document.getElementById('perc90Expense').textContent = data.percentile_90?.expense?.toFixed(2) || '0';
        document.getElementById('perc90Total').textContent = data.percentile_90?.amount?.toFixed(2) || '0';
        
    } catch (error) {
        console.error('Error loading analytics:', error);
    }
}

async function loadTransactions() {
    const from = document.getElementById('fromDate').value;
    const to = document.getElementById('toDate').value;
    const offset = currentPage * limit;
    
    try {
        const response = await fetch(`/api/v1/transactions?from=${from}T00:00:00Z&to=${to}T23:59:59Z&limit=${limit}&offset=${offset}`);
        const transactions = await response.json();
        
        if (transactions.length > 0) {
            currentTransactions = [...currentTransactions, ...transactions];
            displayTransactions(transactions);
            currentPage++;
            document.getElementById('loadMoreBtn').style.display = 'block';
        } else {
            document.getElementById('loadMoreBtn').style.display = 'none';
        }
        
    } catch (error) {
        console.error('Error loading transactions:', error);
    }
}

function displayTransactions(transactions) {
    const container = document.getElementById('transactionsList');
    
    transactions.forEach(transaction => {
        const div = document.createElement('div');
        div.className = `transaction-item ${transaction.type}`;
        div.dataset.id = transaction.id;
        
        div.innerHTML = `
            <div class="transaction-info">
                <div class="transaction-amount">$${transaction.amount.toFixed(2)}</div>
                <div class="transaction-details">
                    <span class="transaction-category">${transaction.category}</span>
                    <span>${transaction.description || 'No description'}</span>
                    <span>${new Date(transaction.created_at).toLocaleString()}</span>
                </div>
            </div>
            <div class="transaction-actions">
                ${editingId === transaction.id ? `
                    <input type="text" class="edit-input" id="edit-category-${transaction.id}" value="${transaction.category}" placeholder="Category">
                    <input type="text" class="edit-input" id="edit-desc-${transaction.id}" value="${transaction.description || ''}" placeholder="Description">
                    <button class="btn btn-save" onclick="saveEdit('${transaction.id}')">Save</button>
                ` : `
                    <button class="btn btn-edit" onclick="startEdit('${transaction.id}')">Edit</button>
                `}
                <button class="btn btn-delete" onclick="deleteTransaction('${transaction.id}')">Delete</button>
            </div>
        `;
        
        container.appendChild(div);
    });
}

function startEdit(id) {
    editingId = id;
    loadData();
}

async function saveEdit(id) {
    const category = document.getElementById(`edit-category-${id}`).value;
    const description = document.getElementById(`edit-desc-${id}`).value;
    const userId = '123e4567-e89b-12d3-a456-426614174000'; // заглушка(предполагается получение от сервиса авторизации)
    
    try {
        const response = await fetch(`/api/v1/transactions/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                category: category,
                description: description,
                user_id: userId
            })
        });
        
        if (response.ok) {
            editingId = null;
            currentPage = 0;
            currentTransactions = [];
            document.getElementById('transactionsList').innerHTML = '';
            await Promise.all([
                loadAnalytics(),
                loadTransactions()
            ]);
        }
    } catch (error) {
        console.error('Error updating transaction:', error);
    }
}

async function createTransaction() {
    const amount = parseFloat(document.getElementById('newAmount').value);
    const type = document.getElementById('newType').value;
    const category = document.getElementById('newCategory').value;
    const description = document.getElementById('newDescription').value;
    
    if (!amount || !category) {
        alert('Please fill in amount and category');
        return;
    }
    
    const userId = '123e4567-e89b-12d3-a456-426614174000'; // заглушка(предполагается получение от сервиса авторизации)
    
    try {
        const response = await fetch('/api/v1/transactions', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                user_id: userId,
                amount: amount,
                type: type,
                category: category,
                description: description
            })
        });
        
        if (response.ok) {
            document.getElementById('newAmount').value = '';
            document.getElementById('newType').value = 'income';
            document.getElementById('newCategory').value = '';
            document.getElementById('newDescription').value = '';
            
            currentPage = 0;
            currentTransactions = [];
            document.getElementById('transactionsList').innerHTML = '';
            await Promise.all([
                loadAnalytics(),
                loadTransactions()
            ]);
        }
    } catch (error) {
        console.error('Error creating transaction:', error);
    }
}

async function deleteTransaction(id) {
    try {
        const response = await fetch(`/api/v1/transactions/${id}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            currentPage = 0;
            currentTransactions = [];
            document.getElementById('transactionsList').innerHTML = '';
            await Promise.all([
                loadAnalytics(),
                loadTransactions()
            ]);
        }
    } catch (error) {
        console.error('Error deleting transaction:', error);
    }
}

function loadMore() {
    loadTransactions();
}