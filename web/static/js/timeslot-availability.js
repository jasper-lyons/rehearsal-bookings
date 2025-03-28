async function fetchAvailability(date) {
    try {
        const response = await fetch(`/rooms?day=${date}`);
        if (!response.ok) throw new Error('Failed to fetch availability');
        const data = await response.json();
        return data.rooms; // Return the rooms array
    } catch (error) {
        console.error(error);
        return [];
    }
}

// Function to set the availability of each timeslot based on the results of rooms API
async function setAvailability() {
    const timeSlots = document.querySelectorAll('.time-slot');
    const datePicker = document.getElementById('date-input');
    const rooms = await fetchAvailability(datePicker.value); // ✅ Await the result
    const room1 = rooms.find(r => r.name === "Room 1"); // ✅ Find Room 1
    if (!room1) {
        console.error('Room 1 not found!');
        return;
    }

    const room2 = rooms.find(r => r.name === "Room 2"); // ✅ Find Room 1
    if (!room2) {
        console.error('Room 2 not found!');
        return;
    }

    timeSlots.forEach(slot => {
        slot.classList.remove('unavailable');
        const slotTime = slot.dataset.time;
        const slotRoom = slot.dataset.room;
        const timeLabel = formatHour(slotTime);

        // Disable slot if the room's availability is false
        if (slotRoom === "room1" && !room1.availability[timeLabel]) {
            slot.classList.add('unavailable');
            return;
        }

        if (slotRoom === "room2" && !room2.availability[timeLabel]) {
            slot.classList.add('unavailable');
            return;
        }

    });
}

// Function to set the availability of each timeslot based on the business logic
// used on the user side of bookings, rather than admin
async function setBookableSlots() {
    const timeSlots = document.querySelectorAll('.time-slot');
    const datePicker = document.getElementById('date-input');

    let date = new Date(datePicker.value);
    let isWeekday = date.getDay() >= 1 && date.getDay() <= 5; // Monday = 1, Friday = 5

    timeSlots.forEach(slot => {
        const slotTime = slot.dataset.time;

        if (datePicker.value === datePicker.min && slotTime < new Date().getHours() + 2) {
            slot.classList.add('unavailable');
        }

        // If the selected date is a weekday and the slot is before 12pm, disable it
        if (isWeekday && slotTime < 12) {
            slot.classList.add('unavailable');
        }

        // If the selected date is a weekday and the session type is solo and the slot is after 5pm, disable it
        if (isWeekday && document.getElementById('session-type').value === "solo" && slotTime > 17) {
            slot.classList.add('unavailable');
        }
    });
}
