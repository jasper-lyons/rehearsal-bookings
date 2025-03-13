// Initialize variables for time slot selection
const timeSlots = document.querySelectorAll('.time-slot');
let startSlot = null;
let endSlot = null;
let selectedRoom = null;

// Function to clear the selection and reset the variables 
// This is needed so that users can't book more than one room or time-range
// This function is also called by updateDatePicker()
function clearSelection() {
    // Remove the selection class from all slots
    timeSlots.forEach(slot => {
        slot.classList.remove('selected','start-slot','end-slot');
    });

    // Reset the variables
    startSlot = null;
    endSlot = null;
    selectedRoom = null;
    // Reset hidden inputs
    document.getElementById('start-time').value = null;
    document.getElementById('end-time').value = null;
    document.getElementById('room').value = null;

    // Trigger the change event on the end-time input to update the price
    document.getElementById('end-time').dispatchEvent(new Event('change'));
    // Remove the enabled class from the book now button if exists
    const bookNowButton = document.getElementById('book-now');
    if (bookNowButton) {
        bookNowButton.classList.remove('enabled');
    }
    setAvailability();
}

// Function for initialising the first time slot (used on odd clicks, 1, 3 etc...)
function selectFirstSlot(slot) {
    clearSelection();
    startSlot = slot;
    selectedRoom = slot.dataset.room;
    slot.classList.add('selected', 'grabbing', 'start-slot');
}

// Function to select slots between start and end
function selectSlotsBetween(start, end) {
    const startTime = parseInt(start.dataset.time, 10);
    const endTime = parseInt(end.dataset.time, 10);
    const room = start.dataset.room;

    timeSlots.forEach(slot => {
        const slotTime = parseInt(slot.dataset.time, 10);
        if (
            slot.dataset.room === room &&
            slotTime >= Math.min(startTime, endTime) &&
            slotTime <= Math.max(startTime, endTime)
        ) {
            slot.classList.add('selected');
        }
    });
}

// Add event listeners for clicking on time slots
timeSlots.forEach(slot => {
    slot.addEventListener('click', () => {
        if (!slot.classList.contains('unavailable')) {
            if (!startSlot) {
                // First click: set the start slot
                selectFirstSlot(slot);
            } else if (!endSlot && slot.dataset.room === selectedRoom) {
                // Second click: set the end slot and select all in between
                endSlot = slot;
                slot.classList.add('end-slot');
                selectSlotsBetween(startSlot, endSlot);

                const startTimeValue = parseInt(startSlot.dataset.time, 10);
                const endTimeValue = parseInt(endSlot.dataset.time, 10);
                const formattedStartTime = formatHour(Math.min(startTimeValue, endTimeValue));
                const formattedEndTime = formatHour(Math.max(startTimeValue, endTimeValue) + 1);
                const formattedRoom = selectedRoom === 'room1' ? 'Room 1' : 'Room 2';

                // Save the data to the hidden inputs
                document.getElementById('start-time').value = formattedStartTime;
                document.getElementById('end-time').value = formattedEndTime;
                document.getElementById('room').value = formattedRoom;

                // Trigger the change event on the end-time input to update the price
                document.getElementById('end-time').dispatchEvent(new Event('change'));
                // Enable book now button if exists
                const bookNowButton = document.getElementById('book-now');
                if (bookNowButton) {
                    bookNowButton.classList.add('enabled');
                }                // print the selection for user to see room/time details
                updatePrice();
                populateSummary();

                timeSlots.forEach(slot => slot.classList.remove('grabbing'));
            } else {
                // Third click: clear the selection and start again
                selectFirstSlot(slot);
            }
        }
    });

    // On mouse over, if dragging, add the slots affected to the selectedSlots array
    slot.addEventListener('mouseover', () => {
        if (startSlot && !endSlot) {
            if (selectedRoom === slot.dataset.room) {
                slot.classList.add('grabbing');						
            }
        }
    });
});

// Populate Summary function
function populateSummary() {
    // Get the values from the form
    const room = document.getElementById('room').value;
    const date = document.getElementById('date-input').value;
    const startTime = document.getElementById('start-time').value;
    const endTime = document.getElementById('end-time').value;
    const price = document.getElementById('price').textContent;

    // Update both Step 2 and Step 3 summaries using classes
    const summaryRooms = document.querySelectorAll('.summary-room');
    const summaryDate= document.querySelectorAll('.summary-date');
    const summaryTimes = document.querySelectorAll('.summary-time');
    const summaryPrices = document.querySelectorAll('.summary-price');

    // Update the summary content for all elements with these classes
    summaryRooms.forEach(element => element.innerHTML = `<p>Room: <strong>${room}</strong></p>`);
    summaryDate.forEach(element => element.innerHTML = `<p>Date: <strong>${date}</strong></p>`);
    summaryTimes.forEach(element => element.innerHTML = `<p>Time: <strong>${startTime} - ${endTime}</strong></p>`);
    summaryPrices.forEach(element => element.innerHTML = `<p>Price: <strong>${price}</strong></p>`);
};

// set the availability for the initial date
setAvailability();
