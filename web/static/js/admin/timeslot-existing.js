// Function to select slots between start and end
function selectedSlots() {
    const originalBookingDate = document.getElementById('original-booking-date').textContent
    const newBookingDate = document.getElementById('date-input').value
    console.log(originalBookingDate === newBookingDate)
    if (originalBookingDate === newBookingDate) {
        const startTime = parseInt(document.getElementById('original-start-time').textContent, 10);
        const endTime = parseInt(document.getElementById('original-end-time').textContent, 10);
        const room = document.getElementById('room').value === 'Room 1' ? 'room1' : 'room2';

        timeSlots.forEach(slot => {
            const slotTime = parseInt(slot.dataset.time, 10);
            if (
                slot.dataset.room === room &&
                slotTime >= startTime &&
                slotTime < endTime
            ){
                slot.classList.remove('unavailable');
                slot.classList.add('current-booking');
            }
        });
    }   
}

// set the availability for the initial date
setAvailability().then(() => {
    selectedSlots();
});
