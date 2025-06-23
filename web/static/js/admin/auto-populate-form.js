// Populate form fields when a payment method is selected
function internalBookingsOveride(selectElement) {
    const selectedOption = selectElement.options[selectElement.selectedIndex];
    const cymbalsElement = document.getElementById('cymbals');

    // If 'Internal Bookings', 'Weekly Regulars', or 'Last Minute Cancellation' is selected,
    // set the status to 'paid' and revised price to 0
    switch (selectedOption.value) {
        case 'internal':
        case 'regulars':
            document.getElementById('status').value = 'paid';
            document.getElementById('revised-price').value = 0;
            break;
        case 'last-minute-cancellation':
            if (cymbalsElement.checked) {
                document.getElementById('status').value = 'unpaid';
                document.getElementById('revised-price').value = 3;
            }
            break;
        default:
            break;       
    }
}

// Filter users in the dropdown based on search input
function filterUsers() {
    const searchInput = document.getElementById('user-search').value.toLowerCase();
    const dropdown = document.getElementById('user-select');
    const options = dropdown.options;

    for (let i = 0; i < options.length; i++) {
        const optionText = options[i].textContent.toLowerCase();
        options[i].style.display = optionText.includes(searchInput) ? '' : 'none';
    }
}

// Populate form fields when a user is selected
function populateUserDetails(selectElement) {
    const selectedOption = selectElement.options[selectElement.selectedIndex];
    document.getElementById('name').value = selectedOption.dataset.name || '';
    document.getElementById('email').value = selectedOption.dataset.email || '';
    document.getElementById('phone').value = selectedOption.dataset.phone || '';
}