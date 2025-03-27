function selectButton(button) {
    // Deselect all buttons
    const buttons = document.querySelectorAll('.toggle-button');
    const infoBox = document.getElementById('info-box');
    buttons.forEach(btn => btn.classList.remove('selected'));

    // Select the clicked button
    button.classList.add('selected');

    if (infoBox) {
        // Update and show the information box content
        if (button.dataset.value === 'band') {
            infoBox.innerHTML = `
                    <p>Rehearsal session for up to six people.</p>
					<p class="explainer-text">£12 per hour for bookings under 4 hours, £10 per hour for 4+ hours, £100 for all day bookings (9+ hours) </p>
                    `;
        } else if (button.dataset.value === 'solo') {
            infoBox.innerHTML = `
                    <p align=right>Rehearsal session for one person</p>
					<p align=right class="explainer-text">£6.50 per hour. Solo bookings are not available on weekday evenings after 6pm</p>
                    `;
        }
    }

    // Save the selected value to type input
    document.getElementById('session-type').value = button.getAttribute('data-value');
    clearSelection(); // clear the time slot selection when session type changes
}