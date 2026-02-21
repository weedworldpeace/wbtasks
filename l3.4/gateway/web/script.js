class ImageProcessor {
     constructor() {
        this.apiBaseUrl = 'http://localhost:8080/api/v1';
        this.currentTaskId = null;
        this.init();
    }

    init() {
        this.bindElements();
        this.bindEvents();
    }

    bindElements() {
        this.uploadForm = document.getElementById('uploadForm');
        this.imageInput = document.getElementById('imageInput');
        this.fileName = document.getElementById('fileName');
        this.uploadBtn = document.getElementById('uploadBtn');
        this.uploadResult = document.getElementById('uploadResult');
        this.taskId = document.getElementById('taskId');
        this.copyIdBtn = document.getElementById('copyIdBtn');
        this.checkTaskId = document.getElementById('checkTaskId');
        this.checkBtn = document.getElementById('checkBtn');
        this.loadingIndicator = document.getElementById('loadingIndicator');
        this.errorMessage = document.getElementById('errorMessage');
        this.resultsSection = document.getElementById('resultsSection');
        this.taskStatus = document.getElementById('taskStatus');
        this.imagesGrid = document.getElementById('imagesGrid');
        this.deleteBtn = document.getElementById('deleteBtn');
    }

        bindEvents() {
        this.imageInput.addEventListener('change', (e) => this.handleFileSelect(e));
        this.uploadForm.addEventListener('submit', (e) => this.handleUpload(e));
        this.checkBtn.addEventListener('click', () => this.checkStatus());
        this.checkTaskId.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.checkStatus();
        });
        this.copyIdBtn.addEventListener('click', () => this.copyTaskId());
        
        if (this.deleteBtn) {
            this.deleteBtn.addEventListener('click', () => this.deleteTask());
        }
    }

    handleFileSelect(event) {
        const file = event.target.files[0];
        if (file) {
            this.fileName.textContent = file.name;
            this.uploadBtn.disabled = false;
            
            if (file.size > 10 * 1024 * 1024) {
                this.showError('Файл слишком большой. Максимальный размер 10MB');
                this.uploadBtn.disabled = true;
            }
            
            if (!file.type.startsWith('image/')) {
                this.showError('Пожалуйста, выберите изображение');
                this.uploadBtn.disabled = true;
            }
        } else {
            this.fileName.textContent = 'Файл не выбран';
            this.uploadBtn.disabled = true;
        }
    }

    async handleUpload(event) {
        event.preventDefault();
        
        const file = this.imageInput.files[0];
        if (!file) {
            this.showError('Выберите файл для загрузки');
            return;
        }

        const formData = new FormData();
        formData.append('image', file);

        this.showLoading();
        this.hideError();
        this.hideResults();

        try {
            const response = await fetch(`${this.apiBaseUrl}/upload`, {
                method: 'POST',
                body: formData
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || 'Ошибка загрузки');
            }

            this.showUploadResult(data.id);
        } catch (error) {
            this.showError(error.message);
        } finally {
            this.hideLoading();
        }
    }

    async deleteTask() {
        if (!this.currentTaskId) return;
        
        this.showLoading();
        
        try {
            const response = await fetch(`${this.apiBaseUrl}/image/${this.currentTaskId}`, {
                method: 'DELETE'
            });
            
            if (!response.ok) {
                throw new Error('Ошибка удаления');
            }
            
            this.resultsSection.classList.add('hidden');
            this.deleteBtn.classList.add('hidden');
            this.checkTaskId.value = '';
            
            alert('✅ Задача успешно удалена');
            
        } catch (error) {
            alert('❌ Ошибка: ' + error.message);
        } finally {
            this.hideLoading();
        }
    }

    async checkStatus() {
        const taskId = this.checkTaskId.value.trim();
        
        if (!taskId) {
            this.showError('Введите ID задачи');
            return;
        }

        const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
        if (!uuidRegex.test(taskId)) {
            this.showError('Неверный формат ID');
            return;
        }

        this.showLoading();
        this.hideError();

        try {
            const res = await fetch(`${this.apiBaseUrl}/image/${taskId}`);
            const data = await res.json();

            if (!res.ok) {
                console.log('Ошибка при получении данных:', data);
                if (data.error == 'pending') {
                    throw new Error('Изображения еще не готовы. Пожалуйста, попробуйте позже.');
                } else {
                    throw new Error('Ошибка получения данных');
                }
            }

            this.currentTaskId = taskId;
            this.displayResults(data);
            
        } catch (error) {
            this.showError(error.message);
        } finally {
            this.hideLoading();
        }
    }

    showUploadResult(taskId) {
        this.taskId.textContent = taskId;
        this.uploadResult.classList.remove('hidden');
        this.checkTaskId.value = taskId;
    }

    displayResults(data) {
        this.resultsSection.classList.remove('hidden');
        
        this.deleteBtn.classList.remove('hidden');
        
        this.taskStatus.textContent = data.status;
        this.taskStatus.className = `status-badge ${data.status}`;
        
        this.imagesGrid.innerHTML = '';
        
        const images = [
            { type: 'original', title: 'Оригинал', url: data.original_url },
            { type: 'processed', title: 'С водяным знаком', url: data.watermarked_url },
            { type: 'thumbnail', title: 'Миниатюра', url: data.thumbnail_url }
        ];

        images.forEach(img => {
            if (img.url) {
                const card = this.createImageCard(img.title, img.url);
                this.imagesGrid.appendChild(card);
            }
        });
    }

    createImageCard(title, url) {
        const card = document.createElement('div');
        card.className = 'image-card';
        
        card.innerHTML = `
            <h3>${title}</h3>
            <img src="${url}" alt="${title}" loading="lazy">
            <a href="${url}" target="_blank" download>Скачать</a>
        `;
        
        return card;
    }

    copyTaskId() {
        const taskId = this.taskId.textContent;
        navigator.clipboard.writeText(taskId).then(() => {
            this.copyIdBtn.textContent = '✓';
            setTimeout(() => {
                this.copyIdBtn.textContent = '📋';
            }, 2000);
        });
    }

    showLoading() {
        this.loadingIndicator.classList.remove('hidden');
    }

    hideLoading() {
        this.loadingIndicator.classList.add('hidden');
    }

    showError(message) {
        this.errorMessage.textContent = message;
        this.errorMessage.classList.remove('hidden');
    }

    hideError() {
        this.errorMessage.classList.add('hidden');
    }

    hideResults() {
        this.resultsSection.classList.add('hidden');
    }
}

document.addEventListener('DOMContentLoaded', () => {
    new ImageProcessor();
});