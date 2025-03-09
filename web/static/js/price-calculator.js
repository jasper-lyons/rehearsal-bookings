function updatePrice() {
    let price = document.getElementById('price')
    let session_type = document.getElementById('session-type')
    let duration = parseInt(document.getElementById('duration').value, 10)
    let cymbals = document.getElementById('cymbals')
    let priceValue = 0.00

    if (selectedRoom) {
        switch (session_type.value) {
        case 'solo':
            priceValue = 6.50 * duration
            break;
        case 'band':
            if (duration > 9) {
                priceValue = 100.00
            } else if (duration > 3) {
                priceValue = 10.00 * duration
            } else {
                priceValue = 12.00 * duration
            }
            break;
        }
    }

    if (cymbals.checked) {
        priceValue += 3.00
    }

    price.textContent = `£${(priceValue).toFixed(2)}`
}

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

window.addEventListener('load', function () {
    let session_type = document.getElementById('session-type')
    let duration = document.getElementById('duration')
    let buttons = document.querySelectorAll('.toggle-button');
    let cymbals = document.getElementById('cymbals')

    // event listeners for input changes
    session_type.addEventListener('change', updatePrice)
    duration.addEventListener('change', updatePrice)
    cymbals.addEventListener('change', updatePrice)
    cymbals.addEventListener('change', populateSummary)

    // event listeners for toggle buttons
    buttons.forEach(button => {
        button.addEventListener('click', function (e) {
            e.preventDefault(); // Prevent form submission
            selectButton(button);
            session_type.dispatchEvent(new Event('change')); // Trigger price update
        });
    });

    updatePrice()
})
