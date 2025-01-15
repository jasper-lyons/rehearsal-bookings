// Initialize the current date and define the range
const currentDate = new Date();
const minDate = new Date(currentDate); 	// Start of the range (today)
const maxDate = new Date(currentDate); // End of the range (90 days from today)
maxDate.setDate(maxDate.getDate() + 90);

// define the selectedDate - by default this is the current date
let selectedDate = new Date(currentDate);

// Update the displayed date and button states
function updateDatePicker() {
    const datePicker = document.getElementById('date-input');

    datePicker.min = minDate.toLocaleDateString('en-CA');
    datePicker.max = maxDate.toLocaleDateString('en-CA');
    datePicker.value = selectedDate.toLocaleDateString('en-CA');
    document.getElementById('prev-day').disabled = selectedDate <= minDate;
    document.getElementById('next-day').disabled = selectedDate >= maxDate;
}

// Add event listeners for buttons
document.getElementById('prev-day').addEventListener('click', () => {
    if (selectedDate > minDate) {
        selectedDate.setDate(selectedDate.getDate() - 1); // Move one day back
        updateDatePicker();
    }
});

document.getElementById('next-day').addEventListener('click', () => {
    if (selectedDate < maxDate) {
        selectedDate.setDate(selectedDate.getDate() + 1); // Move one day forward
        updateDatePicker();
    }
});

// Add event listener for the date picker input
document.getElementById('date-input').addEventListener('change', (e) => {
    const newDate = new Date(e.target.value);
    if (newDate >= minDate && newDate <= maxDate) {
        selectedDate = newDate;
        updateDatePicker();
    }
});

// Initialize the display and hidden input with today's date
updateDatePicker();
