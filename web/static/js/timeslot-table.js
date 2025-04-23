// Generate time slots dynamically
const startHour = 10; // 12 AM
const endHour = 22; // 10 PM (Last booking to finish at 11pm)

function formatHour(hour_value)  {
    return `${hour_value.toString().padStart(2, '0')}:00`;
}

for (let hour = startHour; hour <= endHour; hour++) {
    const formattedStartHour = formatHour(hour);
    const formattedEndHour = formatHour(hour + 1);
    if (document.getElementById("rec-room-timeslot-table")) {
        document.write(`
            <tr>
                <td>${formattedStartHour}-${formattedEndHour}</td>
                <td class="time-slot" data-room="Room 1" data-time="${hour}"></td>
                <td class="time-slot" data-room="Room 2" data-time="${hour}"></td>
                <td class="time-slot" data-room="Rec Room" data-time="${hour}"></td>
            </tr>
        `);
    } else {
        document.write(`
            <tr>
                <td>${formattedStartHour}-${formattedEndHour}</td>
                <td class="time-slot" data-room="Room 1" data-time="${hour}"></td>
                <td class="time-slot" data-room="Room 2" data-time="${hour}"></td>
            </tr>
        `);
    }
}
