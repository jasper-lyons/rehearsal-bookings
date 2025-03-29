const options = { weekday: 'short', day: '2-digit', month: 'short', year: 'numeric' };

// Add event listeners for buttons
document.getElementById('prev-day').addEventListener('click', () => {
    if(document.getElementById('original-booking-date')) {
        selectedDate.setDate(selectedDate.getDate() - 1);
        updateDatePicker(animated = true, timebound = false);
    } else {
        if (selectedDate > minDate) {
            selectedDate.setDate(selectedDate.getDate() - 1); // Move one day back
            updateDatePicker(animated = true);
        }
    }
    document.getElementById('formatted-date').textContent = new Intl.DateTimeFormat('en-GB', options).format(new Date(datePicker.value));
    clearSelection(); // clear the time slot selection when date changes
});

document.getElementById('next-day').addEventListener('click', () => {
    if(document.getElementById('original-booking-date')) {
        selectedDate.setDate(selectedDate.getDate() + 1);
        updateDatePicker(animated = true, timebound = false);
    } else {
        if (selectedDate < maxDate) {
            selectedDate.setDate(selectedDate.getDate() + 1); // Move one day back
            updateDatePicker(animated = true);
        }
    }
    document.getElementById('formatted-date').textContent = new Intl.DateTimeFormat('en-GB', options).format(new Date(datePicker.value));
    clearSelection(); // clear the time slot selection when date changes
    
});

// Add event listener for the date picker input
document.getElementById('date-input').addEventListener('change', (e) => {
    const newDate = new Date(e.target.value);
    if(document.getElementById('original-booking-date')) {
        selectedDate = newDate;
        updateDatePicker(animated = false, timebound = false);
    } else {
        if (newDate >= minDate && newDate <= maxDate) {
            selectedDate = newDate;
            updateDatePicker();
        } else if (newDate > maxDate) {
            alert("Bookings can only be made up to 60 days in advance!")
            selectedDate = new Date(maxDate);
            updateDatePicker();
        }
        else if (newDate < minDate) {
            alert("Please select a date in the future!")
            selectedDate = new Date(minDate);
            updateDatePicker();
        }
    }
    document.getElementById('formatted-date').textContent = new Intl.DateTimeFormat('en-GB', options).format(new Date(datePicker.value));
    clearSelection(); // clear the time slot selection when date changes
});

// Initialize the display and hidden input with today's date
if(document.getElementById('original-booking-date')) {
    updateDatePicker(animated = false, timebound = false);
} else {
    updateDatePicker();
}
