window.addEventListener('load', function () {
	for (let method of ["POST", "DELETE", "PUT"]) {
		let forms = document.querySelectorAll(`[method="${method}"]`)

		console.log(forms)

		for (form of forms) {
			let onsubmit = form.onsubmit
			form.onsubmit = undefined

			form.addEventListener('submit', async function (e) {
				e.preventDefault();

				if (onsubmit && onsubmit.call(this, e) === false)
					return

				let response = await fetch(e.target.action, {
					method: method,
					body: JSON.stringify(Object.fromEntries(new FormData(e.target))),
					headers: {
						'Content-Type': 'application/json'
					}
				})

				if (response.redirected) {
					window.location.replace(response.url)
				}
			})
		}
	}
})
