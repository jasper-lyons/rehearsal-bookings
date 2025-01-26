// Generate time slots dynamically
const startHour = 10; // 12 AM
const endHour = 22; // 10 PM (Last booking to finish at 11pm)

function formatHour(hour_value)  {
    return `${hour_value.toString().padStart(2, '0')}:00`;
}

for (let hour = startHour; hour <= endHour; hour++) {
    const formattedHour = formatHour(hour);
    document.write(`
        <tr>
            <td>${formattedHour}</td>
            <td class="time-slot" data-room="room1" data-time="${hour}"></td>
            <td class="time-slot" data-room="room2" data-time="${hour}"></td>
        </tr>
    `);
}
