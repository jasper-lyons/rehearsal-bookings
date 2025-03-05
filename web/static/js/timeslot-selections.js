// Initialize variables for time slot selection
const timeSlots = document.querySelectorAll('.time-slot');
let startSlot = null;
let endSlot = null;
let selectedRoom = null;

// Function to set the availability of each timeslot. 
// This function is called by updateDatePicker() and selectButton()
function setAvailability() {
    const datePicker = document.getElementById('date-input');
    timeSlots.forEach(slot => {
        // Reset the availability
        slot.classList.remove('unavailable');
        // If the selected date is today and the slot is in the past or within the next 2 hours, disable it
        if (datePicker.value === datePicker.min && slot.dataset.time < new Date().getHours() + 2) {
            slot.classList.add('unavailable');
        }
        // If the selected date is a weekday and the slot is before 12pm, disable it
        if (isWeekday && slot.dataset.time < 12) {
            slot.classList.add('unavailable');
        }
        // If the selected date is a weekday and the session type is solo and the slot is after 5pm, disable it
        if (isWeekday && document.getElementById('session-type').value === "solo" && slot.dataset.time > 17) {
            slot.classList.add('unavailable');
        }
    });
}

// Function to clear the selection and reset the variables 
// This is needed so that users can't book more than one room or time-range
// This function is also called by updateDatePicker()
function clearSelection() {
    // Remove the selection class from all slots
    timeSlots.forEach(slot => {
        slot.classList.remove('selected');
    });

    // Reset the variables
    startSlot = null;
    endSlot = null;
    selectedRoom = null;
    timeslot_output.innerHTML = null;

    // Reset hidden inputs
    document.getElementById('start-time').value = null;
    document.getElementById('end-time').value = null;
    document.getElementById('room').value = null;
    document.getElementById('duration').value = null;

    // Trigger the change event on the duration input to update the price
    document.getElementById('duration').dispatchEvent(new Event('change'));
    // Remove the enabled class from the book now button
    document.getElementById('book-now').classList.remove('enabled')
}

// Function for initialising the first time slot (used on odd clicks, 1, 3 etc...)
function selectFirstSlot(slot) {
    clearSelection();
    startSlot = slot;
    selectedRoom = slot.dataset.room;
    slot.classList.add('selected', 'grabbing');

    // Disable the book now button until a valid selection is made
    document.getElementById('book-now').classList.remove('enabled');
    document.getElementById('duration').value = 0;
    document.getElementById('duration').dispatchEvent(new Event('change'));
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
                document.getElementById('duration').value = Math.abs(startTimeValue - endTimeValue) + 1;

                // Trigger the change event on the duration input to update the price
                document.getElementById('duration').dispatchEvent(new Event('change'));
                document.getElementById('book-now').classList.add('enabled')
                // print the selection for user to see room/time details
                timeslot_output.innerHTML = `
                    <h2>YOUR SELECTION</h2>
                    <h3><strong>${formattedRoom}</strong>: ${document.getElementById('date-input').value}, ${formattedStartTime} - ${formattedEndTime} </h3>
                `;

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

// set the availability for the initial date
setAvailability();
