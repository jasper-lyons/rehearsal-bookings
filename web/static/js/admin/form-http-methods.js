window.addEventListener('load', function () {
	let delete_forms = document.querySelectorAll('[method="DELETE"]')

	for (form of delete_forms) {
		form.addEventListener('submit', async function (e) {
			e.preventDefault();
			let response = await fetch(e.target.action, {
				method: 'DELETE'
			})

			if (response.redirected) {
				window.location.replace(response.url)
			}
		})
	}

	let paid_forms = document.querySelectorAll(".mark-paid")

	for (form of paid_forms) {
		form.addEventListener('submit', async function (e) {
			e.preventDefault();
			let response = await fetch(e.target.action, {
				method: 'PUT',
				body: JSON.stringify({
					status: "paid"
				}),
				headers: {
					'Content-Type': 'application/json'
				}

			})

			if (response.redirected) {
				window.location.replace(response.url)
			}
		})
	}

	let cancelation_forms = document.querySelectorAll(".mark-cancelled")

	for (form of cancelation_forms) {
		form.addEventListener('submit', async function (e) {
			e.preventDefault();
			let response = await fetch(e.target.action, {
				method: 'PUT',
				body: JSON.stringify({
					status: "cancelled"
				}),
				headers: {
					'Content-Type': 'application/json'
				}

			})

			if (response.redirected) {
				window.location.replace(response.url)
			}
		})
	}

})
