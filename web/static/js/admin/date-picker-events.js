function formatDate(date) {
    return date.toISOString().split('T')[0]
}

const urlParams = new URLSearchParams(window.location.search);
const dateParam = urlParams.get('date');
if (dateParam) {
    const newDate = new Date(dateParam);
    if (!isNaN(newDate)) {
        selectedDate = newDate;
    }
}

function updateUrl() {
    const newUrl = new URL(window.location);
    newUrl.searchParams.set('date', formatDate(selectedDate));
    window.history.pushState({}, '', newUrl);
}

// Add event listeners for buttons
document.getElementById('prev-day').addEventListener('click', () => {
    selectedDate.setDate(selectedDate.getDate() - 1); // Move one day back
    updateDatePicker(animated = true, timebound = false).then(() => {
        updateUrl();
        showDateGroup();
        setAvailability();
        location.reload(); // Rerender the page
    });

});

document.getElementById('next-day').addEventListener('click', () => {
    selectedDate.setDate(selectedDate.getDate() + 1); // Move one day forward
    updateDatePicker(animated = true, timebound = false).then(() => {
        updateUrl();
        showDateGroup();
        setAvailability();
        location.reload(); // Rerender the page
    });
});

// Add event listener for the date picker input
document.getElementById('date-input').addEventListener('change', (e) => {
    const newDate = new Date(e.target.value);
    selectedDate = newDate;
    updateDatePicker(animated = false, timebound = false).then(() => {
        updateUrl();
        showDateGroup();
        setAvailability();
        location.reload(); // Rerender the page
    });
});

// Initialize the display and hidden input with today's date
updateDatePicker(animated = false, timebound = false).then(() => {
    updateUrl();
    showDateGroup();
    setAvailability();
});