// Populate Summary function
async function populateSummary() {
    // Get the values from the form
    const room = document.getElementById('room').value;
    const session_type = document.getElementById('session-type').value;
    const dateInput = document.getElementById('date-input').value;
    const formattedDate = new Date(dateInput).toLocaleDateString('en-GB', {
        weekday: 'long',
        day: '2-digit',
        month: 'long',
        year: 'numeric'
    }).replace(/,/g, '');
    const startTime = document.getElementById('start-time').value;
    const endTime = document.getElementById('end-time').value;
    const cymbals = document.getElementById('cymbals').checked? 'Yes' : 'No';
    const price = document.getElementById('price').textContent;

    // Update both Step 2 and Step 3 summaries using classes
    const summaryRooms = document.querySelectorAll('.summary-room');
    const summaryTypes = document.querySelectorAll('.summary-type');
    const summaryDate= document.querySelectorAll('.summary-date');
    const summaryTimes = document.querySelectorAll('.summary-time');
    const summaryCymbals = document.querySelectorAll('.summary-cymbals');
    const summaryPrices = document.querySelectorAll('.summary-price');


    // Update the summary content for all elements with these classes
    summaryRooms.forEach(element => element.innerHTML = `<p>Room: <strong>${room}</strong></p>`);
    summaryTypes.forEach(element => element.innerHTML = `<p>Session Type: <strong>${session_type.charAt(0).toUpperCase() + session_type.slice(1)}</strong></p>`);
    summaryDate.forEach(element => element.innerHTML = `<p>Date: <strong>${formattedDate}</strong></p>`);
    summaryTimes.forEach(element => element.innerHTML = `<p>Time: <strong>${startTime} - ${endTime}</strong></p>`);
    summaryCymbals.forEach(element => element.innerHTML = `<p>Cymbal Hire: <strong>${cymbals}</strong></p>`);
    summaryPrices.forEach(element => element.innerHTML = `<p>Price: <strong>${price}</strong></p>`);
};

async function populateCustomerInfo() {
    // const name = document.getElementById('name').value;
    const email = document.getElementById('email').value;
    const phone = document.getElementById('phone').value;
    // const summaryName = document.querySelectorAll('.');
    const summaryEmail = document.getElementById('customer-email');
    const summaryPhone = document.getElementById('customer-phone');

    // summaryName.forEach(element => element.textContent = `Name: ${name}`);
    summaryEmail.innerHTML = `<p>An email confirmation will be sent to: <strong>${email}</strong></p>`;
    summaryPhone.innerHTML = `<p>On the day, an SMS with the access codes will be sent to: <strong>${phone}</strong></p>`;
}

function adminUpdatedDetails() {
    let summary = document.getElementById('updated-details')
    const room = document.getElementById('room').value;
    const date = document.getElementById('date-input').value;
    const startTime = document.getElementById('start-time').value;
    const endTime = document.getElementById('end-time').value;
    const type = document.getElementById('session-type').value;
    const cymbals = document.getElementById('cymbals').checked? 'Yes' : 'No';
    const rev_price = document.getElementById('revised-price').value || document.getElementById('price').textContent.replace('£','');
    const rev_status = document.getElementById('status').value;
    const customerName = document.getElementById('customer-name').value;
    const customerPhone = document.getElementById('customer-phone').value;
    const customerEmail = document.getElementById('customer-email').value;

    summary.innerHTML = `
        <div class="admin-panel">
            <h2>Booking Details:</h2>
            <p><strong>Room: </strong>${room}</p>
            <p><strong>Date: </strong>${date}</p>
            <p><strong>Booking Time: </strong>${startTime} - ${endTime}</p>
            <p><strong>Session Type: </strong>${type}</p>
            <p><strong>Cymbals: </strong>${cymbals}</p>	
            <p><strong>Price: £</strong>${rev_price}</p>
            <p><strong>Status: </strong>${rev_status}</p>
        </div>
        <div class="admin-panel">
            <h2>Customer Details:</h2>
            <p><strong>Name: </strong>${customerName}</p>
            <p><strong>Phone: </strong>${customerPhone}</p>
            <p><strong>Email: </strong>${customerEmail}</p>
        </div>
        `;
}