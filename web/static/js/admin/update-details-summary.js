function updatedDetails() {
    let summary = document.getElementById('updated-details')
    const room = document.getElementById('room').value;
    const date = document.getElementById('date-input').value;
    const startTime = document.getElementById('start-time').value;
    const endTime = document.getElementById('end-time').value;
    const type = document.getElementById('session-type').value;
    const cymbals = document.getElementById('cymbals').checked? 'Yes' : 'No';
    const rev_price = document.getElementById('revised-price').value || document.getElementById('price').textContent.replace('£','');
    const rev_status = document.getElementById('status').value;
    const customerName = document.getElementById('name').value;
    const customerPhone = document.getElementById('phone').value;
    const customerEmail = document.getElementById('email').value;

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