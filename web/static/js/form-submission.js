window.addEventListener('load', function () {
    let form = document.getElementById('form')
    form.addEventListener('submit', async function (e) {
        e.preventDefault()

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
        SumUpCard.mount({
            id: 'sumup-card',
            checkoutId: checkout.id,
            onResponse: async function (type, body) {
                switch (type) {
                    case "sent":
                        break;
                    case "success":
                        // TODO: Confirm booking
												if (body.status === 'SUCCESS') {
													let confirmBookingResponse = await fetch(`/bookings/${booking.id}/confirm`, {
														method: 'POST',
														body: JSON.stringify(body),
														headers: {
															'Content-Type': 'application/json'
														}
													})

													// TODO: Display success, perhaps redirect?
													document.getElementById('sumup-card').style.display = "none"
													document.getElementById('success').style.display = "block"
												} else {
													// do nothing
												}
                        break;
                    case "error":
                        // TODO: Handle error
                        break;
                }
            }
        })
    })
})
