window.addEventListener('load', function () {
	let forms = document.querySelectorAll('[method="DELETE"]')

	for (form of forms) {
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
})
