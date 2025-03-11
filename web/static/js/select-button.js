function selectButton(button) {
    // Deselect all buttons
    const buttons = document.querySelectorAll('.toggle-button');
    const infoBox = document.getElementById('info-box');
    buttons.forEach(btn => btn.classList.remove('selected'));

    // Select the clicked button
    button.classList.add('selected');

    // Update and show the information box content
    if (button.dataset.value === 'band') {
        infoBox.innerHTML = '<p>Rehearsal session for up to six people</p>';
    } else if (button.dataset.value === 'solo') {
        infoBox.innerHTML = '<p align=right>Rehearsal session for one person</p>';
    }

    // Save the selected value to type input
    document.getElementById('session-type').value = button.getAttribute('data-value');
    clearSelection(); // clear the time slot selection when session type changes
}