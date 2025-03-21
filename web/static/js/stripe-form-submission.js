window.addEventListener('load', function () {
	let stripe = Stripe(window.env.STRIPE_PUBLISHABLE_KEY);
	let elements;

	// Add event listeners for input changes
	let sessionType = document.getElementById('session-type');
	let startTime = document.getElementById('start-time');
	let endTime = document.getElementById('end-time');
	let cymbals = document.getElementById('cymbals');

	sessionType.addEventListener('change', updatePrice);
	startTime.addEventListener('change', updatePrice);
	endTime.addEventListener('change', updatePrice);
	cymbals.addEventListener('change', updatePrice);

	// Initial price update
	updatePrice();

	function updatePrice() {
		fetchPrice().then(function (price) {
			document.getElementById('stripe-submit').textContent = `Pay: £${price.toFixed(2)}`;
		});
	}

	let form = document.getElementById('form');
	form.addEventListener('submit', async function (e) {
		e.preventDefault();

		// Show the payment form
		document.getElementById('stripe-form').style.display = 'block';

		// Create held booking
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
				end_time: document.getElementById('end-time').value,
			}),
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (!bookingResponse.ok) {
			alert("Can't book " + document.getElementById('room').value + " at that time!");
			return;
		}

		let booking = await bookingResponse.json();

		// Get payment intent from server
		let intentResponse = await fetch('/stripe/payment-intents', {
			method: 'POST',
			body: JSON.stringify({
				booking_id: booking.id
			}),
			headers: {
				'Content-Type': 'application/json'
			}
		});

		let intent = await intentResponse.json();

		// Initialize Elements with client secret
		elements = stripe.elements({
			clientSecret: intent.client_secret,
			appearance: {
				theme: 'stripe'
			}
		});

		// Create and mount the Payment Element
		const paymentElement = elements.create('payment');
		paymentElement.mount('#payment-gateway');

		// Handle payment submission
		document.getElementById('stripe-submit').addEventListener('click', async function (e) {
			e.preventDefault()
			const result = await stripe.confirmPayment({
				elements,
				confirmParams: {
					return_url: window.location.origin + '/booking-confirmation',
				},
				redirect: 'if_required',
			});

			if (result.error) {
				// Show error to your customer
				console.error(result.error.message);
				alert('Payment failed: ' + result.error.message);
			} else {
				document.getElementById('stripe-form').style.display = 'none';
			}
		});
	});
});
