window.addEventListener('load', function () {
	let form = document.getElementById('form')
	form.addEventListener('submit', async function (e) {
		e.preventDefault()

		document.getElementById('payment-gateway').style.display = 'block'

		let sumup = SumUpCard.mount({
			id: 'payment-gateway',
			email: document.getElementById('email').value,
			onResponse: async function (type, body) {
				switch (type) {
					case "sent":
						break;
					case "success":
						let confirmBookingResponse = await fetch(`/bookings/${booking.id}/confirm`, {
							method: 'POST',
							body: JSON.stringify(body),
							headers: {
								'Content-Type': 'application/json'
							}
						})

						if (200 <= confirmBookingResponse.status && confirmBookingResponse.status < 300) {
							// TODO: Display success, perhaps redirect?
							document.getElementById('payment-gateway').style.display = "none"
							document.getElementById('success').style.display = "block"
						} else {
							alert("Payment failed!")
						}
						break;
					case "error":
						// TODO: Handle error
						break;
				}
			}
		})

		// create held booking
		let bookingResponse = await fetch('/bookings', {
			method: 'POST',
			body: JSON.stringify({
				type: document.getElementById('session-type').value,
				name: document.getElementById('name').value,
				email: document.getElementById('email').value,
				phone: document.getElementById('phone').value,
				room: document.getElementById('room').value,
				date: document.getElementById('date-input').value,
				start_time: document.getElementById('start-time').value,
				duration: parseInt(document.getElementById('duration').value, 10),
			}),
			headers: {
				'Content-Type': 'application/json'
			}
		})

		if (!bookingResponse.ok) {
			alert("Can't book ", document.getElementById('room').value, " at that time!")
			return
		}

		let booking = await bookingResponse.json()

		// create charge
		let checkoutResponse = await fetch('/sumup/checkouts', {
			method: 'POST',
			body: JSON.stringify({
				amount: parseFloat(document.getElementById('price').textContent.replace('£', '')),
				checkout_reference: `booking-${booking.id}`
			}),
			headers: {
				'Content-Type': 'application/json'
			}
		})

		if (!checkoutResponse.ok) {
			alert("Payment provider failed, please reach out!")
			return
		}

		let checkout = await checkoutResponse.json()

		// open sumup UI
		sumup.update({
			checkoutId: checkout.id
		})
	})
})
