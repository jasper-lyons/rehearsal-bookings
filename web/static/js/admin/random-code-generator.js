// Generate a random 4-digit code
function generateRandomCode() {
    return Math.floor(1000 + Math.random() * 9000); // Generates a number between 1000 and 9999
}

// Set the generated code as the default value of the input field
document.addEventListener('DOMContentLoaded', function () {
    const codeValueInput = document.getElementById('code-value');
    codeValueInput.value = generateRandomCode();
});