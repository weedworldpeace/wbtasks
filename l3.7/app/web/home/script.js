let itemsPage = 0;
let historyPage = 0;
const limit = 20;
let editingItemId = null;

async function handleResponse(response) {
    if (response.status === 401) {
        window.location.href = '/auth';
        return null;
    }
    if (!response.ok && response.status != 403) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
}

document.addEventListener('DOMContentLoaded', function() {
    loadItems();
    loadHistory();
});

function toggleCreate() {
    const form = document.getElementById('createForm');
    form.style.display = form.style.display === 'none' ? 'block' : 'none';
}

async function createItem() {
    const name = document.getElementById('newName').value;
    const description = document.getElementById('newDescription').value;
    const quantity = parseInt(document.getElementById('newQuantity').value);
    const price = parseFloat(document.getElementById('newPrice').value);

    if (!name || !quantity || !price) {
        alert('Please fill all required fields');
        return;
    }

    try {
        const response = await fetch('/api/v1/items', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({ name, description, quantity, price })
        });

        const data = await handleResponse(response);
        if (data === null) return;

        if (response.ok) {
            document.getElementById('newName').value = '';
            document.getElementById('newDescription').value = '';
            document.getElementById('newQuantity').value = '';
            document.getElementById('newPrice').value = '';
            
            itemsPage = 0;
            document.getElementById('itemsList').innerHTML = '';
            loadItems();
            toggleCreate();
        } else {
            alert("no permission")
        }
    } catch (error) {
        console.error('Error:', error);
    }
}

async function loadItems() {
    const offset = itemsPage * limit;
    
    try {
        const response = await fetch(`/api/v1/items?limit=${limit}&offset=${offset}`, {
            credentials: 'include'
        });
        const items = await handleResponse(response);
        if (items === null) return;
        
        if (items.length > 0) {
            displayItems(items);
            itemsPage++;
            document.getElementById('loadMoreItemsBtn').style.display = 'block';
        } else {
            document.getElementById('loadMoreItemsBtn').style.display = 'none';
        }
    } catch (error) {
        console.error('Error:', error);
    }
}

function displayItems(items) {
    const container = document.getElementById('itemsList');
    
    items.forEach(item => {
        const div = document.createElement('div');
        div.className = 'item-row';
        div.dataset.id = item.id;
        
        if (editingItemId === item.id) {
            div.innerHTML = `
                <div class="item-fields">
                    <input type="text" id="edit-name-${item.id}" value="${item.name || ''}" placeholder="Name">
                    <input type="text" id="edit-desc-${item.id}" value="${item.description || ''}" placeholder="Description">
                    <input type="number" id="edit-qty-${item.id}" value="${item.quantity || 0}" placeholder="Quantity">
                    <input type="number" id="edit-price-${item.id}" value="${item.price || 0}" step="0.01" placeholder="Price">
                </div>
                <div class="item-actions">
                    <button class="btn btn-update" onclick="updateItem('${item.id}')">Update</button>
                    <button class="btn btn-delete" onclick="deleteItem('${item.id}')">Delete</button>
                </div>
            `;
        } else {
            div.innerHTML = `
                <div class="item-fields">
                    <span><strong>${item.name}</strong></span>
                    <span>ItemId: ${item.id}</span>
                    <span>${item.description || '-'}</span>
                    <span>Quantity: ${item.quantity}</span>
                    <span>Price: $${item.price?.toFixed(2)}</span>
                    <span>created_at: ${new Date(item.created_at).toLocaleString()}</span>
                    <span>updated_at: ${new Date(item.updated_at).toLocaleString()}</span>
                </div>
                <div class="item-actions">
                    <button class="btn btn-update" onclick="startEdit('${item.id}')">Edit</button>
                    <button class="btn btn-delete" onclick="deleteItem('${item.id}')">Delete</button>
                </div>
            `;
        }
        
        container.appendChild(div);
    });
}

function startEdit(id) {
    editingItemId = id;
    itemsPage = 0;
    document.getElementById('itemsList').innerHTML = '';
    loadItems();
}

async function updateItem(id) {
    const name = document.getElementById(`edit-name-${id}`).value;
    const description = document.getElementById(`edit-desc-${id}`).value;
    const quantity = parseInt(document.getElementById(`edit-qty-${id}`).value);
    const price = parseFloat(document.getElementById(`edit-price-${id}`).value);

    try {
        const response = await fetch(`/api/v1/items/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({ name, description, quantity, price })
        });;

        const data = await handleResponse(response);
        if (data === null) return;

        if (response.ok) {
            editingItemId = null;
            itemsPage = 0;
            document.getElementById('itemsList').innerHTML = '';
            loadItems();
        } else {
            alert("no permission")
        }
    } catch (error) {
        console.error('Error:', error);
    }
}

async function deleteItem(id) {
    try {
        const response = await fetch(`/api/v1/items/${id}`, {
            method: 'DELETE',
            credentials: 'include'
        });

        if (response.ok) {
            itemsPage = 0;
            document.getElementById('itemsList').innerHTML = '';
            loadItems();
            loadHistory();
        } else if (response.status == 403) {
            alert("no permission")
        }
    } catch (error) {
        console.error('Error:', error);
    }
}

async function loadHistory() {
    const offset = historyPage * limit;
    
    try {
        const response = await fetch(`/api/v1/history?limit=${limit}&offset=${offset}`, {
            credentials: 'include'
        });
        const history = await handleResponse(response);
        if (history === null) return;
        
        if (history.length > 0) {
            displayHistory(history);
            historyPage++;
            document.getElementById('loadMoreHistoryBtn').style.display = 'block';
        } else {
            document.getElementById('loadMoreHistoryBtn').style.display = 'none';
        }
    } catch (error) {
        console.error('Error:', error);
    }
}

function displayHistory(history) {
    const container = document.getElementById('historyList');
    container.innerHTML = '';

    history.forEach(entry => {
        const div = document.createElement('div');
        div.className = 'item-row';

        let details = '';
        if (entry.action === 'INSERT' && entry.new_data) {
            const d = entry.new_data;
            details = `
                <div style="color: #27ae60">Inserted:</div>
                <div>Name: ${d.name || ''}</div>
                <div>Desc: ${d.description || ''}</div>
                <div>Qty: ${d.quantity || 0}</div>
                <div>Price: $${d.price?.toFixed(2)}</div>
            `;
        } else if (entry.action === 'UPDATE' && entry.old_data && entry.new_data) {
            const old = entry.old_data;
            const New = entry.new_data;
            details = `
                <div style="color: #e74c3c">Old: ${old.name || ''} | ${old.description || ''} | Qty: ${old.quantity || 0} | $${old.price?.toFixed(2)}</div>
                <div style="color: #27ae60">New: ${New.name || ''} | ${New.description || ''} | Qty: ${New.quantity || 0} | $${New.price?.toFixed(2)}</div>
            `;
        } else if (entry.action === 'DELETE' && entry.old_data) {
            const d = entry.old_data;
            details = `
                <div style="color: #e74c3c">Deleted:</div>
                <div>Name: ${d.name || ''}</div>
                <div>Desc: ${d.description || ''}</div>
                <div>Qty: ${d.quantity || 0}</div>
                <div>Price: $${d.price?.toFixed(2)}</div>
            `;
        } else {
            details = 'No details';
        }

        div.innerHTML = `
            <div class="item-fields" style="flex-direction: column; align-items: flex-start;">
                <div style="display: flex; gap: 20px; margin-bottom: 10px;">
                    <span><strong>${entry.action}</strong></span>
                    <span>Item: ${entry.item_id}</span>
                    <span>Role: ${entry.user_role}</span>
                    <span>${new Date(entry.changed_at).toLocaleString()}</span>
                </div>
                <div style="background: #f5f5f5; padding: 10px; border-radius: 4px; width: 100%;">
                    ${details}
                </div>
            </div>
        `;
        container.appendChild(div);
    });
}

function loadMoreItems() {
    loadItems();
}

function loadMoreHistory() {
    loadHistory();
}