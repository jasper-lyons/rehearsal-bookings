async function sendForm(endPoint, content) {
	let bookingResponse = await fetch(endPoint, content)
	if (!bookingResponse.ok) {
		alert("Can't book ", document.getElementById('room').value, " at that time!")
		return
	}

	if (bookingResponse.ok) {
		document.getElementById('form-container').style.display = 'none';
		document.getElementById('success').style.display = 'block';
		// Update reload prompt event listener
		if (window.beforeUnloadListenerAdded) {
		window.removeEventListener('beforeunload', handleBeforeUnload);
		window.beforeUnloadListenerAdded = false; // Reset the flag
		}
		setTimeout(() => {
			window.location.href = '/admin/bookings';
		}, 1000);
	}

	let booking = await bookingResponse.json()

	return booking
}


function adminCreateBooking() {
	// check if all fields are filled
	let requiredFields = ['session-type', 'name', 'phone', 'room', 'date-input', 'start-time', 
						'status', 'payment-method']
	for (let field of requiredFields) {
		if (!document.getElementById(field).value) {
			alert('Please fill out all fields')
			return
		}
	}

	cymbals = 0
	if (document.getElementById('cymbals').checked) {
		cymbals = 1
	}		

	content = {
		method: 'POST',
		body: JSON.stringify({
			type: document.getElementById('session-type').value,
			name: document.getElementById('name').value,
			email: document.getElementById('email').value || '-',
			phone: document.getElementById('phone').value,
			room: document.getElementById('room').value,
			date: document.getElementById('date-input').value,
			start_time: document.getElementById('start-time').value,
			end_time: document.getElementById('end-time').value,
			cymbals: cymbals,
			revised_price: document.getElementById('revised-price').value,
			status: document.getElementById('status').value,
			payment_method: document.getElementById('payment-method').value,
			booking_notes: document.getElementById('booking-notes').value,
		}),
		headers: {
			'Content-Type': 'application/json'
		}
	}

	let endPoint = '/admin/bookings'

	return sendForm(endPoint, content)
}

function adminUpdateBooking() {
	// check if all fields are filled
	let requiredFields = ['session-type', 'customer-name', 'customer-phone', 'room',
						 'date-input', 'start-time', 'end-time', 'status', 'payment-method']
	for (let field of requiredFields) {
		if (!document.getElementById(field).value) {
			alert('Please fill out all fields')
			return
		}
	}

	cymbals = 0
	if (document.getElementById('cymbals').checked) {
		cymbals = 1
	}

	content = {
		method: 'PUT',
		body: JSON.stringify({
			type: document.getElementById('session-type').value,
			name: document.getElementById('customer-name').value,
			email: document.getElementById('customer-email').value || '-',
			phone: document.getElementById('customer-phone').value,
			room: document.getElementById('room').value,
			date: document.getElementById('date-input').value,
			start_time: document.getElementById('start-time').value,
			end_time: document.getElementById('end-time').value,
			cymbals: cymbals,
			revised_price: document.getElementById('revised-price').value,
			status: document.getElementById('status').value,
			payment_method: document.getElementById('payment-method').value,
			booking_notes: document.getElementById('booking-notes').value,
		}),
		headers: {
			'Content-Type': 'application/json'
		}
	}

	bookingId = document.getElementById('booking-id').textContent
	endPoint = `/admin/bookings/${bookingId}/update`

	return sendForm(endPoint, content)
}
