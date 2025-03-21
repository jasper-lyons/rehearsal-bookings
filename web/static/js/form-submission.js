async function formSubmission(admin=false, update=false) {
	let form = document.getElementById('form')
	form.addEventListener('submit', async function (e) {
		e.preventDefault()

		if (admin) {
			console.log("admin mode")
			// check if all fields are filled
			const name = document.getElementById('name').value.trim();
			const email = document.getElementById('email').value.trim();
			const phone = document.getElementById('phone').value.trim();
		
			let nameRegex = /^[a-zA-Z]+([-' ][a-zA-Z]+)*\s+[a-zA-Z]+([-' ][a-zA-Z]+)*$/;
			let emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
			let phoneRegex = /^\+?\d{11}$/;
		
			if (!nameRegex.test(name))  {
				alert('Please provide your full name.');
				return false;
			}
		
			if (!emailRegex.test(email)) {
				alert('Please enter a valid email address.');
				return false;
			}

			if (!phoneRegex.test(phone)) {
				alert('Please enter a valid phone number.');
				return false;
			}

			// check if all fields are filled
			let requiredFields = ['session-type', 'name','email', 'phone', 'room', 'date-input', 'start-time']
			for (let field of requiredFields) {
				if (!document.getElementById(field).value) {
					alert('Please fill out all fields')
					return
				}
				console.log(document.getElementById(field).value)
			}
		}

		cymbals = 0
		if (document.getElementById('cymbals').checked) {
			cymbals = 1
		}

		let endPoint = '/bookings'
		if (update) {
			bookingId = document.getElementById('booking-id').textContent
			endPoint = `/admin/bookings/${bookingId}/update`
			content = {
				method: 'PUT',
				body: JSON.stringify({
					type: document.getElementById('session-type').value,
					name: document.getElementById('name').value,
					email: document.getElementById('email').value,
					phone: document.getElementById('phone').value,
					room: document.getElementById('room').value,
					date: document.getElementById('date-input').value,
					start_time: document.getElementById('start-time').value,
					end_time: document.getElementById('end-time').value,
					cymbals: cymbals,
					price: document.getElementById('revised-price').value,
					status: document.getElementById('status').value,
				}),
				headers: {
					'Content-Type': 'application/json'
				}
			}
		} else {
			content = {
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
					cymbals: cymbals,
				}),
				headers: {
					'Content-Type': 'application/json'
				}
			}
		}
		
		let bookingResponse = await fetch(endPoint, content)
		if (!bookingResponse.ok) {
			alert("Can't book ", document.getElementById('room').value, " at that time!")
			return
		}

		if (bookingResponse.ok && admin) {
			document.getElementById('form-container').style.display = 'none';
			document.getElementById('success').style.display = 'block';
			setTimeout(() => {
				location.reload();
			}, 1500);
		}

		let booking = await bookingResponse.json()

		return booking
	})
}
