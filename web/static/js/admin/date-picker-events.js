// Add event listeners for buttons
document.getElementById('prev-day').addEventListener('click', () => {
                selectedDate.setDate(selectedDate.getDate() - 1); // Move one day back
        updateDatePicker(animated = true, timebound = false).then(() => {
            showDateGroup();
        });
});

document.getElementById('next-day').addEventListener('click', () => {
        selectedDate.setDate(selectedDate.getDate() + 1); // Move one day forward
        updateDatePicker(animated = true, timebound = false).then(() => {
                showDateGroup();
            });
});

// Add event listener for the date picker input
document.getElementById('date-input').addEventListener('change', (e) => {
        const newDate = new Date(e.target.value);
        selectedDate = newDate;
        updateDatePicker(animated = false, timebound = false).then(() => {
                showDateGroup();
            });
});

// Initialize the display and hidden input with today's date
updateDatePicker(animated = false, timebound = false).then(() => {
        showDateGroup();
});