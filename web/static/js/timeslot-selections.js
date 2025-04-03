// Initialize variables for time slot selection
const timeSlots = document.querySelectorAll('.time-slot');
let startSlot = null;
let endSlot = null;
let selectedRoom = null;

// Function to clear the selection and reset the variables 
// This is needed so that users can't book more than one room or time-range
function clearSelection() {

    // Remove the selection class from all slots
    timeSlots.forEach(slot => {
        slot.classList.remove('selected','start-slot','end-slot','current-booking');
        slot.textContent = ''
    });

    // Reset the variables
    startSlot = null;
    endSlot = null;
    selectedRoom = null;


    // Reset hidden inputs
    if (document.getElementById('original-booking-date')) {
        document.getElementById('start-time').value = document.getElementById('original-start-time').textContent;
        document.getElementById('end-time').value = document.getElementById('original-end-time').textContent;
        document.getElementById('room').value = document.getElementById('original-room').textContent;
    } else {
        document.getElementById('start-time').value = null;
        document.getElementById('end-time').value = null;
        document.getElementById('room').value = null;
    }

    // Trigger the change event on the end-time input to update the price
    document.getElementById('end-time').dispatchEvent(new Event('change'));
    // Remove the enabled class from the book now button if exists
    const bookNowButton = document.getElementById('book-now');
    if (bookNowButton) {
        bookNowButton.classList.remove('enabled');
    }

    
    // set the availability for the initial date
    setAvailability().then(() => {
        if (document.getElementById('original-booking-date')) {
            selectedSlots();
        }

        if (!document.getElementById('payment-method')) {
            setBookableSlots();
        }
    });
}

function saveSlotSelections(start, end, room) {
    const startTimeValue = parseInt(start.dataset.time, 10);
    const endTimeValue = parseInt(end.dataset.time, 10);
    const formattedStartTime = formatHour(Math.min(startTimeValue, endTimeValue));
    const formattedEndTime = formatHour(Math.max(startTimeValue, endTimeValue) + 1);
    const formattedRoom = room === 'room1' ? 'Room 1' : 'Room 2';

    if (endTimeValue < startTimeValue) {
        end.textContent = formattedStartTime
        start.textContent = formattedEndTime
    } else if (endTimeValue === startTimeValue) {  
        end.textContent = `${formattedStartTime} - ${formattedEndTime}`;
    } else {
        start.textContent = formattedStartTime
        end.textContent = formattedEndTime
    }

    // Save the data to the hidden inputs
    document.getElementById('start-time').value = formattedStartTime;
    document.getElementById('end-time').value = formattedEndTime;
    document.getElementById('room').value = formattedRoom;

    // Trigger the change event on the end-time input to update the price
    document.getElementById('end-time').dispatchEvent(new Event('change'));

    // Add the beforeunload listener
    if (!window.beforeUnloadListenerAdded) {
        window.beforeUnloadListenerAdded = true; // Set the property to true after adding the listener
        window.addEventListener('beforeunload', handleBeforeUnload);
    }       

    updatePrice();
}

// Function for initialising the first time slot (used on odd clicks, 1, 3 etc...)
function selectFirstSlot(slot) {
    clearSelection();
    startSlot = slot;
    endSlot = slot;
    selectedRoom = slot.dataset.room;
    slot.classList.add('selected', 'start-slot');
    saveSlotSelections(startSlot, endSlot, selectedRoom);
}

// Function to select slots between start and end
function selectSlotsBetween(start, end) {
    const startTime = parseInt(start.dataset.time, 10);
    const endTime = parseInt(end.dataset.time, 10);
    const room = start.dataset.room;
    const filteredSlots = Array.from(timeSlots).filter(slot => {
        const slotTime = parseInt(slot.dataset.time, 10);
        return (
            slot.dataset.room === room &&
            slotTime >= Math.min(startTime, endTime) &&
            slotTime <= Math.max(startTime, endTime)
        );
    });

    const hasUnavailableSlots = filteredSlots.some(slot => slot.classList.contains('unavailable'));

    if (hasUnavailableSlots) {
        alert("can't book over already booked timeslots!");
        selectFirstSlot(start)
        return;
    }

    filteredSlots.forEach(slot => {
        slot.classList.add('selected');
    });
}

// Add event listeners for clicking on time slots
timeSlots.forEach(slot => {
    slot.addEventListener('click', () => {
        if (!slot.classList.contains('unavailable') || slot.classList.contains('current-booking')) {
            if (!startSlot) {
                // First click: set the start slot
                selectFirstSlot(slot);
                
                // Enable book now button 
                document.getElementById('book-now').classList.add('enabled');
            } else if (startSlot == endSlot && slot.dataset.room === selectedRoom) {
                // Second click: set the end slot and select all in between
                endSlot = slot;
                slot.classList.add('end-slot');
                selectSlotsBetween(startSlot, endSlot);
                saveSlotSelections(startSlot, endSlot, selectedRoom);
            } else {
                // Third click: clear the selection and start again
                selectFirstSlot(slot);

                // Enable book now button 
                document.getElementById('book-now').classList.add('enabled');
            }
        }
    });
});
    
// set the availability for the initial date
setAvailability().then(() => {
    if (!document.getElementById('payment-method')) {
        setBookableSlots();
    }
});
