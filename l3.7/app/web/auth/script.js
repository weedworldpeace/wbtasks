async function getToken() {
    const role = document.getElementById('roleSelect').value;
    const button = document.querySelector('.auth-btn');
    const message = document.getElementById('message');
    
    button.classList.add('loading');
    button.textContent = 'Getting token...';
    message.style.display = 'none';
    
    try {
        const response = await fetch(`/api/v1/auth?role=${role}`, {
            method: 'GET',
            credentials: 'include'
        });
        
        if (response.ok) {
            message.className = 'message success';
            message.textContent = 'Token received! Redirecting to home...';
            
            setTimeout(() => {
                window.location.href = '/home';
            }, 1000);
        } else {
            const error = await response.json();
            throw new Error(error.error || 'Failed to get token');
        }
    } catch (error) {
        message.className = 'message error';
        message.textContent = error.message || 'Error getting token';
    } finally {
        button.classList.remove('loading');
        button.textContent = 'Get Access';
    }
}