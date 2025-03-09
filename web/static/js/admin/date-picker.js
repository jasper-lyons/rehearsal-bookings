// Initialize the current date and define the range
const currentDate = new Date();

// define the selectedDate - by default this is the current date
let selectedDate = new Date(currentDate);

// Update the displayed date and button states
function updateDatePicker(animated = false) {
    const datePicker = document.getElementById('date-input');
    datePicker.value = selectedDate.toLocaleDateString('en-CA');
    const dateGroups = document.querySelectorAll(".date-group");

    // Show the correct date group
    dateGroups.forEach(group => {
        if (group.getAttribute("data-date") === datePicker.value) {
            group.style.display = "block";
        } else {
            group.style.display = "none";
        }
    });

    // Trigger animation
    if (animated) {
        datePicker.classList.add('animate');
        setTimeout(() => datePicker.classList.remove('animate'), 150); // Match the animation duration
    }
}

// Add event listeners for buttons
document.getElementById('prev-day').addEventListener('click', () => {
        selectedDate.setDate(selectedDate.getDate() - 1); // Move one day back
        updateDatePicker(animated = true);
});

document.getElementById('next-day').addEventListener('click', () => {
        selectedDate.setDate(selectedDate.getDate() + 1); // Move one day forward
        updateDatePicker(animated = true);
});

// Add event listener for the date picker input
document.getElementById('date-input').addEventListener('change', (e) => {
    const newDate = new Date(e.target.value);
        selectedDate = newDate;
        updateDatePicker();
});

// Initialize the display and hidden input with today's date
document.addEventListener("DOMContentLoaded", function() {
    updateDatePicker();
});