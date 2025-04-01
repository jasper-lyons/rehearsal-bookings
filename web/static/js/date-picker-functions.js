// Initialize the current date and define the range
const currentDate = new Date();
const minDate = new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate()); // Start of the range (today)
const maxDate = new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate()); // End of the range (60 days from today)
maxDate.setDate(maxDate.getDate() + 60);

// define the selectedDate - by default this is the current date
let selectedDate = new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate()); 

if (document.getElementById('original-booking-date')) {
    selectedDate = new Date(document.getElementById('original-booking-date').textContent);
}

const datePicker = document.getElementById('date-input');

// Update the displayed date and button states
async function updateDatePicker(animated = false, timebound = true) {
    if (timebound) {
        datePicker.min = minDate.toLocaleDateString('en-CA');
        datePicker.max = maxDate.toLocaleDateString('en-CA');
        datePicker.value = selectedDate.toLocaleDateString('en-CA');
        document.getElementById('prev-day').disabled = selectedDate <= minDate;
        document.getElementById('next-day').disabled = selectedDate >= maxDate;
    } else {
        datePicker.value = selectedDate.toLocaleDateString('en-CA');
        document.getElementById('prev-day').disabled = false;
        document.getElementById('next-day').disabled = false;
    }

    // Trigger animation
    if (animated) {
        datePicker.classList.add('animate');
        setTimeout(() => datePicker.classList.remove('animate'), 150); // Match the animation duration   
    }
}


// function to show the correct date group based on selectedDate
function showDateGroup() {
    const dateGroups = document.querySelectorAll(".date-group");

    // Show the correct date group
    dateGroups.forEach(group => {
        if (group.getAttribute("data-date") === datePicker.value) {
            group.style.display = "block";
        } else {
            group.style.display = "none";
        }
    });

}