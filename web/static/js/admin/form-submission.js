window.addEventListener('load', function () {
	let form = document.getElementById('form')
	form.addEventListener('submit', async function (e) {
		e.preventDefault()
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
		let requiredFields = ['session-type', 'name','email', 'phone', 'room', 'date-input', 'start-time', 'duration']
		for (let field of requiredFields) {
			if (!document.getElementById(field).value) {
				alert('Please fill out all fields')
				return
			}
			console.log(document.getElementById(field).value)
		}
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
	})
})
