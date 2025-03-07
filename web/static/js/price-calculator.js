function updatePrice() {
    let price = document.getElementById('price')
    let session_type = document.getElementById('session-type')
    let duration = document.getElementById('duration')

    if (selectedRoom) {
        switch (session_type.value) {
        case 'solo':
            price.textContent = `£${(parseInt(duration.value, 10) * 6.50).toFixed(2)}`
            break;
        case 'band':
            if (duration.value > 9) {
                price.textContent = `£100.00`
            } else if (duration.value > 3) {
                price.textContent = `£${(parseInt(duration.value, 10) * 10.00).toFixed(2)}`
            } else {
                price.textContent = `£${(parseInt(duration.value, 10) * 12.00).toFixed(2)}`
            }
            break;
        }
    } else {
        price.textContent = '£0.00'
    }
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

    // event listeners for input changes
    session_type.addEventListener('change', updatePrice)
    duration.addEventListener('change', updatePrice)

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
