// State management
let selectedFile = null;

// DOM elements
const fileInput = document.getElementById('fileInput');
const uploadBtn = document.getElementById('uploadBtn');
const uploadSection = document.getElementById('uploadSection');
const selectedFileName = document.getElementById('selectedFileName');
const filesContainer = document.getElementById('filesContainer');
const refreshBtn = document.getElementById('refreshBtn');
const statusMessage = document.getElementById('statusMessage');

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    loadFiles();
    setupEventListeners();
});

// Event listeners
function setupEventListeners() {
    fileInput.addEventListener('change', handleFileSelect);
    uploadBtn.addEventListener('click', handleUpload);
    refreshBtn.addEventListener('click', loadFiles);

    // Drag and drop
    uploadSection.addEventListener('dragover', handleDragOver);
    uploadSection.addEventListener('dragleave', handleDragLeave);
    uploadSection.addEventListener('drop', handleDrop);
}

// File selection
function handleFileSelect(event) {
    const file = event.target.files[0];
    if (file) {
        selectedFile = file;
        showSelectedFile(file.name);
    }
}

function showSelectedFile(filename) {
    selectedFileName.textContent = `Selected: ${filename}`;
    selectedFileName.style.display = 'block';
    uploadBtn.style.display = 'inline-block';
}

// Drag and drop handlers
function handleDragOver(event) {
    event.preventDefault();
    uploadSection.classList.add('dragover');
}

function handleDragLeave(event) {
    event.preventDefault();
    uploadSection.classList.remove('dragover');
}

function handleDrop(event) {
    event.preventDefault();
    uploadSection.classList.remove('dragover');

    const file = event.dataTransfer.files[0];
    if (file) {
        selectedFile = file;
        fileInput.files = event.dataTransfer.files;
        showSelectedFile(file.name);
    }
}

// Upload functionality
async function handleUpload() {
    if (!selectedFile) {
        showStatus('Please select a file first', 'error');
        return;
    }

    try {
        uploadBtn.disabled = true;
        uploadBtn.textContent = 'Uploading...';

        // Step 1: Get presigned URL
        const presignResp = await fetch(`/api/upload?filename=${encodeURIComponent(selectedFile.name)}`);

        if (!presignResp.ok) {
            throw new Error('Failed to get upload URL');
        }

        const { url } = await presignResp.json();

        console.log(url);

        // Step 2: Upload directly to S3/MinIO
        const uploadResp = await fetch(url, {
            method: 'PUT',
            body: selectedFile
        });

        if (!uploadResp.ok) {
            throw new Error('Upload failed');
        }

        showStatus(`Successfully uploaded ${selectedFile.name}`, 'success');

        // Reset form
        selectedFile = null;
        fileInput.value = '';
        selectedFileName.style.display = 'none';
        uploadBtn.style.display = 'none';

        // Reload file list
        await loadFiles();

    } catch (error) {
        console.error('Upload error:', error);
        showStatus(`Upload failed: ${error.message}`, 'error');
    } finally {
        uploadBtn.disabled = false;
        uploadBtn.textContent = 'Upload File';
    }
}

// Load and display files
async function loadFiles() {
    try {
        filesContainer.innerHTML = `
            <div class="loading">
                <div class="spinner"></div>
                <p>Loading files...</p>
            </div>
        `;

        const response = await fetch('/api/list?prefix=');

        if (!response.ok) {
            throw new Error('Failed to load files');
        }

        const files = await response.json();
        displayFiles(files);

    } catch (error) {
        console.error('Error loading files:', error);
        filesContainer.innerHTML = `
            <div class="empty-state">
                <p>Failed to load files. Please try again.</p>
            </div>
        `;
    }
}

function displayFiles(files) {
    if (!files || files.length === 0) {
        filesContainer.innerHTML = `
            <div class="empty-state">
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
                <p>No files yet. Upload your first file to get started!</p>
            </div>
        `;
        return;
    }

    const filesList = files.map(filename => `
        <li class="file-item">
            <span class="file-name">${escapeHtml(filename)}</span>
            <button class="download-btn" onclick="downloadFile('${escapeHtml(filename)}')">
                Download
            </button>
        </li>
    `).join('');

    filesContainer.innerHTML = `<ul class="files-list">${filesList}</ul>`;
}

// Download functionality
async function downloadFile(filename) {
    try {
        showStatus(`Preparing ${filename}...`, 'success');

        // 1. Get presigned download URL from your server
        const response = await fetch(`/api/download?filename=${encodeURIComponent(filename)}`);

        if (!response.ok) {
            throw new Error('Failed to get download URL');
        }

        const { url } = await response.json();

        // 2. Fetch the actual file content as a Blob
        // This bypasses the browser's default "open in tab" behavior for images
        const fileResponse = await fetch(url);
        if (!fileResponse.ok) throw new Error('File download failed');
        
        const blob = await fileResponse.blob();

        // 3. Create a temporary object URL and trigger download
        const blobUrl = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = blobUrl;
        link.download = filename; // This attribute forces the download
        
        document.body.appendChild(link);
        link.click();
        
        // 4. Cleanup
        document.body.removeChild(link);
        window.URL.revokeObjectURL(blobUrl);

        showStatus(`${filename} downloaded successfully`, 'success');

    } catch (error) {
        console.error('Download error:', error);
        showStatus(`Download failed: ${error.message}`, 'error');
    }
}

// Status message
function showStatus(message, type) {
    statusMessage.textContent = message;
    statusMessage.className = `status-message ${type}`;

    setTimeout(() => {
        statusMessage.className = 'status-message';
    }, 5000);
}

// Utility function to escape HTML
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, m => map[m]);
}
