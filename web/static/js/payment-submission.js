window.addEventListener('load', function () {
	let form = document.getElementById('form')
	form.addEventListener('submit', async function (e) {
		e.preventDefault()

		let sumup = SumUpCard.mount({
			id: 'sumup-card',
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
							document.getElementById('sumup-card').style.display = "none"
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

		let booking = formSubmission();

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
