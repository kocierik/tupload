document.addEventListener('DOMContentLoaded', () => {
    const dropZone = document.getElementById('dropZone');
    const fileInput = document.getElementById('fileInput');
    const resultContainer = document.getElementById('resultContainer');
    const downloadLink = document.getElementById('downloadLink');

    // Handle drag and drop
    dropZone.addEventListener('dragover', (e) => {
        e.preventDefault();
        dropZone.classList.add('dragover');
    });

    dropZone.addEventListener('dragleave', () => {
        dropZone.classList.remove('dragover');
    });

    dropZone.addEventListener('drop', (e) => {
        e.preventDefault();
        dropZone.classList.remove('dragover');
        const file = e.dataTransfer.files[0];
        if (file) {
            uploadFile(file);
        }
    });

    // Handle file input
    dropZone.addEventListener('click', () => {
        fileInput.click();
    });

    fileInput.addEventListener('change', (e) => {
        const file = e.target.files[0];
        if (file) {
            uploadFile(file);
        }
    });

    // Upload file function
    function uploadFile(file) {
        const formData = new FormData();
        formData.append('file', file);

        fetch('/upload', {
            method: 'PUT',
            body: file
        })
        .then(response => response.text())
        .then(text => {
            // Extract download URL from response
            const match = text.match(/wget (https:\/\/[^\s]+)/);
            if (match) {
                const url = match[1];
                downloadLink.value = `wget ${url}`;
                resultContainer.style.display = 'block';
                dropZone.style.display = 'none';
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Error uploading file. Please try again.');
        });
    }
});

// Copy link function
function copyLink() {
    const downloadLink = document.getElementById('downloadLink');
    downloadLink.select();
    document.execCommand('copy');
    
    const button = document.querySelector('.download-link button');
    const originalText = button.textContent;
    button.textContent = 'Copied!';
    setTimeout(() => {
        button.textContent = originalText;
    }, 2000);
} 